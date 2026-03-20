package payment

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/middleware"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/logger"
	"github.com/t-line/backend/internal/pkg/response"
)

// CallbackVerifier verifies wechat pay callback signatures.
type CallbackVerifier interface {
	VerifyCallback(body []byte, signature string) error
}

type Handler struct {
	svc              *Service
	callbackVerifier CallbackVerifier
}

func NewHandler(svc *Service, callbackVerifier CallbackVerifier) *Handler {
	return &Handler{svc: svc, callbackVerifier: callbackVerifier}
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
	// Read raw body for signature verification
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.L.Errorw("wechat callback: failed to read body", "error", err)
		response.ServerError(c, "读取请求体失败")
		return
	}

	// Verify callback signature
	signature := c.GetHeader("Wechatpay-Signature")
	if h.callbackVerifier != nil {
		if err := h.callbackVerifier.VerifyCallback(rawBody, signature); err != nil {
			logger.L.Errorw("wechat callback: signature verification failed", "error", err)
			response.Unauthorized(c, "回调签名验证失败")
			return
		}
	}

	// Parse callback data from raw body (body stream already consumed)
	var req WechatCallbackReq
	if err := json.Unmarshal(rawBody, &req); err != nil {
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
