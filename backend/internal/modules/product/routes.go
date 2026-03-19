package product

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(public, admin *gin.RouterGroup) {
	// Public user routes
	public.GET("/products", h.ListProducts)
	public.GET("/products/:id", h.GetProduct)

	// Admin routes
	admin.GET("/products", h.AdminListProducts)
	admin.POST("/products", h.AdminCreateProduct)
	admin.PUT("/products/:id", h.AdminUpdateProduct)
	admin.DELETE("/products/:id", h.AdminDeleteProduct)
	admin.PUT("/products/:id/stock", h.AdminUpdateStock)
}
