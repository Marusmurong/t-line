package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

const (
	slotLockPrefix = "booking:lock:"
	slotLockTTL    = 10 * time.Minute
)

// VenuePricer abstracts venue price lookup.
type VenuePricer interface {
	GetTimeSlotPrice(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) (decimal.Decimal, error)
}

// OrderCreator abstracts order creation from order module.
type OrderCreator interface {
	CreateOrderForBooking(ctx context.Context, userID int64, bookingID int64, venueName string, date, startTime, endTime string, amount decimal.Decimal) (int64, error)
}

type Service struct {
	repo         *Repository
	rdb          *redis.Client
	venuePricer  VenuePricer
	orderCreator OrderCreator
}

func NewService(repo *Repository, rdb *redis.Client) *Service {
	return &Service{
		repo: repo,
		rdb:  rdb,
	}
}

func (s *Service) SetVenuePricer(vp VenuePricer) {
	s.venuePricer = vp
}

func (s *Service) SetOrderCreator(oc OrderCreator) {
	s.orderCreator = oc
}

// CreateBooking: check slot -> Redis SETNX lock -> create booking -> create order
func (s *Service) CreateBooking(ctx context.Context, userID int64, req CreateBookingReq) (*BookingResp, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	// Check if slot is already booked
	booked, err := s.repo.IsSlotBooked(ctx, req.VenueID, date, req.StartTime, req.EndTime)
	if err != nil {
		return nil, apperrors.ErrInternal
	}
	if booked {
		return nil, apperrors.ErrSlotUnavailable
	}

	// Redis SETNX lock to prevent race conditions
	lockKey := fmt.Sprintf("%s%d:%s:%s:%s", slotLockPrefix, req.VenueID, req.Date, req.StartTime, req.EndTime)
	ok, err := s.rdb.SetNX(ctx, lockKey, userID, slotLockTTL).Result()
	if err != nil {
		return nil, apperrors.ErrInternal
	}
	if !ok {
		return nil, apperrors.ErrSlotConflict
	}

	// Calculate duration and price
	duration := calculateDuration(req.StartTime, req.EndTime)
	price := decimal.Zero
	if s.venuePricer != nil {
		price, err = s.venuePricer.GetTimeSlotPrice(ctx, req.VenueID, date, req.StartTime, req.EndTime)
		if err != nil {
			_ = s.rdb.Del(ctx, lockKey)
			return nil, err
		}
	}

	totalAmount := price.Mul(duration)

	booking := &Booking{
		UserID:        userID,
		VenueID:       req.VenueID,
		Date:          date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		DurationHours: duration,
		Status:        StatusPending,
		TotalAmount:   totalAmount,
	}

	if err := s.repo.Create(ctx, booking); err != nil {
		_ = s.rdb.Del(ctx, lockKey)
		return nil, apperrors.ErrInternal
	}

	resp := ToBookingResp(booking)
	return &resp, nil
}

// CancelBooking: release lock -> check waitlist -> trigger refund
func (s *Service) CancelBooking(ctx context.Context, userID, bookingID int64, reason string) error {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return apperrors.ErrBookingNotFound
	}

	if booking.UserID != userID {
		return apperrors.ErrForbidden
	}

	if booking.Status == StatusCancelled {
		return apperrors.ErrBookingCancelled
	}

	if booking.Status != StatusPending && booking.Status != StatusConfirmed {
		return apperrors.New(40203, "当前状态不可取消")
	}

	now := time.Now()
	booking.Status = StatusCancelled
	booking.CancelReason = reason
	booking.CancelledAt = &now

	if err := s.repo.Update(ctx, booking); err != nil {
		return apperrors.ErrInternal
	}

	// Release Redis lock
	lockKey := fmt.Sprintf("%s%d:%s:%s:%s", slotLockPrefix, booking.VenueID,
		booking.Date.Format("2006-01-02"), booking.StartTime, booking.EndTime)
	_ = s.rdb.Del(ctx, lockKey)

	// Notify first waitlisted user
	go s.notifyFirstWaitlisted(context.Background(), booking.VenueID, booking.Date, booking.StartTime, booking.EndTime)

	return nil
}

// GetBooking returns a booking by ID.
func (s *Service) GetBooking(ctx context.Context, userID, bookingID int64) (*BookingResp, error) {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, apperrors.ErrBookingNotFound
	}

	if booking.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	resp := ToBookingResp(booking)
	return &resp, nil
}

// ListBookings returns user's bookings (upcoming or history).
func (s *Service) ListBookings(ctx context.Context, userID int64, upcoming bool, offset, limit int) ([]BookingResp, int64, error) {
	bookings, total, err := s.repo.ListByUserID(ctx, userID, upcoming, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]BookingResp, 0, len(bookings))
	for i := range bookings {
		result = append(result, ToBookingResp(&bookings[i]))
	}
	return result, total, nil
}

// IsSlotBooked implements venue.BookingChecker interface.
func (s *Service) IsSlotBooked(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) (bool, error) {
	return s.repo.IsSlotBooked(ctx, venueID, date, startTime, endTime)
}

// notifyFirstWaitlisted notifies the first waiting user when a slot becomes available.
func (s *Service) notifyFirstWaitlisted(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) {
	waitlist, err := s.repo.GetFirstWaiting(ctx, venueID, date, startTime, endTime)
	if err != nil {
		return // No one waiting
	}

	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)
	waitlist.Status = WaitlistNotified
	waitlist.NotifiedAt = &now
	waitlist.ExpiresAt = &expiresAt

	_ = s.repo.UpdateWaitlistEntry(ctx, waitlist)
	// TODO: send notification via notify module
}

func calculateDuration(startTime, endTime string) decimal.Decimal {
	start, err1 := time.Parse("15:04", startTime)
	end, err2 := time.Parse("15:04", endTime)
	if err1 != nil || err2 != nil {
		return decimal.NewFromFloat(1.0)
	}
	hours := end.Sub(start).Hours()
	if hours <= 0 {
		hours = 1.0
	}
	return decimal.NewFromFloat(hours)
}
