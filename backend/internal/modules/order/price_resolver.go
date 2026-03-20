package order

import (
	"context"

	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

// ProductInfoProvider provides product name and price (with optional SKU).
type ProductInfoProvider interface {
	GetProductNameAndPrice(ctx context.Context, productID int64, skuID *int64) (string, decimal.Decimal, error)
}

// ActivityInfoProvider provides activity name and price.
type ActivityInfoProvider interface {
	GetActivityNameAndPrice(ctx context.Context, activityID int64) (string, decimal.Decimal, error)
}

// CompositePriceResolver implements PriceResolver by delegating to module-specific providers.
// Note: booking orders are created via CreateOrderDirect and do not use this resolver.
type CompositePriceResolver struct {
	product  ProductInfoProvider
	activity ActivityInfoProvider
}

func NewCompositePriceResolver(
	product ProductInfoProvider,
	activity ActivityInfoProvider,
) *CompositePriceResolver {
	return &CompositePriceResolver{
		product:  product,
		activity: activity,
	}
}

func (r *CompositePriceResolver) ResolvePrice(ctx context.Context, itemType string, itemID int64, skuID *int64) (string, decimal.Decimal, error) {
	switch itemType {
	case TypeProduct:
		return r.product.GetProductNameAndPrice(ctx, itemID, skuID)
	case TypeActivity:
		return r.activity.GetActivityNameAndPrice(ctx, itemID)
	default:
		return "", decimal.Zero, apperrors.New(40001, "不支持的商品类型: "+itemType)
	}
}
