package model

import (
	"time"
)

// Product 商品模型
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Category    string    `gorm:"type:varchar(100);index;not null" json:"category"` // 分类索引
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	ImageURL    string    `gorm:"type:varchar(512)" json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
