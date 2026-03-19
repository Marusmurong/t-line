package product

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Product struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	Category      string          `gorm:"size:20;not null" json:"category"`
	SubCategory   string          `gorm:"size:40" json:"sub_category"`
	Name          string          `gorm:"size:128;not null" json:"name"`
	Description   string          `gorm:"type:text" json:"description"`
	CoverImage    string          `gorm:"size:512" json:"cover_image"`
	Images        datatypes.JSON  `gorm:"type:jsonb;default:'[]'" json:"images"`
	Price         decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	OriginalPrice decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"original_price"`
	Stock         int             `gorm:"default:0" json:"stock"`
	SalesCount    int             `gorm:"default:0" json:"sales_count"`
	Status        int             `gorm:"default:0" json:"status"` // 0=off-shelf, 1=on-shelf
	Attributes    datatypes.JSON  `gorm:"type:jsonb;default:'{}'" json:"attributes"`
	SortOrder     int             `gorm:"default:0" json:"sort_order"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (Product) TableName() string { return "products" }

type ProductSKU struct {
	ID         int64           `gorm:"primaryKey" json:"id"`
	ProductID  int64           `json:"product_id"`
	Name       string          `gorm:"size:128" json:"name"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Stock      int             `gorm:"default:0" json:"stock"`
	Attributes datatypes.JSON  `gorm:"type:jsonb;default:'{}'" json:"attributes"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func (ProductSKU) TableName() string { return "product_skus" }
