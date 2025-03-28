package utils

import (
	"errors"
	"time"
)

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

// 计算从 startDate 到今天已经是第几期账单
func CalculateBillingCycles(startDate, endDate time.Time) (uint, error) {
	// 如果今天早于开始日期，则返回 0
	if endDate.Before(startDate) {
		return 0, nil
	}
	// 计算起始日期和今天的年份和月份差
	yearsDiff := endDate.Year() - startDate.Year()
	monthsDiff := endDate.Month() - startDate.Month()
	cycles := yearsDiff*12 + int(monthsDiff)
	// 确保当前周期是否已满
	startDay := startDate.Day()
	if endDate.Day() >= startDay {
		cycles++ // 包括当前周期
	}
	if cycles < 0 {
		return 0, errors.New("cycles below zero")

	}
	return uint(cycles), nil
}
