package order

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

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	resp, err := h.svc.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) ListOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	params := pagination.Parse(c)
	status := c.Query("status")

	orders, total, err := h.svc.ListUserOrders(c.Request.Context(), userID, status, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, orders, total, params.Page, params.PageSize)
}

func (h *Handler) GetOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的订单ID")
		return
	}

	resp, err := h.svc.GetOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的订单ID")
		return
	}

	var req CancelOrderReq
	_ = c.ShouldBindJSON(&req)

	if err := h.svc.CancelOrder(c.Request.Context(), userID, orderID, req.Reason); err != nil {
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
