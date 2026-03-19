package venue

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

type BookingChecker interface {
	IsSlotBooked(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) (bool, error)
}

type Service struct {
	repo           *Repository
	bookingChecker BookingChecker
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetBookingChecker(bc BookingChecker) {
	s.bookingChecker = bc
}

// --- User-facing ---

func (s *Service) ListVenues(ctx context.Context, offset, limit int) ([]VenueResp, int64, error) {
	active := 1
	venues, total, err := s.repo.List(ctx, &active, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]VenueResp, 0, len(venues))
	for i := range venues {
		result = append(result, ToVenueResp(&venues[i]))
	}
	return result, total, nil
}

func (s *Service) GetVenue(ctx context.Context, id int64) (*VenueResp, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}
	resp := ToVenueResp(v)
	return &resp, nil
}

func (s *Service) GetAvailability(ctx context.Context, venueID int64, date time.Time) ([]TimeSlotResp, error) {
	v, err := s.repo.GetByID(ctx, venueID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}
	if v.Status != 1 {
		return nil, apperrors.ErrVenueUnavailable
	}

	dayType := getDayType(date)
	rules, err := s.repo.ListActiveTimeRules(ctx, venueID, dayType)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	blockedTimes, err := s.repo.ListBlockedTimes(ctx, venueID, date)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	slots := make([]TimeSlotResp, 0, len(rules))
	for _, rule := range rules {
		available := true

		// Check blocked times
		for _, bt := range blockedTimes {
			if isTimeOverlap(date, rule.StartTime, rule.EndTime, bt.StartAt, bt.EndAt) {
				available = false
				break
			}
		}

		// Check existing bookings
		if available && s.bookingChecker != nil {
			booked, checkErr := s.bookingChecker.IsSlotBooked(ctx, venueID, date, rule.StartTime, rule.EndTime)
			if checkErr == nil && booked {
				available = false
			}
		}

		slots = append(slots, TimeSlotResp{
			StartTime: rule.StartTime,
			EndTime:   rule.EndTime,
			Price:     rule.Price.StringFixed(2),
			Available: available,
			DayType:   rule.DayType,
		})
	}

	return slots, nil
}

// --- Admin ---

func (s *Service) AdminListVenues(ctx context.Context, status *int, offset, limit int) ([]VenueResp, int64, error) {
	venues, total, err := s.repo.List(ctx, status, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]VenueResp, 0, len(venues))
	for i := range venues {
		result = append(result, ToVenueResp(&venues[i]))
	}
	return result, total, nil
}

func (s *Service) AdminCreateVenue(ctx context.Context, req CreateVenueReq) (*VenueResp, error) {
	v := &Venue{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Facilities:  req.Facilities,
		Status:      1,
		SortOrder:   req.SortOrder,
	}

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToVenueResp(v)
	return &resp, nil
}

func (s *Service) AdminUpdateVenue(ctx context.Context, id int64, req UpdateVenueReq) (*VenueResp, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if req.Name != nil {
		v.Name = *req.Name
	}
	if req.Type != nil {
		v.Type = *req.Type
	}
	if req.Description != nil {
		v.Description = *req.Description
	}
	if req.CoverImage != nil {
		v.CoverImage = *req.CoverImage
	}
	if req.Facilities != nil {
		v.Facilities = *req.Facilities
	}
	if req.Status != nil {
		v.Status = *req.Status
	}
	if req.SortOrder != nil {
		v.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToVenueResp(v)
	return &resp, nil
}

func (s *Service) AdminDeleteVenue(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperrors.ErrRecordNotFound
	}
	return s.repo.Delete(ctx, id)
}

// Time rules

func (s *Service) AdminCreateTimeRule(ctx context.Context, venueID int64, req CreateTimeRuleReq) (*VenueTimeSlotRule, error) {
	if _, err := s.repo.GetByID(ctx, venueID); err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rule := &VenueTimeSlotRule{
		VenueID:        venueID,
		DayType:        req.DayType,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Price:          price,
		MemberDiscount: req.MemberDiscount,
		IsActive:       isActive,
	}

	if err := s.repo.CreateTimeRule(ctx, rule); err != nil {
		return nil, apperrors.ErrInternal
	}

	return rule, nil
}

func (s *Service) AdminUpdateTimeRule(ctx context.Context, ruleID int64, req UpdateTimeRuleReq) (*VenueTimeSlotRule, error) {
	rule, err := s.repo.GetTimeRuleByID(ctx, ruleID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if req.DayType != nil {
		rule.DayType = *req.DayType
	}
	if req.StartTime != nil {
		rule.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		rule.EndTime = *req.EndTime
	}
	if req.Price != nil {
		price, pErr := decimal.NewFromString(*req.Price)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		rule.Price = price
	}
	if req.MemberDiscount != nil {
		rule.MemberDiscount = *req.MemberDiscount
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateTimeRule(ctx, rule); err != nil {
		return nil, apperrors.ErrInternal
	}

	return rule, nil
}

func (s *Service) AdminDeleteTimeRule(ctx context.Context, ruleID int64) error {
	if _, err := s.repo.GetTimeRuleByID(ctx, ruleID); err != nil {
		return apperrors.ErrRecordNotFound
	}
	return s.repo.DeleteTimeRule(ctx, ruleID)
}

func (s *Service) AdminListTimeRules(ctx context.Context, venueID int64) ([]VenueTimeSlotRule, error) {
	return s.repo.ListTimeRules(ctx, venueID)
}

// GetTimeSlotPrice returns price for a given venue, date, start/end time
func (s *Service) GetTimeSlotPrice(ctx context.Context, venueID int64, date time.Time, startTime, endTime string) (decimal.Decimal, error) {
	dayType := getDayType(date)
	rules, err := s.repo.ListActiveTimeRules(ctx, venueID, dayType)
	if err != nil {
		return decimal.Zero, apperrors.ErrInternal
	}

	for _, rule := range rules {
		if rule.StartTime == startTime && rule.EndTime == endTime {
			return rule.Price, nil
		}
	}

	return decimal.Zero, apperrors.ErrSlotUnavailable
}

// helpers

func getDayType(date time.Time) string {
	weekday := date.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return "weekend"
	}
	return "weekday"
}

func isTimeOverlap(date time.Time, slotStart, slotEnd string, blockStart, blockEnd time.Time) bool {
	loc := date.Location()
	dateStr := date.Format("2006-01-02")

	slotStartTime, err1 := time.ParseInLocation("2006-01-02 15:04:05", dateStr+" "+slotStart+":00", loc)
	slotEndTime, err2 := time.ParseInLocation("2006-01-02 15:04:05", dateStr+" "+slotEnd+":00", loc)
	if err1 != nil || err2 != nil {
		return false
	}

	return slotStartTime.Before(blockEnd) && slotEndTime.After(blockStart)
}
