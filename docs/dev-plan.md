# T-Line Sports Club 项目任务分解与开发排期

> 版本：v1.0
> 日期：2026-03-19
> 状态：规划阶段
> 团队规模：1-2 名全栈开发

---

## 已明确的决策

- **架构**: 模块化单体 (Modular Monolith)，单 Go 二进制，内部按业务域分模块
- **仓库策略**: Monorepo，三端代码统一管理 (backend / miniapp / admin)
- **技术栈**: Go Gin + PostgreSQL + Redis / uni-app Vue3 / Vue3 Arco Design Pro
- **认证**: JWT + 微信 OAuth，Access Token 2h + Refresh Token 30d
- **支付**: 组合支付（余额优先 + 微信补差），冻结-确认模式
- **部署**: Docker + Docker Compose，初期单机 4C8G
- **缓存**: Redis 管理场地时段锁（SETNX）、候补队列（Sorted Set）、热数据缓存

---

## 整体规划概述

### 项目目标

构建 T-Line 网球馆 SaaS 系统，包含用户端微信小程序、管理端 Web 后台、Go 后端 API 三端，实现场地预订（含候补）、商品售卖、活动管理、教务系统、数据统计等核心功能。MVP 目标支撑 1000 用户以内的日常运营。

### 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22+ / Gin 1.10+ / GORM v2 / PostgreSQL 16+ / Redis 7+ |
| 用户端 | uni-app + Vue 3.4+ / TypeScript / Pinia |
| 管理端 | Vue 3.4+ / Arco Design Pro 2.x / TypeScript / Pinia |
| 基础设施 | Docker / Docker Compose / Nginx / 阿里云 OSS |
| 外部服务 | 微信开放平台 / 微信支付 / 短信服务 |

### 主要阶段

| 阶段 | 名称 | 预估工期 | MVP 必要性 |
|------|------|---------|-----------|
| Phase 1 | 基础设施 + 用户系统 | 12 天 | **必须** |
| Phase 2 | 核心业务 -- 场地预订 | 18 天 | **必须** |
| Phase 3 | 商品 + 活动 | 14 天 | **必须**（活动为核心） |
| Phase 4 | 教务系统 | 12 天 | 可延后 |
| Phase 5 | 数据统计 + 优化 | 8 天 | 部分必须 |
| Phase 6 | 部署上线 | 6 天 | **必须** |
| **总计** | | **约 70 天** | |

> 按 1 名全栈开发计算。若 2 人并行（前后端分工），可压缩至 40-45 天。

---

## 详细任务分解

### Phase 1：基础设施 + 用户系统（12 天）

**目标**: 搭建三端项目脚手架，完成数据库初始化、用户认证、基础中间件，三端可联调跑通登录流程。

#### 1.1 后端基础设施

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P1-001 | Go 项目初始化 | backend | 初始化 Go module，建立 `cmd/server/main.go` 入口；按架构文档创建 `internal/` 目录结构（config / server / middleware / pkg / modules）；配置 Makefile | 无 | 0.5 |
| P1-002 | 配置管理 | backend | 实现 `internal/config/config.go`：YAML + 环境变量加载（Viper）；创建 `configs/config.yaml` 和 `config.example.yaml`；包含数据库、Redis、JWT、微信、OSS 等配置项 | P1-001 | 0.5 |
| P1-003 | 日志系统 | backend | 实现 `internal/pkg/logger/logger.go`：Zap 初始化、日志等级配置、Lumberjack 日志轮转、请求 ID 上下文传递 | P1-002 | 0.5 |
| P1-004 | 数据库初始化 | backend | 实现 `internal/pkg/database/postgres.go`：GORM 初始化、连接池配置、自动迁移开关；实现 `internal/pkg/database/redis.go`：Redis 客户端封装、通用缓存/锁操作 | P1-002 | 1 |
| P1-005 | 统一响应与错误码 | backend | 实现 `internal/pkg/response/response.go`：统一成功/分页/错误响应格式；实现 `internal/pkg/errors/codes.go`：错误码常量定义（40000-50099 各范围） | P1-001 | 0.5 |
| P1-006 | 基础中间件 | backend | 实现 `internal/middleware/` 下所有中间件：CORS、Panic Recovery、请求日志（含请求 ID）、限流（令牌桶，基于 Redis） | P1-003, P1-004 | 1 |
| P1-007 | HTTP 服务与路由 | backend | 实现 `internal/server/server.go`：Gin 引擎初始化、优雅关闭；实现 `internal/server/router.go`：总路由注册，挂载中间件链，健康检查 `/api/v1/health` | P1-006 | 0.5 |
| P1-008 | 参数校验与分页工具 | backend | 实现 `internal/pkg/validator/validator.go`：Gin binding 自定义校验器；实现 `internal/pkg/pagination/pagination.go`：统一分页参数解析与响应构造 | P1-005 | 0.5 |
| P1-009 | Docker 开发环境 | backend | 编写根目录 `docker-compose.yml`：PostgreSQL 16 + Redis 7 + Go API 服务（热重载 air）；Makefile 增加 `make dev` 一键启动 | P1-007 | 0.5 |

**涉及文件**:
- `backend/cmd/server/main.go`
- `backend/internal/config/config.go`
- `backend/internal/pkg/logger/logger.go`
- `backend/internal/pkg/database/postgres.go`
- `backend/internal/pkg/database/redis.go`
- `backend/internal/pkg/response/response.go`
- `backend/internal/pkg/errors/codes.go`
- `backend/internal/middleware/*.go`
- `backend/internal/server/server.go`
- `backend/internal/server/router.go`
- `backend/internal/pkg/validator/validator.go`
- `backend/internal/pkg/pagination/pagination.go`
- `backend/configs/config.yaml`
- `backend/configs/config.example.yaml`
- `backend/Dockerfile`
- `backend/Makefile`
- `backend/go.mod`
- `docker-compose.yml`

#### 1.2 数据库迁移 -- 用户域

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P1-010 | 用户域数据库迁移 | backend | 编写 `migrations/000001_init_users.up.sql`：创建 users / wallets / wallet_transactions / points 四张表及索引；编写对应 down 迁移；编写 `scripts/migrate.sh` | P1-004 | 0.5 |

