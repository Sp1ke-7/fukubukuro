package model

import (
	"time"
)

// CouponRecord 用户领取优惠券的记录
type CouponRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(64);not null;index" json:"user_id"`
	CouponID  uint      `gorm:"not null;index" json:"coupon_id"`
	Status    string    `gorm:"type:varchar(32);not null;default:'unused'" json:"status"` // unused, used, expired
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
