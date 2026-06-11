package service

import (
	"context"
	"fmt"
	"time"

	"fukubukuro/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ClaimCoupon 用户领取优惠券
func ClaimCoupon(db *gorm.DB, rdb *redis.Client, couponID uint, userID string) error {
	ctx := context.Background()
	lockKey := fmt.Sprintf("coupon_lock:%d", couponID)

	// 1. 用 Redis SETNX 抢分布式锁
	locked, err := rdb.SetNX(ctx, lockKey, userID, 5*time.Second).Result()
	if err != nil {
		return fmt.Errorf("抢锁失败: %w", err)
	}
	if !locked {
		return fmt.Errorf("优惠券太火爆，请稍后重试")
	}
	defer rdb.Del(ctx, lockKey) // 函数结束时释放锁

	// 2. 在事务中扣库存 + 记录领取
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 扣减优惠券库存（行锁，防超领）
	result := tx.Model(&model.Coupon{}).
		Where("id = ? AND stock > 0", couponID).
		Update("stock", gorm.Expr("stock - 1"))
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("优惠券已抢光")
	}

	// 记录领取关系
	couponRecord := model.CouponRecord{
		UserID:   userID,
		CouponID: couponID,
	}
	if err := tx.Create(&couponRecord).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录领取失败: %w", err)
	}

	return tx.Commit().Error
}
