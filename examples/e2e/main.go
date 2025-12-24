package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/serpapi/serpapi-go"
)

func main() {
	apiKey := os.Getenv("SERPAPI_API_KEY")
	if apiKey == "" {
		log.Fatal("SERPAPI_API_KEY environment variable is not set")
	}

	client := serpapi.NewClient(apiKey)
	ctx := context.Background()

	params := map[string]string{
		"engine": "google",
		"q":      "Coffee",
	}

	fmt.Println("--- Searching Page 1 ---")
	results, err := client.Search(ctx, params)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Print some organic results
	if organic, ok := results["organic_results"].([]interface{}); ok && len(organic) > 0 {
		fmt.Printf("Top Result Page 1: %v\n", organic[0].(map[string]interface{})["title"])
	}

	// Demonstrate Pagination
	nextParams := results.NextPageParams()
	if nextParams != nil {
		fmt.Println("\n--- Searching Page 2 ---")
		nextResults, err := client.Search(ctx, nextParams)
		if err != nil {
			log.Fatalf("Search failed: %v", err)
		}
		if organic, ok := nextResults["organic_results"].([]interface{}); ok && len(organic) > 0 {
			fmt.Printf("Top Result Page 2: %v\n", organic[0].(map[string]interface{})["title"])
		}
	} else {
		fmt.Println("No next page available.")
	}

	// Print full results for the first page
	fmt.Println("\n--- Full Results (Page 1) ---")
	pretty, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(pretty))
}
