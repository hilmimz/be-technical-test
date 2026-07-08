package domain

type Product struct {
	ID        int64   `json:"id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CreateProductRequest struct {
	SKU   string  `json:"sku" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Qty   int     `json:"qty" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type CreateProductResponse struct {
	ID        int64   `json:"id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type ProductRepository interface {
}

type ProductUsecase interface {
	CreateProduct(req *CreateProductRequest) (*CreateProductResponse, error)
}