**涉及文件**:
- `backend/migrations/000001_init_users.up.sql`
- `backend/migrations/000001_init_users.down.sql`
- `backend/scripts/migrate.sh`

#### 1.3 用户认证模块

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P1-011 | JWT 工具 | backend | 实现 `internal/pkg/jwt/jwt.go`：Token 签发（Access + Refresh）、解析、刷新逻辑；JWT Payload 包含 user_id / role / member_level | P1-002 | 0.5 |
| P1-012 | JWT 认证中间件 | backend | 实现 `internal/middleware/auth.go`：从 Header 提取 Bearer Token，验证并注入用户信息到 Context；实现 `internal/middleware/rbac.go`：基于角色的权限校验 | P1-011 | 0.5 |
| P1-013 | 微信 SDK 封装 | backend | 实现 `internal/integration/wechat/auth.go`：code2session 换取 openid/session_key；手机号解密（encrypted_data + iv） | P1-002 | 1 |
| P1-014 | 短信 SDK 封装 | backend | 实现 `internal/integration/sms/sms.go`：发送验证码、验证码校验（Redis 存储 5 分钟过期） | P1-004 | 0.5 |
| P1-015 | Auth 模块 -- Model + Repository | backend | 实现 `internal/modules/auth/model.go`：User / Wallet / WalletTransaction / Points GORM 模型；实现 `repository.go`：用户 CRUD、按 openid/phone 查询、钱包操作 | P1-010 | 1 |
| P1-016 | Auth 模块 -- Service | backend | 实现 `internal/modules/auth/service.go`：微信登录（查找/创建用户 + 签发 Token）、手机号登录、Token 刷新、获取/更新个人信息、钱包查询 | P1-015, P1-013, P1-014 | 1 |
| P1-017 | Auth 模块 -- Handler + Routes | backend | 实现 `internal/modules/auth/handler.go`：各 API Handler；实现 `dto.go`：请求/响应 DTO 定义与校验标签；实现 `routes.go`：路由注册 | P1-016 | 0.5 |

**涉及文件**:
- `backend/internal/pkg/jwt/jwt.go`
- `backend/internal/middleware/auth.go`
- `backend/internal/middleware/rbac.go`
- `backend/internal/integration/wechat/auth.go`
- `backend/internal/integration/sms/sms.go`
- `backend/internal/modules/auth/model.go`
- `backend/internal/modules/auth/repository.go`
- `backend/internal/modules/auth/service.go`
- `backend/internal/modules/auth/handler.go`
- `backend/internal/modules/auth/dto.go`
- `backend/internal/modules/auth/routes.go`

#### 1.4 小程序项目初始化

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P1-018 | uni-app 项目搭建 | miniapp | 使用 HBuilderX / CLI 创建 uni-app Vue3 + TypeScript 项目；配置 `manifest.json`（AppID）、`pages.json`（TabBar 路由）、`vite.config.ts`；引入 Pinia、scss 全局变量 | 无 | 0.5 |
| P1-019 | 请求层封装 | miniapp | 实现 `src/api/request.ts`：uni.request 封装、拦截器（自动 token 注入、401 自动刷新、统一错误提示）；定义 `src/types/api.d.ts` 基础响应类型 | P1-018 | 0.5 |
| P1-020 | 通用组件 | miniapp | 实现基础通用组件：`t-navbar`（自定义导航栏，适配状态栏）、`t-tabbar`（底部 Tab）、`t-loading`（加载态）、`t-empty`（空态） | P1-018 | 1 |
| P1-021 | 微信登录流程 | miniapp | 实现 `src/api/auth.ts`；实现 `src/composables/useAuth.ts`：wx.login 获取 code、调用后端登录接口、Token 存储（Storage）、自动登录；实现 `src/stores/user.ts` | P1-019, P1-017 | 1 |
| P1-022 | 个人中心页面 | miniapp | 实现 `src/pages/profile/index.vue`：头像昵称展示、菜单入口（我的预约/我的订单/我的课程/钱包/会员/设置）；未登录态引导登录 | P1-021, P1-020 | 0.5 |

**涉及文件**:
- `miniapp/src/main.ts`
- `miniapp/src/App.vue`
- `miniapp/src/manifest.json`
- `miniapp/src/pages.json`
- `miniapp/src/uni.scss`
- `miniapp/src/api/request.ts`
- `miniapp/src/api/auth.ts`
- `miniapp/src/types/api.d.ts`
- `miniapp/src/composables/useAuth.ts`
- `miniapp/src/stores/user.ts`
- `miniapp/src/components/t-navbar/index.vue`
- `miniapp/src/components/t-tabbar/index.vue`
- `miniapp/src/components/t-loading/index.vue`
- `miniapp/src/components/t-empty/index.vue`
- `miniapp/src/pages/profile/index.vue`
- `miniapp/package.json`
- `miniapp/tsconfig.json`
- `miniapp/vite.config.ts`

#### 1.5 管理端项目初始化

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P1-023 | Arco Design Pro 项目搭建 | admin | 使用 Arco Design Pro 脚手架创建项目；配置路由结构、布局（侧边栏 + 顶栏）、Pinia 状态管理 | 无 | 0.5 |
| P1-024 | 请求层 + 登录页 | admin | 实现 `src/api/request.ts`：Axios 封装、拦截器、Token 管理；实现 `src/views/login/index.vue`：管理员账号密码登录；实现 `src/router/guards.ts`：路由守卫 | P1-023, P1-017 | 1 |

**涉及文件**:
- `admin/src/main.ts`
- `admin/src/App.vue`
- `admin/src/api/request.ts`
- `admin/src/views/login/index.vue`
- `admin/src/layouts/default.vue`
- `admin/src/layouts/login.vue`
- `admin/src/router/index.ts`
- `admin/src/router/guards.ts`
- `admin/src/stores/user.ts`
- `admin/src/stores/app.ts`
- `admin/package.json`
- `admin/tsconfig.json`
- `admin/vite.config.ts`

