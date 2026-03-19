package product

import (
	"context"

	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// --- User-facing ---

func (s *Service) ListProducts(ctx context.Context, category *string, offset, limit int) ([]ProductResp, int64, error) {
	onShelf := 1
	products, total, err := s.repo.List(ctx, category, &onShelf, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]ProductResp, 0, len(products))
	for i := range products {
		result = append(result, ToProductResp(&products[i]))
	}
	return result, total, nil
}

func (s *Service) GetProduct(ctx context.Context, id int64) (*ProductResp, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrProductNotFound
	}

	skus, err := s.repo.ListSKUs(ctx, id)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToProductDetailResp(p, skus)
	return &resp, nil
}

// --- Admin ---

func (s *Service) AdminListProducts(ctx context.Context, category *string, status *int, offset, limit int) ([]ProductResp, int64, error) {
	products, total, err := s.repo.List(ctx, category, status, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]ProductResp, 0, len(products))
	for i := range products {
		result = append(result, ToProductResp(&products[i]))
	}
	return result, total, nil
}

func (s *Service) AdminCreateProduct(ctx context.Context, req CreateProductReq) (*ProductResp, error) {
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	originalPrice := decimal.Zero
	if req.OriginalPrice != "" {
		originalPrice, err = decimal.NewFromString(req.OriginalPrice)
		if err != nil {
			return nil, apperrors.ErrInvalidParams
		}
	}

	p := &Product{
		Category:      req.Category,
		SubCategory:   req.SubCategory,
		Name:          req.Name,
		Description:   req.Description,
		CoverImage:    req.CoverImage,
		Images:        req.Images,
		Price:         price,
		OriginalPrice: originalPrice,
		Stock:         req.Stock,
		Attributes:    req.Attributes,
		SortOrder:     req.SortOrder,
		Status:        0, // default off-shelf
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToProductResp(p)
	return &resp, nil
}

func (s *Service) AdminUpdateProduct(ctx context.Context, id int64, req UpdateProductReq) (*ProductResp, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrProductNotFound
	}

	if req.Category != nil {
		p.Category = *req.Category
	}
	if req.SubCategory != nil {
		p.SubCategory = *req.SubCategory
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.CoverImage != nil {
		p.CoverImage = *req.CoverImage
	}
	if req.Images != nil {
		p.Images = *req.Images
	}
	if req.Price != nil {
		price, pErr := decimal.NewFromString(*req.Price)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		p.Price = price
	}
	if req.OriginalPrice != nil {
		origPrice, pErr := decimal.NewFromString(*req.OriginalPrice)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		p.OriginalPrice = origPrice
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	if req.Attributes != nil {
		p.Attributes = *req.Attributes
	}
	if req.SortOrder != nil {
		p.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		p.Status = *req.Status
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToProductResp(p)
	return &resp, nil
}

func (s *Service) AdminDeleteProduct(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperrors.ErrProductNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) AdminUpdateStock(ctx context.Context, id int64, delta int) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperrors.ErrProductNotFound
	}

	if err := s.repo.UpdateStock(ctx, id, delta); err != nil {
		return apperrors.ErrStockInsufficient
	}
	return nil
}
