package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

const maxWaitlistSize = 10

// JoinWaitlist adds a user to the waitlist for a specific slot.
func (s *Service) JoinWaitlist(ctx context.Context, userID int64, req JoinWaitlistReq) (*WaitlistResp, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	// Check current waitlist size
	position, err := s.repo.GetWaitlistPosition(ctx, req.VenueID, date, req.StartTime, req.EndTime)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	if position >= maxWaitlistSize {
		return nil, apperrors.ErrWaitlistFull
	}

	// Add to Redis sorted set for fast lookup
	redisKey := fmt.Sprintf("waitlist:%d:%s:%s:%s", req.VenueID, req.Date, req.StartTime, req.EndTime)
	score := float64(time.Now().UnixMilli())
	_ = s.rdb.ZAdd(ctx, redisKey, redis.Z{Score: score, Member: userID}).Err()

	entry := &BookingWaitlist{
		UserID:    userID,
		VenueID:   req.VenueID,
		Date:      date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Position:  position + 1,
		Status:    WaitlistWaiting,
	}

	if err := s.repo.CreateWaitlistEntry(ctx, entry); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToWaitlistResp(entry)
	return &resp, nil
}

// HandleWaitlistTimeout handles expired waitlist notifications.
// This would be called by a scheduler to mark expired entries and notify the next user.
func (s *Service) HandleWaitlistTimeout(_ context.Context, _ int64) error {
	return nil
}
