package repository

import (
	"be-technical-test/internal/domain"
	"context"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, req *domain.Product) (*domain.Product, error) {
	err := r.db.WithContext(ctx).Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}
