package activity

import (
	"context"
	"time"

	"github.com/t-line/backend/internal/pkg/logger"
)

// AutoCancelUnderfilledActivities checks activities past their cancel_check_at
// time and cancels those with fewer participants than min_participants.
// This method is designed to be called by a scheduler/cron job.
func (s *Service) AutoCancelUnderfilledActivities(ctx context.Context) error {
	activities, err := s.repo.ListPendingCancelCheck(ctx, time.Now())
	if err != nil {
		logger.L.Errorf("auto-cancel: failed to list activities: %v", err)
		return err
	}

	for _, a := range activities {
		if a.CurrentParticipants < a.MinParticipants {
			a.Status = "cancelled"
			if updateErr := s.repo.Update(ctx, &a); updateErr != nil {
				logger.L.Errorf("auto-cancel: failed to cancel activity %d: %v", a.ID, updateErr)
				continue
			}
			logger.L.Infof("auto-cancel: activity %d cancelled (participants: %d, min: %d)",
				a.ID, a.CurrentParticipants, a.MinParticipants)
		} else {
			// Enough participants, confirm the activity
			a.Status = "confirmed"
			if updateErr := s.repo.Update(ctx, &a); updateErr != nil {
				logger.L.Errorf("auto-cancel: failed to confirm activity %d: %v", a.ID, updateErr)
				continue
			}
			logger.L.Infof("auto-cancel: activity %d confirmed (participants: %d, min: %d)",
				a.ID, a.CurrentParticipants, a.MinParticipants)
		}
	}

	return nil
}
