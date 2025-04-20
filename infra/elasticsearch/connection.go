package elasticsearch

import (
	"fmt"
	"product_search/global"

	"github.com/elastic/go-elasticsearch/v8"
)

func InitElasticsearch() {
	config := global.Config.ElasticSearchConfig

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.Address},
		Username:  config.Username,
		Password:  config.Password,
	})

	if err != nil {
		fmt.Printf("Error creating the client: %s", err)
	}
	global.ESClient = client
	fmt.Printf("Elasticsearch client initialized successfully: %s", config.Address)
}
