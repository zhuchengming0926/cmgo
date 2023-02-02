package map_use

import (
	"time"
)

// 内存map使用，用户并发读写，通过管道来解除并发问题
type Processor struct {
	keyMap           map[string]int64      // 数据
	backupTicker     *time.Ticker          // 备份ticker
	backupInterval   time.Duration         // 多久备份一次
	backupJob        BackupJob             // 备份回调函数，因为需要用到groupInfo，所以用回调函数
	updateSignalChan chan UpdateChanSignal // 更新keyMap的管道
	stopChannel      chan bool             // 停止信号
}

type UpdateChanSignal struct {
	Key  string
	Time int64
}

// BackupJob 备份回调函数
type BackupJob func(keyMap map[string]int64, start time.Time)

// 创建内存map数据处理器
func New(backUpInterVal time.Duration, backupJob BackupJob) *Processor {
	if backUpInterVal <= 0 {
		return nil
	}
	adGP := &Processor{
		backupInterval:   backUpInterVal,
		updateSignalChan: make(chan UpdateChanSignal),
		stopChannel:      make(chan bool),
		backupJob:        backupJob,
		keyMap:           make(map[string]int64),
	}
	return adGP
}

func (adgp *Processor) Start() {
	adgp.backupTicker = time.NewTicker(adgp.backupInterval)
	adgp.start()
}

func (adgp *Processor) start() {
	for {
		select {
		case <-adgp.backupTicker.C:
			adgp.backupTickHandler(false)
		case updateSignal := <-adgp.updateSignalChan:
			adgp.updateKeyMap(&updateSignal)
		case <-adgp.stopChannel:
			adgp.backupTicker.Stop()
			adgp.backupTickHandler(true)
			return
		}
	}
}

func (adgp *Processor) Stop() {
	adgp.stopChannel <- true
}

func (adgp *Processor) backupTickHandler(isStopSignal bool) {
	start := time.Now()
	// 这里需要深拷贝一份，不然会有并发问题
	backKeyMap := make(map[string]int64, len(adgp.keyMap))
	for key, val := range adgp.keyMap {
		backKeyMap[key] = val
	}

	//为了不阻塞正常消费流程，平常备份走子协程方式
	if !isStopSignal {
		go adgp.backupJob(backKeyMap, start)
	} else {
		adgp.backupJob(backKeyMap, start)
	}
}

// 这里保证并发安全
func (adgp *Processor) UpdateKeyMap(updateSignal UpdateChanSignal) {
	adgp.updateSignalChan <- updateSignal
}

func (adgp *Processor) updateKeyMap(updateSignal *UpdateChanSignal) {
	adgp.keyMap[updateSignal.Key] = updateSignal.Time
}

// 非并发安全，只有加载数据时可用
func (adgp *Processor) LoadKey(key string, time int64) {
	adgp.keyMap[key] = time
}

// 非并发安全，只有在start前可用
func (adgp *Processor) KeyMapLen() int {
	return len(adgp.keyMap)
}
