package scheduler

import (
	"context"
	"time"

	"github.com/t-line/backend/internal/pkg/logger"
)

// expireUnpaidOrders closes orders that have exceeded their payment deadline,
// releases associated booking slots and product stock.
func (s *Scheduler) expireUnpaidOrders() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find expired pending orders
	var expiredIDs []int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT id FROM orders
		WHERE status = 'pending'
		  AND expires_at IS NOT NULL
		  AND expires_at < NOW()
	`).Scan(&expiredIDs).Error
	if err != nil {
		logger.L.Errorw("scheduler: failed to find expired orders", "error", err)
		return
	}

	if len(expiredIDs) == 0 {
		return
	}

	logger.L.Infow("scheduler: processing expired orders", "count", len(expiredIDs))

	for _, id := range expiredIDs {
		s.cancelExpiredOrder(ctx, id)
	}
}

func (s *Scheduler) cancelExpiredOrder(ctx context.Context, orderID int64) {
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

	// Cancel the order using optimistic locking on status
	result := tx.Exec(`
		UPDATE orders
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = ? AND status = 'pending'
	`, orderID)
	if result.Error != nil {
		tx.Rollback()
		logger.L.Errorw("scheduler: failed to cancel order", "order_id", orderID, "error", result.Error)
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return
	}

	// Get order type to determine what to release
	var orderType string
	if err := tx.Raw(`SELECT type FROM orders WHERE id = ?`, orderID).Scan(&orderType).Error; err != nil {
		tx.Rollback()
		logger.L.Errorw("scheduler: failed to get order type", "order_id", orderID, "error", err)
		return
	}

	// Release resources based on order type
	switch orderType {
	case "booking":
		// Cancel associated booking
		tx.Exec(`
			UPDATE bookings
			SET status = 'cancelled', cancel_reason = '支付超时自动取消', cancelled_at = NOW(), updated_at = NOW()
			WHERE order_id = ? AND status IN ('pending', 'confirmed')
		`, orderID)

	case "product":
		// Restore product stock
		tx.Exec(`
			UPDATE product_skus ps
			SET stock = ps.stock + oi.quantity
			FROM order_items oi
			WHERE oi.order_id = ?
			  AND oi.item_type = 'product'
			  AND oi.sku_id IS NOT NULL
			  AND ps.id = oi.sku_id
		`, orderID)

	case "activity":
		// Cancel activity registration and decrement participant count
		tx.Exec(`
			UPDATE activity_registrations
			SET status = 'cancelled'
			WHERE order_id = ? AND status = 'registered'
		`, orderID)
		tx.Exec(`
			UPDATE activities a
			SET current_participants = current_participants - sub.cnt
			FROM (
				SELECT activity_id, COUNT(*) AS cnt
				FROM activity_registrations
				WHERE order_id = ? AND status = 'cancelled'
				GROUP BY activity_id
			) sub
			WHERE a.id = sub.activity_id
			  AND a.current_participants >= sub.cnt
		`, orderID)
	}

	if err := tx.Commit().Error; err != nil {
		logger.L.Errorw("scheduler: failed to commit order cancellation", "order_id", orderID, "error", err)
	}
}
