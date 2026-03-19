package payment

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Payment struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	PaymentNo     string          `gorm:"size:64;uniqueIndex;not null" json:"payment_no"`
	OrderID       int64           `gorm:"not null" json:"order_id"`
	UserID        int64           `gorm:"not null" json:"user_id"`
	Method        string          `gorm:"size:20;not null" json:"method"`
	Amount        decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	BalanceAmount decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"balance_amount"`
	WechatAmount  decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"wechat_amount"`
	WechatTradeNo string          `gorm:"size:64" json:"wechat_trade_no"`
	Status        string          `gorm:"size:20;not null;default:pending" json:"status"`
	PaidAt        *time.Time      `json:"paid_at"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (Payment) TableName() string { return "payments" }

// Payment status constants
const (
	PayStatusPending  = "pending"
	PayStatusSuccess  = "success"
	PayStatusFailed   = "failed"
	PayStatusRefunded = "refunded"
)

// Payment method constants
const (
	MethodBalance = "balance"
	MethodWechat  = "wechat"
	MethodCombo   = "combo"
)

type Coupon struct {
	ID              int64           `gorm:"primaryKey" json:"id"`
	Name            string          `gorm:"size:64;not null" json:"name"`
	Type            string          `gorm:"size:20;not null" json:"type"`
	Value           decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"value"`
	MinAmount       decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"min_amount"`
	ApplicableTypes datatypes.JSON  `gorm:"type:jsonb;default:'[]'" json:"applicable_types"`
	TotalCount      int             `json:"total_count"`
	UsedCount       int             `gorm:"default:0" json:"used_count"`
	StartAt         time.Time       `json:"start_at"`
	EndAt           time.Time       `json:"end_at"`
	Status          int             `gorm:"default:1" json:"status"`
}

func (Coupon) TableName() string { return "coupons" }

type UserCoupon struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	UserID      int64     `gorm:"not null" json:"user_id"`
	CouponID    int64     `gorm:"not null" json:"coupon_id"`
	Status      int       `gorm:"default:0" json:"status"`
	UsedOrderID *int64    `json:"used_order_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`

	Coupon *Coupon `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

func (UserCoupon) TableName() string { return "user_coupons" }

// UserCoupon status constants
const (
	CouponUnused  = 0
	CouponUsed    = 1
	CouponExpired = 2
)
