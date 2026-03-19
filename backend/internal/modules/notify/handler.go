package notify

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

func (h *Handler) ListNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)
	params := pagination.Parse(c)

	notifications, total, err := h.svc.ListNotifications(c.Request.Context(), userID, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, notifications, total, params.Page, params.PageSize)
}

func (h *Handler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的通知ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.svc.MarkRead(c.Request.Context(), id, userID); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *Handler) MarkAllRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.svc.MarkAllRead(c.Request.Context(), userID); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *Handler) UnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, err := h.svc.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, gin.H{"count": count})
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
