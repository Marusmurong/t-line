package stats

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/response"
)

// Handler handles stats-related HTTP requests for admin panel.
type Handler struct {
	svc *Service
}

// NewHandler creates a new stats handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GetDashboard returns aggregated dashboard data.
func (h *Handler) GetDashboard(c *gin.Context) {
	data, err := h.svc.GetDashboard(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, data)
}

// GetRevenueStats returns revenue trend and composition for a date range.
func (h *Handler) GetRevenueStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "start_date 和 end_date 为必填参数")
		return
	}

	data, err := h.svc.GetRevenueStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, data)
}

// GetVenueUsageStats returns venue usage heatmap and per-venue stats.
func (h *Handler) GetVenueUsageStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "start_date 和 end_date 为必填参数")
		return
	}

	data, err := h.svc.GetVenueUsageStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, data)
}

// GetUserStats returns user growth and member distribution.
func (h *Handler) GetUserStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "start_date 和 end_date 为必填参数")
		return
	}

	data, err := h.svc.GetUserStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, data)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
