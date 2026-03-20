package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/t-line/backend/internal/pkg/logger"
)

// aggregateDailyStats pre-computes the previous day's statistics
// and stores them in Redis for fast dashboard reads.
func (s *Scheduler) aggregateDailyStats() {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	yesterday := time.Now().AddDate(0, 0, -1)
	dateStr := yesterday.Format("2006-01-02")

	logger.L.Infow("scheduler: aggregating daily stats", "date", dateStr)

	stats, err := s.computeDailyStats(ctx, dateStr)
	if err != nil {
		logger.L.Errorw("scheduler: failed to compute daily stats", "date", dateStr, "error", err)
		return
	}

	data, err := json.Marshal(stats)
	if err != nil {
		logger.L.Errorw("scheduler: failed to marshal daily stats", "error", err)
		return
	}

	cacheKey := fmt.Sprintf("stats:daily:%s", dateStr)
	if err := s.rdb.Set(ctx, cacheKey, data, 30*24*time.Hour).Err(); err != nil {
		logger.L.Errorw("scheduler: failed to cache daily stats", "error", err)
		return
	}

	logger.L.Infow("scheduler: daily stats aggregated successfully", "date", dateStr)
}

type dailyStats struct {
	Date           string          `json:"date"`
	Revenue        decimal.Decimal `json:"revenue"`
	OrderCount     int             `json:"order_count"`
	NewUsers       int             `json:"new_users"`
	ActiveUsers    int             `json:"active_users"`
	BookingCount   int             `json:"booking_count"`
	VenueUsageRate float64         `json:"venue_usage_rate"`
}

func (s *Scheduler) computeDailyStats(ctx context.Context, dateStr string) (*dailyStats, error) {
	stats := &dailyStats{Date: dateStr}

	// Daily revenue
	var revenue *decimal.Decimal
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(pay_amount), 0) FROM orders
		WHERE paid_at >= ?::date AND paid_at < ?::date + INTERVAL '1 day'
		  AND status NOT IN ('cancelled', 'refunded')
	`, dateStr, dateStr).Scan(&revenue).Error; err != nil {
		return nil, fmt.Errorf("revenue query: %w", err)
	}
	if revenue != nil {
		stats.Revenue = *revenue
	}

	// Daily order count
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM orders
		WHERE created_at >= ?::date AND created_at < ?::date + INTERVAL '1 day'
	`, dateStr, dateStr).Scan(&stats.OrderCount).Error; err != nil {
		return nil, fmt.Errorf("order count query: %w", err)
	}

	// New users
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM users
		WHERE created_at >= ?::date AND created_at < ?::date + INTERVAL '1 day'
	`, dateStr, dateStr).Scan(&stats.NewUsers).Error; err != nil {
		return nil, fmt.Errorf("new users query: %w", err)
	}

	// Active users (placed orders that day)
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT user_id) FROM orders
		WHERE created_at >= ?::date AND created_at < ?::date + INTERVAL '1 day'
	`, dateStr, dateStr).Scan(&stats.ActiveUsers).Error; err != nil {
		return nil, fmt.Errorf("active users query: %w", err)
	}

	// Booking count
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM bookings
		WHERE date = ?::date AND status NOT IN ('cancelled')
	`, dateStr).Scan(&stats.BookingCount).Error; err != nil {
		return nil, fmt.Errorf("booking count query: %w", err)
	}

	// Venue usage rate
	var booked, total int
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM bookings
		WHERE date = ?::date AND status NOT IN ('cancelled')
	`, dateStr).Scan(&booked)
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM venue_time_slot_rules
		WHERE is_active = true
	`).Scan(&total)
	// Multiply by active venues
	var venueCount int
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM venues WHERE status = 1`).Scan(&venueCount)
	totalSlots := total * venueCount
	if totalSlots > 0 {
		stats.VenueUsageRate = float64(booked) / float64(totalSlots)
	}

	return stats, nil
}
