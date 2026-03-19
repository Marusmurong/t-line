package payment

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Payment operations

func (r *Repository) CreatePayment(ctx context.Context, p *Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repository) GetPaymentByNo(ctx context.Context, paymentNo string) (*Payment, error) {
	var p Payment
	if err := r.db.WithContext(ctx).Where("payment_no = ?", paymentNo).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) GetPaymentByOrderID(ctx context.Context, orderID int64) (*Payment, error) {
	var p Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) UpdatePaymentStatus(ctx context.Context, paymentID int64, status string, wechatTradeNo string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == PayStatusSuccess {
		now := time.Now()
		updates["paid_at"] = &now
	}
	if wechatTradeNo != "" {
		updates["wechat_trade_no"] = wechatTradeNo
	}
	return r.db.WithContext(ctx).Model(&Payment{}).Where("id = ?", paymentID).Updates(updates).Error
}

// Coupon operations

func (r *Repository) GetCouponByID(ctx context.Context, id int64) (*Coupon, error) {
	var c Coupon
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListUserCoupons(ctx context.Context, userID int64, status *int) ([]UserCoupon, error) {
	var coupons []UserCoupon
	query := r.db.WithContext(ctx).Preload("Coupon").Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	err := query.Order("created_at DESC").Find(&coupons).Error
	return coupons, err
}

func (r *Repository) ListAvailableCoupons(ctx context.Context, userID int64, orderType string, amount decimal.Decimal) ([]UserCoupon, error) {
	var coupons []UserCoupon
	now := time.Now()

	err := r.db.WithContext(ctx).
		Preload("Coupon").
		Joins("JOIN coupons ON coupons.id = user_coupons.coupon_id").
		Where("user_coupons.user_id = ? AND user_coupons.status = 0 AND user_coupons.expires_at > ?", userID, now).
		Where("coupons.min_amount <= ?", amount).
		Where("coupons.status = 1").
		Find(&coupons).Error

	return coupons, err
}

func (r *Repository) UseUserCoupon(ctx context.Context, userCouponID, orderID int64) error {
	return r.db.WithContext(ctx).
		Model(&UserCoupon{}).
		Where("id = ? AND status = 0", userCouponID).
		Updates(map[string]interface{}{
			"status":        CouponUsed,
			"used_order_id": orderID,
		}).Error
}

// Transaction helper

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}
