package domain

import (
	"be-technical-test/pkg/errs"
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID        uint64          `json:"id"`
	SKU       string          `json:"sku"`
	Name      string          `json:"name"`
	Qty       int             `json:"qty"`
	Price     decimal.Decimal `json:"price"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateProductRequest struct {
	SKU   string          `json:"sku" binding:"required"`
	Name  string          `json:"name" binding:"required"`
	Qty   int             `json:"qty" binding:"min=0"`
	Price decimal.Decimal `json:"price" binding:"required"`
}

type UpdateProductRequest struct {
	Name  string          `json:"name" binding:"required"`
	Qty   int             `json:"qty" binding:"min=0"`
	Price decimal.Decimal `json:"price" binding:"required"`
}

type PurchaseProductRequest struct {
	Qty int `json:"qty" binding:"required,min=1"`
}

type ProductRepository interface {
	FindBySKU(ctx context.Context, sku string) (*Product, error)
	Create(ctx context.Context, req *Product) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	DeleteBySKU(ctx context.Context, sku string) error
	Update(ctx context.Context, product *Product) error
	DecrementStockNaive(ctx context.Context, sku string, qty int) error
	DecrementStock(ctx context.Context, sku string, qty int) (int64, error)
}

type StockRepository interface {
	SetStock(ctx context.Context, sku string, qty int) error
	DecrementStock(ctx context.Context, sku string, qty int) (int, error)
	IncrementStock(ctx context.Context, sku string, qty int) error
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, *errs.Error)
	GetAllProducts(ctx context.Context) ([]*Product, *errs.Error)
	GetProductBySKU(ctx context.Context, sku string) (*Product, *errs.Error)
	DeleteProduct(ctx context.Context, sku string) *errs.Error
	UpdateProduct(ctx context.Context, sku string, req *UpdateProductRequest) (*Product, *errs.Error)
	PurchaseProduct(ctx context.Context, sku string, req *PurchaseProductRequest) *errs.Error
}
