package booking

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(auth *gin.RouterGroup) {
	auth.POST("/bookings", h.CreateBooking)
	auth.GET("/bookings", h.ListBookings)
	auth.GET("/bookings/:id", h.GetBooking)
	auth.POST("/bookings/:id/cancel", h.CancelBooking)
	auth.POST("/bookings/waitlist", h.JoinWaitlist)
}
