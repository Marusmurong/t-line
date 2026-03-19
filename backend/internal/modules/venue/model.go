package venue

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Venue struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:64;not null" json:"name"`
	Type        string         `gorm:"size:20;not null" json:"type"`
	Description string         `gorm:"type:text" json:"description"`
	CoverImage  string         `gorm:"size:512" json:"cover_image"`
	Facilities  datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"facilities"`
	Status      int            `gorm:"default:1" json:"status"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Venue) TableName() string { return "venues" }

type VenueTimeSlotRule struct {
	ID             int64           `gorm:"primaryKey" json:"id"`
	VenueID        int64           `json:"venue_id"`
	DayType        string          `gorm:"size:20;not null" json:"day_type"`
	StartTime      string          `gorm:"type:time;not null" json:"start_time"`
	EndTime        string          `gorm:"type:time;not null" json:"end_time"`
	Price          decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`
	MemberDiscount datatypes.JSON  `gorm:"type:jsonb;default:'{}'" json:"member_discount"`
	IsActive       bool            `gorm:"default:true" json:"is_active"`
}

func (VenueTimeSlotRule) TableName() string { return "venue_time_slot_rules" }

type VenueBlockedTime struct {
	ID      int64     `gorm:"primaryKey" json:"id"`
	VenueID int64     `json:"venue_id"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Reason  string    `gorm:"size:256" json:"reason"`
}

func (VenueBlockedTime) TableName() string { return "venue_blocked_times" }
