package handler

import (
	"net/http"
	"strconv"

	"fukubukuro/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CartHandler struct {
	RDB *redis.Client
}

// AddToCart 添加商品到购物车
func (h *CartHandler) AddToCart(c *gin.Context) {
	// 从鉴权中间件获取用户ID
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(string)

	var req struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	if err := service.AddToCart(h.RDB, userID, req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加购物车失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已添加到购物车"})
}

// GetCart 查看购物车
func (h *CartHandler) GetCart(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(string)

	cart, err := service.GetCart(h.RDB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取购物车失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

// UpdateCart 修改商品数量
func (h *CartHandler) UpdateCart(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(string)

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品ID"})
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if err := service.UpdateCartQuantity(h.RDB, userID, uint(productID), req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新购物车失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "购物车已更新"})
}

// RemoveFromCart 从购物车移除商品
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(string)

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品ID"})
		return
	}

	if err := service.RemoveFromCart(h.RDB, userID, uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除购物车失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已从购物车移除"})
}
