package activity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Activity struct {
	ID                   int64           `gorm:"primaryKey" json:"id"`
	Title                string          `gorm:"size:128;not null" json:"title"`
	Type                 string          `gorm:"size:20;not null" json:"type"`
	Description          string          `gorm:"type:text" json:"description"`
	CoverImage           string          `gorm:"size:512" json:"cover_image"`
	VenueID              *int64          `json:"venue_id"`
	StartAt              time.Time       `json:"start_at"`
	EndAt                time.Time       `json:"end_at"`
	RegistrationDeadline time.Time       `json:"registration_deadline"`
	MinParticipants      int             `gorm:"default:0" json:"min_participants"`
	MaxParticipants      int             `gorm:"default:0" json:"max_participants"`
	CurrentParticipants  int             `gorm:"default:0" json:"current_participants"`
	Price                decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	LevelRequirement     string          `gorm:"size:20;default:'all'" json:"level_requirement"`
	Status               string          `gorm:"size:20;default:'draft'" json:"status"`
	CancelCheckAt        *time.Time      `json:"cancel_check_at"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

func (Activity) TableName() string { return "activities" }

type ActivityRegistration struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	ActivityID int64     `json:"activity_id"`
	UserID     int64     `json:"user_id"`
	OrderID    *int64    `json:"order_id"`
	Status     string    `gorm:"size:20;default:'registered'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ActivityRegistration) TableName() string { return "activity_registrations" }
