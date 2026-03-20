package academic

import (
	"context"
	"fmt"

	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

// CheckConflict validates that a schedule does not conflict with existing
// coach schedules, venue bookings, or approved coach leaves.
// excludeScheduleID allows skipping a specific schedule (useful for updates).
func (s *Service) CheckConflict(ctx context.Context, coachID, venueID int64, date, startTime, endTime string, excludeScheduleID int64) error {
	// 1. Check coach time overlap
	coachConflict, err := s.repo.FindCoachConflict(ctx, coachID, date, startTime, endTime, excludeScheduleID)
	if err != nil {
		return apperrors.ErrInternal
	}
	if coachConflict {
		return apperrors.New(
			apperrors.ErrScheduleConflict.Code,
			fmt.Sprintf("教练在 %s %s-%s 时段已有其他课程", date, startTime, endTime),
		)
	}

	// 2. Check venue time overlap
	venueConflict, err := s.repo.FindVenueConflict(ctx, venueID, date, startTime, endTime, excludeScheduleID)
	if err != nil {
		return apperrors.ErrInternal
	}
	if venueConflict {
		return apperrors.New(
			apperrors.ErrVenueUnavailable.Code,
			fmt.Sprintf("场地在 %s %s-%s 时段已被占用", date, startTime, endTime),
		)
	}

	// 3. Check coach leave
	onLeave, err := s.repo.FindCoachLeaveConflict(ctx, coachID, date)
	if err != nil {
		return apperrors.ErrInternal
	}
	if onLeave {
		return apperrors.New(
			apperrors.ErrCoachUnavailable.Code,
			fmt.Sprintf("教练在 %s 已请假", date),
		)
	}

	return nil
}
