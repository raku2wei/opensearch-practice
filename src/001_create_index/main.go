package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	opensearch "github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

const IndexName = "go-test-index1"

func main() {
	// Initialize the client with SSL/TLS enabled.
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For testing only. Use certificate for validation.
		},
		Addresses: []string{"https://host.docker.internal:9200"},
		Username:  "admin", // For testing only. Don't store credentials in code.
		Password:  "admin",
	})
	if err != nil {
		fmt.Println("cannot initialize", err)
		os.Exit(1)
	}

	// Define index mapping.
	mapping := strings.NewReader(`{
		"settings": {
			"index": {
				"number_of_shards": 3,
				"number_of_replicas": 1
			}
		},
		"mappings": {
			"properties": {
				"tweet_text":{
					"type": "text",
					"analyzer": "kuromoji"
				},
				"user_name":{
					"type": "keyword"
				}
			}
		}
	}`)

	// Create an index with non-default settings.
	createIndex := opensearchapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  mapping,
	}
	createIndexResponse, err := createIndex.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to create index ", err)
		os.Exit(1)
	}
	fmt.Println("creating index", createIndexResponse)
}
