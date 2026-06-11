package model

import (
	"fmt"
	"time"
)

// Order 订单模型
type Order struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"type:varchar(64);not null" json:"user_id"`
	ActivityID string    `gorm:"type:varchar(64);not null" json:"activity_id"`
	Status     string    `gorm:"type:varchar(32);not null;default:'pending'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// 订单状态常量
const (
	OrderStatusPending   = "pending"   // 待支付
	OrderStatusPaid      = "paid"      // 已支付
	OrderStatusShipped   = "shipped"   // 已发货
	OrderStatusCompleted = "completed" // 已完成
	OrderStatusCancelled = "cancelled" // 已取消
)

// Pay 支付订单
func (o *Order) Pay() error {
	if o.Status != OrderStatusPending {
		return fmt.Errorf("订单状态不是待支付，无法支付")
	}
	o.Status = OrderStatusPaid
	return DB.Model(o).Update("status", OrderStatusPaid).Error
}

// Ship 发货
func (o *Order) Ship() error {
	if o.Status != OrderStatusPaid {
		return fmt.Errorf("订单状态不是已支付，无法发货")
	}
	o.Status = OrderStatusShipped
	return DB.Model(o).Update("status", OrderStatusShipped).Error
}