---

### Phase 1 里程碑验收标准

- [ ] `make dev` 一键启动后端 + PostgreSQL + Redis
- [ ] 后端 `/api/v1/health` 返回 200，数据库和 Redis 连接正常
- [ ] 小程序可通过微信授权登录，获取到 JWT Token
- [ ] 小程序个人中心页面展示用户头像和昵称
- [ ] 管理端可通过账号密码登录，进入主界面
- [ ] 三端联调跑通完整登录流程
- [ ] 单元测试覆盖率 >= 80%（auth 模块 service + repository）

---

### Phase 2：核心业务 -- 场地预订（18 天）

**目标**: 实现场地管理、时段规则配置、场地预订（含候补）、订单管理、支付（余额+微信）的完整闭环流程。

#### 2.1 场地管理模块

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P2-001 | 场地域数据库迁移 | backend | 编写 `migrations/000002_init_venues.up.sql`：创建 venues / venue_time_slot_rules / venue_blocked_times 表及索引 | P1-004 | 0.5 |
| P2-002 | Venue 模块 -- Model + Repository | backend | 实现场地 GORM 模型（Venue / VenueTimeSlotRule / VenueBlockedTime）；Repository 实现 CRUD、时段规则查询、可用性查询 | P2-001 | 1 |
| P2-003 | Venue 模块 -- Service | backend | 实现场地 CRUD 业务逻辑；时段规则管理；场地可用性查询（结合 Redis 缓存、blocked_times、已有预订过滤）；Write-through 缓存更新 | P2-002, P1-004 | 1.5 |
| P2-004 | Venue 模块 -- Handler + Routes | backend | 实现用户端 API（场地列表、详情、可用时段查询）和管理端 API（场地 CRUD、时段规则 CRUD、时段视图）；DTO 定义 | P2-003 | 1 |

**涉及文件**:
- `backend/migrations/000002_init_venues.up.sql`
- `backend/migrations/000002_init_venues.down.sql`
- `backend/internal/modules/venue/model.go`
- `backend/internal/modules/venue/repository.go`
- `backend/internal/modules/venue/service.go`
- `backend/internal/modules/venue/handler.go`
- `backend/internal/modules/venue/dto.go`
- `backend/internal/modules/venue/routes.go`

#### 2.2 订单模块

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P2-005 | 订单域数据库迁移 | backend | 创建 orders / order_items / refunds 表及索引 | P1-004 | 0.5 |
| P2-006 | Order 模块 -- Model + Repository | backend | Order / OrderItem / Refund GORM 模型；订单号生成（YYYYMMDDHHmmss + 随机数）；Repository 实现创建、查询、状态更新 | P2-005 | 1 |
| P2-007 | Order 模块 -- 状态机 | backend | 实现 `state_machine.go`：订单状态流转规则（pending -> paid -> used -> completed / cancelled / refunding -> refunded）；禁止非法状态转换 | P2-006 | 0.5 |
| P2-008 | Order 模块 -- Service + Handler | backend | 创建订单、订单列表（分页筛选）、订单详情、取消订单、管理端订单列表/退款审核；DTO + Routes | P2-007 | 1.5 |

**涉及文件**:
- `backend/migrations/000003_init_orders.up.sql`
- `backend/internal/modules/order/model.go`
- `backend/internal/modules/order/repository.go`
- `backend/internal/modules/order/state_machine.go`
- `backend/internal/modules/order/service.go`
- `backend/internal/modules/order/handler.go`
- `backend/internal/modules/order/handler_admin.go`
- `backend/internal/modules/order/dto.go`
- `backend/internal/modules/order/routes.go`

#### 2.3 支付模块

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P2-009 | 支付域数据库迁移 | backend | 创建 payments / coupons / user_coupons 表及索引 | P1-004 | 0.5 |
| P2-010 | 微信支付 SDK | backend | 实现 `internal/integration/wechat/pay.go`：统一下单、支付回调验签、退款申请、订单查询 | P1-002 | 1.5 |
| P2-011 | Payment 模块 -- Model + Repository | backend | Payment / Coupon / UserCoupon GORM 模型；支付流水号生成；Repository | P2-009 | 0.5 |
| P2-012 | Payment 模块 -- 组合支付 Service | backend | 实现 `service_combo.go`：纯余额支付（乐观锁扣减）；组合支付（冻结余额 + 微信补差）；微信回调处理（幂等）；失败解冻回滚；退款处理（余额退回 + 微信退款） | P2-011, P2-010, P1-015 | 2 |
| P2-013 | Payment 模块 -- Handler + Routes | backend | 预支付接口、微信回调接口、支付状态查询、优惠券列表/可用优惠券查询 | P2-012 | 0.5 |

**涉及文件**:
- `backend/migrations/000004_init_payments.up.sql`
- `backend/internal/integration/wechat/pay.go`
- `backend/internal/modules/payment/model.go`
- `backend/internal/modules/payment/repository.go`
- `backend/internal/modules/payment/service.go`
- `backend/internal/modules/payment/service_combo.go`
- `backend/internal/modules/payment/handler.go`
- `backend/internal/modules/payment/dto.go`
- `backend/internal/modules/payment/routes.go`

#### 2.4 预订模块

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P2-014 | 预订域数据库迁移 | backend | 创建 bookings / booking_waitlist 表及索引 | P1-004 | 0.5 |
| P2-015 | Booking 模块 -- Model + Repository | backend | Booking / BookingWaitlist GORM 模型；Repository 实现预订创建、冲突检查、候补队列操作 | P2-014 | 1 |
| P2-016 | Booking 模块 -- 预订 Service | backend | 核心预订流程：Redis 时段检查 -> SETNX 加锁 -> 创建订单 -> 调用支付 -> 确认/回滚；价格计算（时段规则 + 会员折扣 - 优惠券）；取消预订（退款规则：24h/12h/即时） | P2-015, P2-003, P2-008, P2-012 | 2 |
| P2-017 | Booking 模块 -- 候补 Service | backend | 实现 `service_waitlist.go`：加入候补（Redis Sorted Set）；取消触发候补通知；候补确认/超时处理；队列位置查询 | P2-016 | 1 |
| P2-018 | Booking 模块 -- Handler + Routes | backend | 用户端：创建预订、我的预订列表、预订详情、取消预订、加入候补；管理端：预订列表、时段视图 | P2-017 | 1 |

