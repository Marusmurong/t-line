package product

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Product CRUD

func (r *Repository) Create(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Product, error) {
	var p Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) List(ctx context.Context, category *string, status *int, offset, limit int) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.WithContext(ctx).Model(&Product{})
	if category != nil {
		query = query.Where("category = ?", *category)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort_order ASC, id DESC").Offset(offset).Limit(limit).Find(&products).Error
	return products, total, err
}

func (r *Repository) Update(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&Product{}, id).Error
}

// SKU operations

func (r *Repository) ListSKUs(ctx context.Context, productID int64) ([]ProductSKU, error) {
	var skus []ProductSKU
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("id ASC").
		Find(&skus).Error
	return skus, err
}

func (r *Repository) GetSKUByID(ctx context.Context, id int64) (*ProductSKU, error) {
	var sku ProductSKU
	if err := r.db.WithContext(ctx).First(&sku, id).Error; err != nil {
		return nil, err
	}
	return &sku, nil
}

func (r *Repository) CreateSKU(ctx context.Context, sku *ProductSKU) error {
	return r.db.WithContext(ctx).Create(sku).Error
}

func (r *Repository) DeleteSKU(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&ProductSKU{}, id).Error
}

// Stock

func (r *Repository) UpdateStock(ctx context.Context, id int64, delta int) error {
	return r.db.WithContext(ctx).
		Model(&Product{}).
		Where("id = ? AND stock + ? >= 0", id, delta).
		Update("stock", gorm.Expr("stock + ?", delta)).Error
}

// Transaction helper

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
