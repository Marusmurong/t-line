package scheduler

import (
	"context"
	"time"

	"github.com/t-line/backend/internal/pkg/logger"
	"gorm.io/gorm"
)

// checkActivityCancellation evaluates activities past their cancellation check time.
// If minimum participants are not met, the activity is cancelled and all registrants refunded.
// If the threshold is met, the activity is confirmed.
func (s *Scheduler) checkActivityCancellation() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var activityIDs []int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT id FROM activities
		WHERE status = 'registration'
		  AND cancel_check_at IS NOT NULL
		  AND cancel_check_at < NOW()
	`).Scan(&activityIDs).Error
	if err != nil {
		logger.L.Errorw("scheduler: failed to find activities for cancel check", "error", err)
		return
	}

	if len(activityIDs) == 0 {
		return
	}

	logger.L.Infow("scheduler: checking activities for cancellation", "count", len(activityIDs))

	for _, id := range activityIDs {
		s.evaluateActivity(ctx, id)
	}
}

func (s *Scheduler) evaluateActivity(ctx context.Context, activityID int64) {
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

	// Lock and read the activity
	var activity struct {
		ID                  int64
		MinParticipants     int
		CurrentParticipants int
		Status              string
	}
	err := tx.Raw(`
		SELECT id, min_participants, current_participants, status
		FROM activities
		WHERE id = ? AND status = 'registration'
		FOR UPDATE
	`, activityID).Scan(&activity).Error
	if err != nil {
		tx.Rollback()
		logger.L.Errorw("scheduler: failed to lock activity", "activity_id", activityID, "error", err)
		return
	}
	if activity.ID == 0 {
		tx.Rollback()
		return
	}

	if activity.CurrentParticipants < activity.MinParticipants {
		// Insufficient participants: cancel activity
		s.cancelActivity(tx, activityID)
	} else {
		// Sufficient participants: confirm activity
		tx.Exec(`
			UPDATE activities SET status = 'confirmed', updated_at = NOW()
			WHERE id = ?
		`, activityID)
	}

	if err := tx.Commit().Error; err != nil {
		logger.L.Errorw("scheduler: failed to commit activity evaluation", "activity_id", activityID, "error", err)
	}
}

func (s *Scheduler) cancelActivity(tx *gorm.DB, activityID int64) {
	// Cancel the activity
	tx.Exec(`
		UPDATE activities SET status = 'cancelled', updated_at = NOW()
		WHERE id = ?
	`, activityID)

	// Get all registered users' order IDs for full refund
	var registrations []struct {
		ID      int64
		UserID  int64
		OrderID *int64
	}
	tx.Raw(`
		SELECT id, user_id, order_id
		FROM activity_registrations
		WHERE activity_id = ? AND status = 'registered'
	`, activityID).Scan(&registrations)

	// Cancel all registrations
	tx.Exec(`
		UPDATE activity_registrations
		SET status = 'cancelled'
		WHERE activity_id = ? AND status = 'registered'
	`, activityID)

	// Refund all associated paid orders
	for _, reg := range registrations {
		if reg.OrderID == nil {
			continue
		}

		// Mark order as refunding
		tx.Exec(`
			UPDATE orders
			SET status = 'refunded', updated_at = NOW()
			WHERE id = ? AND status = 'paid'
		`, *reg.OrderID)

		// Create refund record
		tx.Exec(`
			INSERT INTO refunds (refund_no, order_id, user_id, amount, balance_refund, wechat_refund, reason, status, created_at)
			SELECT
				'RF' || TO_CHAR(NOW(), 'YYYYMMDD') || LPAD(CAST(nextval('refund_no_seq') AS TEXT), 6, '0'),
				o.id, o.user_id, o.pay_amount, o.balance_paid, o.wechat_paid,
				'活动人数不足自动取消', 'approved', NOW()
			FROM orders o
			WHERE o.id = ? AND o.status = 'refunded'
		`, *reg.OrderID)

		// Refund balance to user wallet if balance was used
		tx.Exec(`
			UPDATE wallets w
			SET balance = balance + o.balance_paid, updated_at = NOW()
			FROM orders o
			WHERE o.id = ?
			  AND w.user_id = o.user_id
			  AND o.balance_paid > 0
		`, *reg.OrderID)
	}
}
