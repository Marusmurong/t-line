package activity

import (
	"strconv"

	"github.com/gin-gonic/gin"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/pagination"
	"github.com/t-line/backend/internal/pkg/response"
)

type AdminHandler struct {
	svc *Service
}

func NewAdminHandler(svc *Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) AdminListActivities(c *gin.Context) {
	params := pagination.Parse(c)

	var typePtr *string
	if t := c.Query("type"); t != "" {
		typePtr = &t
	}

	var statusPtr *string
	if s := c.Query("status"); s != "" {
		statusPtr = &s
	}

	activities, total, err := h.svc.AdminListActivities(c.Request.Context(), typePtr, statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OKWithPage(c, activities, total, params.Page, params.PageSize)
}

func (h *AdminHandler) AdminCreateActivity(c *gin.Context) {
	var req CreateActivityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.AdminCreateActivity(c.Request.Context(), req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) AdminUpdateActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的活动ID")
		return
	}

	var req UpdateActivityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.AdminUpdateActivity(c.Request.Context(), id, req)
	if err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AdminHandler) AdminDeleteActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的活动ID")
		return
	}

	if err := h.svc.AdminDeleteActivity(c.Request.Context(), id); err != nil {
		adminHandleError(c, err)
		return
	}

	response.OK(c, nil)
}

func adminHandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
