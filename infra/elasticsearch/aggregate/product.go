package aggregate

import "product_search/pkg/xtypes"

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	SKU         string      `json:"sku"`
	Description string      `json:"description"`
	Price       int64       `json:"price"`
	Category    []Category  `json:"category"`
	Tags        []string    `json:"tags"`
	Images      []string    `json:"images"`
	Stock       int64       `json:"stock"`
	Brand       Brand       `json:"brand"`
	Attributes  xtypes.JSON `json:"attributes"`
	RatingAvg   float64     `json:"rating_avg"`
	RatingCount int64       `json:"rating_count"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	DeletedAt   string      `json:"deleted_at"`
}
