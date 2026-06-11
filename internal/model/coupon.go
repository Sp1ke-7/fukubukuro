package model

import (
	"time"
)

// Coupon 优惠券模型
type Coupon struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Stock     int       `gorm:"not null;default:0" json:"stock"` // 剩余库存
	Total     int       `gorm:"not null;default:0" json:"total"` // 总库存
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
