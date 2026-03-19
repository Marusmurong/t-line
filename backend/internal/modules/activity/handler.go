package activity

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

// --- User-facing handlers ---

func (h *Handler) ListActivities(c *gin.Context) {
	params := pagination.Parse(c)

	var typePtr *string
	if t := c.Query("type"); t != "" {
		typePtr = &t
	}

	var statusPtr *string
	if s := c.Query("status"); s != "" {
		statusPtr = &s
	}

	activities, total, err := h.svc.ListActivities(c.Request.Context(), typePtr, statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, activities, total, params.Page, params.PageSize)
}

func (h *Handler) GetActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的活动ID")
		return
	}

	resp, err := h.svc.GetActivity(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) Register(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的活动ID")
		return
	}

	userID := middleware.GetUserID(c)
	resp, err := h.svc.Register(c.Request.Context(), id, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) CancelRegistration(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的活动ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.svc.CancelRegistration(c.Request.Context(), id, userID); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
