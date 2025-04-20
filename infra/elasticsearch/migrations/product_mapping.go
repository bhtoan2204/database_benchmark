package migrations

import (
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

const ProductIndex = "products"

var ProductMapping = `
{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "custom_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding", "word_delimiter"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": {
        "type": "text",
        "analyzer": "custom_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "sku": { "type": "keyword" },
      "description": {
        "type": "text",
        "analyzer": "custom_analyzer"
      },
      "price": { "type": "long" },
      "category": {
        "properties": {
          "id": { "type": "keyword" },
          "name": { 
            "type": "text",
            "analyzer": "custom_analyzer",
            "fields": {
              "keyword": {
                "type": "keyword"
              }
            }
          }
        }
      },
      "tags": { "type": "keyword" },
      "images": { "type": "keyword" },
      "stock": { "type": "long" },
      "brand": {
        "properties": {
          "id": { "type": "keyword" },
          "name": {
            "type": "text",
            "analyzer": "custom_analyzer",
            "fields": {
              "keyword": {
                "type": "keyword"
              }
            }
          }
        }
      },
      "attributes": { "type": "object" },
      "rating_avg": { "type": "float" },
      "rating_count": { "type": "long" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" },
      "deleted_at": { "type": "date" }
    }
  }
}`

func CreateProductIndex(client *elasticsearch.Client) error {
	res, err := client.Indices.Exists([]string{ProductIndex})
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		// Index already exists
		return nil
	}

	res, err = client.Indices.Create(
		ProductIndex,
		client.Indices.Create.WithBody(strings.NewReader(ProductMapping)),
	)
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}
