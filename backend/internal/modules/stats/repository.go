package stats

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Repository handles all statistics-related database queries using raw SQL.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new stats repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetTodayRevenue returns the sum of pay_amount for paid orders today.
func (r *Repository) GetTodayRevenue(ctx context.Context) (decimal.Decimal, error) {
	var result *decimal.Decimal
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(pay_amount), 0) AS total
		FROM orders
		WHERE paid_at >= CURRENT_DATE
		  AND paid_at < CURRENT_DATE + INTERVAL '1 day'
		  AND status NOT IN ('cancelled', 'refunded')
	`).Scan(&result).Error
	if err != nil {
		return decimal.Zero, err
	}
	if result == nil {
		return decimal.Zero, nil
	}
	return *result, nil
}

// GetTodayOrderCount returns the count of orders created today.
func (r *Repository) GetTodayOrderCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM orders
		WHERE created_at >= CURRENT_DATE
		  AND created_at < CURRENT_DATE + INTERVAL '1 day'
	`).Scan(&count).Error
	return count, err
}

// GetVenueUsageRate calculates the ratio of booked slots to total available slots for a given date.
func (r *Repository) GetVenueUsageRate(ctx context.Context, date string) (float64, error) {
	var result struct {
		Booked int
		Total  int
	}
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			(SELECT COUNT(*) FROM bookings
			 WHERE date = ? AND status NOT IN ('cancelled')) AS booked,
			(SELECT COUNT(*) FROM venue_time_slot_rules
			 WHERE is_active = true) *
			(SELECT COUNT(*) FROM venues WHERE status = 1) AS total
	`, date).Scan(&result).Error
	if err != nil {
		return 0, err
	}
	if result.Total == 0 {
		return 0, nil
	}
	return float64(result.Booked) / float64(result.Total), nil
}

// GetActiveUsers returns the count of distinct users who placed orders in the last N days.
func (r *Repository) GetActiveUsers(ctx context.Context, days int) (int, error) {
	var count int
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT user_id) FROM orders
		WHERE created_at >= CURRENT_DATE - INTERVAL '1 day' * ?
	`, days).Scan(&count).Error
	return count, err
}

// GetRevenueTrend returns daily revenue between start and end dates.
func (r *Repository) GetRevenueTrend(ctx context.Context, startDate, endDate string) ([]DayRevenue, error) {
	var rows []DayRevenue
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(paid_at, 'YYYY-MM-DD') AS date,
			COALESCE(SUM(pay_amount), 0) AS amount
		FROM orders
		WHERE paid_at >= ?::date
		  AND paid_at < ?::date + INTERVAL '1 day'
		  AND status NOT IN ('cancelled', 'refunded')
		GROUP BY TO_CHAR(paid_at, 'YYYY-MM-DD')
		ORDER BY date
	`, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetRevenueComposition returns revenue grouped by order type for a date range.
func (r *Repository) GetRevenueComposition(ctx context.Context, startDate, endDate string) ([]CategoryRevenue, error) {
	var rows []struct {
		Category string
		Amount   decimal.Decimal
	}
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			type AS category,
			COALESCE(SUM(pay_amount), 0) AS amount
		FROM orders
		WHERE paid_at >= ?::date
		  AND paid_at < ?::date + INTERVAL '1 day'
		  AND status NOT IN ('cancelled', 'refunded')
		GROUP BY type
		ORDER BY amount DESC
	`, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	total := decimal.Zero
	for _, row := range rows {
		total = total.Add(row.Amount)
	}

	result := make([]CategoryRevenue, 0, len(rows))
	for _, row := range rows {
		pct := 0.0
		if total.IsPositive() {
			pct, _ = row.Amount.Div(total).Mul(decimal.NewFromInt(100)).Float64()
		}
		result = append(result, CategoryRevenue{
			Category: row.Category,
			Amount:   row.Amount,
			Percent:  pct,
		})
	}
	return result, nil
}

