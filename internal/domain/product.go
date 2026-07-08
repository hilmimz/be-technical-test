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
	Qty   int             `json:"qty" binding:"required,min=0"`
	Price decimal.Decimal `json:"price" binding:"required"`
}

type ProductRepository interface {
	FindBySKU(ctx context.Context, sku string) (*Product, error)
	CreateProduct(ctx context.Context, req *Product) (*Product, error)
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, *errs.Error)
}
