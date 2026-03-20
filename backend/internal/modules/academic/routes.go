package academic

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(public, auth *gin.RouterGroup) {
	// Public routes - coach listing
	public.GET("/coaches", h.ListCoaches)
	public.GET("/coaches/:id", h.GetCoach)

	// Authenticated user routes
	auth.GET("/my-courses/records", h.ListMyRecords)
	auth.POST("/records/:id/rating", h.AddRating)
}

func (h *AdminHandler) RegisterRoutes(admin *gin.RouterGroup) {
	// Coach management
	admin.GET("/coaches", h.ListCoaches)
	admin.GET("/coaches/:id", h.GetCoach)
	admin.POST("/coaches", h.CreateCoach)
	admin.PUT("/coaches/:id", h.UpdateCoach)
	admin.DELETE("/coaches/:id", h.DeleteCoach)
	admin.GET("/coaches/:id/performance", h.GetCoachPerformance)

	// Schedule management
	admin.GET("/schedules", h.ListSchedules)
	admin.GET("/schedules/:id", h.GetSchedule)
	admin.POST("/schedules", h.CreateSchedule)
	admin.PUT("/schedules/:id", h.UpdateSchedule)
	admin.DELETE("/schedules/:id", h.CancelSchedule)
	admin.PUT("/schedules/:id/substitute", h.SubstituteCoach)
	admin.POST("/schedules/conflict-check", h.ConflictCheck)
	admin.POST("/schedules/batch", h.BatchCreateSchedules)

	// Leave management
	admin.GET("/coach-leaves", h.ListLeaves)
	admin.GET("/coach-leaves/:id", h.GetLeave)
	admin.POST("/coach-leaves", h.CreateLeave)
	admin.PUT("/coach-leaves/:id", h.UpdateLeave)
	admin.DELETE("/coach-leaves/:id", h.DeleteLeave)

	// Student management
	admin.GET("/students", h.ListStudents)
	admin.GET("/students/:id/records", h.ListStudentRecords)
	admin.POST("/records/:id/feedback", h.AddFeedback)
}
