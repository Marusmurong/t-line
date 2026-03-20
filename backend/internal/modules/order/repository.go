package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GenerateOrderNo generates a unique order number.
func GenerateOrderNo() string {
	return time.Now().Format("20060102150405") + uuid.New().String()[:8]
}

// GenerateRefundNo generates a unique refund number.
func GenerateRefundNo() string {
	return "RF" + time.Now().Format("20060102150405") + uuid.New().String()[:6]
}

// CreateOrder creates an order with items in a transaction.
func (r *Repository) CreateOrder(ctx context.Context, order *Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Order, error) {
	var order Order
	err := r.db.WithContext(ctx).Preload("Items").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) GetByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	var order Order
	err := r.db.WithContext(ctx).Preload("Items").Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID int64, status string, offset, limit int) ([]Order, int64, error) {
	var orders []Order
	var total int64

	query := r.db.WithContext(ctx).Model(&Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Items").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	return orders, total, err
}

func (r *Repository) ListAll(ctx context.Context, status, orderType string, offset, limit int) ([]Order, int64, error) {
	var orders []Order
	var total int64

	query := r.db.WithContext(ctx).Model(&Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if orderType != "" {
		query = query.Where("type = ?", orderType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Items").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	return orders, total, err
}

func (r *Repository) UpdateStatus(ctx context.Context, orderID int64, fromStatus, toStatus string) error {
	result := r.db.WithContext(ctx).
		Model(&Order{}).
		Where("id = ? AND status = ?", orderID, fromStatus).
		Updates(map[string]interface{}{
			"status":     toStatus,
			"updated_at": time.Now(),
		})
	if result.RowsAffected == 0 {
		return fmt.Errorf("订单状态变更失败: 当前状态非 %s", fromStatus)
	}
	return result.Error
}

func (r *Repository) UpdateOrderPaid(ctx context.Context, orderID int64, balancePaid, wechatPaid interface{}) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&Order{}).
		Where("id = ? AND status = ?", orderID, StatusPending).
		Updates(map[string]interface{}{
			"status":       StatusPaid,
			"balance_paid": balancePaid,
			"wechat_paid":  wechatPaid,
			"paid_at":      &now,
			"updated_at":   now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("订单支付状态更新失败: 订单可能已被处理")
	}
	return nil
}

// Refund operations

func (r *Repository) CreateRefund(ctx context.Context, refund *Refund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

func (r *Repository) GetRefundByID(ctx context.Context, id int64) (*Refund, error) {
	var refund Refund
	if err := r.db.WithContext(ctx).First(&refund, id).Error; err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *Repository) GetRefundByOrderID(ctx context.Context, orderID int64) (*Refund, error) {
	var refund Refund
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&refund).Error; err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *Repository) UpdateRefund(ctx context.Context, refund *Refund) error {
	return r.db.WithContext(ctx).Save(refund).Error
}

// Transaction helper

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}
