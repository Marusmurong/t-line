package order

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/middleware"
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

func (h *AdminHandler) ListOrders(c *gin.Context) {
	params := pagination.Parse(c)
	status := c.Query("status")
	orderType := c.Query("type")

	orders, total, err := h.svc.ListAllOrders(c.Request.Context(), status, orderType, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, orders, total, params.Page, params.PageSize)
}

func (h *AdminHandler) ReviewRefund(c *gin.Context) {
	refundID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的退款ID")
		return
	}

	var req RefundReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	adminID := middleware.GetUserID(c)
	resp, err := h.svc.ReviewRefund(c.Request.Context(), adminID, refundID, req.Action, req.Remark)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}
