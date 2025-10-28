package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating ES client: %s", err)
	}

	// Test connection
	_, err = client.Info()
	if err != nil {
		log.Fatalf("Error connecting to ES: %s", err)
	}

	ES = client
	log.Println(" Elasticsearch connected")
}
