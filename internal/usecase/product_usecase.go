package usecase

import (
	"be-technical-test/internal/domain"
	"be-technical-test/pkg/errs"
	"be-technical-test/pkg/idgen"
	"context"
	"errors"
	"log"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	productRepo domain.ProductRepository
	stockRepo   domain.StockRepository
}

func NewProductUseCase(productRepository domain.ProductRepository, stockRepo domain.StockRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepository,
		stockRepo:   stockRepo,
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

	// Sync stock to Redis
	if err := u.stockRepo.SetStock(ctx, prod.SKU, prod.Qty); err != nil {
		log.Printf("warning: failed to sync stock to redis for sku=%s: %v", prod.SKU, err)
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

	// Sync stock to Redis
	if err := u.stockRepo.SetStock(ctx, product.SKU, product.Qty); err != nil {
		log.Printf("warning: failed to sync stock to redis for sku=%s: %v", product.SKU, err)
	}

	return product, nil
}

func (u *ProductUseCase) PurchaseProduct(ctx context.Context, sku string, req *domain.PurchaseProductRequest) *errs.Error {
	result, err := u.stockRepo.DecrementStock(ctx, sku, req.Qty)
	if err != nil {
		return errs.Internal("failed to check stock", err)
	}

	switch result {
	case 1:
		if _, err := u.productRepo.DecrementStock(ctx, sku, req.Qty); err != nil {
			return errs.Internal("failed to update stock", err)
		}
		return nil
	case 0:
		return errs.BadRequest("insufficient stock", nil)
	default:
		// Fallback to database
		rowsAffected, err := u.productRepo.DecrementStock(ctx, sku, req.Qty)
		if err != nil {
			if rollbackErr := u.stockRepo.IncrementStock(ctx, sku, req.Qty); rollbackErr != nil {
				log.Printf("CRITICAL: failed to rollback redis stock after db update failure sku=%s qty=%d: %v", sku, req.Qty, rollbackErr)
			}
			return errs.Internal("failed to update stock", err)
		}

		if rowsAffected == 0 {
			_, err := u.productRepo.FindBySKU(ctx, sku)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errs.NotFound("product not found", nil)
				}
				return errs.Internal("failed to get product info", err)
			}
			return errs.BadRequest("insufficient stock", nil)
		}

		latestProduct, err := u.productRepo.FindBySKU(ctx, sku)
		if err == nil {
			// Sync stock to Redis
			if err := u.stockRepo.SetStock(ctx, sku, latestProduct.Qty); err != nil {
				log.Printf("warning: failed to sync stock to redis for sku=%s: %v", sku, err)
			}
		}

		return nil
	}
}
