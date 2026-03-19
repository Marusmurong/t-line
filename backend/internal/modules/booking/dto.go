package booking

import "time"

// --- Request DTOs ---

type CreateBookingReq struct {
	VenueID   int64  `json:"venue_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type CancelBookingReq struct {
	Reason string `json:"reason" binding:"omitempty,max=256"`
}

type JoinWaitlistReq struct {
	VenueID   int64  `json:"venue_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

// --- Response DTOs ---

type BookingResp struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	VenueID       int64      `json:"venue_id"`
	VenueName     string     `json:"venue_name,omitempty"`
	Date          string     `json:"date"`
	StartTime     string     `json:"start_time"`
	EndTime       string     `json:"end_time"`
	DurationHours string     `json:"duration_hours"`
	Status        string     `json:"status"`
	TotalAmount   string     `json:"total_amount"`
	OrderID       *int64     `json:"order_id"`
	CancelReason  string     `json:"cancel_reason,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type WaitlistResp struct {
	ID        int64      `json:"id"`
	VenueID   int64      `json:"venue_id"`
	Date      string     `json:"date"`
	StartTime string     `json:"start_time"`
	EndTime   string     `json:"end_time"`
	Position  int        `json:"position"`
	Status    string     `json:"status"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func ToBookingResp(b *Booking) BookingResp {
	return BookingResp{
		ID:            b.ID,
		UserID:        b.UserID,
		VenueID:       b.VenueID,
		Date:          b.Date.Format("2006-01-02"),
		StartTime:     b.StartTime,
		EndTime:       b.EndTime,
		DurationHours: b.DurationHours.StringFixed(1),
		Status:        b.Status,
		TotalAmount:   b.TotalAmount.StringFixed(2),
		OrderID:       b.OrderID,
		CancelReason:  b.CancelReason,
		CancelledAt:   b.CancelledAt,
		CreatedAt:     b.CreatedAt,
	}
}

func ToWaitlistResp(w *BookingWaitlist) WaitlistResp {
	return WaitlistResp{
		ID:        w.ID,
		VenueID:   w.VenueID,
		Date:      w.Date.Format("2006-01-02"),
		StartTime: w.StartTime,
		EndTime:   w.EndTime,
		Position:  w.Position,
		Status:    w.Status,
		ExpiresAt: w.ExpiresAt,
		CreatedAt: w.CreatedAt,
	}
}
