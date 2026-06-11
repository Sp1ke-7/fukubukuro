package service

import (
	"fukubukuro/internal/model"

	"gorm.io/gorm"
)

// CreateProduct 创建商品
func CreateProduct(db *gorm.DB, product *model.Product) error {
	return db.Create(product).Error
}

// GetProductByID 根据ID查询商品详情
func GetProductByID(db *gorm.DB, id uint) (*model.Product, error) {
	var product model.Product
	err := db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// ListProducts 分页查询商品列表，支持按分类筛选
func ListProducts(db *gorm.DB, category string, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := db.Model(&model.Product{})
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 先统计总数
	query.Count(&total)

	// 再分页查询
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&products).Error

	return products, total, err
}
