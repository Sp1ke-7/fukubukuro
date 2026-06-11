package service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"fukubukuro/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// PlaceOrder 从购物车下单
func PlaceOrder(rdb *redis.Client, db *gorm.DB, userID string) ([]model.Order, error) {
	// 1. 获取购物车数据
	cart, err := GetCart(rdb, userID)
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败: %w", err)
	}
	if len(cart) == 0 {
		return nil, fmt.Errorf("购物车为空")
	}

	// 2. 开启数据库事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var orders []model.Order
	for productIDStr, qty := range cart {
		productID, _ := strconv.ParseUint(productIDStr, 10, 64)
		if qty <= 0 {
			continue
		}

		// 3. 扣减商品库存（行锁）
		result := tx.Model(&model.Product{}).
			Where("id = ? AND stock >= ?", productID, qty).
			Update("stock", gorm.Expr("stock - ?", qty))
		if result.Error != nil || result.RowsAffected == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("商品 %d 库存不足", productID)
		}

		// 4. 生成订单
		order := model.Order{
			UserID:     userID,
			ActivityID: "normal", // 普通订单，非秒杀
			Status:     model.OrderStatusPending,
		}
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("生成订单失败: %w", err)
		}
		orders = append(orders, order)
		log.Printf("订单生成成功: OrderID=%d, UserID=%s, ProductID=%d, Quantity=%d", order.ID, userID, productID, qty)
	}

	// 5. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 6. 清空购物车
	cartKey := fmt.Sprintf("cart:%s", userID)
	rdb.Del(context.Background(), cartKey)

	return orders, nil
}
