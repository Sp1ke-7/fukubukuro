package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// AddToCart 添加商品到购物车
func AddToCart(rdb *redis.Client, userID string, productID uint, quantity int) error {
	ctx := context.Background()
	cartKey := fmt.Sprintf("cart:%s", userID)
	return rdb.HSet(ctx, cartKey, strconv.FormatUint(uint64(productID), 10), quantity).Err()
}

// GetCart 获取购物车所有商品及其数量
func GetCart(rdb *redis.Client, userID string) (map[string]int, error) {
	ctx := context.Background()
	cartKey := fmt.Sprintf("cart:%s", userID)
	items, err := rdb.HGetAll(ctx, cartKey).Result()
	if err != nil {
		return nil, err
	}

	cart := make(map[string]int, len(items))
	for productID, qtyStr := range items {
		qty, _ := strconv.Atoi(qtyStr)
		cart[productID] = qty
	}
	return cart, nil
}

// RemoveFromCart 从购物车移除商品
func RemoveFromCart(rdb *redis.Client, userID string, productID uint) error {
	ctx := context.Background()
	cartKey := fmt.Sprintf("cart:%s", userID)
	return rdb.HDel(ctx, cartKey, strconv.FormatUint(uint64(productID), 10)).Err()
}

// UpdateCartQuantity 更新购物车商品数量
func UpdateCartQuantity(rdb *redis.Client, userID string, productID uint, quantity int) error {
	ctx := context.Background()
	cartKey := fmt.Sprintf("cart:%s", userID)
	return rdb.HSet(ctx, cartKey, strconv.FormatUint(uint64(productID), 10), quantity).Err()
}
