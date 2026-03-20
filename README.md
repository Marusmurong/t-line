# T-Line Sports Club

网球馆智能管理系统 — 微信小程序 + Web 管理后台 + Go 后端 API

整合会员管理、场地预约、课程管理、商品销售、活动报名等核心功能，打造"预约-通行-服务-结算"闭环生态。

## 系统架构

```
┌────────────────────┐     ┌──────────────────────┐
│  微信小程序 (C端)    │     │  Web 管理后台 (B端)    │
│  uni-app + Vue3     │     │  Vue3 + Arco Design   │
└─────────┬──────────┘     └───────────┬──────────┘
          │  HTTPS                     │  HTTPS
          └────────────┬───────────────┘
                       ▼
              ┌─────────────────┐
              │   Nginx 反向代理  │
              └────────┬────────┘
                       ▼
              ┌─────────────────┐
              │  Go (Gin) API   │
              │  12 个业务模块    │
              └───┬────────┬────┘
                  │        │
          ┌───────▼──┐ ┌───▼──────┐
          │PostgreSQL │ │  Redis   │
          │  16+      │ │  7+      │
          └──────────┘ └──────────┘
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 用户端 | uni-app + Vue 3 + TypeScript + Pinia |
| 管理端 | Vue 3 + Arco Design Pro + TypeScript + Pinia |
| 后端 | Go 1.22+ / Gin / GORM v2 |
| 数据库 | PostgreSQL 16 + Redis 7 |
| 部署 | Docker + Docker Compose + Nginx |

## 项目结构

```
t-line/
├── backend/                 # Go 后端 API
│   ├── cmd/server/          # 入口
│   ├── configs/             # 配置文件
│   ├── internal/
│   │   ├── config/          # 配置加载
│   │   ├── server/          # HTTP 服务 + 路由
│   │   ├── middleware/      # JWT认证/RBAC/限流/CORS/日志
│   │   ├── pkg/             # 共享工具 (response/errors/jwt/pagination)
│   │   ├── integration/     # 外部服务 (微信/短信)
│   │   ├── scheduler/       # 定时任务引擎
│   │   └── modules/         # 业务模块
│   │       ├── auth/        #   用户认证 + 会员 + 钱包
│   │       ├── venue/       #   场地管理 + 时段规则
│   │       ├── booking/     #   预订 + 候补队列
│   │       ├── order/       #   订单 + 状态机
│   │       ├── payment/     #   组合支付 (余额+微信)
│   │       ├── product/     #   商品 (课程/球具/服务)
│   │       ├── activity/    #   活动 + 自动取消
│   │       ├── academic/    #   教务 (排期/教练/学员)
│   │       ├── notify/      #   通知消息
│   │       └── stats/       #   数据统计
│   ├── migrations/          # 数据库迁移 (9个版本, 24+张表)
│   ├── Dockerfile
│   └── Makefile
│
├── miniapp/                 # 微信小程序 (uni-app)
│   └── src/
│       ├── pages/           # 15 个页面
│       ├── api/             # 7 个 API 模块
│       ├── stores/          # Pinia 状态管理
│       ├── composables/     # 组合式函数
│       ├── components/      # 通用组件
│       └── types/           # TypeScript 类型
│
├── admin/                   # Web 管理后台
│   └── src/
│       ├── views/           # 10 个管理页面
│       ├── api/             # 6 个 API 模块
│       ├── layouts/         # 侧边栏布局
│       ├── router/          # 路由 + 守卫
│       └── stores/          # Pinia 状态管理
│
├── docs/                    # 项目文档
│   ├── architecture.md      # 系统架构设计
│   ├── dev-plan.md          # 开发计划
│   ├── visual-design-user.html   # 小程序视觉稿
│   └── visual-design-admin.html  # 管理端视觉稿
│
└── docker-compose.yml       # 本地开发环境
```

## 功能模块

### 用户端 (微信小程序)

| 模块 | 功能 |
|------|------|
| 场地预订 | 日历+时间轴选择、候补订场、取消退款 |
| 商城 | 课程购买(单次/课包)、球具(现货/预定)、服务(租赁/穿线) |
| 活动 | 畅打/团课/赛事/主题活动报名 |
| 支付 | 余额优先+微信补差组合支付、优惠券 |
| 个人中心 | 会员等级(普通/银/金/钻)、钱包、订单、课程记录 |

### 管理端 (Web 后台)

| 模块 | 功能 |
|------|------|
| 总览 | 营收/订单/使用率/活跃用户、收入趋势、最近订单 |
| 场地管理 | 场地 CRUD、时段规则、时间网格视图 |
| 订单管理 | 多条件筛选、退款审核、异常处理 |
| 商品管理 | 课程/球具/服务 CRUD、上下架、库存 |
| 活动管理 | 发布、自动取消(人数不足)、报名管理 |
| 教务管理 | 课程排期(冲突检测)、教练排班/绩效、学员管理 |
| 数据统计 | 收入趋势、场地热力图、用户增长、会员分布 |

### 后端核心能力

| 能力 | 实现方式 |
|------|---------|
| 防超卖 | Redis SETNX 分布式锁 |
| 候补队列 | Redis Sorted Set (先到先得) |
| 组合支付 | 余额冻结 → 微信补差 → 回调确认/解冻 |
| 排课冲突检测 | SQL 时间重叠查询 (教练+场地+休假) |
| 定时任务 | 候补超时/订单过期/活动自动取消/统计聚合 |

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- Docker & Docker Compose
- 微信开发者工具 (小程序开发)

### 1. 克隆项目

```bash
git clone https://github.com/Marusmurong/t-line.git
cd t-line
```

### 2. 启动数据库

```bash
docker compose up -d
```

这会启动：
- PostgreSQL 16 (端口 5432, 用户 tline/tline123)
- Redis 7 (端口 6379)

### 3. 启动后端

```bash
cd backend

