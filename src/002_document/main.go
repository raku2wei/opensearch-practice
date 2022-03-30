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

	// Add a document to the index.
	document := strings.NewReader(`{
	    "tweet_text": "吾輩は猫である。名前はまだ無い。どこで生れたかとんと見当がつかぬ。",
	    "user_name": "夏目"
	}`)

	docId := "1" // 既存サービスのDB上のidを入れる（Fantiaでいえば、posts.id とか）
	req := opensearchapi.IndexRequest{
		Index:      IndexName,
		DocumentID: docId, // DocumentIDを指定するとPUT(更新)
		Body:       document,
	}
	insertResponse, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to insert document ", err)
		os.Exit(1)
	}
	fmt.Println("### create document")
	fmt.Println(insertResponse)

	// Get document.
	get := opensearchapi.GetRequest{
		Index:      IndexName,
		DocumentID: docId,
	}

	getResponse, err := get.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to search document ", err)
		os.Exit(1)
	}
	fmt.Println("### get document")
	fmt.Println(getResponse)

	// Delete the document.
	delete := opensearchapi.DeleteRequest{
		Index:      IndexName,
		DocumentID: docId,
	}

	deleteResponse, err := delete.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to delete document ", err)
		os.Exit(1)
	}
	fmt.Println("### deleting document")
	fmt.Println(deleteResponse)
}
