package mq

import (
	"context"
	"encoding/json"
	"log"

	"fukubukuro/internal/model"

	"github.com/redis/go-redis/v9"
)

// StartOrderConsumer 启动订单消费者 Worker
func StartOrderConsumer(rdb *redis.Client) {
	go func() {
		ctx := context.Background()
		streamName := "order_stream"
		consumerGroup := "order_group"
		consumerName := "consumer-1"

		// 1. 创建消费组（如果不存在）
		if err := rdb.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err(); err != nil {
			// 消费组已存在不是错误，忽略
			if err.Error() != "BUSYGROUP Consumer Group name already exists" {
				log.Printf("创建消费组失败: %v", err)
			}
		}

		for {
			// 2. 从 Stream 读取新消息（阻塞等待）
			msgs, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: consumerName,
				Streams:  []string{streamName, ">"},
				Block:    0, // 阻塞直到有新消息
			}).Result()
			if err != nil {
				log.Printf("读取消息失败: %v", err)
				continue
			}

			// 3. 处理取到的消息
			for _, msg := range msgs[0].Messages {
				orderData, ok := msg.Values["order"]
				if !ok {
					log.Printf("消息格式错误: %v", msg.ID)
					rdb.XAck(ctx, streamName, consumerGroup, msg.ID)
					continue
				}

				var orderMsg OrderMessage
				if err := json.Unmarshal([]byte(orderData.(string)), &orderMsg); err != nil {
					log.Printf("反序列化订单消息失败: %v", err)
					rdb.XAck(ctx, streamName, consumerGroup, msg.ID)
					continue
				}

				// 4. 生成订单写入数据库
				order := model.Order{
					UserID:     orderMsg.UserID,
					ActivityID: orderMsg.ActivityID,
					Status:     "pending",
				}
				if err := model.DB.Create(&order).Error; err != nil {
					log.Printf("订单生成失败: %v", err)
				} else {
					log.Printf("订单生成成功: OrderID=%d, UserID=%s, ActivityID=%s", order.ID, orderMsg.UserID, orderMsg.ActivityID)
				}

				// 5. 确认消息处理完成
				rdb.XAck(ctx, streamName, consumerGroup, msg.ID)
			}
		}
	}()
}
