package order

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	OrderNo        string          `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID         int64           `gorm:"not null" json:"user_id"`
	Type           string          `gorm:"size:20;not null" json:"type"`
	Status         string          `gorm:"size:20;not null;default:pending" json:"status"`
	TotalAmount    decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	DiscountAmount decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"discount_amount"`
	PayAmount      decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"pay_amount"`
	BalancePaid    decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"balance_paid"`
	WechatPaid     decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"wechat_paid"`
	CouponID       *int64          `json:"coupon_id"`
	Remark         string          `gorm:"size:512" json:"remark"`
	PaidAt         *time.Time      `json:"paid_at"`
	ExpiresAt      *time.Time      `json:"expires_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	ID         int64           `gorm:"primaryKey" json:"id"`
	OrderID    int64           `gorm:"not null" json:"order_id"`
	ItemType   string          `gorm:"size:20;not null" json:"item_type"`
	ItemID     int64           `gorm:"not null" json:"item_id"`
	ItemName   string          `gorm:"size:128;not null" json:"item_name"`
	SkuID      *int64          `json:"sku_id"`
	Quantity   int             `gorm:"default:1" json:"quantity"`
	UnitPrice  decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total_price"`
}

func (OrderItem) TableName() string { return "order_items" }

type Refund struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	RefundNo      string          `gorm:"size:32;uniqueIndex;not null" json:"refund_no"`
	OrderID       int64           `gorm:"not null" json:"order_id"`
	UserID        int64           `gorm:"not null" json:"user_id"`
	Amount        decimal.Decimal `gorm:"type:decimal(12,2);not null" json:"amount"`
	BalanceRefund decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"balance_refund"`
	WechatRefund  decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"wechat_refund"`
	Reason        string          `gorm:"size:512" json:"reason"`
	Status        string          `gorm:"size:20;not null;default:pending" json:"status"`
	ReviewedBy    *int64          `json:"reviewed_by"`
	ReviewedAt    *time.Time      `json:"reviewed_at"`
	CompletedAt   *time.Time      `json:"completed_at"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (Refund) TableName() string { return "refunds" }

// Order status constants
const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusUsed      = "used"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
	StatusRefunding = "refunding"
	StatusRefunded  = "refunded"
)

// Order type constants
const (
	TypeBooking  = "booking"
	TypeProduct  = "product"
	TypeActivity = "activity"
	TypeRecharge = "recharge"
)
