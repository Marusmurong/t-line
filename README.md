<div align="center">

<img src="docs/logo.jpg" alt="T-Line Sports Club" width="200" />

# T-Line Sports Club

**网球馆智能管理系统**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://vuejs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat-square&logo=redis&logoColor=white)](https://redis.io)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)](LICENSE)

微信小程序 + Web 管理后台 + Go 后端 API

整合会员管理、场地预约、课程管理、商品销售、活动报名等核心功能
打造 **"预约 - 通行 - 服务 - 结算"** 闭环生态

</div>

---

## 📸 产品预览

### 用户端 (微信小程序)

<div align="center">
<img src="docs/preview-banner.png" alt="小程序预览" width="100%" />
</div>

> 💡 完整交互稿：浏览器打开 `docs/visual-design-user.html`

### 管理后台

<div align="center">
<img src="docs/preview-admin.png" alt="管理后台预览" width="100%" />
</div>

> 💡 完整交互稿：浏览器打开 `docs/visual-design-admin.html`

---

## 🏗 系统架构

```
┌────────────────────┐       ┌──────────────────────┐
│  微信小程序 (C端)    │       │  Web 管理后台 (B端)    │
│  uni-app + Vue3     │       │  Vue3 + Arco Design   │
└─────────┬──────────┘       └───────────┬──────────┘
          │  HTTPS                       │  HTTPS
          └──────────────┬───────────────┘
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

---

## ✨ 功能亮点

<table>
<tr>
<td width="50%">

### 🎾 用户端 (微信小程序)

| 模块 | 功能 |
|:-----|:-----|
| **场地预订** | 日历+时间轴选择、候补订场、取消退款 |
| **商城** | 课程(单次/课包)、球具(现货/预定)、服务(租赁/穿线) |
| **活动** | 畅打/团课/赛事/主题活动报名 |
| **支付** | 余额+微信组合支付、优惠券 |
| **个人中心** | 会员(普通/银/金/钻)、钱包、订单、课程 |

</td>
<td width="50%">

### 🖥 管理后台 (Web)

| 模块 | 功能 |
|:-----|:-----|
| **总览** | 营收/订单/使用率、收入趋势、最近订单 |
| **场地管理** | CRUD、时段规则、时间网格视图 |
| **订单管理** | 多条件筛选、退款审核 |
| **教务管理** | 排期冲突检测、教练排班/绩效、学员 |
| **数据统计** | 收入趋势、热力图、用户增长、会员分布 |

</td>
</tr>
</table>

### 🔧 后端核心能力

| 能力 | 实现方式 |
|:-----|:---------|
| 🔒 防超卖 | Redis `SETNX` 分布式锁 |
| 📋 候补队列 | Redis `Sorted Set` (先到先得) |
| 💰 组合支付 | 余额冻结 → 微信补差 → 回调确认/解冻 |
| 📅 排课冲突检测 | SQL 时间重叠查询 (教练 + 场地 + 休假) |
| ⏰ 定时任务 | 候补超时 / 订单过期 / 活动自动取消 / 统计聚合 |

---

## 📁 项目结构

```
t-line/
├── backend/                 # Go 后端 API
│   ├── cmd/server/          #   入口
│   ├── configs/             #   配置
│   ├── internal/
│   │   ├── middleware/      #   JWT/RBAC/限流/CORS
│   │   ├── scheduler/       #   定时任务引擎
│   │   └── modules/         #   业务模块 (12个)
│   │       ├── auth/        #     用户认证 + 会员 + 钱包
│   │       ├── venue/       #     场地 + 时段规则
│   │       ├── booking/     #     预订 + 候补
│   │       ├── order/       #     订单 + 状态机
│   │       ├── payment/     #     组合支付
│   │       ├── product/     #     商品
│   │       ├── activity/    #     活动 + 自动取消
│   │       ├── academic/    #     教务 (排期/教练/学员)
│   │       ├── notify/      #     通知
│   │       └── stats/       #     数据统计
│   └── migrations/          #   数据库迁移 (9版本, 24+表)
│
├── miniapp/                 # 微信小程序 (uni-app Vue3)
│   └── src/
│       ├── pages/           #   15 个页面
│       ├── api/             #   7 个 API 模块
│       └── stores/          #   Pinia 状态管理
│
├── admin/                   # Web 管理后台 (Vue3 + Arco)
│   └── src/
│       ├── views/           #   10 个管理页面
│       ├── api/             #   6 个 API 模块
│       └── layouts/         #   侧边栏布局
│
├── docs/                    # 项目文档 + 视觉稿
└── docker-compose.yml       # 本地开发环境
```

---

## 🚀 快速开始

### 环境要求

| 工具 | 版本 | 用途 |
|------|------|------|
| Go | 1.22+ | 后端 |
| Node.js | 18+ | 前端 |
| Docker | latest | 数据库 |
| 微信开发者工具 | latest | 小程序调试 |

### 1️⃣ 克隆项目

```bash
git clone https://github.com/Marusmurong/t-line.git
cd t-line
```

### 2️⃣ 启动数据库

```bash
docker compose up -d
```

> PostgreSQL 16 (5432) + Redis 7 (6379)

### 3️⃣ 启动后端

```bash
cd backend

