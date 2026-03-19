package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/middleware"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) PreparePayment(c *gin.Context) {
	var req PreparePayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	// Get user's wechat openID from context (stored during auth)
	openID, _ := c.Get("wechat_openid")
	openIDStr, _ := openID.(string)

	resp, err := h.svc.PreparePayment(c.Request.Context(), userID, req, openIDStr)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) WechatCallback(c *gin.Context) {
	var req WechatCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	if err := h.svc.HandleWechatCallback(c.Request.Context(), req); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *Handler) ListCoupons(c *gin.Context) {
	userID := middleware.GetUserID(c)

	coupons, err := h.svc.ListUserCoupons(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, coupons)
}

func (h *Handler) ListAvailableCoupons(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderType := c.Query("order_type")
	amount := c.Query("amount")

	if orderType == "" || amount == "" {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "order_type 和 amount 参数必填")
		return
	}

	coupons, err := h.svc.ListAvailableCoupons(c.Request.Context(), userID, orderType, amount)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, coupons)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
