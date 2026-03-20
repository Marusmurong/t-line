package academic

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/middleware"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/pagination"
	"github.com/t-line/backend/internal/pkg/response"
)

// ============================================================
// User-facing handler
// ============================================================

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListCoaches(c *gin.Context) {
	params := pagination.Parse(c)

	var statusPtr *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			statusPtr = &v
		}
	}

	coaches, total, err := h.svc.ListCoaches(c.Request.Context(), statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, coaches, total, params.Page, params.PageSize)
}

func (h *Handler) GetCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的教练ID")
		return
	}

	resp, err := h.svc.GetCoach(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) ListMyRecords(c *gin.Context) {
	params := pagination.Parse(c)
	userID := middleware.GetUserID(c)

	records, total, err := h.svc.ListMyRecords(c.Request.Context(), userID, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, records, total, params.Page, params.PageSize)
}

func (h *Handler) AddRating(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的记录ID")
		return
	}

	var req RatingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	resp, err := h.svc.AddRating(c.Request.Context(), id, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

// ============================================================
// Admin handler
// ============================================================

type AdminHandler struct {
	svc *Service
}

func NewAdminHandler(svc *Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// --- Coach admin ---

func (h *AdminHandler) ListCoaches(c *gin.Context) {
	params := pagination.Parse(c)

	var statusPtr *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			statusPtr = &v
		}
	}

	coaches, total, err := h.svc.ListCoaches(c.Request.Context(), statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OKWithPage(c, coaches, total, params.Page, params.PageSize)
}

func (h *AdminHandler) GetCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的教练ID")
		return
	}

	resp, err := h.svc.GetCoach(c.Request.Context(), id)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) CreateCoach(c *gin.Context) {
	var req CreateCoachReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.CreateCoach(c.Request.Context(), req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) UpdateCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的教练ID")
		return
	}

	var req UpdateCoachReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.UpdateCoach(c.Request.Context(), id, req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) DeleteCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的教练ID")
		return
	}

	if err := h.svc.DeleteCoach(c.Request.Context(), id); err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *AdminHandler) GetCoachPerformance(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的教练ID")
		return
	}

	resp, err := h.svc.GetCoachPerformance(c.Request.Context(), id)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

// --- Schedule admin ---

func (h *AdminHandler) ListSchedules(c *gin.Context) {
	params := pagination.Parse(c)

	var coachIDPtr *int64
	if v := c.Query("coach_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			coachIDPtr = &id
		}
	}

	var venueIDPtr *int64
	if v := c.Query("venue_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			venueIDPtr = &id
		}
	}

	var datePtr *string
	if d := c.Query("date"); d != "" {
		datePtr = &d
	}

	var statusPtr *string
	if s := c.Query("status"); s != "" {
		statusPtr = &s
	}

	schedules, total, err := h.svc.ListSchedules(c.Request.Context(), coachIDPtr, venueIDPtr, datePtr, statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OKWithPage(c, schedules, total, params.Page, params.PageSize)
}

func (h *AdminHandler) GetSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的排课ID")
		return
	}

	resp, err := h.svc.GetSchedule(c.Request.Context(), id)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) CreateSchedule(c *gin.Context) {
	var req CreateScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的排课ID")
		return
	}

	var req UpdateScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.UpdateSchedule(c.Request.Context(), id, req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) CancelSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的排课ID")
		return
	}

	if err := h.svc.CancelSchedule(c.Request.Context(), id); err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *AdminHandler) SubstituteCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的排课ID")
		return
	}

	var req SubstituteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.SubstituteCoach(c.Request.Context(), id, req.SubstituteCoachID)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) ConflictCheck(c *gin.Context) {
	var req ConflictCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	if err := h.svc.CheckConflict(c.Request.Context(), req.CoachID, req.VenueID, req.Date, req.StartTime, req.EndTime, 0); err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, gin.H{"conflict": false, "message": "无冲突"})
}

func (h *AdminHandler) BatchCreateSchedules(c *gin.Context) {
	var req BatchScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	results, err := h.svc.BatchCreateSchedules(c.Request.Context(), req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, results)
}

// --- Leave admin ---

func (h *AdminHandler) ListLeaves(c *gin.Context) {
	params := pagination.Parse(c)

	var coachIDPtr *int64
	if v := c.Query("coach_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			coachIDPtr = &id
		}
	}

	var statusPtr *string
	if s := c.Query("status"); s != "" {
		statusPtr = &s
	}

	leaves, total, err := h.svc.ListLeaves(c.Request.Context(), coachIDPtr, statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OKWithPage(c, leaves, total, params.Page, params.PageSize)
}

func (h *AdminHandler) GetLeave(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的休假ID")
		return
	}

	resp, err := h.svc.GetLeave(c.Request.Context(), id)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) CreateLeave(c *gin.Context) {
	var req CreateLeaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.CreateLeave(c.Request.Context(), req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) UpdateLeave(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的休假ID")
		return
	}

	var req UpdateLeaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.UpdateLeave(c.Request.Context(), id, req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) DeleteLeave(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的休假ID")
		return
	}

	if err := h.svc.DeleteLeave(c.Request.Context(), id); err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, nil)
}

// --- Student admin ---

func (h *AdminHandler) ListStudents(c *gin.Context) {
	params := pagination.Parse(c)

	items, total, err := h.svc.ListStudents(c.Request.Context(), params.Offset(), params.PageSize)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OKWithPage(c, items, total, params.Page, params.PageSize)
}

func (h *AdminHandler) ListStudentRecords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的学员ID")
		return
	}

	params := pagination.Parse(c)

	records, total, err := h.svc.ListStudentRecords(c.Request.Context(), id, params.Offset(), params.PageSize)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OKWithPage(c, records, total, params.Page, params.PageSize)
}

func (h *AdminHandler) AddFeedback(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的记录ID")
		return
	}

	var req FeedbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.AddFeedback(c.Request.Context(), id, req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

// --- Error handling ---

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}

func adminHandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
