package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"product_search/infra/elasticsearch/aggregate"
	"product_search/infra/elasticsearch/config"
	"product_search/infra/elasticsearch/repository"

	"github.com/gin-gonic/gin"
)

type ProductAPI struct {
	productRepo *repository.ProductRepository
}

func NewProductAPI(productRepo *repository.ProductRepository) *ProductAPI {
	return &ProductAPI{
		productRepo: productRepo,
	}
}

func main() {
	// Initialize Elasticsearch client
	client, err := config.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}

	// Initialize repository
	productRepo := repository.NewProductRepository(client)

	// Initialize API
	api := NewProductAPI(productRepo)

	// Setup Gin router
	r := gin.Default()
	r.Use(gin.Recovery())

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	v1 := r.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("/", api.CreateProduct)
			products.POST("/bulk", api.BulkCreateProducts)
			products.GET("/search", api.SearchProducts)
		}
	}

	// Start server
	log.Fatal(r.Run(":8080"))
}

func (api *ProductAPI) CreateProduct(c *gin.Context) {
	var product aggregate.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set timestamps
	now := time.Now().UTC().Format(time.RFC3339)
	product.CreatedAt = now
	product.UpdatedAt = now

	// Index single product
	ctx := context.Background()
	err := api.productRepo.BulkIndex(ctx, []aggregate.Product{product})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (api *ProductAPI) BulkCreateProducts(c *gin.Context) {
	var products []aggregate.Product
	if err := c.ShouldBindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set timestamps for all products
	now := time.Now().UTC().Format(time.RFC3339)
	for i := range products {
		products[i].CreatedAt = now
		products[i].UpdatedAt = now
	}

	// Bulk index products
	ctx := context.Background()
	err := api.productRepo.BulkIndex(ctx, products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Products created successfully",
		"count":   len(products),
	})
}

func (api *ProductAPI) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	page := 1
	pageSize := 10

	params := repository.SearchProductParams{
		Query:    query,
		Page:     page,
		PageSize: pageSize,
	}

	ctx := context.Background()
	products, total, err := api.productRepo.Search(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"products": products,
	})
}
