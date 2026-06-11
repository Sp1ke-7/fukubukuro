package handler

import (
	"net/http"
	"strconv"

	"fukubukuro/internal/model"
	"fukubukuro/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

// GetProduct 获取商品详情
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品ID"})
		return
	}

	product, err := service.GetProductByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// ListProducts 商品列表（分页 + 分类筛选）
func (h *ProductHandler) ListProducts(c *gin.Context) {
	category := c.Query("category")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := service.ListProducts(h.DB, category, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询商品列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products":  products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateProduct 创建商品（管理后台用）
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}
	if err := service.CreateProduct(h.DB, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}
	c.JSON(http.StatusCreated, product)
}
