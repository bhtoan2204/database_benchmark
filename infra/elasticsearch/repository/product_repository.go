package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"product_search/infra/elasticsearch/aggregate"

	"github.com/elastic/go-elasticsearch/v8"
)

type ProductRepository struct {
	client *elasticsearch.Client
}

func NewProductRepository(client *elasticsearch.Client) *ProductRepository {
	return &ProductRepository{client: client}
}

type SearchProductParams struct {
	Query       string
	CategoryIDs []string
	BrandIDs    []string
	MinPrice    *int64
	MaxPrice    *int64
	Tags        []string
	Page        int
	PageSize    int
	SortBy      string
	SortOrder   string
}

func (r *ProductRepository) Search(ctx context.Context, params SearchProductParams) ([]aggregate.Product, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	// Build search query
	query := map[string]interface{}{
		"from": (params.Page - 1) * params.PageSize,
		"size": params.PageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
	}

	boolQuery := query["query"].(map[string]interface{})["bool"].(map[string]interface{})
	must := boolQuery["must"].([]map[string]interface{})

	// Add text search if query is provided
	if params.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  params.Query,
				"fields": []string{"name^3", "description", "brand.name^2", "category.name^2"},
				"type":   "best_fields",
			},
		})
	}

	// Add filters
	if len(params.CategoryIDs) > 0 {
		must = append(must, map[string]interface{}{
			"terms": map[string]interface{}{
				"category.id": params.CategoryIDs,
			},
		})
	}

	if len(params.BrandIDs) > 0 {
		must = append(must, map[string]interface{}{
			"terms": map[string]interface{}{
				"brand.id": params.BrandIDs,
			},
		})
	}

	if len(params.Tags) > 0 {
		must = append(must, map[string]interface{}{
			"terms": map[string]interface{}{
				"tags": params.Tags,
			},
		})
	}

	// Add price range if provided
	if params.MinPrice != nil || params.MaxPrice != nil {
		priceRange := map[string]interface{}{}
		if params.MinPrice != nil {
			priceRange["gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			priceRange["lte"] = *params.MaxPrice
		}
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}

	// Add sorting
	if params.SortBy != "" {
		order := "desc"
		if params.SortOrder != "" {
			order = strings.ToLower(params.SortOrder)
		}

		query["sort"] = []map[string]interface{}{
			{
				params.SortBy: map[string]interface{}{
					"order": order,
				},
			},
		}
	}

	// Execute search
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("products"),
		r.client.Search.WithBody(strings.NewReader(mustJSON(query))),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source aggregate.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	products := make([]aggregate.Product, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		products[i] = hit.Source
	}

	return products, result.Hits.Total.Value, nil
}

func (r *ProductRepository) BulkIndex(ctx context.Context, products []aggregate.Product) error {
	var buf strings.Builder
	for _, product := range products {
		// Add metadata
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "products",
				"_id":    product.ID,
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}

		// Add document
		if err := json.NewEncoder(&buf).Encode(product); err != nil {
			return err
		}
	}

	res, err := r.client.Bulk(
		strings.NewReader(buf.String()),
		r.client.Bulk.WithContext(ctx),
		r.client.Bulk.WithIndex("products"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch bulk indexing error: %s", res.String())
	}

	return nil
}

func mustJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
