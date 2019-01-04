package utils

import (
	"fmt"
	"strings"
	"time"
)

// 中国时区对象缓存
var chinaTimeLocation *time.Location

//DateFormat 格式化时间
func DateFormat(t1 time.Time, format string) (result string) {
	tmpFormat := strings.Replace(format, "yyyy", "2006", 1)
	tmpFormat = strings.Replace(tmpFormat, "MM", "01", 1)
	tmpFormat = strings.Replace(tmpFormat, "dd", "02", 1)
	tmpFormat = strings.Replace(tmpFormat, "HH", "15", 1)
	tmpFormat = strings.Replace(tmpFormat, "mm", "04", 1)
	tmpFormat = strings.Replace(tmpFormat, "ss", "05", 1)
	return t1.Format(tmpFormat)
}

/*
字符串格式化时间Long（yyyy-MM-dd HH:mm:ss）
@t1 时间
*/
func LongStringToTime(t1 string) time.Time {
	t2, _ := time.Parse("2006-01-02 15:04:05", t1)
	return t2
}

//GetNowString 获取当前时间（返回字符串）
func GetNowString() string {
	return GetDateTimeString(time.Now().Local())
}

func GetDateTimeString(time time.Time) string {
	return time.Local().Format("2006-01-02 15:04:05")
}

//GetNow 获取当前时间（返回时间类型）
func GetNow() time.Time {
	return time.Now().Local()
}

/*
字符串格式化为时间Long string (yyyy-MM-dd HH::mm::ss)
@t1 时间
*/
func LongStringToShortString(t1 string, format string) string {
	tmpTime, _ := time.Parse("2006-01-02 15:04:05", t1)
	return DateFormat(tmpTime, format)
}

/*
时间格式化字符串Short（yyyy-MM-dd）
@t1 时间
*/
func ShortTimeToString(t1 time.Time) string {
	return t1.Format("2006-01-02")
}

/*
字符串格式化时间Short（yyyy-MM-dd）
@t1 时间
*/
func ShortStringToTime(t1 string) time.Time {
	t2, _ := time.Parse("2006-01-02", t1)
	return t2
}

// 格式化UTC时间为正常类型 当formatType=1时格式化为yyyy-MM-dd 当formatType=2时格式化为yyyy-MM-dd HH:ii:ss
func FormatUtcDateToDate(utcDate string, formatType int) string {

	if len(utcDate) >= 19 {
		if formatType == 1 {
			return utcDate[0:10]
		} else if formatType == 2 {
			return utcDate[0:19]
		}

	}
	return utcDate
}

// ParseChinaTime 使用中国东八区时区解析时间，参数与time.Parse()相同
func ParseChinaTime(layout string, value string) (t time.Time, err error) {
	if chinaTimeLocation == nil {
		if chinaTimeLocation, err = time.LoadLocation("Asia/Chongqing"); err != nil {
			return
		}
	}
	t, err = time.ParseInLocation(layout, value, chinaTimeLocation)
	return
}

/*
时间格式化字符串Long（yyyy-MM-dd HH:mm:ss）
@t1 时间
*/
func LongTimeToString(t1 time.Time) string {
	return t1.Format("2006-01-02 15:04:05")
}

// 转换时间戳
func GetTimeParse(times string) int64 {
	if "" == times {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04:05", times, loc)
	return parse.Unix()
}

/*
时间加上分钟（参数为string类型）
@t1 时间string类型
@minute 分钟
*/
func AddTimeStringToMinute(t1 string, minute int) string {
	t2 := LongStringToTime(t1) //string转为time格式
	m1, _ := time.ParseDuration(fmt.Sprint(minute, "m"))
	t3 := t2.Add(m1)
	return LongTimeToString(t3)
}

//获取两个时间相差分钟数
func GetMinuteDiffer(start_time, end_time string) int64 {
	var minute int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil {
		diff := t2.Unix() - t1.Unix()
		minute = diff / 60
		if minute < 0 {
			minute = -minute
		}
		return minute
	} else {
		return minute
	}
}

//时间戳转字符串
func TimestampToString(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05")
}

//获取时间格式HH:mm:ss time=6:00:00
func TimeShortString(timeHour string) string {
	loc, _ := time.LoadLocation("Local")
	formatTimeStr := time.Now().Format("2006-01-02")
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", formatTimeStr+timeHour, loc)
	un := t1.Unix()
	str := time.Unix(un, 0).Format("15:04:05")
	return str
}

// GetWeek 获取星期
func GetWeek(t string) (weekDay string) {
	t1, _ := time.Parse("2006-01-02", t)
	week := t1.Weekday().String()
	week = strings.ToLower(week)
	if week == "monday" {
		weekDay = "星期一 "
	} else if week == "tuesday" {
		weekDay = "星期二 "
	} else if week == "wednesday" {
		weekDay = "星期三 "
	} else if week == "thursday" {
		weekDay = "星期四 "
	} else if week == "friday" {
		weekDay = "星期五 "
	} else if week == "saturday" {
		weekDay = "星期六 "
	} else if week == "sunday" {
		weekDay = "星期日 "
	}
	return
}