# 复制配置
cp configs/config.example.yaml configs/config.yaml
# 编辑 config.yaml，填入微信/短信/OSS 等配置

# 安装依赖
go mod tidy

# 执行数据库迁移
# 方式一：使用 golang-migrate
migrate -path migrations -database "postgres://tline:tline123@localhost:5432/tline?sslmode=disable" up

# 方式二：手动执行 SQL
psql -h localhost -U tline -d tline -f migrations/000001_init_users.up.sql
# ... 依次执行所有 up.sql

# 启动服务 (热重载)
make dev

# 或直接运行
make run
```

后端启动后访问 http://localhost:8080/api/v1/health 验证。

### 4. 启动管理端

```bash
cd admin
npm install
npm run dev
```

访问 http://localhost:3000 进入管理后台。

### 5. 启动小程序

```bash
cd miniapp
npm install
npm run dev:mp-weixin
```

用微信开发者工具打开 `miniapp/dist/dev/mp-weixin` 目录。

## API 概览

### 用户端 API

```
POST   /api/v1/auth/wechat-login     # 微信登录
POST   /api/v1/auth/phone-login      # 手机号登录
GET    /api/v1/auth/profile           # 个人信息

GET    /api/v1/venues                 # 场地列表
GET    /api/v1/venues/:id/availability # 可用时段

POST   /api/v1/bookings              # 创建预订
POST   /api/v1/bookings/waitlist     # 加入候补

POST   /api/v1/orders                # 创建订单
POST   /api/v1/payments/prepare      # 预支付

GET    /api/v1/products              # 商品列表
GET    /api/v1/activities            # 活动列表
POST   /api/v1/activities/:id/register # 活动报名

GET    /api/v1/coaches               # 教练列表
GET    /api/v1/wallet                # 钱包信息
GET    /api/v1/notifications         # 通知列表
```

### 管理端 API

```
GET    /api/v1/admin/dashboard       # 仪表盘
CRUD   /api/v1/admin/venues          # 场地管理
CRUD   /api/v1/admin/products        # 商品管理
CRUD   /api/v1/admin/activities      # 活动管理
GET    /api/v1/admin/orders          # 订单管理
CRUD   /api/v1/admin/coaches         # 教练管理
CRUD   /api/v1/admin/schedules       # 课程排期
POST   /api/v1/admin/schedules/conflict-check # 冲突检测
GET    /api/v1/admin/stats/revenue   # 收入统计
GET    /api/v1/admin/stats/venue-usage # 场地使用率
GET    /api/v1/admin/stats/users     # 用户统计
```

## 配置说明

复制 `backend/configs/config.example.yaml` 为 `config.yaml`，按需修改：

| 配置项 | 说明 | 必填 |
|--------|------|------|
| `database.*` | PostgreSQL 连接信息 | 是 |
| `redis.*` | Redis 连接信息 | 是 |
| `jwt.secret` | JWT 签名密钥 (生产环境务必修改) | 是 |
| `wechat.app_id/app_secret` | 微信小程序凭证 | 小程序登录必填 |
| `wechat.mch_id/mch_api_key` | 微信支付商户凭证 | 微信支付必填 |
| `sms.*` | 短信服务配置 (阿里云/腾讯云) | 手机号登录必填 |
| `oss.*` | 对象存储配置 | 图片上传必填 |

## 数据库迁移

项目包含 9 个迁移版本，覆盖 24+ 张表：

| 版本 | 内容 |
|------|------|
| 000001 | users, wallets, wallet_transactions, points |
| 000002 | venues, venue_time_slot_rules, venue_blocked_times |
| 000003 | orders, order_items, refunds |
| 000004 | payments, coupons, user_coupons |
| 000005 | bookings, booking_waitlist |
| 000006 | products, product_skus |
| 000007 | activities, activity_registrations |
| 000008 | notifications |
| 000009 | coaches, course_schedules, coach_leaves, student_course_records |

## 生产部署

### Docker 部署

```bash
# 构建后端镜像
cd backend
docker build -t tline-api .

# 运行
docker run -d \
  -p 8080:8080 \
  -e DATABASE_HOST=your-db-host \
  -e REDIS_HOST=your-redis-host \
  -e JWT_SECRET=your-production-secret \
  tline-api
```

### Nginx 配置示例

```nginx
server {
    listen 443 ssl;
    server_name api.t-line.com;

    ssl_certificate     /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 微信支付回调 (无需认证)
    location /api/v1/payments/wechat-callback {
        proxy_pass http://127.0.0.1:8080;
    }
}

server {
    listen 443 ssl;
    server_name admin.t-line.com;

    ssl_certificate     /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    root /var/www/tline-admin/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
    }
}
```

### 小程序发布

1. 在 `miniapp/src/manifest.json` 中填入微信小程序 AppID
2. 运行 `npm run build:mp-weixin`
3. 微信开发者工具打开 `dist/build/mp-weixin`
4. 上传代码 → 提交审核 → 发布

## License

MIT
