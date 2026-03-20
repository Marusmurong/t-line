package booking

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

// Booking operations

func (r *Repository) Create(ctx context.Context, b *Booking) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Booking, error) {
	var b Booking
	if err := r.db.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID int64, upcoming bool, offset, limit int) ([]Booking, int64, error) {
	var bookings []Booking
	var total int64

	query := r.db.WithContext(ctx).Model(&Booking{}).Where("user_id = ?", userID)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if upcoming {
		query = query.Where("date >= ? AND status IN ?", today, []string{StatusPending, StatusConfirmed})
	} else {
		query = query.Where("date < ? OR status IN ?", today, []string{StatusCancelled, StatusCompleted, StatusNoShow})
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderDir := "date ASC, start_time ASC"
	if !upcoming {
		orderDir = "date DESC, start_time DESC"
	}

	err := query.Order(orderDir).Offset(offset).Limit(limit).Find(&bookings).Error
	return bookings, total, err
}

func (r *Repository) Update(ctx context.Context, b *Booking) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&Booking{}, id).Error
}

// GetByOrderID finds a booking by its associated order ID.
func (r *Repository) GetByOrderID(ctx context.Context, orderID int64) (*Booking, error) {
	var b Booking
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) IsSlotBooked(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Booking{}).
		Where("venue_id = ? AND date = ? AND start_time = ? AND end_time = ? AND status IN ?",
			venueID, date, startTime, endTime, []string{StatusPending, StatusConfirmed}).
		Count(&count).Error
	return count > 0, err
}

// Waitlist operations

func (r *Repository) CreateWaitlistEntry(ctx context.Context, w *BookingWaitlist) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *Repository) GetWaitlistPosition(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&BookingWaitlist{}).
		Where("venue_id = ? AND date = ? AND start_time = ? AND end_time = ? AND status = ?",
			venueID, date, startTime, endTime, WaitlistWaiting).
		Count(&count).Error
	return int(count), err
}

func (r *Repository) GetFirstWaiting(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) (*BookingWaitlist, error) {
	var w BookingWaitlist
	err := r.db.WithContext(ctx).
		Where("venue_id = ? AND date = ? AND start_time = ? AND end_time = ? AND status = ?",
			venueID, date, startTime, endTime, WaitlistWaiting).
		Order("position ASC").
		First(&w).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *Repository) UpdateWaitlistEntry(ctx context.Context, w *BookingWaitlist) error {
	return r.db.WithContext(ctx).Save(w).Error
}

// Transaction helper

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
