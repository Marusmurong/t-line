package venue

import (
	"time"

	"gorm.io/datatypes"
)

// --- Request DTOs ---

type CreateVenueReq struct {
	Name        string         `json:"name" binding:"required,max=64"`
	Type        string         `json:"type" binding:"required,oneof=indoor_hard outdoor_hard indoor_clay outdoor_clay"`
	Description string         `json:"description"`
	CoverImage  string         `json:"cover_image" binding:"omitempty,max=512"`
	Facilities  datatypes.JSON `json:"facilities"`
	SortOrder   int            `json:"sort_order"`
}

type UpdateVenueReq struct {
	Name        *string         `json:"name" binding:"omitempty,max=64"`
	Type        *string         `json:"type" binding:"omitempty,oneof=indoor_hard outdoor_hard indoor_clay outdoor_clay"`
	Description *string         `json:"description"`
	CoverImage  *string         `json:"cover_image" binding:"omitempty,max=512"`
	Facilities  *datatypes.JSON `json:"facilities"`
	Status      *int            `json:"status" binding:"omitempty,oneof=0 1 2"`
	SortOrder   *int            `json:"sort_order"`
}

type CreateTimeRuleReq struct {
	DayType        string         `json:"day_type" binding:"required,oneof=weekday weekend holiday"`
	StartTime      string         `json:"start_time" binding:"required"`
	EndTime        string         `json:"end_time" binding:"required"`
	Price          string         `json:"price" binding:"required"`
	MemberDiscount datatypes.JSON `json:"member_discount"`
	IsActive       *bool          `json:"is_active"`
}

type UpdateTimeRuleReq struct {
	DayType        *string         `json:"day_type" binding:"omitempty,oneof=weekday weekend holiday"`
	StartTime      *string         `json:"start_time"`
	EndTime        *string         `json:"end_time"`
	Price          *string         `json:"price"`
	MemberDiscount *datatypes.JSON `json:"member_discount"`
	IsActive       *bool           `json:"is_active"`
}

type CreateBlockedTimeReq struct {
	StartAt time.Time `json:"start_at" binding:"required"`
	EndAt   time.Time `json:"end_at" binding:"required"`
	Reason  string    `json:"reason"`
}

// --- Response DTOs ---

type VenueResp struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	CoverImage  string         `json:"cover_image"`
	Facilities  datatypes.JSON `json:"facilities"`
	Status      int            `json:"status"`
	SortOrder   int            `json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type TimeSlotResp struct {
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Price          string `json:"price"`
	MemberPrice    string `json:"member_price,omitempty"`
	Available      bool   `json:"available"`
	DayType        string `json:"day_type"`
}

type TimeGridSlot struct {
	VenueID   int64  `json:"venue_id"`
	VenueName string `json:"venue_name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	BookingID *int64 `json:"booking_id,omitempty"`
	UserName  string `json:"user_name,omitempty"`
}

func ToVenueResp(v *Venue) VenueResp {
	return VenueResp{
		ID:          v.ID,
		Name:        v.Name,
		Type:        v.Type,
		Description: v.Description,
		CoverImage:  v.CoverImage,
		Facilities:  v.Facilities,
		Status:      v.Status,
		SortOrder:   v.SortOrder,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}
}
