package auth

import (
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

func (h *Handler) WeChatLogin(c *gin.Context) {
	var req WeChatLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.WeChatLogin(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) PhoneLogin(c *gin.Context) {
	var req PhoneLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.PhoneLogin(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) PasswordLogin(c *gin.Context) {
	var req PasswordLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.PasswordLogin(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) SendSMSCode(c *gin.Context) {
	var req SendSMSReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	if err := h.svc.SendSMSCode(c.Request.Context(), req.Phone); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	tokenPair, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, tokenPair)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	resp, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) GetWallet(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.svc.GetWallet(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) GetWalletTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	params := pagination.Parse(c)

	wallet, err := h.svc.repo.GetWalletByUserID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	txs, total, err := h.svc.repo.ListWalletTransactions(c.Request.Context(), wallet.ID, params.Offset(), params.PageSize)
	if err != nil {
		response.ServerError(c, "获取交易记录失败")
		return
	}

	response.OKWithPage(c, txs, total, params.Page, params.PageSize)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.BadRequest(c, appErr.Code, appErr.Message)
		return
	}
	response.ServerError(c, "服务器内部错误")
}
