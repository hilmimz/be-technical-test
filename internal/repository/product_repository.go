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

func (r *ProductRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, req *domain.Product) (*domain.Product, error) {
	err := r.db.WithContext(ctx).Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *ProductRepository) DeleteBySKU(ctx context.Context, sku string) error {
	tx := r.db.WithContext(ctx).Where("sku = ?", sku).Delete(&domain.Product{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}
