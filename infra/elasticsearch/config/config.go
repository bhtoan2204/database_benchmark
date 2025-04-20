package config

import (
	"crypto/tls"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearchClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return elasticsearch.NewClient(cfg)
}
