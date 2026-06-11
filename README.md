# 「福袋」社交化电商平台 (Fukubukuro)

基于 Go 语言开发的社交化电商平台，聚焦高并发场景下的秒杀抢购与订单处理，覆盖用户中心、商品系统、秒杀福袋、购物车、订单、优惠券六大核心电商模块。

## 技术栈

- **语言**：Go 1.25
- **框架**：Gin
- **数据库**：MySQL (GORM)
- **缓存/锁/消息队列**：Redis (go-redis)
- **部署**：Docker + Docker Compose（多阶段构建）

## 核心功能与技术亮点

| 模块 | 核心功能 | 技术亮点 |
|:---|:---|:---|
| **用户中心** | 注册/登录、JWT鉴权、个人信息管理 | 基于Redis的Token主动失效机制，实现单点登录互踢 |
| **商品系统** | 商品列表、详情、分类检索 | MySQL复合索引优化，GORM Preload解决N+1查询问题 |
| **秒杀福袋** | 定时开抢、库存扣减、订单生成 | **Redis + Lua脚本保证原子性**，防止超卖；消息队列异步处理订单 |
| **购物车与订单** | 加入购物车、下单、订单状态流转 | 购物车数据Redis缓存，订单状态机设计，保证数据最终一致性 |
| **优惠券系统** | 领券、用券、过期回收 | 优惠券库存扣减的高并发处理，使用Redis分布式锁防止超领 |
| **部署与监控** | Docker Compose一键部署 | Prometheus + Grafana监控QPS、订单处理延迟、数据库连接池 |

## 快速启动

docker compose up -d

服务启动后访问 `http://localhost:8080/ping`，返回 `{"message":"fukubukuro"}` 即部署成功。

## API 文档

| 方法 | 路径 | 说明 |
|:---|:---|:---|
| POST | /api/register | 用户注册 |
| POST | /api/login | 用户登录 |
| GET | /api/products | 商品列表（支持分类筛选与分页） |
| GET | /api/products/:id | 商品详情 |
| POST | /api/seckill | 秒杀福袋（需登录） |
| POST | /api/cart | 添加购物车（需登录） |
| GET | /api/cart | 查看购物车（需登录） |
| PUT | /api/cart/:product_id | 修改购物车商品数量（需登录） |
| DELETE | /api/cart/:product_id | 删除购物车商品（需登录） |
| POST | /api/orders | 下单（需登录） |
| GET | /api/orders | 订单列表（需登录） |
| GET | /api/orders/:id | 订单详情（需登录） |
| POST | /api/coupons/:id/claim | 领取优惠券（需登录） |