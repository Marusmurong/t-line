package auth

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(public, auth *gin.RouterGroup) {
	// public routes
	public.POST("/auth/wechat-login", h.WeChatLogin)
	public.POST("/auth/phone-login", h.PhoneLogin)
	public.POST("/auth/password-login", h.PasswordLogin)
	public.POST("/auth/sms-code", h.SendSMSCode)
	public.POST("/auth/refresh", h.RefreshToken)

	// authenticated routes
	auth.GET("/auth/profile", h.GetProfile)
	auth.PUT("/auth/profile", h.UpdateProfile)
	auth.GET("/wallet", h.GetWallet)
	auth.GET("/wallet/transactions", h.GetWalletTransactions)
}
