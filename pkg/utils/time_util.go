package utils

import (
	"errors"
	"fmt"
	"time"
)

// 定义时间格式化常量
const (
	FullTimeFmt       = "2006-01-02 15:04:05"
	ZoneTimeFmt       = "2006-01-02T15:04:05Z"
	DayTimeFmt        = "2006-01-02"
	MonthDayTimeFmt   = "01-02 15:04:05"
	HourMinuteTimeFmt = "15:04"
)

// GetUTCTime 获取UTC时间
func GetUTCTime() time.Time {
	return time.Now().UTC()
}

// GetLocalTime 获取当前时间
func GetLocalTime() time.Time {
	return time.Now()
}

// GetLocalTimeStr 获取本地当前时间的字符串表示，格式为 FullTimeFmt
func GetLocalTimeStr() string {
	return time.Now().Format(FullTimeFmt)
}

// GetLocalTimestamp 获取本地当前时间的时间戳，返回int64格式
func GetLocalTimestamp() int64 {
	return time.Now().Unix()
}

// GetLocalTimestampMilli 获取本地当前时间的毫秒级时间戳，返回int64格式
func GetLocalTimestampMilli() int64 {
	return time.Now().UnixMilli()
}

// TimeParseForZone 将符合 ZoneTimeFmt 格式的字符串解析为时间
func TimeParseForZone(s string) (time.Time, error) {
	tm, err := time.Parse(ZoneTimeFmt, s)
	if err != nil {
		return time.Time{}, errors.New("Time parsing failed: " + err.Error())
	}
	return time.Unix(tm.Unix(), 0), nil
}

// TimeToStr 将时间按照指定格式转换为字符串
func TimeToStr(tm time.Time, fmt string) string {
	return tm.Format(fmt)
}

// StrToTime 将符合指定格式的字符串解析为时间
func StrToTime(fmt, str string) (time.Time, error) {
	tm, err := time.Parse(fmt, str)
	if err != nil {
		return time.Time{}, errors.New("Time parsing failed: " + err.Error())
	}
	return tm, nil
}

// TimestampToFmtStr 将时间戳按照指定格式转换为字符串
func TimestampToFmtStr(timestamp int64, fmt string) string {
	return time.Unix(timestamp, 0).Format(fmt)
}

// FmtStrToTimestamp 将格式化的时间字符串转换为时间戳
func FmtStrToTimestamp(timeStr string, fmt string) (int64, error) {
	// 解析时间字符串
	tm, err := time.Parse(fmt, timeStr)
	if err != nil {
		return 0, err // 返回错误信息
	}
	// 转换为Unix时间戳（秒）
	return tm.Unix(), nil
}

// TimestampToTime 将时间戳转换为time.Time类型
// timestamp为秒级时间戳
func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// GetMonthStartEnd 获取指定时间所在月的开始和结束时间
func GetMonthStartEnd(tm time.Time) (time.Time, time.Time) {
	monthStartDay := tm.AddDate(0, 0, -tm.Day()+1)
	monthStartTime := time.Date(monthStartDay.Year(), monthStartDay.Month(), monthStartDay.Day(), 0, 0, 0, 0, tm.Location())
	monthEndDay := monthStartTime.AddDate(0, 1, -1)
	monthEndTime := time.Date(monthEndDay.Year(), monthEndDay.Month(), monthEndDay.Day(), 23, 59, 59, 0, tm.Location())
	return monthStartTime, monthEndTime
}

// IsTimestampToday 判断秒级时间戳是否为今天（本地时区）
func IsTimestampToday(timestamp int64) bool {
	// 获取今天的日期（年、月、日）
	year, month, day := time.Now().Local().Date()

	// 计算今天 00:00:00 和 23:59:59 的时间戳
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
	endOfDay := time.Date(year, month, day, 23, 59, 59, 0, time.Local).Unix()

	// 判断目标时间戳是否在今天的范围内
	return timestamp >= startOfDay && timestamp <= endOfDay
}

// IsTimestampYesterday 判断秒级时间戳是否为昨天（本地时区）
func IsTimestampYesterday(timestamp int64) bool {
	now := time.Now().Local()

	// 昨天 00:00:00
	startOfYesterday := time.Date(
		now.Year(), now.Month(), now.Day()-1,
		0, 0, 0, 0, now.Location(),
	)

	// 昨天 23:59:59
	endOfYesterday := time.Date(
		now.Year(), now.Month(), now.Day()-1,
		23, 59, 59, 0, now.Location(),
	)

	// 判断目标时间戳是否在昨天的范围内
	return timestamp >= startOfYesterday.Unix() && timestamp <= endOfYesterday.Unix()
}

