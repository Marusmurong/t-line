package stats

import (
	"time"

	"github.com/shopspring/decimal"
)

// DashboardData aggregates all dashboard metrics for the admin panel.
type DashboardData struct {
	TodayRevenue       decimal.Decimal   `json:"today_revenue"`
	TodayOrders        int               `json:"today_orders"`
	VenueUsageRate     float64           `json:"venue_usage_rate"`
	ActiveUsers        int               `json:"active_users"`
	RevenueTrend       []DayRevenue      `json:"revenue_trend"`
	RevenueComposition []CategoryRevenue `json:"revenue_composition"`
	RecentOrders       []RecentOrder     `json:"recent_orders"`
}

// DayRevenue represents revenue for a single day.
type DayRevenue struct {
	Date   string          `json:"date"`
	Amount decimal.Decimal `json:"amount"`
}

// CategoryRevenue represents revenue breakdown by order type.
type CategoryRevenue struct {
	Category string          `json:"category"`
	Amount   decimal.Decimal `json:"amount"`
	Percent  float64         `json:"percent"`
}

// RecentOrder is a simplified order record for dashboard display.
type RecentOrder struct {
	OrderNo   string          `json:"order_no"`
	User      string          `json:"user"`
	Type      string          `json:"type"`
	Amount    decimal.Decimal `json:"amount"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

// VenueUsageStat represents usage for a specific venue/hour/date combination.
type VenueUsageStat struct {
	VenueID   int64  `json:"venue_id"`
	VenueName string `json:"venue_name"`
	Date      string `json:"date"`
	Hour      int    `json:"hour"`
	IsUsed    bool   `json:"is_used"`
}

// HeatmapCell represents an aggregated heatmap data point (weekday x hour).
type HeatmapCell struct {
	Weekday  int     `json:"weekday"`
	Hour     int     `json:"hour"`
	Count    int     `json:"count"`
	Rate     float64 `json:"rate"`
}

// UserGrowthStat tracks daily new and cumulative user counts.
type UserGrowthStat struct {
	Date     string `json:"date"`
	NewUsers int    `json:"new_users"`
	Total    int    `json:"total"`
}

// MemberDistribution shows the breakdown of users by member level.
type MemberDistribution struct {
	Level   string  `json:"level"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

// CoachPerformance ranks coaches by completed lessons and rating.
type CoachPerformance struct {
	CoachID        int64           `json:"coach_id"`
	CoachTitle     string          `json:"coach_title"`
	LessonCount    int             `json:"lesson_count"`
	StudentCount   int             `json:"student_count"`
	AverageRating  decimal.Decimal `json:"average_rating"`
	TotalRevenue   decimal.Decimal `json:"total_revenue"`
}

// RevenueStatsResp is the response for revenue statistics API.
type RevenueStatsResp struct {
	Trend       []DayRevenue      `json:"trend"`
	Composition []CategoryRevenue `json:"composition"`
	Total       decimal.Decimal   `json:"total"`
}

// VenueUsageStatsResp is the response for venue usage API.
type VenueUsageStatsResp struct {
	Heatmap    []HeatmapCell      `json:"heatmap"`
	UsageRate  float64            `json:"usage_rate"`
	ByVenue    []VenueUsageDetail `json:"by_venue"`
}

// VenueUsageDetail shows usage rate per venue.
type VenueUsageDetail struct {
	VenueID   int64   `json:"venue_id"`
	VenueName string  `json:"venue_name"`
	UsageRate float64 `json:"usage_rate"`
}

// UserStatsResp is the response for user statistics API.
type UserStatsResp struct {
	Growth       []UserGrowthStat     `json:"growth"`
	Distribution []MemberDistribution `json:"distribution"`
	TotalUsers   int                  `json:"total_users"`
}