**涉及文件**:
- `backend/migrations/000005_init_bookings.up.sql`
- `backend/internal/modules/booking/model.go`
- `backend/internal/modules/booking/repository.go`
- `backend/internal/modules/booking/service.go`
- `backend/internal/modules/booking/service_waitlist.go`
- `backend/internal/modules/booking/handler.go`
- `backend/internal/modules/booking/handler_admin.go`
- `backend/internal/modules/booking/dto.go`
- `backend/internal/modules/booking/routes.go`

#### 2.5 定时任务

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P2-019 | 定时任务调度器 | backend | 实现 `internal/scheduler/scheduler.go`：go-cron 初始化、任务注册 | P1-007 | 0.5 |
| P2-020 | 订单超时关闭 | backend | 实现 `scheduler/order_expire.go`：每分钟扫描 `orders WHERE status='pending' AND expires_at < now()`，自动关闭并释放锁/库存 | P2-019, P2-008 | 0.5 |
| P2-021 | 候补超时检查 | backend | 实现 `scheduler/booking_expire.go`：每分钟扫描已通知但超时的候补记录，取消并通知下一位 | P2-019, P2-017 | 0.5 |

**涉及文件**:
- `backend/internal/scheduler/scheduler.go`
- `backend/internal/scheduler/order_expire.go`
- `backend/internal/scheduler/booking_expire.go`

#### 2.6 小程序 -- 预订前端

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P2-022 | 首页 | miniapp | 实现 `src/pages/home/index.vue`：Banner 轮播、快捷入口（预订/活动/商城）、今日可用场地概览、近期活动推荐 | P1-020 | 1 |
| P2-023 | 场地预订页 | miniapp | 实现 `src/pages/booking/index.vue`：日期选择器（横向滚动周历）、场地列表卡片、时间轴视图（显示各时段状态：可订/已订/候补/维护）；实现 `src/api/venue.ts` + `src/api/booking.ts` | P2-004, P1-019 | 2 |
| P2-024 | 预订确认页 | miniapp | 实现 `src/pages/booking/confirm.vue`：预订信息摘要、价格明细（原价/会员折扣/优惠券）、优惠券选择弹窗、支付方式选择（余额/微信/组合）、提交预订 | P2-023 | 1 |
| P2-025 | 支付流程 | miniapp | 实现 `src/composables/usePayment.ts`：调用预支付接口 -> wx.requestPayment 调起微信支付 -> 轮询支付结果 -> 跳转结果页；支付成功/失败/超时处理 | P2-024, P2-013 | 1 |
| P2-026 | 我的预约页 | miniapp | 实现 `src/pages/my-bookings/index.vue`：Tab 切换（即将到来/历史记录）、预订卡片（场地信息/时间/状态/取消按钮）；取消预订弹窗确认 + 退款规则提示 | P2-018 | 1 |
| P2-027 | 订单模块页面 | miniapp | 实现 `src/pages/orders/index.vue`：Tab 筛选（全部/待支付/已支付/已完成/退款）；`src/pages/orders/detail.vue`：订单详情、支付/取消操作 | P2-018, P2-025 | 1 |

**涉及文件**:
- `miniapp/src/pages/home/index.vue`
- `miniapp/src/pages/booking/index.vue`
- `miniapp/src/pages/booking/confirm.vue`
- `miniapp/src/pages/my-bookings/index.vue`
- `miniapp/src/pages/orders/index.vue`
- `miniapp/src/pages/orders/detail.vue`
- `miniapp/src/api/venue.ts`
- `miniapp/src/api/booking.ts`
- `miniapp/src/api/order.ts`
- `miniapp/src/api/payment.ts`
- `miniapp/src/composables/useBooking.ts`
- `miniapp/src/composables/usePayment.ts`
- `miniapp/src/composables/usePagination.ts`
- `miniapp/src/types/venue.d.ts`
- `miniapp/src/types/booking.d.ts`
- `miniapp/src/types/order.d.ts`
- `miniapp/src/components/t-card/index.vue`
- `miniapp/src/components/t-price/index.vue`
- `miniapp/src/components/t-badge/index.vue`
- `miniapp/src/utils/date.ts`
- `miniapp/src/utils/price.ts`

#### 2.7 管理端 -- 场地与订单

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P2-028 | 通用组件 | admin | 实现 `SearchForm`（搜索筛选表单封装）、`DataTable`（表格 + 分页封装）、`StatusTag`（状态标签）、`PriceDisplay`（金额展示）；实现 `useTable.ts` + `useCrud.ts` composable | P1-024 | 1 |
| P2-029 | 场地管理页面 | admin | `src/views/venue/list.vue`：场地列表表格（名称/类型/状态/操作）；`form.vue`：场地新增/编辑表单（含图片上传）；`time-rules.vue`：时段定价规则配置（工作日/周末/节假日）；`time-grid.vue`：日历 + 时段格子视图 | P2-028, P2-004 | 2 |
| P2-030 | 订单管理页面 | admin | `src/views/order/list.vue`：订单列表（筛选：类型/状态/日期范围/关键词）；`detail.vue`：订单详情、退款审核操作 | P2-028, P2-008 | 1 |

**涉及文件**:
- `admin/src/components/SearchForm/index.vue`
- `admin/src/components/DataTable/index.vue`
- `admin/src/components/StatusTag/index.vue`
- `admin/src/components/PriceDisplay/index.vue`
- `admin/src/composables/useTable.ts`
- `admin/src/composables/useCrud.ts`
- `admin/src/views/venue/list.vue`
- `admin/src/views/venue/form.vue`
- `admin/src/views/venue/time-rules.vue`
- `admin/src/views/venue/time-grid.vue`
- `admin/src/views/order/list.vue`
- `admin/src/views/order/detail.vue`
- `admin/src/api/venue.ts`
- `admin/src/api/order.ts`
- `admin/src/router/routes/venue.ts`
- `admin/src/router/routes/order.ts`

