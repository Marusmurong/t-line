package booking

import (
	"time"

	"github.com/shopspring/decimal"
)

type Booking struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	UserID        int64           `gorm:"not null" json:"user_id"`
	VenueID       int64           `gorm:"not null" json:"venue_id"`
	Date          time.Time       `gorm:"type:date;not null" json:"date"`
	StartTime     string          `gorm:"type:time;not null" json:"start_time"`
	EndTime       string          `gorm:"type:time;not null" json:"end_time"`
	DurationHours decimal.Decimal `gorm:"type:decimal(3,1);not null" json:"duration_hours"`
	Status        string          `gorm:"size:20;not null;default:pending" json:"status"`
	TotalAmount   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	OrderID       *int64          `json:"order_id"`
	CancelReason  string          `gorm:"size:256" json:"cancel_reason"`
	CancelledAt   *time.Time      `json:"cancelled_at"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (Booking) TableName() string { return "bookings" }

// Booking status constants
const (
	StatusPending    = "pending"
	StatusConfirmed  = "confirmed"
	StatusWaitlisted = "waitlisted"
	StatusCancelled  = "cancelled"
	StatusCompleted  = "completed"
	StatusNoShow     = "no_show"
)

type BookingWaitlist struct {
	ID         int64      `gorm:"primaryKey" json:"id"`
	BookingID  *int64     `json:"booking_id"`
	UserID     int64      `gorm:"not null" json:"user_id"`
	VenueID    int64      `gorm:"not null" json:"venue_id"`
	Date       time.Time  `gorm:"type:date;not null" json:"date"`
	StartTime  string     `gorm:"type:time;not null" json:"start_time"`
	EndTime    string     `gorm:"type:time;not null" json:"end_time"`
	Position   int        `gorm:"not null" json:"position"`
	Status     string     `gorm:"size:20;not null;default:waiting" json:"status"`
	NotifiedAt *time.Time `json:"notified_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (BookingWaitlist) TableName() string { return "booking_waitlist" }

// Waitlist status constants
const (
	WaitlistWaiting   = "waiting"
	WaitlistNotified  = "notified"
	WaitlistConfirmed = "confirmed"
	WaitlistExpired   = "expired"
	WaitlistCancelled = "cancelled"
)
