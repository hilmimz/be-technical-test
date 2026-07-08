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

	// Generate ID using sonnyflake ID generator
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

	created, err := u.productRepo.Create(ctx, prod)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errs.Conflict("sku already exists", nil)
		}
		return nil, errs.Internal("failed to create product", err)
	}
	return created, nil
}

func (u *ProductUseCase) GetAllProducts(ctx context.Context) ([]*domain.Product, *errs.Error) {
	products, err := u.productRepo.GetAll(ctx)
	if err != nil {
		return nil, errs.Internal("failed to get all products", err)
	}
	return products, nil
}

func (u *ProductUseCase) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, *errs.Error) {
	product, err := u.productRepo.FindBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("product not found", nil)
		}
		return nil, errs.Internal("failed to get product by sku", err)
	}
	return product, nil
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, sku string) *errs.Error {
	err := u.productRepo.DeleteBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NotFound("product not found", nil)
		}
		return errs.Internal("failed to delete product", err)
	}
	return nil
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, sku string, req *domain.UpdateProductRequest) (*domain.Product, *errs.Error) {
	if req.Qty < 0 {
		return nil, errs.BadRequest("qty cannot be negative", nil)
	}
	if req.Price.LessThanOrEqual(decimal.Zero) {
		return nil, errs.BadRequest("price must be greater than 0", nil)
	}
	product, err := u.productRepo.FindBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("product not found", nil)
		}
		return nil, errs.Internal("failed to get product by sku", err)
	}
	product.Name = req.Name
	product.Qty = req.Qty
	product.Price = req.Price
	err = u.productRepo.Update(ctx, product)
	if err != nil {
		return nil, errs.Internal("failed to update product", err)
	}
	return product, nil
}

func (u *ProductUseCase) PurchaseProduct(ctx context.Context, sku string, req *domain.PurchaseProductRequest) (*domain.Product, *errs.Error) {
	product, err := u.productRepo.FindBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("product not found", nil)
		}
		return nil, errs.Internal("failed to get product", err)
	}

	if product.Qty < req.Qty {
		return nil, errs.BadRequest("insufficient stock", nil)
	}

	product.Qty -= req.Qty

	if err := u.productRepo.DecrementStock(ctx, sku, req.Qty); err != nil {
		return nil, errs.Internal("failed to update stock", err)
	}

	return product, nil
}
