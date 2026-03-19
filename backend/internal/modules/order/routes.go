package order

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(auth *gin.RouterGroup) {
	auth.POST("/orders", h.CreateOrder)
	auth.GET("/orders", h.ListOrders)
	auth.GET("/orders/:id", h.GetOrder)
	auth.POST("/orders/:id/cancel", h.CancelOrder)
}

func (h *AdminHandler) RegisterRoutes(admin *gin.RouterGroup) {
	admin.GET("/orders", h.ListOrders)
	admin.POST("/refunds/:id/review", h.ReviewRefund)
}