// GetTimestamp30DaysAgo 返回30天前的该时刻的时间戳（秒级）
// 参数：currentTimestamp - 当前时刻的秒级时间戳
// 返回：30天前的秒级时间戳
func GetTimestamp30DaysAgo(currentTimestamp int64) int64 {
	// 将秒级时间戳转为Time对象
	tm := time.Unix(currentTimestamp, 0)

	// 减去30天（24小时×30天）
	thirtyDaysAgo := tm.Add(-30 * 24 * time.Hour)

	// 返回秒级时间戳
	return thirtyDaysAgo.Unix()
}

// GetTimestampNDaysLater 返回指定天数后的该时刻的时间戳（秒级）
// 参数：
//
//	currentTimestamp - 当前时刻的秒级时间戳
//	days - 要增加的天数（正数表示未来，负数表示过去）
//
// 返回：指定天数后的秒级时间戳
func GetTimestampNDaysLater(currentTimestamp int64, days int) int64 {
	// 将秒级时间戳转为Time对象
	tm := time.Unix(currentTimestamp, 0)

	// 加上指定天数
	futureTime := tm.Add(time.Duration(days) * 24 * time.Hour)

	// 返回秒级时间戳
	return futureTime.Unix()
}

// GetTomorrowZeroTimestamp 获取明天 00:00:00 的时间戳（秒级）
func GetTomorrowZeroTimestamp() int64 {
	now := time.Now()

	// 获取明天的日期
	tomorrow := now.Add(24 * time.Hour)

	// 构造明天零点的时间
	tomorrowZero := time.Date(
		tomorrow.Year(),
		tomorrow.Month(),
		tomorrow.Day(),
		0, 0, 0, 0, // 时、分、秒、纳秒
		tomorrow.Location(),
	)

	return tomorrowZero.Unix()
}

// IsSameDay 判断两个秒级时间戳是否为同一天（本地时区）
func IsSameDay(timestamp1, timestamp2 int64) bool {
	// 将时间戳转换为本地时间
	time1 := time.Unix(timestamp1, 0).Local()
	time2 := time.Unix(timestamp2, 0).Local()

	// 比较年、月、日是否相同
	return time1.Year() == time2.Year() &&
		time1.Month() == time2.Month() &&
		time1.Day() == time2.Day()
}

// GetYesterdayStartEnd 获取昨天 00:00:00 和 23:59:59 的时间戳（秒级）
func GetYesterdayStartEnd() (int64, int64) {
	now := time.Now().Local()

	// 昨天日期
	yesterday := now.AddDate(0, 0, -1)

	// 昨天 00:00:00
	startOfYesterday := time.Date(
		yesterday.Year(),
		yesterday.Month(),
		yesterday.Day(),
		0, 0, 0, 0,
		yesterday.Location(),
	).Unix()

	// 昨天 23:59:59
	endOfYesterday := time.Date(
		yesterday.Year(),
		yesterday.Month(),
		yesterday.Day(),
		23, 59, 59, 0,
		yesterday.Location(),
	).Unix()

	return startOfYesterday, endOfYesterday
}

// GetTimestampDayStartEnd 获取指定秒级时间戳所在日期的 00:00:00 和 23:59:59 的时间戳（秒级）
func GetTimestampDayStartEnd(timestamp int64) (int64, int64) {
	// 将秒级时间戳转换为 time.Time 类型，并使用本地时区
	date := time.Unix(timestamp, 0).Local()

	// 指定日期 00:00:00
	startOfDay := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		date.Location(),
	).Unix()

	// 指定日期 23:59:59
	endOfDay := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		23, 59, 59, 0,
		date.Location(),
	).Unix()

	return startOfDay, endOfDay
}

// IsDateLater 判断时间戳A的日期是否晚于时间戳B的日期
func IsDateLater(timestampA, timestampB int64) bool {
	// 将时间戳转换为日期（忽略具体时间）
	dateA := time.Unix(timestampA, 0).UTC().Truncate(24 * time.Hour)
	dateB := time.Unix(timestampB, 0).UTC().Truncate(24 * time.Hour)

	// 比较日期
	return dateA.After(dateB)
}

// FormatDuration 格式化时间差值
func FormatDuration(d time.Duration) string {
	// 转换为秒数
	seconds := int(d.Seconds()) % 60
	// 转换为分钟数
	minutes := int(d.Minutes()) % 60
	// 转换为小时数
	hours := int(d.Hours())

	// 根据时间长度返回不同格式
	switch {
	case hours > 0:
		return fmt.Sprintf("%d时%d分%d秒", hours, minutes, seconds)
	case minutes > 0:
		return fmt.Sprintf("%d分%d秒", minutes, seconds)
	default:
		return fmt.Sprintf("%d秒", seconds)
	}
}