---

### Phase 2 里程碑验收标准

- [ ] 管理员可在后台创建场地、配置时段定价规则
- [ ] 用户在小程序可查看场地可用时段，选择时段下单
- [ ] 余额支付和微信支付完整流程跑通
- [ ] 组合支付（余额+微信）流程跑通
- [ ] 预订成功后场地时段状态正确更新
- [ ] 取消预订退款规则生效（24h/12h/即时）
- [ ] 候补排队功能：加入候补 -> 有人取消 -> 通知候补 -> 确认/超时
- [ ] 订单超时（30 分钟）自动关闭
- [ ] 管理端场地时段视图正确展示
- [ ] 管理端订单列表筛选和退款审核功能正常

---

### Phase 3：商品 + 活动（14 天）

**目标**: 实现商品管理（课程/球具/服务）和活动管理（畅打/团课/赛事）模块，含活动自动取消逻辑。

#### 3.1 商品模块

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P3-001 | 商品域数据库迁移 | backend | 创建 products / product_skus 表及索引 | P1-004 | 0.5 |
| P3-002 | 文件上传模块 | backend | 实现 `internal/pkg/upload/oss.go`：阿里云 OSS / 腾讯 COS 上传封装、图片压缩、签名 URL 生成 | P1-002 | 1 |
| P3-003 | Product 模块 -- 后端完整实现 | backend | Model + Repository + Service + Handler + Routes：商品 CRUD、分类/子分类筛选、上下架管理、SKU 管理、库存扣减（Redis 预扣 + DB 最终一致） | P3-001, P3-002, P2-008 | 2 |
| P3-004 | 商城页面 -- 小程序 | miniapp | `src/pages/shop/index.vue`：分类 Tab 切换（课程/球具/服务）、商品卡片网格；`detail.vue`：商品详情、SKU 选择、加入购买 | P3-003, P2-025 | 1.5 |
| P3-005 | 商品管理 -- 管理端 | admin | `src/views/product/list.vue`：商品列表（分类筛选/上下架/库存）；`form.vue`：商品编辑（富文本描述、多图上传、SKU 配置） | P3-003, P2-028 | 1.5 |

**涉及文件**:
- `backend/migrations/000006_init_products.up.sql`
- `backend/internal/pkg/upload/oss.go`
- `backend/internal/modules/product/model.go`
- `backend/internal/modules/product/repository.go`
- `backend/internal/modules/product/service.go`
- `backend/internal/modules/product/handler.go`
- `backend/internal/modules/product/handler_admin.go`
- `backend/internal/modules/product/dto.go`
- `backend/internal/modules/product/routes.go`
- `miniapp/src/pages/shop/index.vue`
- `miniapp/src/pages/shop/detail.vue`
- `miniapp/src/api/product.ts`
- `miniapp/src/types/product.d.ts`
- `admin/src/views/product/list.vue`
- `admin/src/views/product/form.vue`
- `admin/src/api/product.ts`
- `admin/src/router/routes/product.ts`

#### 3.2 活动模块

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P3-006 | 活动域数据库迁移 | backend | 创建 activities / activity_registrations 表及索引 | P1-004 | 0.5 |
| P3-007 | Activity 模块 -- 后端完整实现 | backend | Model + Repository + Service + Handler：活动 CRUD、报名（Redis 计数器原子递增 + 上限检测）、取消报名退款、活动状态流转（draft -> published -> registration -> confirmed/cancelled -> completed） | P3-006, P2-012, P2-003 | 2.5 |
| P3-008 | Activity 模块 -- 自动取消 | backend | 实现 `service_auto_cancel.go`：低于最低人数自动取消逻辑；`scheduler/activity_cancel.go`：定时扫描到期活动；批量退款、释放场地、通知用户 | P3-007, P2-019 | 1 |
| P3-009 | 通知模块 -- 基础实现 | backend | 实现 `internal/modules/notify/`：Notification GORM 模型、站内消息 CRUD、标记已读/全部已读；实现 `internal/integration/wechat/message.go`：微信模板消息/订阅消息推送 | P1-013, P1-004 | 1 |
| P3-010 | 活动页面 -- 小程序 | miniapp | `src/pages/activity/index.vue`：活动列表（类型筛选 Tab、状态标签、报名进度条）；`detail.vue`：活动详情、报名按钮、报名人数/剩余名额、水平要求展示 | P3-007, P2-025 | 1.5 |
| P3-011 | 钱包页面 -- 小程序 | miniapp | `src/pages/wallet/index.vue`：余额展示、充值入口、流水记录列表（分页加载）；`src/api/wallet.ts` | P1-016, P2-025 | 1 |
| P3-012 | 活动管理 -- 管理端 | admin | `src/views/activity/list.vue`：活动列表（类型/状态筛选、报名人数统计）；`form.vue`：活动编辑（场地选择、时间设置、人数设置、自动取消配置） | P3-007, P2-028 | 1 |

**涉及文件**:
- `backend/migrations/000007_init_activities.up.sql`
- `backend/migrations/000008_init_notifications.up.sql`
- `backend/internal/modules/activity/model.go`
- `backend/internal/modules/activity/repository.go`
- `backend/internal/modules/activity/service.go`
- `backend/internal/modules/activity/service_auto_cancel.go`
- `backend/internal/modules/activity/handler.go`
- `backend/internal/modules/activity/handler_admin.go`
- `backend/internal/modules/activity/dto.go`
- `backend/internal/modules/activity/routes.go`
- `backend/internal/scheduler/activity_cancel.go`
- `backend/internal/modules/notify/model.go`
- `backend/internal/modules/notify/repository.go`
- `backend/internal/modules/notify/service.go`
- `backend/internal/modules/notify/handler.go`
- `backend/internal/modules/notify/routes.go`
- `backend/internal/integration/wechat/message.go`
- `miniapp/src/pages/activity/index.vue`
- `miniapp/src/pages/activity/detail.vue`
- `miniapp/src/pages/wallet/index.vue`
- `miniapp/src/api/activity.ts`
- `miniapp/src/api/wallet.ts`
- `admin/src/views/activity/list.vue`
- `admin/src/views/activity/form.vue`
- `admin/src/api/activity.ts`
- `admin/src/router/routes/activity.ts`

