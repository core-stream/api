# api

[![Go Reference](https://pkg.go.dev/badge/github.com/core-stream/api.svg)](https://pkg.go.dev/github.com/core-stream/api)

Go client library for the core.stream API.

## Installation

```bash
go get github.com/core-stream/api
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	corestream "github.com/core-stream/api"
)

func main() {
	client, err := corestream.NewClient("your-api-token")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Search streams
	results, err := client.SearchStreams(ctx, "gaming", 1, 10, "week")
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results.Results {
		fmt.Printf("%s by %s\n", result.Title, result.UserDisplayName)
	}
}
```
