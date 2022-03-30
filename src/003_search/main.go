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

	// Add documents to the index.
	documents := strings.NewReader(`
		{ "index":{ "_id" : "1" } }
		{"tweet_text":"カレーが好き", "user_name":"三井"}
		{ "index":{ "_id" : "2" } }
		{"tweet_text":"カツカレーが食べたい", "user_name":"三井"}
		{ "index":{ "_id" : "3" } }
		{"tweet_text":"カツ丼が好き", "user_name":"匿名"}
		{ "index":{ "_id" : "4" } }
		{"tweet_text":"秒速五㌢㍍で落ちているサーバー", "user_name":"匿名"}
		{ "index":{ "_id" : "5" } }
		{"tweet_text":"東京都は日本の首都である。", "user_name":"匿名"}
		{ "index":{ "_id" : "6" } }
		{"tweet_text":"京都は昔は日本の首都だった。", "user_name":"匿名"}`)

	req := opensearchapi.BulkRequest{
		Index: IndexName,
		Body:  documents,
	}
	bulkInsertResponse, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to insert document ", err)
		os.Exit(1)
	}
	fmt.Println("### bulk insert document")
	fmt.Println(bulkInsertResponse)

	// Search for the document.
	content := strings.NewReader(`{
	   "size": 5,
	   "query": {
			"match": {
				"tweet_text": {
			  		"query": "カレー",
			  		"operator": "and"
				}
		  	}
	  	}
	}`)

	search := opensearchapi.SearchRequest{
		Body: content,
	}

	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to search document ", err)
		os.Exit(1)
	}
	fmt.Println("### search document")
	fmt.Println(searchResponse)

	// Delete documents.
	deleteDocuments := strings.NewReader(`
		{ "delete":{ "_id" : "1" } }
		{ "delete":{ "_id" : "2" } }
		{ "delete":{ "_id" : "3" } }
		{ "delete":{ "_id" : "4" } }
		{ "delete":{ "_id" : "5" } }
		{ "delete":{ "_id" : "6" } }`)

	deleteReq := opensearchapi.BulkRequest{
		Index: IndexName,
		Body:  deleteDocuments,
	}
	bulkDeleteResponse, err := deleteReq.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to delete document ", err)
		os.Exit(1)
	}
	fmt.Println("### bulk delete document")
	fmt.Println(bulkDeleteResponse)
}
