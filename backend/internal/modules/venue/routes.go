package venue

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(public, auth, admin *gin.RouterGroup) {
	// Public / authenticated user routes
	public.GET("/venues", h.ListVenues)
	public.GET("/venues/:id", h.GetVenue)
	auth.GET("/venues/:id/availability", h.GetAvailability)

	// Admin routes
	admin.GET("/venues", h.AdminListVenues)
	admin.POST("/venues", h.AdminCreateVenue)
	admin.PUT("/venues/:id", h.AdminUpdateVenue)
	admin.DELETE("/venues/:id", h.AdminDeleteVenue)
	admin.GET("/venues/time-grid", h.AdminGetTimeGrid)
	admin.GET("/venues/:id/time-rules", h.AdminListTimeRules)
	admin.POST("/venues/:id/time-rules", h.AdminCreateTimeRule)
	admin.PUT("/venues/:id/time-rules/:ruleId", h.AdminUpdateTimeRule)
	admin.DELETE("/venues/:id/time-rules/:ruleId", h.AdminDeleteTimeRule)
}