---

### Phase 3 里程碑验收标准

- [ ] 管理员可创建/编辑/上下架商品（含多 SKU）
- [ ] 用户可浏览商品并完成购买支付
- [ ] 库存扣减正确（Redis 预扣 + DB 一致）
- [ ] 管理员可发布活动，配置自动取消条件
- [ ] 用户可报名活动并完成支付
- [ ] 活动报名人数达到上限后无法继续报名
- [ ] 活动人数不足自动取消，批量退款正确执行
- [ ] 站内通知功能正常（已读/未读/全部已读）
- [ ] 微信模板消息推送正常（预订确认/活动取消等场景）

---

### Phase 4：教务系统（12 天）

**目标**: 实现教练管理、课程排期（含冲突检测）、重复课程、教练休假/代课、学员课程记录与教练反馈。

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P4-001 | 教务域数据库迁移 | backend | 创建 coaches / course_schedules / coach_leaves / student_course_records 表及索引 | P1-004 | 0.5 |
| P4-002 | Academic 模块 -- Model + Repository | backend | GORM 模型定义；Repository：教练 CRUD、排期 CRUD（含批量创建）、休假管理、学员记录查询 | P4-001 | 1 |
| P4-003 | Academic 模块 -- 冲突检测 Service | backend | 实现 `service_conflict.go`：教练冲突检测（同时段已有课程）、场地冲突检测（同时段已被占用）、教练休假检测；支持"试算"和"确认"两种模式 | P4-002, P2-003 | 1.5 |
| P4-004 | Academic 模块 -- 排期 Service | backend | RRULE 解析（展开重复规则为日期列表）；批量创建排期（同 recurrence_group_id）；代课指定；课程取消/批量取消整组 | P4-003 | 1.5 |
| P4-005 | Academic 模块 -- 教练绩效 & 学员记录 | backend | 教练绩效统计（课时数/学员数/评分）；学员出勤记录；教练课后反馈；学员评价 | P4-002 | 1 |
| P4-006 | Academic 模块 -- Handler + Routes | backend | 管理端 API：教练 CRUD、排期 CRUD、冲突检测、休假审批、学员列表/记录；用户端 API：我的课程列表、课程详情、评价教练 | P4-004, P4-005 | 1 |
| P4-007 | 我的课程 -- 小程序 | miniapp | `src/pages/my-courses/index.vue`：课程日历视图、课程卡片（教练/场地/时间/状态）、课后评价入口 | P4-006 | 1 |
| P4-008 | 教练管理 -- 管理端 | admin | `src/views/academic/coaches.vue`：教练列表 + 新增/编辑；`coach-detail.vue`：教练详情 + 绩效统计图表 | P4-006, P2-028 | 1 |
| P4-009 | 课程排期 -- 管理端 | admin | `src/views/academic/schedules.vue`：周视图课表（横轴时间、纵轴场地/教练）、拖拽创建/修改排期、冲突高亮提示、重复规则配置弹窗 | P4-006, P2-028 | 2 |
| P4-010 | 学员管理 -- 管理端 | admin | `src/views/academic/students.vue`：学员列表（搜索/筛选）、学员课程记录、出勤统计 | P4-006, P2-028 | 0.5 |

**涉及文件**:
- `backend/migrations/000009_init_academic.up.sql`
- `backend/internal/modules/academic/model.go`
- `backend/internal/modules/academic/repository.go`
- `backend/internal/modules/academic/service.go`
- `backend/internal/modules/academic/service_conflict.go`
- `backend/internal/modules/academic/handler.go`
- `backend/internal/modules/academic/dto.go`
- `backend/internal/modules/academic/routes.go`
- `miniapp/src/pages/my-courses/index.vue`
- `admin/src/views/academic/schedules.vue`
- `admin/src/views/academic/coaches.vue`
- `admin/src/views/academic/coach-detail.vue`
- `admin/src/views/academic/students.vue`
- `admin/src/api/academic.ts`
- `admin/src/router/routes/academic.ts`

---

### Phase 4 里程碑验收标准

- [ ] 管理员可添加教练并配置基本信息
- [ ] 课程排期创建时冲突检测正确（教练/场地/休假）
- [ ] 重复课程批量创建正常，可按组取消
- [ ] 教练休假审批流程正常
- [ ] 代课指定功能正常
- [ ] 学员课程记录出勤标记正常
- [ ] 教练课后反馈功能正常
- [ ] 小程序"我的课程"页面展示正确
- [ ] 管理端课程排期周视图展示正确

---

### Phase 5：数据统计 + 优化（8 天）

**目标**: 实现管理端仪表盘、各维度数据统计、性能优化。

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P5-001 | 统计聚合定时任务 | backend | 实现 `scheduler/stats_aggregate.go`：每日凌晨聚合前一天数据写入 daily_stat / revenue_stat / venue_usage_stat；仪表盘数据缓存到 Redis（5 分钟 TTL） | P2-019 | 1 |
| P5-002 | Stats 模块 -- 后端实现 | backend | Model + Repository + Service + Handler：仪表盘聚合数据（今日收入/订单数/新用户/场地使用率）、收入构成分析（按订单类型/时间范围）、场地使用率（热力图数据）、用户画像统计 | P5-001 | 2 |
| P5-003 | 仪表盘 -- 管理端 | admin | `src/views/dashboard/index.vue`：数据卡片（今日收入/订单/用户/场地使用率）、收入趋势折线图、场地使用热力图、近期订单列表、待办事项（待审核退款/待确认活动） | P5-002, P2-028 | 2 |
| P5-004 | 数据统计页 -- 管理端 | admin | `src/views/stats/index.vue`：多维度统计（时间范围选择器）、收入构成饼图、场地使用率柱状图、用户增长曲线、教练绩效排行 | P5-002 | 1 |
| P5-005 | 性能优化 | backend | 数据库慢查询分析和索引优化；Redis 缓存命中率检查；API 响应时间 P99 优化；N+1 查询消除（GORM Preload）；分页查询优化 | P2-018, P3-007 | 1 |
| P5-006 | 小程序性能优化 | miniapp | 首屏加载优化（分包加载、图片懒加载、骨架屏）；接口聚合（首页一次请求获取所有数据）；列表虚拟滚动；本地缓存策略 | P2-022 | 1 |

