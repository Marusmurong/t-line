package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authmod "github.com/t-line/backend/internal/modules/auth"
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
}