# 配置
cp configs/config.example.yaml configs/config.yaml
# 编辑 config.yaml 填入微信/短信/OSS 配置

# 依赖 + 迁移 + 启动
go mod tidy
migrate -path migrations -database "postgres://tline:tline123@localhost:5432/tline?sslmode=disable" up
make dev
```

验证：http://localhost:8080/api/v1/health

### 4️⃣ 启动管理端

```bash
cd admin
npm install
npm run dev
```

访问：http://localhost:3000

### 5️⃣ 启动小程序

```bash
cd miniapp
npm install
npm run dev:mp-weixin
```

微信开发者工具打开 `miniapp/dist/dev/mp-weixin`

---

## 🔌 API 概览

<details>
<summary><b>用户端 API</b> (点击展开)</summary>

```
# 认证
POST   /api/v1/auth/wechat-login        微信登录
POST   /api/v1/auth/phone-login         手机号登录
GET    /api/v1/auth/profile              个人信息

# 场地预订
GET    /api/v1/venues                    场地列表
GET    /api/v1/venues/:id/availability   可用时段
POST   /api/v1/bookings                 创建预订
POST   /api/v1/bookings/waitlist        加入候补

# 订单支付
POST   /api/v1/orders                   创建订单
POST   /api/v1/payments/prepare         预支付

# 商品活动
GET    /api/v1/products                 商品列表
GET    /api/v1/activities               活动列表
POST   /api/v1/activities/:id/register  活动报名

# 教务
GET    /api/v1/coaches                  教练列表
GET    /api/v1/my-courses/records       我的课程记录

# 其他
GET    /api/v1/wallet                   钱包信息
GET    /api/v1/notifications            通知列表
GET    /api/v1/coupons                  优惠券
```
</details>

<details>
<summary><b>管理端 API</b> (点击展开)</summary>

```
# 总览
GET    /api/v1/admin/dashboard               仪表盘

# 场地
CRUD   /api/v1/admin/venues                  场地管理
GET    /api/v1/admin/venues/time-grid        时段视图

# 业务
CRUD   /api/v1/admin/products                商品管理
CRUD   /api/v1/admin/activities              活动管理
GET    /api/v1/admin/orders                  订单管理
POST   /api/v1/admin/refunds/:id/review     退款审核

# 教务
CRUD   /api/v1/admin/coaches                 教练管理
CRUD   /api/v1/admin/schedules               课程排期
POST   /api/v1/admin/schedules/conflict-check 冲突检测
GET    /api/v1/admin/students                学员管理

# 统计
GET    /api/v1/admin/stats/revenue           收入统计
GET    /api/v1/admin/stats/venue-usage       场地使用率
GET    /api/v1/admin/stats/users             用户统计
```
</details>

---

## ⚙️ 配置说明

复制 `backend/configs/config.example.yaml` 为 `config.yaml`：

| 配置项 | 说明 | 必填 |
|:-------|:-----|:----:|
| `database.*` | PostgreSQL 连接 | ✅ |
| `redis.*` | Redis 连接 | ✅ |
| `jwt.secret` | JWT 签名密钥 | ✅ |
| `wechat.app_id` / `app_secret` | 微信小程序凭证 | 小程序 |
| `wechat.mch_id` / `mch_api_key` | 微信支付商户凭证 | 支付 |
| `sms.*` | 短信服务 (阿里云/腾讯云) | 短信登录 |
| `oss.*` | 对象存储 | 图片上传 |

---

## 🗄 数据库

9 个迁移版本，24+ 张表：

| 版本 | 表 |
|:-----|:---|
| `000001` | users, wallets, wallet_transactions, points |
| `000002` | venues, venue_time_slot_rules, venue_blocked_times |
| `000003` | orders, order_items, refunds |
| `000004` | payments, coupons, user_coupons |
| `000005` | bookings, booking_waitlist |
| `000006` | products, product_skus |
| `000007` | activities, activity_registrations |
| `000008` | notifications |
| `000009` | coaches, course_schedules, coach_leaves, student_course_records |

---

## 🐳 生产部署

### Docker

```bash
cd backend
docker build -t tline-api .

docker run -d -p 8080:8080 \
  -e DATABASE_HOST=your-db-host \
  -e REDIS_HOST=your-redis-host \
  -e JWT_SECRET=your-production-secret \
  tline-api
```

### Nginx

```nginx
server {
    listen 443 ssl;
    server_name api.t-line.com;

    ssl_certificate     /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

server {
    listen 443 ssl;
    server_name admin.t-line.com;

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

1. `miniapp/src/manifest.json` 填入 AppID
2. `npm run build:mp-weixin`
3. 微信开发者工具上传 → 审核 → 发布

---

<div align="center">

**Built with ❤️ for T-Line Sports Club**

[架构文档](docs/architecture.md) · [开发计划](docs/dev-plan.md) · [视觉稿-用户端](docs/visual-design-user.html) · [视觉稿-管理端](docs/visual-design-admin.html)

</div>
