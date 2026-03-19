# T-Line Sports Club 系统架构设计文档

> 版本：v1.0  
> 日期：2026-03-19  
> 状态：设计阶段  

---

## 目录

1. [系统架构总览](#1-系统架构总览)
2. [后端服务划分](#2-后端服务划分)
3. [数据库设计](#3-数据库设计)
4. [API 设计规范](#4-api-设计规范)
5. [项目目录结构](#5-项目目录结构)
6. [关键业务流程](#6-关键业务流程)
7. [架构决策记录](#7-架构决策记录)

---

## 1. 系统架构总览

### 1.1 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端层                                 │
│  ┌──────────────────┐          ┌──────────────────────────┐     │
│  │  微信小程序 (C端)  │          │  Web 管理后台 (B端)       │     │
│  │  uni-app + Vue3   │          │  Vue3 + Arco Design Pro  │     │
│  │  375px 适配        │          │  1440px 桌面端            │     │
│  └────────┬─────────┘          └────────────┬─────────────┘     │
└───────────┼──────────────────────────────────┼──────────────────┘
            │  HTTPS                           │  HTTPS
┌───────────┼──────────────────────────────────┼──────────────────┐
│           ▼                                  ▼                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Nginx 反向代理                         │    │
│  │           SSL 终止 / 负载均衡 / 静态资源                   │    │
│  └───────────────────────┬─────────────────────────────────┘    │
│                          │                                      │
│  ┌───────────────────────▼─────────────────────────────────┐    │
│  │              Go (Gin) REST API 服务                       │    │
│  │  ┌─────────┐ ┌─────────┐ ┌──────────┐ ┌─────────────┐  │    │
│  │  │ 用户模块 │ │ 场地模块 │ │ 商品模块  │ │  教务模块    │  │    │
│  │  ├─────────┤ ├─────────┤ ├──────────┤ ├─────────────┤  │    │
│  │  │ 活动模块 │ │ 订单模块 │ │ 支付模块  │ │  统计模块    │  │    │
│  │  └─────────┘ └─────────┘ └──────────┘ └─────────────┘  │    │
│  │  ┌─────────────────────────────────────────────────────┐│    │
│  │  │        共享层: 认证/日志/错误处理/中间件/定时任务      ││    │
│  │  └─────────────────────────────────────────────────────┘│    │
│  └───────┬─────────────┬──────────────┬────────────────────┘    │
│          │             │              │                          │
│  ┌───────▼───────┐ ┌──▼────────┐ ┌──▼───────────────┐          │
│  │  PostgreSQL   │ │   Redis    │ │  对象存储 (OSS)   │          │
│  │  主数据库      │ │  缓存/队列  │ │  图片/文件         │          │
│  └───────────────┘ └───────────┘ └──────────────────┘          │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                  外部服务集成                              │    │
│  │  微信开放平台 │ 微信支付 │ 短信服务 │ 微信模板消息         │    │
│  └─────────────────────────────────────────────────────────┘    │
│                         服务端层                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 技术选型

| 层级 | 技术 | 版本 | 选型理由 |
|------|------|------|---------|
| **用户端** | uni-app + Vue3 | Vue 3.4+ | 跨端能力强，一套代码编译微信小程序；生态成熟 |
| **管理端** | Vue3 + Arco Design Pro | Vue 3.4+ / Arco 2.x | Arco 组件库企业级质量，表格/表单/图表齐全；字节跳动开源维护活跃 |
| **后端框架** | Go + Gin | Go 1.22+ / Gin 1.10+ | 高并发性能优异，编译型语言部署简单，Gin 生态成熟 |
| **ORM** | GORM | v2 | Go 生态最流行 ORM，支持迁移/关联/钩子，文档完善 |
| **数据库** | PostgreSQL | 16+ | ACID 事务可靠，JSON 支持好，适合复杂业务查询 |
| **缓存** | Redis | 7+ | 场地时段锁定、会话管理、热数据缓存、简单消息队列 |
| **对象存储** | 阿里云 OSS / 腾讯 COS | - | 图片上传（场地图、商品图、头像），CDN 加速 |
| **认证** | JWT + 微信 OAuth | - | 无状态认证，小程序端微信授权登录 |
| **定时任务** | go-cron | v2 | 候补超时、活动自动取消、数据统计聚合 |
| **日志** | Zap + Lumberjack | - | 高性能结构化日志，支持日志轮转 |
| **部署** | Docker + Docker Compose | - | 初期单机部署，后续可扩展 K8s |

### 1.3 部署架构

**初期 (MVP / 1000 用户以内)**

```
┌─────────────────────────────────────────┐
│           云服务器 (4C8G)                 │
│  ┌───────────┐  ┌────────────────────┐  │
│  │  Nginx     │  │  Go API (单实例)    │  │
│  │  :80/:443  │──│  :8080             │  │
│  └───────────┘  └────────┬───────────┘  │
│                          │               │
│  ┌───────────┐  ┌────────▼───────────┐  │
│  │  Redis     │  │  PostgreSQL        │  │
│  │  :6379     │  │  :5432             │  │
│  └───────────┘  └────────────────────┘  │
└─────────────────────────────────────────┘

管理后台: 静态文件部署 (Nginx / CDN)
小程序: 微信平台托管
```

**扩展阶段 (10000+ 用户)**

```
                   ┌─────────────┐
                   │  CDN / SLB   │
                   └──────┬──────┘
                          │
              ┌───────────┼───────────┐
              │           │           │
        ┌─────▼─────┐ ┌──▼──────┐ ┌─▼────────┐
        │ Go API x2  │ │ Go API  │ │ Go API   │
        │ (场地/订单) │ │ (支付)  │ │ (其他)   │
        └─────┬─────┘ └──┬──────┘ └─┬────────┘
              │           │          │
        ┌─────▼───────────▼──────────▼────┐
        │     PostgreSQL (主从)            │
        │     Redis Cluster               │
        └─────────────────────────────────┘
```

---

## 2. 后端服务划分

采用 **模块化单体 (Modular Monolith)** 架构。单一 Go 二进制，内部按业务域严格分模块，模块间通过接口（interface）通信，为未来拆分微服务预留能力。

### 2.1 模块总览

```
┌──────────────────────────────────────────────────────────┐
│                       API 网关层                          │
│  路由注册 / 中间件链 / 请求验证 / 响应格式化               │
├──────────────────────────────────────────────────────────┤
│                       业务模块层                          │
│                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │  auth     │ │  venue    │ │  booking  │ │  order    │   │
│  │  用户认证  │ │  场地管理  │ │  预订管理  │ │  订单管理  │   │
│  ├──────────┤ ├──────────┤ ├──────────┤ ├──────────┤   │
│  │  payment  │ │  product  │ │  activity │ │  academic │   │
│  │  支付模块  │ │  商品管理  │ │  活动管理  │ │  教务管理  │   │
│  ├──────────┤ ├──────────┤                              │
│  │  stats    │ │  notify   │                              │
│  │  数据统计  │ │  通知消息  │                              │
│  └──────────┘ └──────────┘                              │
├──────────────────────────────────────────────────────────┤
│                       共享基础层                          │
│  config / logger / database / cache / errors / middleware │
│  scheduler / upload / wechat / sms                       │
└──────────────────────────────────────────────────────────┘
```

### 2.2 各模块职责

#### auth (用户认证模块)
- **职责**: 微信登录/手机号注册、JWT 签发与刷新、用户信息 CRUD、会员等级管理、积分管理、钱包余额管理
- **核心实体**: User, UserProfile, MemberLevel, Wallet, WalletTransaction, Points
- **对外接口**: 提供用户信息查询接口供其他模块调用
- **依赖**: wechat SDK, sms SDK, cache (会话)

#### venue (场地管理模块)
- **职责**: 场地信息 CRUD、场地状态管理、时段规则配置（工作日/周末/节假日定价）、场地可用性查询
- **核心实体**: Venue, VenueTimeSlotRule, VenueBlockedTime
- **对外接口**: 提供场地可用性检查接口供 booking/academic 模块调用
- **依赖**: cache (场地状态缓存)

#### booking (预订管理模块)
- **职责**: 场地预订创建、候补排队、预订取消与退款规则、预订状态流转、候补超时自动退款
- **核心实体**: Booking, BookingWaitlist
- **对外接口**: 消费 venue 的可用性接口，调用 order 创建订单，调用 payment 发起支付
- **依赖**: venue, order, payment, scheduler (候补超时), cache (时段锁)

#### order (订单管理模块)
- **职责**: 统一订单创建与管理、订单状态机（待支付/已支付/已使用/已完成/已取消/退款中/已退款）、退款审核、订单查询与筛选
- **核心实体**: Order, OrderItem, Refund
- **对外接口**: 提供创建订单、查询订单、更新状态接口
- **依赖**: payment (退款)

#### payment (支付模块)
- **职责**: 组合支付逻辑（余额优先 + 微信补差）、微信支付接口对接、退款处理、优惠券核销、支付回调验证
- **核心实体**: Payment, PaymentTransaction, Coupon, CouponUsage
- **对外接口**: 提供统一支付接口、退款接口
- **依赖**: auth (余额扣减), wechat pay SDK, order (回调更新)

#### product (商品管理模块)
- **职责**: 课程商品（单次/课包）、球具（现货/预定/邮寄）、服务（租赁/穿线）的 CRUD、上下架管理、库存管理、SKU 管理
- **核心实体**: Product, ProductCategory, ProductSKU, Inventory
- **对外接口**: 提供商品查询、库存扣减接口
- **依赖**: order (购买流程), upload (商品图片)

#### activity (活动管理模块)
- **职责**: 活动发布与编辑、活动类型管理（畅打/团课/赛事/主题活动）、报名与取消、最低人数自动取消、活动状态流转
- **核心实体**: Activity, ActivityRegistration, ActivityType
- **对外接口**: 提供活动查询、报名接口
- **依赖**: venue (场地占用), order, payment, scheduler (自动取消), notify (通知)

#### academic (教务管理模块)
- **职责**: 课程排期（含冲突检测: 教练+场地+时间）、重复周期管理、教练排班与休假/代课、学员课程记录、缺课管理、教练课后反馈、教练绩效统计
- **核心实体**: CourseSchedule, Coach, CoachSchedule, CoachLeave, StudentCourseRecord, CoachFeedback
- **冲突检测**: 检查同一时段教练是否已有课程、场地是否已被占用
- **依赖**: venue (场地可用性), auth (教练/学员用户), notify (课程通知)

#### stats (数据统计模块)
- **职责**: 仪表盘数据聚合、场地使用率统计（含热力图数据）、收入构成分析、用户画像、教练绩效排行、趋势分析
- **核心实体**: DailyStat, RevenueStat, VenueUsageStat
- **依赖**: 只读访问其他模块的数据，通过 scheduler 定期聚合

#### notify (通知消息模块)
- **职责**: 微信模板消息/订阅消息推送、短信通知、站内消息管理
- **核心实体**: Notification, NotificationTemplate
- **依赖**: wechat SDK, sms SDK

### 2.3 模块间依赖关系

```
auth ◄──────────── booking, order, payment, academic, activity
venue ◄─────────── booking, academic, activity
order ◄─────────── booking, product, activity
payment ◄───────── booking, order, activity
product ◄───────── order
scheduler ◄──────── booking (候补超时), activity (自动取消), stats (聚合)
notify ◄────────── booking, activity, academic, payment
```

**原则**: 依赖方向单向，禁止循环依赖。模块间通过 Go interface 解耦，上层模块依赖下层接口定义而非具体实现。

### 2.4 共享基础组件

| 组件 | 职责 |
|------|------|
| `config` | 配置加载（env / yaml），环境变量管理 |
| `logger` | Zap 结构化日志，请求 ID 追踪 |
| `database` | GORM 初始化，连接池配置，迁移管理 |
| `cache` | Redis 客户端封装，通用缓存/锁操作 |
| `errors` | 统一错误码定义，错误响应格式 |
| `middleware` | JWT 认证、RBAC 鉴权、请求日志、限流、CORS、Panic 恢复 |
| `scheduler` | 定时任务管理器（候补超时、活动取消、统计聚合） |
| `upload` | 对象存储封装，图片上传/压缩 |
| `wechat` | 微信登录 SDK、微信支付 SDK、模板消息 SDK |
| `sms` | 短信验证码发送 |
| `validator` | 请求参数校验封装 |
| `response` | 统一 API 响应格式封装 |

---

## 3. 数据库设计

### 3.1 核心表设计

#### 用户域 (auth)

**users** - 用户表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| phone | VARCHAR(20) UNIQUE | 手机号 |
| nickname | VARCHAR(64) | 昵称 |
| avatar_url | VARCHAR(512) | 头像 URL |
| wechat_openid | VARCHAR(128) UNIQUE | 微信 OpenID |
| wechat_unionid | VARCHAR(128) | 微信 UnionID |
| utr_rating | DECIMAL(4,2) | UTR 网球等级 |
| member_level | SMALLINT DEFAULT 0 | 会员等级 (0普通/1银/2金/3钻) |
| member_expires_at | TIMESTAMPTZ | 会员到期时间 |
| role | VARCHAR(20) DEFAULT 'user' | 角色 (user/coach/admin/super_admin) |
| status | SMALLINT DEFAULT 1 | 状态 (0禁用/1正常) |
| created_at | TIMESTAMPTZ | 创建时间 |
| updated_at | TIMESTAMPTZ | 更新时间 |

**wallets** - 钱包表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| user_id | BIGINT UNIQUE FK | 用户 ID |
| balance | DECIMAL(12,2) DEFAULT 0 | 余额（分为单位存储更安全，此处元） |
| frozen_amount | DECIMAL(12,2) DEFAULT 0 | 冻结金额 |
| total_recharged | DECIMAL(12,2) DEFAULT 0 | 累计充值 |
| version | INT DEFAULT 0 | 乐观锁版本号 |
| updated_at | TIMESTAMPTZ | 更新时间 |

**wallet_transactions** - 钱包流水表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| wallet_id | BIGINT FK | 钱包 ID |
| type | VARCHAR(20) | 类型 (recharge/consume/refund/freeze/unfreeze) |
| amount | DECIMAL(12,2) | 金额（正为入账，负为出账） |
| balance_after | DECIMAL(12,2) | 交易后余额 |
| related_order_id | BIGINT | 关联订单 ID |
| remark | VARCHAR(256) | 备注 |
| created_at | TIMESTAMPTZ | 创建时间 |

**points** - 积分表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| user_id | BIGINT UNIQUE FK | 用户 ID |
| balance | INT DEFAULT 0 | 当前积分 |
| total_earned | INT DEFAULT 0 | 累计获得 |
| total_spent | INT DEFAULT 0 | 累计消费 |
| updated_at | TIMESTAMPTZ | 更新时间 |

#### 场地域 (venue)

**venues** - 场地表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| name | VARCHAR(64) | 场地名称 |
| type | VARCHAR(20) | 类型 (indoor_hard/outdoor_hard/indoor_clay/outdoor_clay) |
| description | TEXT | 描述 |
| cover_image | VARCHAR(512) | 封面图 |
| facilities | JSONB | 设施信息 (灯光/空调/更衣室等) |
| status | SMALLINT DEFAULT 1 | 状态 (0停用/1正常/2维护中) |
| sort_order | INT DEFAULT 0 | 排序 |
| created_at | TIMESTAMPTZ | 创建时间 |
| updated_at | TIMESTAMPTZ | 更新时间 |

**venue_time_slot_rules** - 时段定价规则表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| venue_id | BIGINT FK | 场地 ID |
| day_type | VARCHAR(20) | 日期类型 (weekday/weekend/holiday) |
| start_time | TIME | 开始时间 |
| end_time | TIME | 结束时间 |
| price | DECIMAL(10,2) | 单价（元/小时） |
| member_discount | JSONB | 会员折扣 {"silver":0.95,"gold":0.8,"diamond":0.7} |
| is_active | BOOLEAN DEFAULT TRUE | 是否启用 |

**venue_blocked_times** - 场地屏蔽时间表（维护/特殊关闭）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| venue_id | BIGINT FK | 场地 ID |
| start_at | TIMESTAMPTZ | 开始时间 |
| end_at | TIMESTAMPTZ | 结束时间 |
| reason | VARCHAR(256) | 原因 |

#### 预订域 (booking)

**bookings** - 预订表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| user_id | BIGINT FK | 用户 ID |
| venue_id | BIGINT FK | 场地 ID |
| date | DATE | 预订日期 |
| start_time | TIME | 开始时间 |
| end_time | TIME | 结束时间 |
| duration_hours | DECIMAL(3,1) | 时长（小时） |
| status | VARCHAR(20) | 状态 (pending/confirmed/waitlisted/cancelled/completed/no_show) |
| total_amount | DECIMAL(10,2) | 总金额 |
| order_id | BIGINT FK | 关联订单 ID |
| cancel_reason | VARCHAR(256) | 取消原因 |
| cancelled_at | TIMESTAMPTZ | 取消时间 |
| created_at | TIMESTAMPTZ | 创建时间 |
| updated_at | TIMESTAMPTZ | 更新时间 |

**booking_waitlist** - 候补队列表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| booking_id | BIGINT FK | 预订 ID |
| user_id | BIGINT FK | 用户 ID |
| venue_id | BIGINT FK | 场地 ID |
| date | DATE | 日期 |
| start_time | TIME | 开始时间 |
| end_time | TIME | 结束时间 |
| position | INT | 队列位置 |
| status | VARCHAR(20) | 状态 (waiting/notified/confirmed/expired/cancelled) |
| notified_at | TIMESTAMPTZ | 通知时间 |
| expires_at | TIMESTAMPTZ | 超时时间（通知后 N 分钟） |
| created_at | TIMESTAMPTZ | 创建时间 |

#### 订单域 (order)

**orders** - 订单表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| order_no | VARCHAR(32) UNIQUE | 订单号 (YYYYMMDDHHmmss + 随机数) |
| user_id | BIGINT FK | 用户 ID |
| type | VARCHAR(20) | 订单类型 (booking/product/activity/recharge) |
| status | VARCHAR(20) | 状态 (pending/paid/used/completed/cancelled/refunding/refunded) |
| total_amount | DECIMAL(12,2) | 订单总额 |
| discount_amount | DECIMAL(12,2) DEFAULT 0 | 优惠金额 |
| pay_amount | DECIMAL(12,2) | 实付金额 |
| balance_paid | DECIMAL(12,2) DEFAULT 0 | 余额支付部分 |
| wechat_paid | DECIMAL(12,2) DEFAULT 0 | 微信支付部分 |
| coupon_id | BIGINT | 使用的优惠券 ID |
| remark | VARCHAR(512) | 备注 |
| paid_at | TIMESTAMPTZ | 支付时间 |
| expires_at | TIMESTAMPTZ | 支付超时时间 (未支付订单 30 分钟过期) |
| created_at | TIMESTAMPTZ | 创建时间 |
| updated_at | TIMESTAMPTZ | 更新时间 |

**order_items** - 订单明细表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| order_id | BIGINT FK | 订单 ID |
| item_type | VARCHAR(20) | 类型 (booking/product/activity/recharge) |
| item_id | BIGINT | 关联 ID (booking_id / product_id / activity_id) |
| item_name | VARCHAR(128) | 商品名称（快照） |
| sku_id | BIGINT | SKU ID |
| quantity | INT DEFAULT 1 | 数量 |
| unit_price | DECIMAL(10,2) | 单价 |
| total_price | DECIMAL(10,2) | 小计 |

**refunds** - 退款表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| refund_no | VARCHAR(32) UNIQUE | 退款单号 |
| order_id | BIGINT FK | 订单 ID |
| user_id | BIGINT FK | 用户 ID |
| amount | DECIMAL(12,2) | 退款金额 |
| balance_refund | DECIMAL(12,2) DEFAULT 0 | 退回余额部分 |
| wechat_refund | DECIMAL(12,2) DEFAULT 0 | 退回微信部分 |
| reason | VARCHAR(512) | 退款原因 |
| status | VARCHAR(20) | 状态 (pending/approved/rejected/processing/completed/failed) |
| reviewed_by | BIGINT | 审核人 |
| reviewed_at | TIMESTAMPTZ | 审核时间 |
| completed_at | TIMESTAMPTZ | 完成时间 |
| created_at | TIMESTAMPTZ | 创建时间 |

#### 支付域 (payment)

**payments** - 支付记录表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| payment_no | VARCHAR(64) UNIQUE | 支付流水号 |
| order_id | BIGINT FK | 订单 ID |
| user_id | BIGINT FK | 用户 ID |
| method | VARCHAR(20) | 支付方式 (balance/wechat/combo) |
| amount | DECIMAL(12,2) | 支付金额 |
| balance_amount | DECIMAL(12,2) DEFAULT 0 | 余额部分 |
| wechat_amount | DECIMAL(12,2) DEFAULT 0 | 微信部分 |
| wechat_trade_no | VARCHAR(64) | 微信交易号 |
| status | VARCHAR(20) | 状态 (pending/success/failed/refunded) |
| paid_at | TIMESTAMPTZ | 支付时间 |
| created_at | TIMESTAMPTZ | 创建时间 |

**coupons** - 优惠券表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| name | VARCHAR(64) | 优惠券名称 |
| type | VARCHAR(20) | 类型 (fixed/percent) |
| value | DECIMAL(10,2) | 面值/折扣率 |
| min_amount | DECIMAL(10,2) DEFAULT 0 | 最低消费 |
| applicable_types | JSONB | 适用订单类型 ["booking","product"] |
| total_count | INT | 发放总量 |
| used_count | INT DEFAULT 0 | 已使用数量 |
| start_at | TIMESTAMPTZ | 生效时间 |
| end_at | TIMESTAMPTZ | 失效时间 |
| status | SMALLINT DEFAULT 1 | 状态 |

**user_coupons** - 用户优惠券表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| user_id | BIGINT FK | 用户 ID |
| coupon_id | BIGINT FK | 优惠券 ID |
| status | SMALLINT DEFAULT 0 | 0未使用/1已使用/2已过期 |
| used_order_id | BIGINT | 使用时的订单 ID |
| expires_at | TIMESTAMPTZ | 过期时间 |
| created_at | TIMESTAMPTZ | 领取时间 |

#### 商品域 (product)

**products** - 商品表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| category | VARCHAR(20) | 分类 (course/equipment/service) |
| sub_category | VARCHAR(40) | 子分类 (single_course/course_pack/racket/ball/rental/stringing) |
| name | VARCHAR(128) | 商品名称 |
| description | TEXT | 描述 |
| cover_image | VARCHAR(512) | 封面图 |
| images | JSONB | 详情图列表 |
| price | DECIMAL(10,2) | 价格 |
| original_price | DECIMAL(10,2) | 原价 |
| stock | INT DEFAULT 0 | 库存 (-1 表示无限) |
| sales_count | INT DEFAULT 0 | 销量 |
| status | SMALLINT DEFAULT 0 | 0下架/1上架 |
| attributes | JSONB | 扩展属性 (课包次数/配送方式等) |
| sort_order | INT DEFAULT 0 | 排序 |
| created_at | TIMESTAMPTZ | 创建时间 |
| updated_at | TIMESTAMPTZ | 更新时间 |

**product_skus** - SKU 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| product_id | BIGINT FK | 商品 ID |
| name | VARCHAR(128) | SKU 名称 (如 "红色/42码") |
| price | DECIMAL(10,2) | SKU 价格 |
| stock | INT DEFAULT 0 | SKU 库存 |
| attributes | JSONB | SKU 属性 |

#### 活动域 (activity)

**activities** - 活动表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| title | VARCHAR(128) | 活动标题 |
| type | VARCHAR(20) | 类型 (open_play/group_class/tournament/themed) |
| description | TEXT | 描述 |
| cover_image | VARCHAR(512) | 封面图 |
| venue_id | BIGINT FK | 场地 ID |
| start_at | TIMESTAMPTZ | 开始时间 |
| end_at | TIMESTAMPTZ | 结束时间 |
| registration_deadline | TIMESTAMPTZ | 报名截止时间 |
| min_participants | INT | 最低人数 (低于自动取消) |
| max_participants | INT | 最大人数 |
| current_participants | INT DEFAULT 0 | 当前报名人数 |
| price | DECIMAL(10,2) | 报名费 |
| level_requirement | VARCHAR(20) | 水平要求 (beginner/intermediate/advanced/all) |
| status | VARCHAR(20) | 状态 (draft/published/registration/confirmed/ongoing/completed/cancelled) |
| cancel_check_at | TIMESTAMPTZ | 自动取消检查时间（报名截止后 N 小时） |
| created_at | TIMESTAMPTZ | 创建时间 |
| updated_at | TIMESTAMPTZ | 更新时间 |

**activity_registrations** - 活动报名表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| activity_id | BIGINT FK | 活动 ID |
| user_id | BIGINT FK | 用户 ID |
| order_id | BIGINT FK | 订单 ID |
| status | VARCHAR(20) | 状态 (registered/cancelled/refunded) |
| created_at | TIMESTAMPTZ | 报名时间 |

#### 教务域 (academic)

**coaches** - 教练表（扩展 users 表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| user_id | BIGINT UNIQUE FK | 关联用户 ID |
| title | VARCHAR(64) | 头衔/职称 |
| specialties | JSONB | 专长 ["正手","发球","比赛策略"] |
| bio | TEXT | 简介 |
| hourly_rate | DECIMAL(10,2) | 课时费 |
| rating | DECIMAL(3,2) DEFAULT 5.0 | 评分 |
| student_count | INT DEFAULT 0 | 学员数 |
| status | SMALLINT DEFAULT 1 | 状态 |
| created_at | TIMESTAMPTZ | 创建时间 |

**course_schedules** - 课程排期表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| coach_id | BIGINT FK | 教练 ID |
| venue_id | BIGINT FK | 场地 ID |
| student_id | BIGINT FK | 学员 ID (私教) |
| product_id | BIGINT FK | 关联课程商品 ID |
| date | DATE | 上课日期 |
| start_time | TIME | 开始时间 |
| end_time | TIME | 结束时间 |
| type | VARCHAR(20) | 类型 (private/group) |
| status | VARCHAR(20) | 状态 (scheduled/ongoing/completed/cancelled/substituted) |
| recurrence_rule | VARCHAR(128) | 重复规则 (RRULE 格式，如 "FREQ=WEEKLY;BYDAY=MO,WE") |
| recurrence_group_id | UUID | 重复组 ID (同一批重复课程) |
| substitute_coach_id | BIGINT | 代课教练 ID |
| created_at | TIMESTAMPTZ | 创建时间 |

**coach_leaves** - 教练休假表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| coach_id | BIGINT FK | 教练 ID |
| start_date | DATE | 开始日期 |
| end_date | DATE | 结束日期 |
| reason | VARCHAR(256) | 原因 |
| status | VARCHAR(20) | 状态 (pending/approved/rejected) |
| created_at | TIMESTAMPTZ | 创建时间 |

**student_course_records** - 学员课程记录表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| schedule_id | BIGINT FK | 排期 ID |
| student_id | BIGINT FK | 学员 ID |
| coach_id | BIGINT FK | 教练 ID |
| attendance | VARCHAR(20) | 出勤 (present/absent/late/excused) |
| coach_feedback | TEXT | 教练反馈 |
| rating | SMALLINT | 学员评分 (1-5) |
| rating_comment | TEXT | 学员评价 |
| created_at | TIMESTAMPTZ | 创建时间 |

#### 通知域 (notify)

**notifications** - 通知表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL PK | 主键 |
| user_id | BIGINT FK | 用户 ID |
| type | VARCHAR(40) | 类型 (booking_confirmed/payment_success/activity_cancelled/...) |
| title | VARCHAR(128) | 标题 |
| content | TEXT | 内容 |
| is_read | BOOLEAN DEFAULT FALSE | 是否已读 |
| extra_data | JSONB | 扩展数据 (跳转参数等) |
| created_at | TIMESTAMPTZ | 创建时间 |

### 3.2 索引策略

```sql
-- 用户表
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_wechat_openid ON users(wechat_openid);
CREATE INDEX idx_users_member_level ON users(member_level);

-- 预订表 (最高频查询)
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_venue_date ON bookings(venue_id, date);
CREATE INDEX idx_bookings_venue_date_time ON bookings(venue_id, date, start_time, end_time);
CREATE INDEX idx_bookings_status ON bookings(status) WHERE status IN ('pending', 'confirmed', 'waitlisted');

-- 候补表
CREATE INDEX idx_waitlist_venue_date_time ON booking_waitlist(venue_id, date, start_time, end_time);
CREATE INDEX idx_waitlist_status_expires ON booking_waitlist(status, expires_at) WHERE status = 'notified';

-- 订单表
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_type_status ON orders(type, status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_expires_at ON orders(expires_at) WHERE status = 'pending';

-- 课程排期 (冲突检测核心)
CREATE INDEX idx_schedules_coach_date ON course_schedules(coach_id, date);
CREATE INDEX idx_schedules_venue_date ON course_schedules(venue_id, date);
CREATE INDEX idx_schedules_coach_date_time ON course_schedules(coach_id, date, start_time, end_time);
CREATE INDEX idx_schedules_venue_date_time ON course_schedules(venue_id, date, start_time, end_time);
CREATE INDEX idx_schedules_recurrence ON course_schedules(recurrence_group_id);

-- 活动表
CREATE INDEX idx_activities_status ON activities(status);
CREATE INDEX idx_activities_type_status ON activities(type, status);
CREATE INDEX idx_activities_start_at ON activities(start_at);
CREATE INDEX idx_activity_reg_activity_user ON activity_registrations(activity_id, user_id);

-- 商品表
CREATE INDEX idx_products_category_status ON products(category, status);
CREATE INDEX idx_products_sort ON products(sort_order, id);

-- 通知表
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
```

### 3.3 Redis 缓存策略

| Key 模式 | 数据类型 | TTL | 用途 |
|---------|---------|-----|------|
| `venue:availability:{venue_id}:{date}` | Hash (time_slot -> status) | 24h | 场地某天各时段可用状态 |
| `venue:lock:{venue_id}:{date}:{start}:{end}` | String (user_id) | 10min | 预订时段分布式锁（防超卖） |
| `booking:waitlist:{venue_id}:{date}:{start}:{end}` | Sorted Set (user_id -> timestamp) | 7d | 候补队列排序 |
| `user:session:{token}` | Hash | 7d | 用户会话信息 |
| `user:info:{user_id}` | JSON String | 30min | 用户基础信息缓存 |
| `product:detail:{product_id}` | JSON String | 10min | 商品详情缓存 |
| `product:stock:{product_id}:{sku_id}` | String (count) | - | 库存计数器（原子操作） |
| `activity:detail:{activity_id}` | JSON String | 5min | 活动详情缓存 |
| `activity:participants:{activity_id}` | String (count) | - | 活动报名人数计数器 |
| `stats:dashboard` | JSON String | 5min | 管理端仪表盘数据 |
| `stats:venue_usage:{date}` | Hash | 1h | 场地使用率数据 |
| `rate_limit:{user_id}:{api}` | String (count) | 1min | API 限流计数 |
| `sms:code:{phone}` | String (code) | 5min | 短信验证码 |
| `order:expire:check` | Sorted Set | - | 待过期订单检查队列 |

**缓存更新策略**:
- **场地可用性**: Write-through。预订成功/取消时同步更新缓存和数据库。
- **商品库存**: Redis 预扣减 + 数据库最终一致。下单时 Redis DECRBY 原子扣减，支付成功后写库，支付超时回滚。
- **活动人数**: Redis INCR 原子递增，定期与数据库同步。
- **其他数据**: Cache-aside 模式，读时加载，写时删除。

---

## 4. API 设计规范

### 4.1 URL 规范

**基础路径**: `/api/v1`

**命名规则**:
- 使用小写字母和连字符（kebab-case）
- 资源名使用复数
- 嵌套资源最多两级

```
# 用户端 API
POST   /api/v1/auth/wechat-login          # 微信登录
POST   /api/v1/auth/phone-login            # 手机号登录
POST   /api/v1/auth/sms-code               # 发送验证码
GET    /api/v1/auth/profile                 # 获取个人信息
PUT    /api/v1/auth/profile                 # 更新个人信息

GET    /api/v1/venues                       # 场地列表
GET    /api/v1/venues/:id                   # 场地详情
GET    /api/v1/venues/:id/availability      # 场地可用时段 (?date=2026-03-19)

POST   /api/v1/bookings                     # 创建预订
GET    /api/v1/bookings                     # 我的预订列表
GET    /api/v1/bookings/:id                 # 预订详情
POST   /api/v1/bookings/:id/cancel          # 取消预订
POST   /api/v1/bookings/waitlist            # 加入候补

POST   /api/v1/orders                       # 创建订单
GET    /api/v1/orders                       # 我的订单列表
GET    /api/v1/orders/:id                   # 订单详情
POST   /api/v1/orders/:id/cancel            # 取消订单

POST   /api/v1/payments/prepare             # 预支付 (获取微信支付参数)
POST   /api/v1/payments/wechat-callback     # 微信支付回调
GET    /api/v1/payments/:id                 # 支付状态查询

GET    /api/v1/products                     # 商品列表 (?category=course&sub_category=course_pack)
GET    /api/v1/products/:id                 # 商品详情

GET    /api/v1/activities                   # 活动列表 (?type=tournament&status=registration)
GET    /api/v1/activities/:id               # 活动详情
POST   /api/v1/activities/:id/register      # 活动报名

GET    /api/v1/wallet                       # 钱包信息
POST   /api/v1/wallet/recharge              # 充值
GET    /api/v1/wallet/transactions           # 流水记录

GET    /api/v1/coupons                      # 我的优惠券
GET    /api/v1/coupons/available             # 可用优惠券 (?order_type=booking&amount=280)

GET    /api/v1/notifications                # 通知列表
PUT    /api/v1/notifications/:id/read       # 标记已读
PUT    /api/v1/notifications/read-all       # 全部已读

# 管理端 API
GET    /api/v1/admin/dashboard              # 仪表盘数据
GET    /api/v1/admin/stats/revenue           # 收入统计
GET    /api/v1/admin/stats/venue-usage       # 场地使用率
GET    /api/v1/admin/stats/users             # 用户画像

CRUD   /api/v1/admin/venues                 # 场地管理
GET    /api/v1/admin/venues/time-grid        # 时段视图 (?date=2026-03-19)
CRUD   /api/v1/admin/venues/:id/time-rules   # 时段规则管理

CRUD   /api/v1/admin/activities              # 活动管理
CRUD   /api/v1/admin/products                # 商品管理
GET    /api/v1/admin/orders                  # 订单列表
POST   /api/v1/admin/refunds/:id/review     # 退款审核

CRUD   /api/v1/admin/coaches                 # 教练管理
GET    /api/v1/admin/coaches/:id/performance  # 教练绩效
CRUD   /api/v1/admin/schedules               # 课程排期
POST   /api/v1/admin/schedules/conflict-check # 冲突检测
CRUD   /api/v1/admin/coach-leaves             # 教练休假
GET    /api/v1/admin/students                 # 学员列表
GET    /api/v1/admin/students/:id/records     # 学员课程记录
```

### 4.2 统一响应格式

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "1号场"
  }
}
```

**分页响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 156,
      "total_pages": 8
    }
  }
}
```

**错误响应**:
```json
{
  "code": 40001,
  "message": "该时段已被预订",
  "data": null
}
```

**错误码规范**:

| 范围 | 含义 |
|------|------|
| 0 | 成功 |
| 40000-40099 | 通用参数错误 |
| 40100-40199 | 认证/授权错误 |
| 40200-40299 | 预订相关错误 |
| 40300-40399 | 支付相关错误 |
| 40400-40499 | 商品相关错误 |
| 40500-40599 | 活动相关错误 |
| 40600-40699 | 教务相关错误 |
| 50000-50099 | 服务器内部错误 |

### 4.3 认证方案

```
┌──────────┐     POST /auth/wechat-login      ┌──────────┐
│  小程序   │  ────────────────────────────►    │  后端     │
│          │  {code, encrypted_data, iv}        │          │
│          │                                    │ 1. code → 微信服务器换 openid
│          │     {access_token, refresh_token}  │ 2. 查找/创建用户
│          │  ◄────────────────────────────     │ 3. 签发 JWT
└──────────┘                                    └──────────┘
```

- **Access Token**: JWT，有效期 2 小时，放 Header `Authorization: Bearer <token>`
- **Refresh Token**: 不透明令牌，有效期 30 天，存 Redis
- **Token 刷新**: `POST /api/v1/auth/refresh`，用 refresh_token 换新 access_token
- **JWT Payload**: `{user_id, role, member_level, exp, iat}`

### 4.4 分页与筛选规范

**分页参数**:
```
GET /api/v1/orders?page=1&page_size=20
```
- `page`: 页码，从 1 开始，默认 1
- `page_size`: 每页条数，默认 20，最大 100

**筛选参数**:
```
GET /api/v1/orders?type=booking&status=paid&start_date=2026-03-01&end_date=2026-03-19
```

**排序参数**:
```
GET /api/v1/products?sort_by=price&sort_order=asc
```
- `sort_by`: 排序字段
- `sort_order`: asc / desc，默认 desc

**搜索参数**:
```
GET /api/v1/admin/orders?keyword=张伟
```

---

## 5. 项目目录结构

### 5.1 仓库策略: Monorepo

采用 **Monorepo** 管理三端代码，理由：
- 团队规模小（预计 3-5 人），统一管理更高效
- 共享类型定义、API 接口约定、lint 规则
- 统一 CI/CD 流程
- 原子性变更（后端接口和前端调用同步修改）

```
T-Line/
├── README.md
├── Makefile                        # 顶层构建命令
├── docker-compose.yml              # 本地开发环境
├── .github/
│   └── workflows/                  # CI/CD
│       ├── backend.yml
│       ├── miniapp.yml
│       └── admin.yml
│
├── docs/                           # 项目文档
│   ├── architecture.md             # 本文档
│   ├── api-spec.md                 # API 详细规范
│   ├── database-erd.md             # ER 图
│   ├── visual-design-user.html     # 用户端视觉稿
│   └── visual-design-admin.html    # 管理端视觉稿
│
├── backend/                        # Go 后端
│   └── (见 5.2)
│
├── miniapp/                        # 用户端小程序
│   └── (见 5.3)
│
└── admin/                          # 管理端 Web
    └── (见 5.4)
```

### 5.2 后端 Go 项目目录

```
backend/
├── cmd/
│   └── server/
│       └── main.go                     # 入口
├── configs/
│   ├── config.yaml                     # 配置模板
│   └── config.example.yaml             # 示例配置
├── internal/
│   ├── config/                         # 配置加载
│   │   └── config.go
│   ├── server/                         # HTTP 服务启动
│   │   ├── server.go
│   │   └── router.go                   # 总路由注册
│   │
│   ├── middleware/                      # 中间件
│   │   ├── auth.go                     # JWT 认证
│   │   ├── rbac.go                     # 角色鉴权
│   │   ├── logger.go                   # 请求日志
│   │   ├── ratelimit.go                # 限流
│   │   ├── cors.go                     # CORS
│   │   └── recovery.go                 # Panic 恢复
│   │
│   ├── pkg/                            # 共享包
│   │   ├── database/                   # 数据库初始化
│   │   │   ├── postgres.go
│   │   │   └── redis.go
│   │   ├── response/                   # 统一响应
│   │   │   └── response.go
│   │   ├── errors/                     # 错误码
│   │   │   └── codes.go
│   │   ├── logger/                     # 日志
│   │   │   └── logger.go
│   │   ├── validator/                  # 参数校验
│   │   │   └── validator.go
│   │   ├── jwt/                        # JWT 工具
│   │   │   └── jwt.go
│   │   ├── pagination/                 # 分页工具
│   │   │   └── pagination.go
│   │   └── upload/                     # 文件上传
│   │       └── oss.go
│   │
│   ├── integration/                    # 外部服务集成
│   │   ├── wechat/                     # 微信 SDK
│   │   │   ├── auth.go                 # 登录
│   │   │   ├── pay.go                  # 支付
│   │   │   └── message.go             # 模板消息
│   │   └── sms/                        # 短信
│   │       └── sms.go
│   │
│   ├── scheduler/                      # 定时任务
│   │   ├── scheduler.go                # 调度器
│   │   ├── booking_expire.go           # 候补超时检查
│   │   ├── order_expire.go             # 订单超时关闭
│   │   ├── activity_cancel.go          # 活动自动取消
│   │   └── stats_aggregate.go          # 统计聚合
│   │
│   └── modules/                        # 业务模块
│       ├── auth/
│       │   ├── handler.go              # HTTP Handler (控制器)
│       │   ├── service.go              # 业务逻辑
│       │   ├── repository.go           # 数据访问
│       │   ├── model.go                # 数据模型 (GORM)
│       │   ├── dto.go                  # 请求/响应 DTO
│       │   └── routes.go              # 模块路由注册
│       │
│       ├── venue/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   ├── dto.go
│       │   └── routes.go
│       │
│       ├── booking/
│       │   ├── handler.go
│       │   ├── handler_admin.go        # 管理端 Handler
│       │   ├── service.go
│       │   ├── service_waitlist.go     # 候补逻辑
│       │   ├── repository.go
│       │   ├── model.go
│       │   ├── dto.go
│       │   └── routes.go
│       │
│       ├── order/
│       │   ├── handler.go
│       │   ├── handler_admin.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   ├── dto.go
│       │   ├── state_machine.go        # 订单状态机
│       │   └── routes.go
│       │
│       ├── payment/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── service_combo.go        # 组合支付逻辑
│       │   ├── repository.go
│       │   ├── model.go
│       │   ├── dto.go
│       │   └── routes.go
│       │
│       ├── product/
│       │   ├── handler.go
│       │   ├── handler_admin.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   ├── dto.go
│       │   └── routes.go
│       │
│       ├── activity/
│       │   ├── handler.go
│       │   ├── handler_admin.go
│       │   ├── service.go
│       │   ├── service_auto_cancel.go  # 自动取消逻辑
│       │   ├── repository.go
│       │   ├── model.go
│       │   ├── dto.go
│       │   └── routes.go
│       │
│       ├── academic/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── service_conflict.go     # 冲突检测逻辑
│       │   ├── repository.go
│       │   ├── model.go
│       │   ├── dto.go
│       │   └── routes.go
│       │
│       ├── stats/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   └── routes.go
│       │
│       └── notify/
│           ├── handler.go
│           ├── service.go
│           ├── repository.go
│           ├── model.go
│           └── routes.go
│
├── migrations/                         # 数据库迁移
│   ├── 000001_init_users.up.sql
│   ├── 000001_init_users.down.sql
│   ├── 000002_init_venues.up.sql
│   └── ...
│
├── scripts/                            # 运维脚本
│   ├── seed.go                         # 种子数据
│   └── migrate.sh                      # 迁移脚本
│
├── Dockerfile
├── go.mod
├── go.sum
└── Makefile
```

### 5.3 用户端 uni-app 项目目录

```
miniapp/
├── src/
│   ├── App.vue
│   ├── main.ts
│   ├── manifest.json                  # uni-app 配置
│   ├── pages.json                     # 页面路由配置
│   ├── uni.scss                       # 全局样式变量
│   │
│   ├── pages/                         # 页面（对应 Tab 页 + 子页面）
│   │   ├── home/
│   │   │   └── index.vue              # 首页
│   │   ├── booking/
│   │   │   ├── index.vue              # 场地预订 (日历+时间轴)
│   │   │   └── confirm.vue            # 预订确认
│   │   ├── activity/
│   │   │   ├── index.vue              # 活动列表
│   │   │   └── detail.vue             # 活动详情
│   │   ├── shop/
│   │   │   ├── index.vue              # 商城首页
│   │   │   └── detail.vue             # 商品详情
│   │   ├── profile/
│   │   │   └── index.vue              # 个人中心
│   │   ├── my-bookings/
│   │   │   └── index.vue              # 我的预约
│   │   ├── my-courses/
│   │   │   └── index.vue              # 我的课程
│   │   ├── wallet/
│   │   │   └── index.vue              # 我的钱包
│   │   ├── orders/
│   │   │   ├── index.vue              # 订单列表
│   │   │   └── detail.vue             # 订单详情
│   │   ├── member/
│   │   │   └── index.vue              # 会员中心
│   │   └── settings/
│   │       └── index.vue              # 设置
│   │
│   ├── components/                    # 通用组件
│   │   ├── t-navbar/                  # 自定义导航栏
│   │   ├── t-tabbar/                  # 底部 TabBar
│   │   ├── t-card/                    # 卡片
│   │   ├── t-badge/                   # 徽章
│   │   ├── t-empty/                   # 空状态
│   │   ├── t-price/                   # 价格展示
│   │   └── t-loading/                 # 加载状态
│   │
│   ├── composables/                   # 组合式函数
│   │   ├── useAuth.ts                 # 认证逻辑
│   │   ├── useBooking.ts              # 预订逻辑
│   │   ├── usePayment.ts              # 支付逻辑
│   │   └── usePagination.ts           # 分页加载
│   │
│   ├── stores/                        # Pinia 状态管理
│   │   ├── user.ts                    # 用户状态
│   │   ├── cart.ts                    # 购物状态
│   │   └── app.ts                     # 应用全局状态
│   │
│   ├── api/                           # API 请求层
│   │   ├── request.ts                 # 请求封装 (拦截器/token 刷新)
│   │   ├── auth.ts
│   │   ├── venue.ts
│   │   ├── booking.ts
│   │   ├── order.ts
│   │   ├── payment.ts
│   │   ├── product.ts
│   │   ├── activity.ts
│   │   └── wallet.ts
│   │
│   ├── utils/                         # 工具函数
│   │   ├── date.ts                    # 日期处理
│   │   ├── price.ts                   # 价格格式化
│   │   └── storage.ts                 # 本地存储
│   │
│   ├── types/                         # TypeScript 类型
│   │   ├── api.d.ts                   # API 响应类型
│   │   ├── user.d.ts
│   │   ├── venue.d.ts
│   │   ├── booking.d.ts
│   │   ├── order.d.ts
│   │   └── product.d.ts
│   │
│   └── static/                        # 静态资源
│       ├── images/
│       └── icons/
│
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .env.example
```

### 5.4 管理端 Web 项目目录

```
admin/
├── src/
│   ├── App.vue
│   ├── main.ts
│   ├── router/                        # 路由
│   │   ├── index.ts
│   │   ├── routes/
│   │   │   ├── dashboard.ts
│   │   │   ├── venue.ts
│   │   │   ├── order.ts
│   │   │   ├── activity.ts
│   │   │   ├── product.ts
│   │   │   ├── academic.ts
│   │   │   ├── stats.ts
│   │   │   └── settings.ts
│   │   └── guards.ts                  # 路由守卫
│   │
│   ├── layouts/                       # 布局
│   │   ├── default.vue                # 侧边栏+顶栏+内容区
│   │   └── login.vue                  # 登录页布局
│   │
│   ├── views/                         # 页面视图
│   │   ├── dashboard/
│   │   │   └── index.vue              # 总览仪表盘
│   │   ├── venue/
│   │   │   ├── list.vue               # 场地列表
│   │   │   ├── time-grid.vue          # 时段视图
│   │   │   ├── form.vue               # 场地编辑
│   │   │   └── time-rules.vue         # 时段规则
│   │   ├── order/
│   │   │   ├── list.vue               # 订单列表
│   │   │   └── detail.vue             # 订单详情
│   │   ├── activity/
│   │   │   ├── list.vue               # 活动列表
│   │   │   └── form.vue               # 活动编辑
│   │   ├── product/
│   │   │   ├── list.vue               # 商品列表
│   │   │   └── form.vue               # 商品编辑
│   │   ├── academic/
│   │   │   ├── schedules.vue          # 课程排期
│   │   │   ├── coaches.vue            # 教练管理
│   │   │   ├── coach-detail.vue       # 教练详情/绩效
│   │   │   └── students.vue           # 学员管理
│   │   ├── stats/
│   │   │   └── index.vue              # 数据统计
│   │   ├── settings/
│   │   │   └── index.vue              # 系统设置
│   │   └── login/
│   │       └── index.vue              # 登录页
│   │
│   ├── components/                    # 通用组件
│   │   ├── SearchForm/                # 搜索筛选表单
│   │   ├── DataTable/                 # 数据表格封装
│   │   ├── StatusTag/                 # 状态标签
│   │   ├── PriceDisplay/             # 金额展示
│   │   └── ChartCard/                # 图表卡片
│   │
│   ├── composables/                   # 组合式函数
│   │   ├── useTable.ts                # 表格数据+分页
│   │   ├── useCrud.ts                 # CRUD 操作
│   │   └── usePermission.ts           # 权限检查
│   │
│   ├── stores/                        # Pinia
│   │   ├── user.ts
│   │   └── app.ts
│   │
│   ├── api/                           # API 层
│   │   ├── request.ts
│   │   ├── dashboard.ts
│   │   ├── venue.ts
│   │   ├── order.ts
│   │   ├── activity.ts
│   │   ├── product.ts
│   │   ├── academic.ts
│   │   └── stats.ts
│   │
│   ├── types/                         # TypeScript 类型
│   │   └── ...
│   │
│   └── utils/
│       └── ...
│
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .env.example
```

---

## 6. 关键业务流程

### 6.1 场地预订流程（含候补）

```
用户选择场地+日期+时段
        │
        ▼
┌──────────────────┐
│ 查询 Redis 时段状态 │
└────────┬─────────┘
         │
    ┌────▼────┐
    │ 可预订？  │
    └─┬────┬──┘
      │    │
    Yes    No
      │    │
      ▼    ▼
┌──────┐ ┌──────────────┐
│ 加锁  │ │ 是否开启候补？ │
│Redis │ │              │
│SETNX │ └──┬────┬─────┘
│10min │   Yes    No
└──┬───┘    │     │
   │        ▼     ▼
   │  ┌──────┐  返回"已满"
   │  │加入候补│
   │  │队列   │
   │  │(Redis │
   │  │Sorted │
   │  │Set)   │
   │  └──┬───┘
   │     │
   │     ▼
   │  创建候补预订
   │  (status=waitlisted)
   │  预扣款（冻结余额或创建待支付订单）
   │
   ▼
创建订单 (status=pending)
计算价格（时段规则+会员折扣-优惠券）
        │
        ▼
调用支付模块 (→ 见 6.2)
        │
   ┌────▼────┐
   │ 支付成功？│
   └─┬────┬──┘
    Yes    No/超时
     │     │
     ▼     ▼
确认预订    释放锁
更新 Redis  回滚订单
时段状态    (status=cancelled)
删除锁
通知用户

--- 候补流程 ---

当有人取消预订时:
        │
        ▼
检查候补队列 (Redis ZPOPMIN)
        │
   ┌────▼─────┐
   │ 有候补？   │
   └─┬────┬───┘
    Yes    No
     │      │
     ▼      ▼
通知候补用户  恢复时段可预订
设置确认超时
(15分钟)
     │
┌────▼─────┐
│ 用户确认？ │
└─┬────┬───┘
 Yes    超时
  │      │
  ▼      ▼
完成支付  取消候补
确认预订  退款/解冻
         检查下一位候补
```

**关键设计点**:
1. Redis SETNX 实现分布式锁，防止同一时段被重复预订
2. 候补队列用 Redis Sorted Set，score 为加入时间戳，保证先到先得
3. 候补超时由 scheduler 定时扫描 `booking_waitlist` 表中 `expires_at < now()` 的记录
4. 取消退款规则: 24小时前取消全额退款，12-24小时退50%，12小时内不退

### 6.2 组合支付流程（余额+微信）

```
用户确认支付
    │
    ▼
┌──────────────────┐
│ 计算应付金额       │
│ pay_amount =      │
│ total - discount  │
│ - coupon          │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ 查询用户余额       │
│ wallet.balance    │
└────────┬─────────┘
         │
    ┌────▼─────────────┐
    │ 余额 >= 应付金额？ │
    └─┬─────────┬──────┘
     Yes         No
      │           │
      ▼           ▼
 纯余额支付    组合支付
      │           │
      ▼           ▼
┌──────────┐ ┌────────────────────────┐
│ BEGIN TX  │ │ BEGIN TX                │
│           │ │                        │
│ 扣减余额   │ │ 1. 冻结全部余额          │
│ (乐观锁    │ │    wallet.frozen +=     │
│  version)  │ │    wallet.balance       │
│           │ │    wallet.balance = 0    │
│ 记录流水   │ │                        │
│           │ │ 2. 计算微信支付金额       │
│ 更新订单   │ │    wechat_amount =      │
│ status=   │ │    pay - balance        │
│ paid      │ │                        │
│           │ │ 3. 记录流水              │
│ COMMIT    │ │                        │
└──────────┘ │ 4. 创建 payment 记录     │
             │    status=pending        │
             │                        │
             │ COMMIT                  │
             └────────┬───────────────┘
                      │
                      ▼
              调用微信支付统一下单 API
              获取 prepay_id
                      │
                      ▼
              返回支付参数给小程序
              小程序调起微信支付
                      │
                 ┌────▼──────┐
                 │ 微信回调    │
                 └─┬─────┬───┘
                 成功    失败
                  │      │
                  ▼      ▼
           ┌──────────┐ ┌────────────┐
           │ BEGIN TX  │ │ 解冻余额    │
           │           │ │ frozen →   │
           │ 确认冻结   │ │ balance    │
           │ 扣除      │ │            │
           │ frozen -= │ │ payment    │
           │ balance   │ │ status=    │
           │ 部分      │ │ failed     │
           │           │ │            │
           │ 更新 order│ │ 更新 order  │
           │ 更新      │ │ status=    │
           │ payment   │ │ cancelled  │
           │ status=   │ │            │
           │ success   │ │ COMMIT     │
           │           │ └────────────┘
           │ COMMIT    │
           └──────────┘
```

**关键设计点**:
1. 钱包余额使用乐观锁 (`version` 字段) 防止并发扣减
2. 组合支付时先冻结余额，微信支付成功后再正式扣除，失败则解冻
3. 微信支付回调需做幂等处理（根据 payment_no 去重）
4. 回调超时时通过主动查询微信订单状态补偿

### 6.3 活动自动取消流程

```
管理员发布活动
(设置: min_participants, registration_deadline, cancel_check_at)
        │
        ▼
活动进入报名状态
(status = registration)
        │
        ├─── 用户报名 ──→ current_participants++
        │                  Redis INCR activity:participants:{id}
        │
        ▼
到达 cancel_check_at 时间
(Scheduler 定时扫描)
        │
        ▼
┌──────────────────────────┐
│ current_participants     │
│ >= min_participants ?    │
└─┬──────────────┬────────┘
  Yes             No
   │               │
   ▼               ▼
活动确认成立      ┌────────────────────┐
status =         │ 自动取消活动         │
confirmed        │ status = cancelled  │
                 │                    │
通知所有          │ 遍历所有报名记录     │
参与者            │ registration.status │
活动确认          │ = cancelled         │
                 │                    │
                 │ 批量创建退款记录     │
                 │ 异步处理退款:        │
                 │ - 余额部分退回余额   │
                 │ - 微信部分调微信退款  │
                 │                    │
                 │ 通知所有报名用户     │
                 │ "活动因人数不足取消"  │
                 │                    │
                 │ 释放场地占用         │
                 └────────────────────┘
```

**关键设计点**:
1. Scheduler 每分钟扫描 `activities WHERE status='registration' AND cancel_check_at <= now()`
2. 批量退款异步处理，避免单次操作耗时过长
3. 活动取消需同步释放 venue_blocked_times 或 booking 占用
4. 通知使用微信模板消息/订阅消息

### 6.4 课程排期冲突检测

```
管理员创建课程排期
(coach_id, venue_id, date, start_time, end_time, recurrence_rule)
        │
        ▼
┌──────────────────────────────────┐
│ 若有重复规则 (recurrence_rule)     │
│ 展开未来 N 周的具体日期列表         │
│ (解析 RRULE → date_list)          │
└──────────────────┬───────────────┘
                   │
        ┌──────────▼──────────┐
        │ 对每个日期执行冲突检测 │
        └──────────┬──────────┘
                   │
    ┌──────────────┼──────────────┐
    │              │              │
    ▼              ▼              ▼
教练冲突检测    场地冲突检测     教练休假检测
    │              │              │
    ▼              ▼              ▼
SELECT EXISTS   SELECT EXISTS   SELECT EXISTS
FROM            FROM            FROM
course_schedules course_schedules coach_leaves
WHERE           WHERE           WHERE
coach_id = ?    venue_id = ?    coach_id = ?
AND date = ?    AND date = ?    AND start_date <= ?
AND status      AND status      AND end_date >= ?
NOT IN          NOT IN          AND status =
('cancelled')   ('cancelled')   'approved'
AND (           AND (
  start < ?end    start < ?end
  AND             AND
  end > ?start    end > ?start
)               )
    │              │              │
    ▼              ▼              ▼
┌──────────────────────────────────┐
│ 汇总所有冲突                       │
│                                  │
│ conflicts = [{                   │
│   date: "2026-03-20",            │
│   type: "coach",                 │
│   detail: "王教练当日14:00-15:00  │
│            已有私教课"             │
│ }, {                             │
│   date: "2026-03-22",            │
│   type: "venue",                 │
│   detail: "3号场当日14:00-16:00   │
│            已有团课"              │
│ }, {                             │
│   date: "2026-03-25",            │
│   type: "leave",                 │
│   detail: "王教练 3/24-3/26 请假" │
│ }]                               │
└──────────────────┬───────────────┘
                   │
          ┌────────▼─────────┐
          │ 存在冲突？         │
          └─┬──────────┬─────┘
           Yes          No
            │            │
            ▼            ▼
    返回冲突详情     批量创建排期记录
    管理员手动       (同一 recurrence_group_id)
    调整后重试       同步更新场地占用
                    通知教练和学员
```

**关键设计点**:
1. 冲突检测使用数据库查询而非内存计算，确保数据一致性
2. 时间段重叠判断: `start1 < end2 AND end1 > start2`
3. 重复课程使用 `recurrence_group_id` 关联，方便批量管理（取消整组、调整）
4. 冲突检测 API 支持 "试算" 模式（不实际创建，仅返回冲突），和 "确认" 模式（通过冲突检测后创建）
5. 索引 `(coach_id, date, start_time, end_time)` 和 `(venue_id, date, start_time, end_time)` 确保查询性能

---

## 7. 架构决策记录

### ADR-001: 采用模块化单体而非微服务

**背景**: 需要决定后端服务架构风格。

**决策**: 采用模块化单体 (Modular Monolith)。

**正面影响**:
- 开发和部署复杂度低，适合初创团队
- 模块间调用无网络开销
- 事务处理简单（单数据库事务）
- 调试方便

**负面影响**:
- 单点故障风险（通过多实例部署缓解）
- 模块间可能出现隐性耦合

**备选方案**:
- **微服务**: 独立部署和扩展，但运维复杂度和基础设施成本过高
- **Serverless**: 冷启动延迟不适合实时预订场景

**迁移路径**: 模块间通过 interface 解耦，未来流量增长后可将高负载模块（如 booking, payment）拆出为独立服务。

### ADR-002: 使用 Redis 管理场地时段锁和候补队列

**背景**: 场地预订需要防止超卖，候补需要有序队列。

**决策**: 使用 Redis SETNX 实现分布式锁，Sorted Set 管理候补队列。

**正面影响**:
- 原子操作，天然防超卖
- Sorted Set 自动排序，先到先得
- 性能高（亚毫秒级）

**负面影响**:
- Redis 宕机时需要降级方案（降级到数据库行锁）
- 需要处理锁续期和过期释放

**备选方案**:
- **PostgreSQL FOR UPDATE**: 可靠但性能瓶颈在高并发场景
- **消息队列**: 过度设计，引入额外组件

### ADR-003: 组合支付采用冻结-确认模式

**背景**: 余额+微信组合支付需要保证原子性。

**决策**: 先冻结余额，微信支付成功后确认扣除，失败则解冻。

**正面影响**:
- 余额不会被误扣
- 中间状态可审计
- 失败自动回滚

**负面影响**:
- 冻结期间余额不可用（通常几十秒，可接受）
- 需要处理异常情况（微信回调延迟、网络超时）

**备选方案**:
- **先扣余额再调微信**: 微信支付失败需退回余额，退回过程中有时间差
- **全部走微信**: 放弃余额功能，用户体验差

### ADR-004: Monorepo 管理三端代码

**背景**: 需要决定代码仓库组织方式。

**决策**: 采用 Monorepo，一个仓库包含 backend / miniapp / admin 三个子项目。

**正面影响**:
- API 类型定义共享，前后端同步
- 统一 CI/CD 流程
- 原子性提交（接口变更+前端适配同一个 PR）
- 代码复用方便

**负面影响**:
- 仓库体积增长较快
- CI 需要智能判断哪个子项目有变更

**备选方案**:
- **多仓库**: 完全独立，但协同成本高，类型定义容易不一致

---

## 附录 A: 性能指标与优化策略

| 场景 | 目标 | 优化手段 |
|------|------|---------|
| 首页加载 | < 1.5s | 接口聚合、Redis 缓存、CDN 静态资源 |
| 场地可用性查询 | < 200ms | Redis Hash 缓存、索引优化 |
| 创建预订 | < 500ms | Redis 锁、异步通知 |
| 支付回调处理 | < 300ms | 幂等设计、最小事务范围 |
| 订单列表查询 | < 500ms | 分页查询、索引覆盖、按需加载 |
| 管理端仪表盘 | < 1s | 预聚合数据（scheduler）、Redis 缓存 |
| 冲突检测 | < 1s | 复合索引、批量查询优化 |

## 附录 B: 安全架构

| 层面 | 措施 |
|------|------|
| 传输层 | 全站 HTTPS，HSTS |
| 认证 | JWT + Refresh Token，微信 OAuth 2.0 |
| 授权 | RBAC (user/coach/admin/super_admin)，中间件统一拦截 |
| 数据安全 | 手机号/身份信息 AES 加密存储，日志脱敏 |
| 支付安全 | 微信支付官方 SDK，回调签名验证，HTTPS 回调 |
| 接口安全 | 请求限流 (IP + User)，参数校验，SQL 注入防护 (GORM 参数化) |
| 审计 | 关键操作日志（支付、退款、权限变更），钱包流水不可删除 |

## 附录 C: 监控与运维

| 维度 | 工具/方案 |
|------|----------|
| 应用日志 | Zap 结构化日志 → ELK / Loki |
| API 监控 | Prometheus + Grafana (请求量/延迟/错误率) |
| 错误追踪 | Sentry |
| 数据库监控 | pg_stat_statements，慢查询日志 |
| 健康检查 | `/api/v1/health` (数据库连接/Redis 连接/磁盘空间) |
| 备份 | PostgreSQL pg_dump 每日全量 + WAL 增量，保留 30 天 |
| 告警 | API 错误率 > 1%、响应 P99 > 2s、数据库连接 > 80% |