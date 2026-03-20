package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	academicmod "github.com/t-line/backend/internal/modules/academic"
	activitymod "github.com/t-line/backend/internal/modules/activity"
	authmod "github.com/t-line/backend/internal/modules/auth"
	bookingmod "github.com/t-line/backend/internal/modules/booking"
	notifymod "github.com/t-line/backend/internal/modules/notify"
	ordermod "github.com/t-line/backend/internal/modules/order"
	paymentmod "github.com/t-line/backend/internal/modules/payment"
	productmod "github.com/t-line/backend/internal/modules/product"
	statsmod "github.com/t-line/backend/internal/modules/stats"
	venuemod "github.com/t-line/backend/internal/modules/venue"
	"github.com/t-line/backend/internal/middleware"
	"github.com/t-line/backend/internal/pkg/response"
)

func (s *Server) setupRoutes() {
	s.engine.Use(
		middleware.Recovery(),
		middleware.CORS(),
		middleware.RequestLogger(),
	)

	// health check
	s.engine.GET("/api/v1/health", func(c *gin.Context) {
		response.OK(c, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	v1 := s.engine.Group("/api/v1")

	// public routes (no auth required)
	public := v1.Group("")
	public.Use(middleware.RateLimit(s.rdb, 60, time.Minute))

	// authenticated routes
	authed := v1.Group("")
	authed.Use(middleware.JWTAuth(s.jwtMgr))

	// admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth(s.jwtMgr), middleware.RequireAdmin())

	// register module routes
	s.authHandler.RegisterRoutes(public, authed)
	s.venueHandler.RegisterRoutes(public, authed, admin)
	s.orderHandler.RegisterRoutes(authed)
	s.orderAdminHandler.RegisterRoutes(admin)
	s.paymentHandler.RegisterRoutes(public, authed)
	s.bookingHandler.RegisterRoutes(authed)
	s.productHandler.RegisterRoutes(public, admin)
	s.activityHandler.RegisterRoutes(public, authed)
	s.activityAdminHandler.RegisterRoutes(admin)
	s.academicHandler.RegisterRoutes(public, authed)
	s.academicAdminHandler.RegisterRoutes(admin)
	s.notifyHandler.RegisterRoutes(authed)
	s.statsHandler.RegisterRoutes(admin)

	// handle 404
	s.engine.NoRoute(func(c *gin.Context) {
		response.Fail(c, http.StatusNotFound, 40400, "接口不存在")
	})
}

func (s *Server) initModules() {
	// Auth module
	authRepo := authmod.NewRepository(s.db)
	authSvc := authmod.NewService(authRepo, s.rdb, s.jwtMgr, s.wechatAuth, s.smsSender)
	s.authHandler = authmod.NewHandler(authSvc)

	// Venue module
	venueRepo := venuemod.NewRepository(s.db)
	venueSvc := venuemod.NewService(venueRepo)
	s.venueHandler = venuemod.NewHandler(venueSvc)

	// Order module
	orderRepo := ordermod.NewRepository(s.db)
	orderSvc := ordermod.NewService(orderRepo)
	s.orderHandler = ordermod.NewHandler(orderSvc)
	s.orderAdminHandler = ordermod.NewAdminHandler(orderSvc)

	// Payment module
	paymentRepo := paymentmod.NewRepository(s.db)
	// Note: WalletOperator and WechatPayer would be injected from actual implementations
	// For now, pass nil - they will be set up when the integration modules are ready
	paymentSvc := paymentmod.NewService(paymentRepo, nil, nil, orderSvc)
	s.paymentHandler = paymentmod.NewHandler(paymentSvc)

	// Booking module
	bookingRepo := bookingmod.NewRepository(s.db)
	bookingSvc := bookingmod.NewService(bookingRepo, s.rdb)
	bookingSvc.SetVenuePricer(venueSvc)
	s.bookingHandler = bookingmod.NewHandler(bookingSvc)

	// Product module
	productRepo := productmod.NewRepository(s.db)
	productSvc := productmod.NewService(productRepo)
	s.productHandler = productmod.NewHandler(productSvc)

	// Activity module
	activityRepo := activitymod.NewRepository(s.db)
	activitySvc := activitymod.NewService(activityRepo)
	s.activityHandler = activitymod.NewHandler(activitySvc)
	s.activityAdminHandler = activitymod.NewAdminHandler(activitySvc)

	// Academic module
	academicRepo := academicmod.NewRepository(s.db)
	academicSvc := academicmod.NewService(academicRepo)
	s.academicHandler = academicmod.NewHandler(academicSvc)
	s.academicAdminHandler = academicmod.NewAdminHandler(academicSvc)

	// Notify module
	notifyRepo := notifymod.NewRepository(s.db)
	notifySvc := notifymod.NewService(notifyRepo)
	s.notifyHandler = notifymod.NewHandler(notifySvc)

	// Stats module
	statsRepo := statsmod.NewRepository(s.db)
	statsSvc := statsmod.NewService(statsRepo)
	s.statsHandler = statsmod.NewHandler(statsSvc)

	// Wire cross-module dependencies
	venueSvc.SetBookingChecker(bookingSvc)
}
