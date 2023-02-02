package timewheel

import (
	"cmgo/lib/logger"
	"container/list"
	"log"
	"time"

	ants "github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

// 用于窗口类型数据更新

// Job 延时任务回调函数
type Job func(int, int, interface{})

// BackupJob 备份回调函数
type BackupJob func(int, []string)

// TaskData 回调函数参数类型

// TimeWheel 时间轮
type TimeWheel struct {
	interval time.Duration // 指针每隔多久往前移动一格
	ticker   *time.Ticker
	slots    []*list.List // 时间轮槽
	// key: 定时器唯一标识 value: 定时器所在的槽, 主要用于删除定时器, 不会出现并发读写，不加锁直接访问
	timer          map[string]int64 //改为int64代表第一次进入时间轮时key所在msg的时间戳(ms)，原先是int代表key所在slot的下标
	totalNum       int32
	currentPos     int       // 当前指针指向哪一个槽
	slotNum        int       // 槽数量
	job            Job       // 定时器回调函数
	addTaskChannel chan Task // 新增任务channel
	stopChannel    chan bool // 停止定时器channel
	backupTicker   *time.Ticker
	backupInterval time.Duration
	backupJob      BackupJob //备份定时器回调函数

	idx           int //标识是哪个时间轮
	goroutinePool *ants.PoolWithFunc
}

// Task 延时任务
type Task struct {
	delay      time.Duration // 延迟时间
	circle     int           // 时间轮需要转动几圈
	key        string        // 定时器唯一标识, 用于删除定时器
	data       interface{}   // 回调函数参数
	timestamp  int64         // 加入时间轮时key所在消息的时间戳
	windowSize int64         // 该key对应特征组窗口大小
}

type poolFuncArg struct {
	Idx        int
	CurrentPos int
	Data       interface{}
	JobFunc    Job
}

func poolFunc(i interface{}) {
	poolFuncArg := i.(poolFuncArg)
	poolFuncArg.JobFunc(poolFuncArg.Idx, poolFuncArg.CurrentPos, poolFuncArg.Data)
}

// New 创建时间轮
func New(interval, backupInterval time.Duration, slotNum, poolSize int, job Job, backupJob BackupJob, satisfySuffixs []uint64) *TimeWheel {
	if interval <= 0 || slotNum <= 0 || job == nil {
		return nil
	}
	tw := &TimeWheel{
		interval:       interval,
		backupInterval: backupInterval,
		slots:          make([]*list.List, slotNum),
		timer:          make(map[string]int64),
		currentPos:     0,
		job:            job,
		backupJob:      backupJob,
		slotNum:        slotNum,
		addTaskChannel: make(chan Task),
		stopChannel:    make(chan bool),
	}
	pool, err := ants.NewPoolWithFunc(poolSize, poolFunc, ants.WithNonblocking(true))
	if err != nil {
		log.Fatalf("newPoolWithFunc faile, err:%v", err)
	}
	tw.goroutinePool = pool

	tw.initSlots()
	return tw
}

// 初始化槽，每个槽指向一个双向链表
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

// Start 启动时间轮
func (tw *TimeWheel) Start(idx int) {
	tw.idx = idx
	tw.ticker = time.NewTicker(tw.interval)
	tw.backupTicker = time.NewTicker(tw.backupInterval)
	tw.start()
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	tw.goroutinePool.Release()
	tw.stopChannel <- true
}

// AddTimer 添加定时器 key为定时器唯一标识
func (tw *TimeWheel) AddTimer(delay time.Duration, key string, timestamp, windowSize int64, data interface{}) {
	if delay < 0 {
		return
	}
	tw.addTaskChannel <- Task{delay: delay, key: key, timestamp: timestamp, windowSize: windowSize, data: data}
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case <-tw.backupTicker.C:
			tw.backupTickHandler(false)
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			tw.backupTicker.Stop()
			tw.backupTickHandler(true)
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l, tw.currentPos)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

func (tw *TimeWheel) backupTickHandler(isStopSignal bool) {
	keys := []string{}
	for key := range tw.timer {
		keys = append(keys, key)
	}
	//为了不阻塞正常时间轮流程，下边起子协程处理
	if !isStopSignal {
		go tw.backupJob(tw.idx, keys)
	} else {
		tw.backupJob(tw.idx, keys)
	}

}

// 扫描链表中过期定时器, 并执行回调函数
func (tw *TimeWheel) scanAndRunTask(l *list.List, currentPos int) {
	tasks := make([]*Task, 0)
	datas := make([]interface{}, 0)
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}
		if time.Now().UnixNano()/1e6-task.timestamp < task.windowSize { //单位都是ms，当前时间距离最老消息间隔小于窗口，不需要推动
			e = e.Next()
			continue
		}
		datas = append(datas, task.data)
		tasks = append(tasks, task)
		next := e.Next()
		l.Remove(e)
		if task.key != "" {
			delete(tw.timer, task.key)
			tw.totalNum -= 1
		}
		e = next
	}
	if len(datas) <= 0 {
		return
	}
	err := tw.goroutinePool.Invoke(poolFuncArg{Idx: tw.idx, CurrentPos: currentPos, Data: datas, JobFunc: tw.job})
	if err != nil {
		//非阻塞模式，invoke失败，需要再将这些任务塞进去
		go func() {
			for i := range tasks {
				tw.AddTimer(time.Minute, tasks[i].key, tasks[i].timestamp, tasks[i].windowSize, tasks[i].data)
			}
		}()
		logger.Error("timewheel invoke failed", zap.Error(err))
	}
}

// 新增任务到链表中
func (tw *TimeWheel) addTask(task *Task) {
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	// 时间轮是特征组维度，更新单个特征不应该刷新时间轮。只有不存在时才会加入新任务
	if task.key != "" {
		if oriTimestamp, ok := tw.timer[task.key]; !ok {
			tw.slots[pos].PushBack(task)
			tw.timer[task.key] = task.timestamp
			tw.totalNum += 1
		} else {
			if task.timestamp < oriTimestamp {
				tw.timer[task.key] = task.timestamp //如果有更老的消息过来，就更新时间轮中key的时间戳
			}
		}
	}
}

// 获取定时器在槽中的位置, 时间轮需要转动的圈数
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle = int(delaySeconds / intervalSeconds / tw.slotNum)
	pos = int(tw.currentPos+delaySeconds/intervalSeconds) % tw.slotNum

	return
}

func (tw *TimeWheel) Len() int32 {
	return tw.totalNum
}
