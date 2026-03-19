package activity

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(public, auth *gin.RouterGroup) {
	// Public routes
	public.GET("/activities", h.ListActivities)
	public.GET("/activities/:id", h.GetActivity)

	// Authenticated user routes
	auth.POST("/activities/:id/register", h.Register)
	auth.POST("/activities/:id/cancel-registration", h.CancelRegistration)
}

func (h *AdminHandler) RegisterRoutes(admin *gin.RouterGroup) {
	admin.GET("/activities", h.AdminListActivities)
	admin.POST("/activities", h.AdminCreateActivity)
	admin.PUT("/activities/:id", h.AdminUpdateActivity)
	admin.DELETE("/activities/:id", h.AdminDeleteActivity)
}
