package order

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock implementations ---

type mockProductInfoProvider struct {
	nameByID  map[int64]string
	priceByID map[int64]decimal.Decimal
	// skuID -> price override
	skuPrices map[int64]decimal.Decimal
	errByID   map[int64]error
}

func (m *mockProductInfoProvider) GetProductNameAndPrice(ctx context.Context, productID int64, skuID *int64) (string, decimal.Decimal, error) {
	if err, ok := m.errByID[productID]; ok {
		return "", decimal.Zero, err
	}
	name, ok := m.nameByID[productID]
	if !ok {
		return "", decimal.Zero, errors.New("product not found")
	}
	price := m.priceByID[productID]
	if skuID != nil {
		if sp, exists := m.skuPrices[*skuID]; exists {
			price = sp
		}
	}
	return name, price, nil
}

type mockActivityInfoProvider struct {
	nameByID  map[int64]string
	priceByID map[int64]decimal.Decimal
	errByID   map[int64]error
}

func (m *mockActivityInfoProvider) GetActivityNameAndPrice(ctx context.Context, activityID int64) (string, decimal.Decimal, error) {
	if err, ok := m.errByID[activityID]; ok {
		return "", decimal.Zero, err
	}
	name, ok := m.nameByID[activityID]
	if !ok {
		return "", decimal.Zero, errors.New("activity not found")
	}
	return name, m.priceByID[activityID], nil
}

// --- Tests ---

func TestResolvePrice_Product(t *testing.T) {
	productMock := &mockProductInfoProvider{
		nameByID:  map[int64]string{1: "Tennis Ball"},
		priceByID: map[int64]decimal.Decimal{1: decimal.NewFromFloat(29.90)},
	}
	resolver := NewCompositePriceResolver(productMock, nil)

	name, price, err := resolver.ResolvePrice(context.Background(), TypeProduct, 1, nil)
	require.NoError(t, err)
	assert.Equal(t, "Tennis Ball", name)
	assert.True(t, decimal.NewFromFloat(29.90).Equal(price), "price should be 29.90, got %s", price)
}

func TestResolvePrice_ProductWithSKU(t *testing.T) {
	skuID := int64(100)
	productMock := &mockProductInfoProvider{
		nameByID:  map[int64]string{1: "Tennis Racket"},
		priceByID: map[int64]decimal.Decimal{1: decimal.NewFromFloat(599.00)},
		skuPrices: map[int64]decimal.Decimal{100: decimal.NewFromFloat(499.00)},
	}
	resolver := NewCompositePriceResolver(productMock, nil)

	name, price, err := resolver.ResolvePrice(context.Background(), TypeProduct, 1, &skuID)
	require.NoError(t, err)
	assert.Equal(t, "Tennis Racket", name)
	assert.True(t, decimal.NewFromFloat(499.00).Equal(price),
		"SKU price should override product price; got %s", price)
}

func TestResolvePrice_Activity(t *testing.T) {
	activityMock := &mockActivityInfoProvider{
		nameByID:  map[int64]string{10: "Weekend Tournament"},
		priceByID: map[int64]decimal.Decimal{10: decimal.NewFromFloat(188.00)},
	}
	resolver := NewCompositePriceResolver(nil, activityMock)

	name, price, err := resolver.ResolvePrice(context.Background(), TypeActivity, 10, nil)
	require.NoError(t, err)
	assert.Equal(t, "Weekend Tournament", name)
	assert.True(t, decimal.NewFromFloat(188.00).Equal(price))
}

func TestResolvePrice_UnknownType(t *testing.T) {
	resolver := NewCompositePriceResolver(nil, nil)

	_, _, err := resolver.ResolvePrice(context.Background(), "unknown_type", 1, nil)
	assert.Error(t, err, "unknown item type should return error")
	assert.Contains(t, err.Error(), "不支持的商品类型")
}

func TestResolvePrice_ProductNotFound(t *testing.T) {
	productMock := &mockProductInfoProvider{
		nameByID:  map[int64]string{},
		priceByID: map[int64]decimal.Decimal{},
		errByID:   map[int64]error{999: errors.New("product not found")},
	}
	resolver := NewCompositePriceResolver(productMock, nil)

	_, _, err := resolver.ResolvePrice(context.Background(), TypeProduct, 999, nil)
	assert.Error(t, err, "non-existent product should return error")
}

func TestResolvePrice_ActivityNotFound(t *testing.T) {
	activityMock := &mockActivityInfoProvider{
		nameByID:  map[int64]string{},
		priceByID: map[int64]decimal.Decimal{},
		errByID:   map[int64]error{999: errors.New("activity not found")},
	}
	resolver := NewCompositePriceResolver(nil, activityMock)

	_, _, err := resolver.ResolvePrice(context.Background(), TypeActivity, 999, nil)
	assert.Error(t, err, "non-existent activity should return error")
}
