/**
 * @File: time.go
 * @Author: zhuchengming
 * @Description:
 * @Date: 2021/2/10 12:05
 */

package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetFormatRequestTime(time time.Time) string {
	return fmt.Sprintf("%d.%d", time.Unix(), time.Nanosecond()/1e3)
}

func GetRequestCost(start, end time.Time) float64 {
	return float64(end.Sub(start).Nanoseconds()/1e4) / 100.0
}

// 将Y-m-d格式时间转化为当天 开始时间 和 结束时间 的unix时间戳
func YmdDateToUnixStartAndEnd(ymdDate string) (startTimeUnix int64, endTimeUnix int64) {
	t, _ := time.Parse("2006-01-02 Z0700 MST", ymdDate+" +0800 CST")
	startTimeUnix = t.Unix()
	endTimeUnix = startTimeUnix + 86399
	return
}

// 获取unix时间戳对应的中文周几
func GetWeekDayFromUnixTimestamp(timestamp int64) string {
	weekMap := map[string]string{
		"Monday":    "周一",
		"Tuesday":   "周二",
		"Wednesday": "周三",
		"Thursday":  "周四",
		"Friday":    "周五",
		"Saturday":  "周六",
		"Sunday":    "周日",
	}
	week := time.Unix(timestamp, 0).Format("Monday")
	weekStr, _ := weekMap[week]
	return weekStr
}

// 获取当天开始时间戳
func GetTodayStartUnix() int64 {
	s := time.Now().Format("2006-01-02") + " 00:00:00"
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	return timestamp.Unix()
}

//获取当天结束时间戳
func GetTodayEndUnix() int64 {
	s := time.Now().Format("2006-01-02") + " 23:59:59"
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	return timestamp.Unix()
}

//时间戳转日期
func FormatUnixToDayTime(unixTimestamp int64) string {
	return time.Unix(unixTimestamp, 0).Format("2006-01-02")
}

//通过给定的时间戳获取当天0点的时间戳
func GetDayStartUnix(unixTimestamp int64) int64 {
	if unixTimestamp == 0 {
		return 0
	}
	s := time.Unix(unixTimestamp, 0).Format("2006-01-02")
	timestamp, _ := time.ParseInLocation("2006-01-02", s, time.Local)
	return timestamp.Unix()
}

//获取中文周几
func GetWeekDayName(timestamp int64) string {
	weekName := [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	weekDay := time.Unix(timestamp, 0).Weekday()

	return weekName[weekDay]
}

//获取当天固定时间时间戳  e.g time = "19:59:23"
func GetUnixByTime(times string) int64 {
	str := time.Now().Format("2006-01-02") + " " + times
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	return timestamp.Unix()
}

// 将unix时间格式化为yyyy-mm-dd HH:ii:ss格式
func GetTimeByUnix(sec int64) string {
	if sec == 0 {
		return ""
	}
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}

//格式化（x岁x个月)
func FormatAge(birthday int64, addMonth bool) string {
	if birthday == 0 {
		return ""
	}
	birthDayStr := time.Unix(birthday, 0).Format("2006-01-02")
	birthDayComp := strings.Split(birthDayStr, "-")
	birthYear, _ := strconv.ParseInt(birthDayComp[0], 10, 64)
	birthMonth, _ := strconv.ParseInt(birthDayComp[1], 10, 64)

	curDayStr := time.Now().Format("2006-01-02")
	curDayComp := strings.Split(curDayStr, "-")
	curYear, _ := strconv.ParseInt(curDayComp[0], 10, 64)
	curMonth, _ := strconv.ParseInt(curDayComp[1], 10, 64)
	var month int64
	var year int64
	var yearPlus int64

	month = curMonth - birthMonth
	if addMonth {
		month += 1
	}

	if month < 0 {
		month += 12
		yearPlus = 1
	}
	year = curYear - yearPlus - birthYear
	if year <= 0 {
		year = 0
	}
	if year > 20 {
		return ""
	}
	outAge := ""
	if year > 0 {
		outAge = fmt.Sprintf("%d岁", year)
	}
	if month > 0 {
		outAge = fmt.Sprintf("%s%d个月", outAge, month)
	}
	return outAge
}

//返回年岁
func FormatAgeYear(birthday int64, addMonth bool) uint64 {
	birthDayStr := time.Unix(birthday, 0).Format("2006-01-02")
	birthDayComp := strings.Split(birthDayStr, "-")
	birthYear, _ := strconv.ParseInt(birthDayComp[0], 10, 64)
	birthMonth, _ := strconv.ParseInt(birthDayComp[1], 10, 64)

	curDayStr := time.Now().Format("2006-01-02")
	curDayComp := strings.Split(curDayStr, "-")
	curYear, _ := strconv.ParseInt(curDayComp[0], 10, 64)
	curMonth, _ := strconv.ParseInt(curDayComp[1], 10, 64)
	var month int64
	var year int64
	var yearPlus int64
	month = curMonth - birthMonth
	if addMonth {
		month += 1
	}
	if month < 0 {
		month += 12
		yearPlus = 1
	}
	year = curYear - yearPlus - birthYear
	if year <= 0 {
		year = 0
	}
	if year > 20 {
		year = 0
	}
	return uint64(year)
}

//获取当天所在周的周一0点时间
func GetDateWeekOfMondayStart(startTime int64) int64 {
	now := time.Unix(startTime, 0)
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStartDate.Unix()
}

//获取当天所在周的周日23：59：59点时间
func GetDateWeekOfSundayEnd(startTime int64) int64 {
	return GetDateWeekOfMondayStart(startTime) + 7*86400 - 1
}

//获取当天所在周的周日00：00：00点时间
func GetDateWeekOfSundayStart(startTime int64) int64 {
	return GetDateWeekOfMondayStart(startTime) + 6*86400
}
