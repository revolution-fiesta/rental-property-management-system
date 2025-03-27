package utils

import "time"

// 计算按天数比例分摊的房租
func CalculateProRatedRent(fullRent float64, date time.Time) float64 {
	year, month, _ := date.Date()
	location := date.Location()
	// 获取当前月份的总天数
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	totalDays := lastOfMonth.Day()
	// 获取当天是这个月的第几天
	currentDay := date.Day()
	// 计算比例并返回费用
	return (float64(currentDay) / float64(totalDays)) * fullRent
}
