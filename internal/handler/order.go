package handler

import (
	"net/http"
	"strconv"

	"fukubukuro/internal/model"
	"fukubukuro/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// PlaceOrder 从购物车下单
func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(string)

	orders, err := service.PlaceOrder(h.RDB, h.DB, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"orders": orders})
}

// ListOrders 查询用户订单列表
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(string)

	var orders []model.Order
	if err := h.DB.Where("user_id = ?", userID).Order("id DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单失败"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetOrder 查询单个订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	var order model.Order
	if err := h.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}
	c.JSON(http.StatusOK, order)
}
