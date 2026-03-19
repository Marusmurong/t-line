package activity

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Activity CRUD

func (r *Repository) Create(ctx context.Context, a *Activity) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Activity, error) {
	var a Activity
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) List(ctx context.Context, actType, status *string, offset, limit int) ([]Activity, int64, error) {
	var activities []Activity
	var total int64

	query := r.db.WithContext(ctx).Model(&Activity{})
	if actType != nil {
		query = query.Where("type = ?", *actType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("start_at DESC").Offset(offset).Limit(limit).Find(&activities).Error
	return activities, total, err
}

func (r *Repository) Update(ctx context.Context, a *Activity) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&Activity{}, id).Error
}

// Registration

func (r *Repository) CreateRegistration(ctx context.Context, reg *ActivityRegistration) error {
	return r.db.WithContext(ctx).Create(reg).Error
}

func (r *Repository) GetRegistration(ctx context.Context, activityID, userID int64) (*ActivityRegistration, error) {
	var reg ActivityRegistration
	err := r.db.WithContext(ctx).
		Where("activity_id = ? AND user_id = ? AND status = 'registered'", activityID, userID).
		First(&reg).Error
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *Repository) UpdateRegistration(ctx context.Context, reg *ActivityRegistration) error {
	return r.db.WithContext(ctx).Save(reg).Error
}

func (r *Repository) IncrementParticipants(ctx context.Context, activityID int64, delta int) error {
	return r.db.WithContext(ctx).
		Model(&Activity{}).
		Where("id = ?", activityID).
		Update("current_participants", gorm.Expr("current_participants + ?", delta)).Error
}

// Auto-cancel support

func (r *Repository) ListPendingCancelCheck(ctx context.Context, before time.Time) ([]Activity, error) {
	var activities []Activity
	err := r.db.WithContext(ctx).
		Where("status IN ('registration', 'published') AND cancel_check_at IS NOT NULL AND cancel_check_at <= ?", before).
		Find(&activities).Error
	return activities, err
}

// Transaction helper

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
