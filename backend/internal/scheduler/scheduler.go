package scheduler

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/t-line/backend/internal/pkg/logger"
	"gorm.io/gorm"
)

// Scheduler manages all periodic background tasks.
type Scheduler struct {
	db   *gorm.DB
	rdb  *redis.Client
	cron *cron.Cron
}

// New creates a new Scheduler instance.
func New(db *gorm.DB, rdb *redis.Client) *Scheduler {
	return &Scheduler{
		db:   db,
		rdb:  rdb,
		cron: cron.New(cron.WithSeconds()),
	}
}

// Start registers and starts all scheduled tasks.
func (s *Scheduler) Start() {
	// Every minute: expire waitlist notifications that timed out
	if _, err := s.cron.AddFunc("0 * * * * *", s.expireBookingWaitlist); err != nil {
		logger.L.Errorw("scheduler: failed to add booking expire task", "error", err)
	}

	// Every minute: close unpaid orders that have expired
	if _, err := s.cron.AddFunc("0 * * * * *", s.expireUnpaidOrders); err != nil {
		logger.L.Errorw("scheduler: failed to add order expire task", "error", err)
	}

	// Every 5 minutes: check activities for auto-cancel or confirmation
	if _, err := s.cron.AddFunc("0 */5 * * * *", s.checkActivityCancellation); err != nil {
		logger.L.Errorw("scheduler: failed to add activity cancel task", "error", err)
	}

	// Every day at 02:00: aggregate previous day's statistics
	if _, err := s.cron.AddFunc("0 0 2 * * *", s.aggregateDailyStats); err != nil {
		logger.L.Errorw("scheduler: failed to add stats aggregate task", "error", err)
	}

	s.cron.Start()
	logger.L.Info("scheduler: started")
}

// Stop gracefully stops all scheduled tasks.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	logger.L.Info("scheduler: stopped")
}
