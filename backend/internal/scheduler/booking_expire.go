package scheduler

import (
	"context"
	"time"

	"github.com/t-line/backend/internal/pkg/logger"
)

// expireBookingWaitlist handles timed-out waitlist notifications.
// It marks expired entries, and promotes the next waiting user in the queue.
func (s *Scheduler) expireBookingWaitlist() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find all notified waitlist entries that have expired
	var expiredIDs []int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT id FROM booking_waitlist
		WHERE status = 'notified'
		  AND expires_at IS NOT NULL
		  AND expires_at < NOW()
	`).Scan(&expiredIDs).Error
	if err != nil {
		logger.L.Errorw("scheduler: failed to find expired waitlist entries", "error", err)
		return
	}

	if len(expiredIDs) == 0 {
		return
	}

	logger.L.Infow("scheduler: processing expired waitlist entries", "count", len(expiredIDs))

	for _, id := range expiredIDs {
		s.processExpiredWaitlistEntry(ctx, id)
	}
}

func (s *Scheduler) processExpiredWaitlistEntry(ctx context.Context, waitlistID int64) {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.L.Errorw("scheduler: failed to begin transaction", "error", tx.Error)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mark the expired entry
	result := tx.Exec(`
		UPDATE booking_waitlist
		SET status = 'expired'
		WHERE id = ? AND status = 'notified'
	`, waitlistID)
	if result.Error != nil {
		tx.Rollback()
		logger.L.Errorw("scheduler: failed to expire waitlist entry", "id", waitlistID, "error", result.Error)
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return
	}

	// Get the expired entry's venue/date/time info to find next in queue
	var entry struct {
		VenueID   int64
		Date      time.Time
		StartTime string
		EndTime   string
	}
	if err := tx.Raw(`
		SELECT venue_id, date, start_time, end_time
		FROM booking_waitlist WHERE id = ?
	`, waitlistID).Scan(&entry).Error; err != nil {
		tx.Rollback()
		logger.L.Errorw("scheduler: failed to get waitlist entry details", "id", waitlistID, "error", err)
		return
	}

	// Find the next waiting user for the same slot
	var nextID *int64
	tx.Raw(`
		SELECT id FROM booking_waitlist
		WHERE venue_id = ? AND date = ? AND start_time = ? AND end_time = ?
		  AND status = 'waiting'
		ORDER BY position ASC
		LIMIT 1
	`, entry.VenueID, entry.Date, entry.StartTime, entry.EndTime).Scan(&nextID)

	if nextID != nil {
		// Notify the next user: set status to notified with a new expiry window
		expiresAt := time.Now().Add(30 * time.Minute)
		tx.Exec(`
			UPDATE booking_waitlist
			SET status = 'notified', notified_at = NOW(), expires_at = ?
			WHERE id = ?
		`, expiresAt, *nextID)
	}

	if err := tx.Commit().Error; err != nil {
		logger.L.Errorw("scheduler: failed to commit waitlist expiry", "id", waitlistID, "error", err)
	}
}
