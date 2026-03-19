package booking

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/middleware"
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

func (h *Handler) CreateBooking(c *gin.Context) {
	var req CreateBookingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	resp, err := h.svc.CreateBooking(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) ListBookings(c *gin.Context) {
	userID := middleware.GetUserID(c)
	params := pagination.Parse(c)

	upcoming := c.DefaultQuery("type", "upcoming") == "upcoming"

	bookings, total, err := h.svc.ListBookings(c.Request.Context(), userID, upcoming, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, bookings, total, params.Page, params.PageSize)
}

func (h *Handler) GetBooking(c *gin.Context) {
	userID := middleware.GetUserID(c)
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的预订ID")
		return
	}

	resp, err := h.svc.GetBooking(c.Request.Context(), userID, bookingID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	userID := middleware.GetUserID(c)
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的预订ID")
		return
	}

	var req CancelBookingReq
	_ = c.ShouldBindJSON(&req)

	if err := h.svc.CancelBooking(c.Request.Context(), userID, bookingID, req.Reason); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *Handler) JoinWaitlist(c *gin.Context) {
	var req JoinWaitlistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	resp, err := h.svc.JoinWaitlist(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
