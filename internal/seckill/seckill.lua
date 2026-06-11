-- seckill.lua
-- 福袋秒杀核心脚本
-- KEYS[1]: 库存 key，如 seckill:stock:1
-- KEYS[2]: 已抢用户集合 key，如 seckill:users:1
-- ARGV[1]: 用户ID
-- ARGV[2]: 每人限购数量

local stock_key = KEYS[1]
local users_key = KEYS[2]
local user_id = ARGV[1]
local limit = tonumber(ARGV[2])

-- 1. 检查用户是否已经抢过
local already = redis.call('SISMEMBER', users_key, user_id)
if already == 1 then
    return -1  -- 用户已抢过，不能重复抢
end

-- 2. 检查库存
local stock = redis.call('GET', stock_key)
if not stock then
    return -2  -- 库存 key 不存在，活动未初始化
end

stock = tonumber(stock)
if stock <= 0 then
    return 0   -- 库存不足
end

-- 3. 原子扣减库存 + 记录用户
redis.call('DECR', stock_key)
redis.call('SADD', users_key, user_id)

return 1  -- 抢购成功