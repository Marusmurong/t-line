package venue

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/pagination"
	"github.com/t-line/backend/internal/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// --- User-facing handlers ---

func (h *Handler) ListVenues(c *gin.Context) {
	params := pagination.Parse(c)

	venues, total, err := h.svc.ListVenues(c.Request.Context(), params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, venues, total, params.Page, params.PageSize)
}

func (h *Handler) GetVenue(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的场地ID")
		return
	}

	resp, err := h.svc.GetVenue(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) GetAvailability(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的场地ID")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "日期格式错误")
		return
	}

	slots, err := h.svc.GetAvailability(c.Request.Context(), id, date)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, slots)
}

// --- Admin handlers ---

func (h *Handler) AdminListVenues(c *gin.Context) {
	params := pagination.Parse(c)

	var statusPtr *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, err := strconv.Atoi(statusStr)
		if err == nil {
			statusPtr = &s
		}
	}

	venues, total, err := h.svc.AdminListVenues(c.Request.Context(), statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, venues, total, params.Page, params.PageSize)
}

func (h *Handler) AdminCreateVenue(c *gin.Context) {
	var req CreateVenueReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.AdminCreateVenue(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) AdminUpdateVenue(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的场地ID")
		return
	}

	var req UpdateVenueReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.AdminUpdateVenue(c.Request.Context(), id, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) AdminDeleteVenue(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的场地ID")
		return
	}

	if err := h.svc.AdminDeleteVenue(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

// --- Admin time rule handlers ---

func (h *Handler) AdminListTimeRules(c *gin.Context) {
	venueID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的场地ID")
		return
	}

	rules, err := h.svc.AdminListTimeRules(c.Request.Context(), venueID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, rules)
}

func (h *Handler) AdminCreateTimeRule(c *gin.Context) {
	venueID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的场地ID")
		return
	}

	var req CreateTimeRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	rule, err := h.svc.AdminCreateTimeRule(c.Request.Context(), venueID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, rule)
}

func (h *Handler) AdminUpdateTimeRule(c *gin.Context) {
	ruleID, err := strconv.ParseInt(c.Param("ruleId"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的规则ID")
		return
	}

	var req UpdateTimeRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	rule, err := h.svc.AdminUpdateTimeRule(c.Request.Context(), ruleID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, rule)
}

func (h *Handler) AdminDeleteTimeRule(c *gin.Context) {
	ruleID, err := strconv.ParseInt(c.Param("ruleId"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的规则ID")
		return
	}

	if err := h.svc.AdminDeleteTimeRule(c.Request.Context(), ruleID); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

// --- Admin time grid handler ---

func (h *Handler) AdminGetTimeGrid(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "日期格式错误")
		return
	}

	// Get all active venues
	active := 1
	venues, _, err := h.svc.repo.List(c.Request.Context(), &active, 0, 100)
	if err != nil {
		response.ServerError(c, "获取场地列表失败")
		return
	}

	var grid []TimeGridSlot
	for _, v := range venues {
		slots, slotErr := h.svc.GetAvailability(c.Request.Context(), v.ID, date)
		if slotErr != nil {
			continue
		}
		for _, slot := range slots {
			status := "available"
			if !slot.Available {
				status = "booked"
			}
			grid = append(grid, TimeGridSlot{
				VenueID:   v.ID,
				VenueName: v.Name,
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
				Status:    status,
			})
		}
	}

	response.OK(c, grid)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
