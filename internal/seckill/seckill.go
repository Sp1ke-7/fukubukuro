package seckill

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/redis/go-redis/v9"
)

//go:embed seckill.lua
var luaScript string

var seckillSHA string

// LoadScript 在服务启动时调用一次，把脚本预加载到Redis
func LoadScript(rdb *redis.Client) error {
	ctx := context.Background()
	sha, err := rdb.ScriptLoad(ctx, luaScript).Result()
	if err != nil {
		return fmt.Errorf("加载秒杀脚本失败: %w", err)
	}
	seckillSHA = sha
	return nil
}

// DoSeckill 执行秒杀
// activityID: 福袋活动ID
// userID: 用户ID
// 返回: 1=成功, 0=库存不足, -1=已抢过, -2=活动不存在
func DoSeckill(rdb *redis.Client, activityID string, userID string, limit int) (int, error) {
	ctx := context.Background()
	stockKey := fmt.Sprintf("seckill:stock:%s", activityID)
	usersKey := fmt.Sprintf("seckill:users:%s", activityID)

	result, err := rdb.EvalSha(ctx, seckillSHA, []string{stockKey, usersKey}, userID, limit).Int()
	if err != nil {
		return 0, fmt.Errorf("执行秒杀脚本失败: %w", err)
	}
	return result, nil
}
