package notify

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, n *Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Notification, error) {
	var n Notification
	if err := r.db.WithContext(ctx).First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID int64, offset, limit int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications).Error
	return notifications, total, err
}

func (r *Repository) MarkRead(ctx context.Context, id, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *Repository) MarkAllRead(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

func (r *Repository) CountUnread(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}
