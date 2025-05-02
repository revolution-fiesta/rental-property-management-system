package runner

import (
	"context"
	"fmt"
	"log/slog"
	"rental-property-management-system/backend/store"
	"rental-property-management-system/backend/utils"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var billingMu sync.Mutex

func StartBillingRunner(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				slog.Info(fmt.Sprintf("Billing Runner scans orders in every %d seconds", 10))

				billingMu.Lock()
				if err := generateBilling(); err != nil {
					slog.Warn(err.Error())
					billingMu.Unlock()
					continue
				}
				billingMu.Unlock()

			case <-ctx.Done():
				slog.Info("billing runner exits")
				return
			}
		}
	}()
}

// WARN: 这个函数的可行性需要严格测试
func generateBilling() error {
	now := time.Now()
	orders := []store.Order{}
	if err := store.GetDB().Find(&orders).Error; err != nil {
		return errors.New("failed to sync orders from database")
	}
	for _, order := range orders {
		if order.RemainingBiilNum == 0 {
			continue
		}
		numCycles, err := utils.CalculateBillingCycles(order.CreatedAt, now)
		if err != nil {
			return errors.Wrapf(err, "failed to calculate billing cycles")
		}
		numPaidCycles := order.TotalTerm - order.RemainingBiilNum
		// 如果已经支付的周期小于当前需要支付的周期, 则生成账单
		// TODO: 这里应该用 Transaction
		if numPaidCycles < numCycles {
			// 根据房间 ID 查询月付价格
			room := store.Room{}
			if err := store.GetDB().Where("id = ?", order.RoomID).Find(&room).Error; err != nil {
				return errors.Wrapf(err, "failed to sync from database")
			}
			// 生成账单
			billing := store.Billing{
				Type:   string(store.BillingTypeMonthlyPayment),
				UserID: order.UserID,
				Price:  room.Price,
				Paid:   false,
				Name:   fmt.Sprintf("%s月付账单 (%d/%d)", room.Name, order.TotalTerm-order.RemainingBiilNum+1, order.TotalTerm),
			}
			if err := store.GetDB().Save(&billing).Error; err != nil {
				return errors.Wrapf(err, "failed to create the bill")
			}
			// 更新账单剩余期数
			order.RemainingBiilNum--
			if err := store.GetDB().Save(&order).Error; err != nil {
				return errors.Wrapf(err, "failed to update remaining cycles number")
			}
		}
	}
	return nil
}
