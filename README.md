# SerpApi Go Client

A robust, idiomatic Go client for [SerpApi](https://serpapi.com).

## Features

- **Full Engine Support**: Works with all SerpApi engines (Google, Bing, eBay, YouTube, etc.).
- **Context Support**: Built-in support for `context.Context` for timeouts and cancellations.
- **Pagination Helper**: Easy results traversal with `NextPageParams()`.
- **Flexible Configuration**: Use functional options for custom HTTP clients or base URLs.
- **Clean API**: Simple, map-based results for maximum flexibility with varying engine schemas.

## Installation

```bash
go get github.com/serpapi/serpapi-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/serpapi/serpapi-go"
)

func main() {
    apiKey := os.Getenv("SERPAPI_API_KEY")
    client := serpapi.NewClient(apiKey)

    params := map[string]string{
        "engine": "google",
        "q":      "Coffee",
    }

    results, err := client.Search(context.Background(), params)
    if err != nil {
        log.Fatalf("Search failed: %v", err)
    }

    // Access results as a map
    if organic, ok := results["organic_results"].([]interface{}); ok {
        for _, r := range organic {
            result := r.(map[string]interface{})
            fmt.Println(result["title"])
        }
    }
}
```

## Advanced Usage

### Pagination

The client provides a helper to fetch subsequent pages easily:

```go
results, _ := client.Search(ctx, params)

// Get parameters for the next page
nextParams := results.NextPageParams()
if nextParams != nil {
    nextResults, _ := client.Search(ctx, nextParams)
    // ... process nextResults
}
```

### Custom Configuration

You can configure the client using functional options:

```go
customHTTPClient := &http.Client{Timeout: 30 * time.Second}

client := serpapi.NewClient(
    apiKey,
    serpapi.WithHTTPClient(customHTTPClient),
    serpapi.WithBaseURL("https://custom.serpapi.proxy"),
)
```

### Context & Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

results, err := client.Search(ctx, params)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
