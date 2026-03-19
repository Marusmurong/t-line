package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/t-line/backend/internal/config"
	"github.com/t-line/backend/internal/integration/sms"
	"github.com/t-line/backend/internal/integration/wechat"
	activitymod "github.com/t-line/backend/internal/modules/activity"
	authmod "github.com/t-line/backend/internal/modules/auth"
	bookingmod "github.com/t-line/backend/internal/modules/booking"
	notifymod "github.com/t-line/backend/internal/modules/notify"
	ordermod "github.com/t-line/backend/internal/modules/order"
	paymentmod "github.com/t-line/backend/internal/modules/payment"
	productmod "github.com/t-line/backend/internal/modules/product"
	venuemod "github.com/t-line/backend/internal/modules/venue"
	"github.com/t-line/backend/internal/pkg/jwt"
	"github.com/t-line/backend/internal/pkg/logger"
	"github.com/t-line/backend/internal/pkg/validator"
	"gorm.io/gorm"
)

type Server struct {
	cfg        *config.Config
	db         *gorm.DB
	rdb        *redis.Client
	engine     *gin.Engine
	jwtMgr     *jwt.Manager
	wechatAuth *wechat.AuthClient
	smsSender  *sms.Sender

	// module handlers
	authHandler            *authmod.Handler
	venueHandler           *venuemod.Handler
	orderHandler           *ordermod.Handler
	orderAdminHandler      *ordermod.AdminHandler
	paymentHandler         *paymentmod.Handler
	bookingHandler         *bookingmod.Handler
	productHandler         *productmod.Handler
	activityHandler        *activitymod.Handler
	activityAdminHandler   *activitymod.AdminHandler
	notifyHandler          *notifymod.Handler
}

func New(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *Server {
	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	validator.Init()

	jwtMgr := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessExpireMin, cfg.JWT.RefreshExpireD)
	wechatAuth := wechat.NewAuthClient(cfg.WeChat.AppID, cfg.WeChat.AppSecret)
	smsSender := sms.NewSender(rdb, cfg.SMS.Provider)

	s := &Server{
		cfg:        cfg,
		db:         db,
		rdb:        rdb,
		engine:     gin.New(),
		jwtMgr:     jwtMgr,
		wechatAuth: wechatAuth,
		smsSender:  smsSender,
	}

	s.initModules()
	s.setupRoutes()

	return s
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.cfg.App.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.L.Infof("server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.L.Info("shutting down server...")
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	logger.L.Info("server stopped")
	return nil
}