// GetRecentOrders returns the most recent orders with user nickname joined.
func (r *Repository) GetRecentOrders(ctx context.Context, limit int) ([]RecentOrder, error) {
	var rows []RecentOrder
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			o.order_no,
			u.nickname AS user,
			o.type,
			o.pay_amount AS amount,
			o.status,
			o.created_at
		FROM orders o
		LEFT JOIN users u ON o.user_id = u.id
		ORDER BY o.created_at DESC
		LIMIT ?
	`, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetVenueUsageHeatmap aggregates booking counts by weekday and hour for a date range.
func (r *Repository) GetVenueUsageHeatmap(ctx context.Context, startDate, endDate string) ([]HeatmapCell, error) {
	var rows []HeatmapCell
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			EXTRACT(DOW FROM date) AS weekday,
			CAST(SPLIT_PART(start_time, ':', 1) AS INTEGER) AS hour,
			COUNT(*) AS count
		FROM bookings
		WHERE date >= ?::date
		  AND date <= ?::date
		  AND status NOT IN ('cancelled')
		GROUP BY weekday, hour
		ORDER BY weekday, hour
	`, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetVenueUsageByVenue returns per-venue usage statistics for a date range.
func (r *Repository) GetVenueUsageByVenue(ctx context.Context, startDate, endDate string) ([]VenueUsageDetail, error) {
	var rows []VenueUsageDetail
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			v.id AS venue_id,
			v.name AS venue_name,
			CASE
				WHEN total_slots.cnt = 0 THEN 0
				ELSE CAST(COALESCE(booked.cnt, 0) AS FLOAT) / total_slots.cnt
			END AS usage_rate
		FROM venues v
		LEFT JOIN (
			SELECT venue_id, COUNT(*) AS cnt
			FROM bookings
			WHERE date >= ?::date AND date <= ?::date
			  AND status NOT IN ('cancelled')
			GROUP BY venue_id
		) booked ON v.id = booked.venue_id
		CROSS JOIN LATERAL (
			SELECT COUNT(*) * ((?::date - ?::date)::int + 1) AS cnt
			FROM venue_time_slot_rules
			WHERE venue_id = v.id AND is_active = true
		) total_slots
		WHERE v.status = 1
		ORDER BY v.sort_order, v.id
	`, startDate, endDate, endDate, startDate).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetUserGrowth returns daily new user count and cumulative total for a date range.
func (r *Repository) GetUserGrowth(ctx context.Context, startDate, endDate string) ([]UserGrowthStat, error) {
	var rows []UserGrowthStat
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(d.date, 'YYYY-MM-DD') AS date,
			COALESCE(daily.cnt, 0) AS new_users,
			(SELECT COUNT(*) FROM users WHERE created_at < d.date + INTERVAL '1 day') AS total
		FROM generate_series(?::date, ?::date, '1 day'::interval) d(date)
		LEFT JOIN (
			SELECT DATE(created_at) AS dt, COUNT(*) AS cnt
			FROM users
			WHERE created_at >= ?::date AND created_at < ?::date + INTERVAL '1 day'
			GROUP BY DATE(created_at)
		) daily ON d.date = daily.dt
		ORDER BY d.date
	`, startDate, endDate, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetMemberDistribution returns the count of users per member level.
func (r *Repository) GetMemberDistribution(ctx context.Context) ([]MemberDistribution, error) {
	var rows []struct {
		Level int
		Count int
	}
	err := r.db.WithContext(ctx).Raw(`
		SELECT member_level AS level, COUNT(*) AS count
		FROM users
		WHERE status = 1
		GROUP BY member_level
		ORDER BY member_level
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	total := 0
	for _, r := range rows {
		total += r.Count
	}

	levelNames := map[int]string{
		0: "普通用户",
		1: "银卡会员",
		2: "金卡会员",
		3: "钻石会员",
	}

	result := make([]MemberDistribution, 0, len(rows))
	for _, row := range rows {
		pct := 0.0
		if total > 0 {
			pct = float64(row.Count) / float64(total) * 100
		}
		name, ok := levelNames[row.Level]
		if !ok {
			name = "未知等级"
		}
		result = append(result, MemberDistribution{
			Level:   name,
			Count:   row.Count,
			Percent: pct,
		})
	}
	return result, nil
}

// GetCoachPerformanceRank returns top coaches ranked by completed lessons.
func (r *Repository) GetCoachPerformanceRank(ctx context.Context, limit int) ([]CoachPerformance, error) {
	var rows []CoachPerformance
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			c.id AS coach_id,
			c.title AS coach_title,
			COALESCE(lessons.cnt, 0) AS lesson_count,
			c.student_count,
			c.rating AS average_rating,
			COALESCE(rev.total, 0) AS total_revenue
		FROM coaches c
		LEFT JOIN (
			SELECT coach_id, COUNT(*) AS cnt
			FROM course_schedules
			WHERE status = 'completed'
			GROUP BY coach_id
		) lessons ON c.id = lessons.coach_id
		LEFT JOIN (
			SELECT cs.coach_id, SUM(o.pay_amount) AS total
			FROM course_schedules cs
			JOIN bookings b ON cs.venue_id = b.venue_id AND cs.date = TO_CHAR(b.date, 'YYYY-MM-DD')
			JOIN orders o ON b.order_id = o.id AND o.status = 'paid'
			GROUP BY cs.coach_id
		) rev ON c.id = rev.coach_id
		WHERE c.status = 1
		ORDER BY lesson_count DESC, average_rating DESC
		LIMIT ?
	`, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetTotalUsers returns the total number of active users.
func (r *Repository) GetTotalUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM users WHERE status = 1
	`).Scan(&count).Error
	return count, err
}

// CacheDailyStats caches daily aggregated stats (placeholder for daily job).
func (r *Repository) GetDailyRevenueTotal(ctx context.Context, date time.Time) (decimal.Decimal, error) {
	var result *decimal.Decimal
	dateStr := date.Format("2006-01-02")
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(pay_amount), 0)
		FROM orders
		WHERE paid_at >= ?::date
		  AND paid_at < ?::date + INTERVAL '1 day'
		  AND status NOT IN ('cancelled', 'refunded')
	`, dateStr, dateStr).Scan(&result).Error
	if err != nil {
		return decimal.Zero, err
	}
	if result == nil {
		return decimal.Zero, nil
	}
	return *result, nil
}
