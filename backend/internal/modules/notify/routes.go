package notify

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(auth *gin.RouterGroup) {
	auth.GET("/notifications", h.ListNotifications)
	auth.GET("/notifications/unread-count", h.UnreadCount)
	auth.PUT("/notifications/:id/read", h.MarkRead)
	auth.PUT("/notifications/read-all", h.MarkAllRead)
}
