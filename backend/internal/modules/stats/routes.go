package stats

import "github.com/gin-gonic/gin"

// RegisterRoutes registers admin stats routes.
func (h *Handler) RegisterRoutes(admin *gin.RouterGroup) {
	admin.GET("/dashboard", h.GetDashboard)

	statsGroup := admin.Group("/stats")
	statsGroup.GET("/revenue", h.GetRevenueStats)
	statsGroup.GET("/venue-usage", h.GetVenueUsageStats)
	statsGroup.GET("/users", h.GetUserStats)
}
