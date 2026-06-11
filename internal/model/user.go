package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 返回时不暴露密码
	TokenVersion int       `gorm:"default:1" json:"-"`                  // 用于单点互踢
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
