# 「福袋」社交化电商平台 (Fukubukuro)

基于 Go 语言开发的社交化电商平台，聚焦高并发场景下的秒杀抢购与订单处理。

## 技术栈

- **语言**：Go 1.25
- **框架**：Gin
- **数据库**：MySQL (GORM)
- **缓存/锁/消息队列**：Redis (go-redis)
- **部署**：Docker + Docker Compose（多阶段构建）

## 核心功能

- **用户中心**：注册/登录、JWT 鉴权
- **商品系统**：商品列表、详情、分类检索（索引优化）
- **秒杀福袋**：Redis + Lua 脚本原子扣库存，防超卖；Redis Streams 消息队列异步处理订单
- **购物车与订单**：Redis Hash 存储购物车，数据库事务下单，订单状态机设计
- **优惠券系统**：Redis 分布式锁（SETNX）防超领
- **容器化部署**：Dockerfile 多阶段构建 + Docker Compose 一键启动

## 快速启动

docker compose up -d

服务启动后访问 `http://localhost:8080/ping`，返回 `{"message":"fukubukuro"}` 即部署成功。

## API 文档

| 方法 | 路径 | 说明 |
|:---|:---|:---|
| POST | /api/register | 用户注册 |
| POST | /api/login | 用户登录 |
| GET | /api/products | 商品列表 |
| GET | /api/products/:id | 商品详情 |
| POST | /api/seckill | 秒杀福袋（需登录） |
| POST | /api/cart | 添加购物车（需登录） |
| GET | /api/cart | 查看购物车（需登录） |
| POST | /api/orders | 下单（需登录） |
| GET | /api/orders | 订单列表（需登录） |
| POST | /api/coupons/:id/claim | 领取优惠券（需登录） |

## 项目亮点

- 秒杀与优惠券采用差异化并发方案：**Lua 原子操作 vs 分布式锁**，根据场景选择最优解
- Redis Streams 实现消息队列，**解耦秒杀接口与数据库写入**
- 分布式锁设 TTL 防死锁，数据库乐观锁兜底，**双重保障数据一致性**
- Docker 多阶段构建，最终镜像仅几十 MB