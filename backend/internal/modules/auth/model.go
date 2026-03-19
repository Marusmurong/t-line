package auth

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID              int64            `gorm:"primaryKey" json:"id"`
	Phone           *string          `gorm:"uniqueIndex;size:20" json:"phone"`
	PasswordHash    string           `gorm:"size:256" json:"-"`
	Nickname        string           `gorm:"size:64" json:"nickname"`
	AvatarURL       string           `gorm:"size:512" json:"avatar_url"`
	Gender          int              `json:"gender"`
	Age             int              `json:"age"`
	WeChatOpenID    *string          `gorm:"uniqueIndex;size:128" json:"-"`
	WeChatUnionID   *string          `gorm:"size:128" json:"-"`
	UTRRating       *decimal.Decimal `gorm:"type:decimal(4,2)" json:"utr_rating"`
	UTRImage        string           `gorm:"size:512" json:"utr_image"`
	BallAge         int              `json:"ball_age"`
	SelfLevel       string           `gorm:"size:20" json:"self_level"`
	MemberLevel     int              `json:"member_level"`
	MemberExpiresAt *time.Time       `json:"member_expires_at"`
	Role            string           `gorm:"size:20;default:user" json:"role"`
	Status          int              `gorm:"default:1" json:"status"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func (User) TableName() string { return "users" }

type Wallet struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	UserID         int64           `gorm:"uniqueIndex" json:"user_id"`
	Balance        decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"balance"`
	FrozenAmount   decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"frozen_amount"`
	TotalRecharged decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"total_recharged"`
	Version        int             `gorm:"default:0" json:"-"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func (Wallet) TableName() string { return "wallets" }

type WalletTransaction struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	WalletID       int64           `json:"wallet_id"`
	Type           string          `gorm:"size:20" json:"type"`
	Amount         decimal.Decimal `gorm:"type:decimal(12,2)" json:"amount"`
	BalanceAfter   decimal.Decimal `gorm:"type:decimal(12,2)" json:"balance_after"`
	RelatedOrderID *int64          `json:"related_order_id"`
	Remark         string          `gorm:"size:256" json:"remark"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (WalletTransaction) TableName() string { return "wallet_transactions" }

type Points struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     int64     `gorm:"uniqueIndex" json:"user_id"`
	Balance    int       `gorm:"default:0" json:"balance"`
	TotalEarned int      `gorm:"default:0" json:"total_earned"`
	TotalSpent int       `gorm:"default:0" json:"total_spent"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Points) TableName() string { return "points" }