**涉及文件**:
- `backend/internal/scheduler/stats_aggregate.go`
- `backend/internal/modules/stats/model.go`
- `backend/internal/modules/stats/repository.go`
- `backend/internal/modules/stats/service.go`
- `backend/internal/modules/stats/handler.go`
- `backend/internal/modules/stats/routes.go`
- `admin/src/views/dashboard/index.vue`
- `admin/src/views/stats/index.vue`
- `admin/src/api/dashboard.ts`
- `admin/src/api/stats.ts`
- `admin/src/components/ChartCard/index.vue`
- `admin/src/router/routes/dashboard.ts`
- `admin/src/router/routes/stats.ts`

---

### Phase 5 里程碑验收标准

- [ ] 仪表盘数据正确展示（收入/订单/用户/场地使用率）
- [ ] 统计图表交互正常（时间范围筛选、数据联动）
- [ ] 场地使用热力图数据正确
- [ ] API P99 响应时间达标（参照架构文档附录 A 指标）
- [ ] 小程序首屏加载 < 1.5s
- [ ] 场地可用性查询 < 200ms

---

### Phase 6：部署上线（6 天）

**目标**: 完成 Docker 化部署、CI/CD 配置、微信小程序审核提交、生产环境验证。

| 编号 | 任务名称 | 模块 | 工作内容 | 前置依赖 | 预估(天) |
|------|---------|------|---------|---------|---------|
| P6-001 | 后端 Docker 化 | backend | 优化 `backend/Dockerfile`（多阶段构建，最终镜像 < 30MB）；生产 `docker-compose.prod.yml`（PostgreSQL + Redis + Go API + Nginx） | P5-005 | 0.5 |
| P6-002 | Nginx 配置 | infra | 编写 Nginx 配置：SSL 终止（Let's Encrypt）、反向代理 API、静态资源缓存策略、GZIP 压缩、安全 Headers（HSTS / CSP） | P6-001 | 0.5 |
| P6-003 | CI/CD 流水线 | infra | 编写 `.github/workflows/backend.yml`：lint + test + build + push image；`miniapp.yml`：lint + build；`admin.yml`：lint + build + deploy to CDN/Nginx；智能触发（只有子项目变更才触发对应流水线） | P6-001 | 1 |
| P6-004 | 数据库生产初始化 | backend | 生产环境数据库迁移脚本；种子数据（管理员账号、初始场地、商品分类）；备份策略配置（pg_dump 每日全量 + WAL 增量） | P6-001 | 0.5 |
| P6-005 | 管理端部署 | admin | 构建生产包（Vite build）；部署到 Nginx / CDN；环境变量配置（API 地址） | P6-002, P6-003 | 0.5 |
| P6-006 | 小程序提审 | miniapp | 微信公众平台配置（服务器域名、业务域名）；小程序审核素材准备（截图、功能描述）；提交审核 | P5-006 | 1 |
| P6-007 | 安全加固 | all | 生产环境安全检查清单：HTTPS 验证、敏感信息加密存储、日志脱敏、限流配置、SQL 注入防护验证、微信回调签名验证 | P6-001 | 0.5 |
| P6-008 | 监控与告警 | infra | 健康检查端点 `/api/v1/health` 完善；Prometheus metrics 暴露；基础告警规则（错误率/响应延迟/DB 连接）；错误追踪（Sentry）接入 | P6-001 | 1 |
| P6-009 | 上线验证 | all | 生产环境全流程冒烟测试：注册登录 -> 预订 -> 支付 -> 取消退款 -> 活动报名 -> 商品购买；压力测试（场地预订并发场景） | P6-006, P6-005 | 0.5 |

**涉及文件**:
- `backend/Dockerfile`
- `docker-compose.prod.yml`
- `nginx/nginx.conf`
- `nginx/ssl/`
- `.github/workflows/backend.yml`
- `.github/workflows/miniapp.yml`
- `.github/workflows/admin.yml`
- `backend/scripts/seed.go`
- `backend/scripts/migrate.sh`
- `backend/scripts/backup.sh`

---

### Phase 6 里程碑验收标准

- [ ] 生产环境 Docker 部署一键启动
- [ ] HTTPS 全站生效，HSTS 配置正确
- [ ] CI/CD 流水线自动触发（push to main）
- [ ] 数据库自动备份任务运行正常
- [ ] 管理端 Web 可通过域名正常访问
- [ ] 微信小程序审核通过并发布
- [ ] 生产环境冒烟测试全部通过
- [ ] 监控告警配置生效

---

## 关键里程碑时间线

```
Phase 1 (Day 1-12)     Phase 2 (Day 13-30)      Phase 3 (Day 31-44)
├─ 脚手架搭建            ├─ 场地+订单+支付          ├─ 商品+活动
├─ 用户认证              ├─ 预订+候补               ├─ 通知模块
├─ 三端联调登录           ├─ 定时任务                ├─ 自动取消
│                       ├─ 小程序预订页面            │
│                       └─ 管理端场地/订单           │
│                                                  │
Phase 4 (Day 45-56)     Phase 5 (Day 57-64)      Phase 6 (Day 65-70)
├─ 教练管理              ├─ 统计聚合                ├─ Docker 部署
├─ 课程排期+冲突检测      ├─ 仪表盘+统计图表         ├─ CI/CD
├─ 学员管理              ├─ 性能优化                ├─ 小程序提审
└─ 排期日历视图           └─ 小程序优化              └─ 上线验证
```

---

## MVP 优先级划分

### 必须 (MVP Core) -- Phase 1 + Phase 2 + Phase 6 核心

