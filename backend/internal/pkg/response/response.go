package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedData struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func OKWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PaginatedData{
			List: list,
			Pagination: Pagination{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}

func Fail(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func BadRequest(c *gin.Context, code int, message string) {
	Fail(c, http.StatusBadRequest, code, message)
}

func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, 40100, message)
}

func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, 40101, message)
}

func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, 40400, message)
}

func ServerError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, 50000, message)
}
