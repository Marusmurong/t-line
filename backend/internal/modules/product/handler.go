package product

import (
	"strconv"

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

func (h *Handler) ListProducts(c *gin.Context) {
	params := pagination.Parse(c)

	var categoryPtr *string
	if cat := c.Query("category"); cat != "" {
		categoryPtr = &cat
	}

	products, total, err := h.svc.ListProducts(c.Request.Context(), categoryPtr, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, products, total, params.Page, params.PageSize)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的商品ID")
		return
	}

	resp, err := h.svc.GetProduct(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

// --- Admin handlers ---

func (h *Handler) AdminListProducts(c *gin.Context) {
	params := pagination.Parse(c)

	var categoryPtr *string
	if cat := c.Query("category"); cat != "" {
		categoryPtr = &cat
	}

	var statusPtr *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, err := strconv.Atoi(statusStr)
		if err == nil {
			statusPtr = &s
		}
	}

	products, total, err := h.svc.AdminListProducts(c.Request.Context(), categoryPtr, statusPtr, params.Offset(), params.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OKWithPage(c, products, total, params.Page, params.PageSize)
}

func (h *Handler) AdminCreateProduct(c *gin.Context) {
	var req CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.AdminCreateProduct(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) AdminUpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的商品ID")
		return
	}

	var req UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	resp, err := h.svc.AdminUpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *Handler) AdminDeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的商品ID")
		return
	}

	if err := h.svc.AdminDeleteProduct(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, nil)
}

func (h *Handler) AdminUpdateStock(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, "无效的商品ID")
		return
	}

	var req UpdateStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, apperrors.ErrInvalidParams.Code, err.Error())
		return
	}

	if err := h.svc.AdminUpdateStock(c.Request.Context(), id, req.Delta); err != nil {
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
