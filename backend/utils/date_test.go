package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestCalculateBillingCycles(t *testing.T) {
	// 2025/4/6 00:00:00:00 使用本地时区
	startTime := time.Date(2025, 3, 4, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2025, 5, 4, 0, 0, 0, 0, time.Local)
	n, _ := CalculateBillingCycles(startTime, endTime)
	fmt.Println(n)
}
