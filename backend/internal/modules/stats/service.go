package stats

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/logger"
)

// Service provides stats aggregation logic for the admin dashboard.
type Service struct {
	repo *Repository
}

// NewService creates a new stats service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetDashboard aggregates all dashboard data into a single response.
func (s *Service) GetDashboard(ctx context.Context) (*DashboardData, error) {
	todayRevenue, err := s.repo.GetTodayRevenue(ctx)
	if err != nil {
		logger.L.Errorw("stats: failed to get today revenue", "error", err)
		return nil, apperrors.ErrInternal
	}

	todayOrders, err := s.repo.GetTodayOrderCount(ctx)
	if err != nil {
		logger.L.Errorw("stats: failed to get today order count", "error", err)
		return nil, apperrors.ErrInternal
	}

	todayStr := time.Now().Format("2006-01-02")
	usageRate, err := s.repo.GetVenueUsageRate(ctx, todayStr)
	if err != nil {
		logger.L.Errorw("stats: failed to get venue usage rate", "error", err)
		return nil, apperrors.ErrInternal
	}

	activeUsers, err := s.repo.GetActiveUsers(ctx, 30)
	if err != nil {
		logger.L.Errorw("stats: failed to get active users", "error", err)
		return nil, apperrors.ErrInternal
	}

	// Revenue trend for the last 7 days
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -6).Format("2006-01-02")

	revenueTrend, err := s.repo.GetRevenueTrend(ctx, startDate, endDate)
	if err != nil {
		logger.L.Errorw("stats: failed to get revenue trend", "error", err)
		return nil, apperrors.ErrInternal
	}

	// Revenue composition for the last 30 days
	compStart := time.Now().AddDate(0, 0, -29).Format("2006-01-02")
	composition, err := s.repo.GetRevenueComposition(ctx, compStart, endDate)
	if err != nil {
		logger.L.Errorw("stats: failed to get revenue composition", "error", err)
		return nil, apperrors.ErrInternal
	}

	recentOrders, err := s.repo.GetRecentOrders(ctx, 10)
	if err != nil {
		logger.L.Errorw("stats: failed to get recent orders", "error", err)
		return nil, apperrors.ErrInternal
	}

	return &DashboardData{
		TodayRevenue:       todayRevenue,
		TodayOrders:        todayOrders,
		VenueUsageRate:     usageRate,
		ActiveUsers:        activeUsers,
		RevenueTrend:       revenueTrend,
		RevenueComposition: composition,
		RecentOrders:       recentOrders,
	}, nil
}

// GetRevenueStats returns revenue trend and composition for a date range.
func (s *Service) GetRevenueStats(ctx context.Context, startDate, endDate string) (*RevenueStatsResp, error) {
	trend, err := s.repo.GetRevenueTrend(ctx, startDate, endDate)
	if err != nil {
		logger.L.Errorw("stats: failed to get revenue trend", "error", err)
		return nil, apperrors.ErrInternal
	}

	composition, err := s.repo.GetRevenueComposition(ctx, startDate, endDate)
	if err != nil {
		logger.L.Errorw("stats: failed to get revenue composition", "error", err)
		return nil, apperrors.ErrInternal
	}

	total := decimal.Zero
	for _, d := range trend {
		total = total.Add(d.Amount)
	}

	return &RevenueStatsResp{
		Trend:       trend,
		Composition: composition,
		Total:       total,
	}, nil
}

// GetVenueUsageStats returns venue usage heatmap and per-venue breakdown.
func (s *Service) GetVenueUsageStats(ctx context.Context, startDate, endDate string) (*VenueUsageStatsResp, error) {
	heatmap, err := s.repo.GetVenueUsageHeatmap(ctx, startDate, endDate)
	if err != nil {
		logger.L.Errorw("stats: failed to get venue heatmap", "error", err)
		return nil, apperrors.ErrInternal
	}

	byVenue, err := s.repo.GetVenueUsageByVenue(ctx, startDate, endDate)
	if err != nil {
		logger.L.Errorw("stats: failed to get venue usage by venue", "error", err)
		return nil, apperrors.ErrInternal
	}

	overallRate, err := s.repo.GetVenueUsageRate(ctx, time.Now().Format("2006-01-02"))
	if err != nil {
		logger.L.Errorw("stats: failed to get overall usage rate", "error", err)
		return nil, apperrors.ErrInternal
	}

	return &VenueUsageStatsResp{
		Heatmap:   heatmap,
		UsageRate: overallRate,
		ByVenue:   byVenue,
	}, nil
}

// GetUserStats returns user growth trend and member level distribution.
func (s *Service) GetUserStats(ctx context.Context, startDate, endDate string) (*UserStatsResp, error) {
	growth, err := s.repo.GetUserGrowth(ctx, startDate, endDate)
	if err != nil {
		logger.L.Errorw("stats: failed to get user growth", "error", err)
		return nil, apperrors.ErrInternal
	}

	distribution, err := s.repo.GetMemberDistribution(ctx)
	if err != nil {
		logger.L.Errorw("stats: failed to get member distribution", "error", err)
		return nil, apperrors.ErrInternal
	}

	totalUsers, err := s.repo.GetTotalUsers(ctx)
	if err != nil {
		logger.L.Errorw("stats: failed to get total users", "error", err)
		return nil, apperrors.ErrInternal
	}

	return &UserStatsResp{
		Growth:       growth,
		Distribution: distribution,
		TotalUsers:   totalUsers,
	}, nil
}
