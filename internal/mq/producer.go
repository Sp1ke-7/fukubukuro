package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// OrderMessage 订单消息体
type OrderMessage struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
}

// SendOrderMessage 秒杀成功后发送订单消息到 Redis Streams
func SendOrderMessage(rdb *redis.Client, userID, activityID string) (string, error) {
	ctx := context.Background()

	msg := OrderMessage{
		UserID:     userID,
		ActivityID: activityID,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("序列化订单消息失败: %w", err)
	}

	// XAdd 追加消息到订单Stream
	msgID, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "order_stream",
		Values: map[string]interface{}{
			"order": string(data),
		},
	}).Result()
	if err != nil {
		return "", fmt.Errorf("发送订单消息失败: %w", err)
	}

	return msgID, nil
}
