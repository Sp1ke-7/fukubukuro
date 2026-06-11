package handler

import (
	"net/http"
	"strconv"

	"fukubukuro/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CouponHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// ClaimCoupon 用户领取优惠券
func (h *CouponHandler) ClaimCoupon(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户标识异常"})
		return
	}

	couponIDStr := c.Param("id")
	couponID, err := strconv.ParseUint(couponIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的优惠券ID"})
		return
	}

	if err := service.ClaimCoupon(h.DB, h.RDB, uint(couponID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "领取成功"})
}
