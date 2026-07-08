package handler

import (
	"be-technical-test/internal/domain"
	"be-technical-test/pkg/response"
	"be-technical-test/pkg/validation"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(productUsecase domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) CreateProductHandler(c *gin.Context) {
	var req domain.CreateProductRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}
		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, errs := h.productUsecase.CreateProduct(ctx, &req)
	if errs != nil {
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	response.ResponseOK(c, http.StatusCreated, "product created successfully", resp)
}

func (h *ProductHandler) GetAllProductsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	products, errs := h.productUsecase.GetAllProducts(ctx)
	if errs != nil {
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	response.ResponseOK(c, http.StatusOK, "products retrieved successfully", products)
}
