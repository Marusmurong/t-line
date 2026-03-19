package product

import (
	"time"

	"gorm.io/datatypes"
)

// --- Request DTOs ---

type CreateProductReq struct {
	Category      string         `json:"category" binding:"required,oneof=course equipment service"`
	SubCategory   string         `json:"sub_category" binding:"omitempty,max=40"`
	Name          string         `json:"name" binding:"required,max=128"`
	Description   string         `json:"description"`
	CoverImage    string         `json:"cover_image" binding:"omitempty,max=512"`
	Images        datatypes.JSON `json:"images"`
	Price         string         `json:"price" binding:"required"`
	OriginalPrice string         `json:"original_price"`
	Stock         int            `json:"stock"`
	Attributes    datatypes.JSON `json:"attributes"`
	SortOrder     int            `json:"sort_order"`
}

type UpdateProductReq struct {
	Category      *string         `json:"category" binding:"omitempty,oneof=course equipment service"`
	SubCategory   *string         `json:"sub_category" binding:"omitempty,max=40"`
	Name          *string         `json:"name" binding:"omitempty,max=128"`
	Description   *string         `json:"description"`
	CoverImage    *string         `json:"cover_image" binding:"omitempty,max=512"`
	Images        *datatypes.JSON `json:"images"`
	Price         *string         `json:"price"`
	OriginalPrice *string         `json:"original_price"`
	Stock         *int            `json:"stock"`
	Attributes    *datatypes.JSON `json:"attributes"`
	SortOrder     *int            `json:"sort_order"`
	Status        *int            `json:"status" binding:"omitempty,oneof=0 1"`
}

type UpdateStockReq struct {
	Delta int `json:"delta" binding:"required"` // positive=add, negative=reduce
}

type CreateSKUReq struct {
	Name       string         `json:"name" binding:"required,max=128"`
	Price      string         `json:"price" binding:"required"`
	Stock      int            `json:"stock"`
	Attributes datatypes.JSON `json:"attributes"`
}

// --- Response DTOs ---

type ProductResp struct {
	ID            int64          `json:"id"`
	Category      string         `json:"category"`
	SubCategory   string         `json:"sub_category"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	CoverImage    string         `json:"cover_image"`
	Images        datatypes.JSON `json:"images"`
	Price         string         `json:"price"`
	OriginalPrice string         `json:"original_price"`
	Stock         int            `json:"stock"`
	SalesCount    int            `json:"sales_count"`
	Status        int            `json:"status"`
	Attributes    datatypes.JSON `json:"attributes"`
	SortOrder     int            `json:"sort_order"`
	SKUs          []SKUResp      `json:"skus,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type SKUResp struct {
	ID         int64          `json:"id"`
	ProductID  int64          `json:"product_id"`
	Name       string         `json:"name"`
	Price      string         `json:"price"`
	Stock      int            `json:"stock"`
	Attributes datatypes.JSON `json:"attributes"`
}

func ToProductResp(p *Product) ProductResp {
	return ProductResp{
		ID:            p.ID,
		Category:      p.Category,
		SubCategory:   p.SubCategory,
		Name:          p.Name,
		Description:   p.Description,
		CoverImage:    p.CoverImage,
		Images:        p.Images,
		Price:         p.Price.StringFixed(2),
		OriginalPrice: p.OriginalPrice.StringFixed(2),
		Stock:         p.Stock,
		SalesCount:    p.SalesCount,
		Status:        p.Status,
		Attributes:    p.Attributes,
		SortOrder:     p.SortOrder,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func ToProductDetailResp(p *Product, skus []ProductSKU) ProductResp {
	resp := ToProductResp(p)
	skuResps := make([]SKUResp, 0, len(skus))
	for _, s := range skus {
		skuResps = append(skuResps, SKUResp{
			ID:         s.ID,
			ProductID:  s.ProductID,
			Name:       s.Name,
			Price:      s.Price.StringFixed(2),
			Stock:      s.Stock,
			Attributes: s.Attributes,
		})
	}
	resp.SKUs = skuResps
	return resp
}
