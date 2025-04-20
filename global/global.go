package global

import (
	"product_search/pkg/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ESClient *elasticsearch.Client
	Config   config.Config
)
