package main

import (
	"be-technical-test/config"
	"be-technical-test/internal/database"
	"be-technical-test/internal/handler"
	"be-technical-test/internal/repository"
	"be-technical-test/internal/usecase"
	"be-technical-test/pkg/validation"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	// Init Database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Register Validators
	validation.RegisterValidators()

	// Init Repository
	productRepo := repository.NewProductRepository(db.DB)

	// Init Usecase
	productUseCase := usecase.NewProductUseCase(productRepo)

	// Init Handlers
	productHandler := handler.NewProductHandler(productUseCase)

	// Setup Router
	router := gin.Default()

	// Product
	products := router.Group("/products")
	products.POST("", productHandler.CreateProductHandler)
	products.GET("", productHandler.GetAllProductsHandler)
	products.GET("/:sku", productHandler.GetProductBySKUHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
