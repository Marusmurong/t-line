package payment

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(public, auth *gin.RouterGroup) {
	// Wechat callback (no auth required)
	public.POST("/payments/wechat-callback", h.WechatCallback)

	// Authenticated routes
	auth.POST("/payments/prepare", h.PreparePayment)
	auth.GET("/coupons", h.ListCoupons)
	auth.GET("/coupons/available", h.ListAvailableCoupons)
}