| 功能 | 说明 |
|------|------|
| 用户注册/登录 | 微信授权 + JWT |
| 场地查看与预订 | 核心业务 |
| 支付（余额 + 微信） | 收入闭环 |
| 订单管理 | 基础订单流转 |
| 管理端场地管理 | 运营基础 |
| 管理端订单管理 | 运营基础 |
| Docker 部署 + 小程序上线 | 可交付 |

### 重要 (MVP Extended) -- Phase 3

| 功能 | 说明 |
|------|------|
| 活动管理 | 网球馆核心运营场景 |
| 商品售卖 | 收入来源之一 |
| 通知消息 | 用户体验关键 |
| 候补排队 | 热门时段用户体验 |

### 可延后 (Post-MVP) -- Phase 4 + Phase 5

| 功能 | 说明 |
|------|------|
| 教务系统 | 复杂但非上线必须 |
| 数据统计仪表盘 | 运营分析，可后补 |
| 性能优化 | 初期用户量小，可迭代 |
| 教练绩效 | 教务的衍生功能 |
| 用户画像 | 数据积累后更有价值 |

---

## 风险清单

| # | 风险 | 影响等级 | 概率 | 应对策略 |
|---|------|---------|------|---------|
| R1 | **微信支付接入复杂** -- 商户号审核周期长（7-15 天）、证书配置易出错 | 高 | 高 | 提前申请微信支付商户号（Phase 1 期间即开始申请）；Phase 2 支付模块先用沙盒环境开发；准备纯余额支付作为降级方案 |
| R2 | **小程序审核不通过** -- 微信审核标准严格，可能多次打回 | 高 | 中 | 提前研读微信审核规范；避免敏感词和非合规功能；准备完善的功能截图和说明；预留 1 周审核缓冲期 |
| R3 | **场地预订并发超卖** -- Redis 锁失效或网络抖动导致同一时段被重复预订 | 高 | 低 | Redis SETNX + 过期时间；数据库 UNIQUE 约束兜底（venue_id + date + start_time + status）；降级方案使用 PostgreSQL FOR UPDATE |
| R4 | **组合支付一致性** -- 余额已冻结但微信支付失败，资金状态不一致 | 高 | 低 | 冻结-确认模式确保原子性；定时任务扫描异常支付记录；微信支付主动查询补偿；完善的流水审计 |
| R5 | **RRULE 重复课程解析复杂** -- 重复规则边界情况多（跨月/跨年/节假日） | 中 | 中 | 使用成熟的 RRULE 解析库（如 go-rrule）；限制展开范围（最多 12 周）；充分的边界测试用例 |
| R6 | **单人全栈开发周期过长** -- 70 天工期对 1 人团队压力大 | 中 | 高 | 严格按 MVP 优先级执行；Phase 4 可延后；UI 复用 Arco/uni-ui 现有组件减少自定义开发；关键模块先实现核心功能，非核心细节迭代补充 |
| R7 | **PostgreSQL 性能瓶颈** -- 预订高峰期数据库压力大 | 中 | 低 | 合理的索引策略（已在架构文档中设计）；Redis 缓存减少数据库查询；慢查询监控和优化；初期单机足够，后续可读写分离 |
| R8 | **微信 API 变更** -- 微信平台接口调整导致功能异常 | 低 | 低 | 微信 SDK 封装在 `internal/integration/wechat/`，隔离变更影响；关注微信开发者社区公告；接口调用增加错误监控 |

---

## 需要进一步明确的问题

### 问题 1：微信支付商户号是否已申请

微信支付商户号审核需要 7-15 个工作日，是 Phase 2 支付模块的硬性前置条件。

**推荐方案**：

- 方案 A：立即申请商户号，Phase 2 开始时应已审核通过，可直接对接
- 方案 B：Phase 2 先用微信支付沙盒环境开发，商户号审核通过后切换生产环境

**等待用户选择**：

```
请选择您偏好的方案，或提供其他建议：
[ ] 方案 A：商户号已申请/立即申请
[ ] 方案 B：先用沙盒开发
[ ] 其他方案：______________________
```

### 问题 2：对象存储选型

架构文档中提到阿里云 OSS 或腾讯 COS，两者功能接近。

**推荐方案**：

- 方案 A：阿里云 OSS -- 如果服务器在阿里云，内网传输免费，CDN 生态完善
- 方案 B：腾讯 COS -- 与微信生态更紧密，小程序内上传可能更顺畅

**等待用户选择**：

```
请选择您偏好的方案，或提供其他建议：
[ ] 方案 A：阿里云 OSS
[ ] 方案 B：腾讯 COS
[ ] 其他方案：______________________
```

### 问题 3：Phase 4 教务系统是否纳入 MVP

教务系统（排期/教练/学员）功能完整但复杂，占 12 天工期。如果网球馆初期不需要教务管理，可以延后到 v1.1 版本。

**推荐方案**：

- 方案 A：纳入 MVP -- 完整交付所有功能，但总工期 70 天
- 方案 B：延后到 v1.1 -- MVP 缩减至 58 天，先上线核心预订+商品+活动功能

**等待用户选择**：

```
请选择您偏好的方案，或提供其他建议：
[ ] 方案 A：纳入 MVP（70 天完整交付）
[ ] 方案 B：延后到 v1.1（58 天快速上线）
[ ] 其他方案：______________________
```

### 问题 4：短信服务提供商选择

手机号登录和通知需要短信服务。

**推荐方案**：

- 方案 A：阿里云短信 -- 与阿里云基础设施统一
- 方案 B：腾讯云短信 -- 与微信生态统一
- 方案 C：初期仅使用微信登录，不做手机号登录 -- 减少开发量和短信成本

**等待用户选择**：

```
请选择您偏好的方案，或提供其他建议：
[ ] 方案 A：阿里云短信
[ ] 方案 B：腾讯云短信
[ ] 方案 C：初期仅微信登录
[ ] 其他方案：______________________
```

---

## 用户反馈区域

请在此区域补充您对整体规划的意见和建议：

```
用户补充内容：

_______________________________________________

_______________________________________________

_______________________________________________
```
