package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"fukubukuro/internal/handler"
	"fukubukuro/internal/model"
	"fukubukuro/internal/mq"
	"fukubukuro/internal/seckill"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	r := gin.Default()

	// 1. 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
	})

	// 2. 加载秒杀 Lua 脚本
	if err := seckill.LoadScript(rdb); err != nil {
		log.Fatalf("加载秒杀脚本失败: %v", err)
	}
	fmt.Println("秒杀脚本加载成功")

	// 3. 连接数据库并自动建表
	if err := model.InitDB(); err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	fmt.Println("数据库连接成功")

	if err := model.AutoMigrate(); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("数据库表结构同步完成")

	// 4. 启动订单消费者 Worker
	mq.StartOrderConsumer(rdb)
	fmt.Println("订单消费者已启动")

	// 5. 用户模块
	userHandler := &handler.UserHandler{DB: model.DB}
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	// 6.商品模块
	productHandler := &handler.ProductHandler{DB: model.DB}
	r.GET("/api/products", productHandler.ListProducts)
	r.GET("/api/products/:id", productHandler.GetProduct)
	r.POST("/api/products", productHandler.CreateProduct)

	// 7. 购物车模块(需要登录)
	cartHandler := &handler.CartHandler{RDB: rdb}
	cartGroup := r.Group("/api/cart", handler.AuthMiddleware())
	{
		cartGroup.POST("", cartHandler.AddToCart)
		cartGroup.GET("", cartHandler.GetCart)
		cartGroup.PUT("/:product_id", cartHandler.UpdateCart)
		cartGroup.DELETE("/:product_id", cartHandler.RemoveFromCart)
	}

	// 8.订单模块（需要登录）
	orderHandler := &handler.OrderHandler{DB: model.DB, RDB: rdb}
	orderGroup := r.Group("/api/orders", handler.AuthMiddleware())
	{
		orderGroup.POST("", orderHandler.PlaceOrder)
		orderGroup.GET("", orderHandler.ListOrders)
		orderGroup.GET("/:id", orderHandler.GetOrder)
	}

	// 9.优惠券模块（需要登录）
	couponHandler := &handler.CouponHandler{DB: model.DB, RDB: rdb}
	r.POST("/api/coupons/:id/claim", handler.AuthMiddleware(), couponHandler.ClaimCoupon)

	// 10. 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "fukubukuro"})
	})

	// 11. 秒杀接口
	r.POST("/api/seckill", handler.AuthMiddleware(), func(c *gin.Context) {
		activityID := c.Query("activity_id")
		// 从鉴权中间件获取当前登录用户的ID
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			return
		}
		userID, ok := userIDVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效"})
			return
		}

		if activityID == "" || userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 activity_id 或 user_id"})
			return
		}

		result, err := seckill.DoSeckill(rdb, activityID, userID, 1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch result {
		case 1:
			// 秒杀成功，发消息到Redis Streams，由Worker异步消费生成订单
			msgID, err := mq.SendOrderMessage(rdb, userID, activityID)
			if err != nil {
				log.Printf("发送订单消息失败: %v", err)
			} else {
				log.Printf("订单消息已发送: MsgID=%s, UserID=%s, ActivityID=%s", msgID, userID, activityID)
			}
			c.JSON(http.StatusOK, gin.H{"message": "抢购成功！"})
		case 0:
			c.JSON(http.StatusOK, gin.H{"message": "库存不足，已抢光"})
		case -1:
			c.JSON(http.StatusOK, gin.H{"message": "您已经抢过了"})
		case -2:
			c.JSON(http.StatusOK, gin.H{"message": "活动不存在"})
		default:
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("未知结果: %d", result)})
		}
	})

	fmt.Println("福袋秒杀系统启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
