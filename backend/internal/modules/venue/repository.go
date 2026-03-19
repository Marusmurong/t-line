package venue

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

// Venue CRUD

func (r *Repository) Create(ctx context.Context, v *Venue) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Venue, error) {
	var v Venue
	if err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *Repository) List(ctx context.Context, status *int, offset, limit int) ([]Venue, int64, error) {
	var venues []Venue
	var total int64

	query := r.db.WithContext(ctx).Model(&Venue{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort_order ASC, id ASC").Offset(offset).Limit(limit).Find(&venues).Error
	return venues, total, err
}

func (r *Repository) Update(ctx context.Context, v *Venue) error {
	return r.db.WithContext(ctx).Save(v).Error
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&Venue{}, id).Error
}

// Time slot rules

func (r *Repository) CreateTimeRule(ctx context.Context, rule *VenueTimeSlotRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *Repository) GetTimeRuleByID(ctx context.Context, id int64) (*VenueTimeSlotRule, error) {
	var rule VenueTimeSlotRule
	if err := r.db.WithContext(ctx).First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *Repository) ListTimeRules(ctx context.Context, venueID int64) ([]VenueTimeSlotRule, error) {
	var rules []VenueTimeSlotRule
	err := r.db.WithContext(ctx).
		Where("venue_id = ?", venueID).
		Order("start_time ASC").
		Find(&rules).Error
	return rules, err
}

func (r *Repository) ListActiveTimeRules(ctx context.Context, venueID int64, dayType string) ([]VenueTimeSlotRule, error) {
	var rules []VenueTimeSlotRule
	err := r.db.WithContext(ctx).
		Where("venue_id = ? AND day_type = ? AND is_active = true", venueID, dayType).
		Order("start_time ASC").
		Find(&rules).Error
	return rules, err
}

func (r *Repository) UpdateTimeRule(ctx context.Context, rule *VenueTimeSlotRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

func (r *Repository) DeleteTimeRule(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&VenueTimeSlotRule{}, id).Error
}

// Blocked times

func (r *Repository) CreateBlockedTime(ctx context.Context, bt *VenueBlockedTime) error {
	return r.db.WithContext(ctx).Create(bt).Error
}

func (r *Repository) ListBlockedTimes(ctx context.Context, venueID int64, date time.Time) ([]VenueBlockedTime, error) {
	var blocked []VenueBlockedTime
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Where("venue_id = ? AND start_at < ? AND end_at > ?", venueID, dayEnd, dayStart).
		Find(&blocked).Error
	return blocked, err
}

func (r *Repository) DeleteBlockedTime(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&VenueBlockedTime{}, id).Error
}

// Transaction helper

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
