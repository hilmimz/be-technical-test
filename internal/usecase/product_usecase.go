package usecase

import (
	"be-technical-test/internal/domain"
	"be-technical-test/pkg/errs"
	"be-technical-test/pkg/idgen"
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	productRepo domain.ProductRepository
}

func NewProductUseCase(productRepository domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepository,
	}
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, *errs.Error) {
	if req.Qty < 0 {
		return nil, errs.BadRequest("qty cannot be negative", nil)
	}
	if req.Price.LessThanOrEqual(decimal.Zero) {
		return nil, errs.BadRequest("price must be greater than 0", nil)
	}

	existing, err := u.productRepo.FindBySKU(ctx, req.SKU)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Internal("failed to check existing sku", err)
	}
	if existing != nil {
		return nil, errs.Conflict("sku already exists", nil)
	}

	// Generate ID using snowflake ID generator
	id, err := idgen.NextID()
	if err != nil {
		return nil, errs.Internal("failed to generate id", err)
	}

	prod := &domain.Product{
		ID:    id,
		SKU:   req.SKU,
		Name:  req.Name,
		Qty:   req.Qty,
		Price: req.Price,
	}

	created, err := u.productRepo.CreateProduct(ctx, prod)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.Conflict("sku already exists", nil)
		}
		return nil, errs.Internal("failed to create product", err)
	}
	return created, nil
}
