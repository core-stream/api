package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	corestream "github.com/core-stream/api"
)

func main() {
	token := os.Getenv("CORESTREAM_API_TOKEN")
	if token == "" {
		log.Fatal("CORESTREAM_API_TOKEN environment variable is required")
	}

	// Create a new client
	client, err := corestream.NewClient(token)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example: List all alerts
	fmt.Println("=== Listing Alerts ===")
	alerts, err := client.ListAlerts(ctx, 1, 10)
	if err != nil {
		log.Printf("Error listing alerts: %v", err)
	} else {
		fmt.Printf("Found %d alerts\n", alerts.Pagination.TotalItems)
		for _, alert := range alerts.Alerts {
			fmt.Printf("  - %s: %s (active: %t)\n", alert.ID, alert.Name, alert.IsActive)
		}
	}

	// Example: Create an alert
	fmt.Println("\n=== Creating Alert ===")
	isActive := true
	newAlert, err := client.CreateAlert(ctx, &corestream.CreateAlertRequest{
		Name:     "Product Mentions",
		Phrases:  []string{"our product", "competitor", "pricing"},
		IsActive: &isActive,
	})
	if err != nil {
		log.Printf("Error creating alert: %v", err)
	} else {
		fmt.Printf("Created alert: %s\n", newAlert.ID)

		// Example: Get the alert we just created
		fmt.Println("\n=== Getting Alert ===")
		alert, err := client.GetAlert(ctx, newAlert.ID)
		if err != nil {
			log.Printf("Error getting alert: %v", err)
		} else {
			fmt.Printf("Alert: %s - Phrases: %v\n", alert.Name, alert.Phrases)
		}

		// Example: Create a webhook for the alert
		fmt.Println("\n=== Creating Webhook ===")
		webhook, err := client.CreateWebhook(ctx, newAlert.ID, &corestream.CreateWebhookRequest{
			URL:    "https://your-domain.com/webhook",
			Secret: "your-secret-key",
		})
		if err != nil {
			log.Printf("Error creating webhook: %v", err)
		} else {
			fmt.Printf("Created webhook: %s\n", webhook.ID)
		}

		// Clean up: Delete the alert
		fmt.Println("\n=== Deleting Alert ===")
		if err := client.DeleteAlert(ctx, newAlert.ID); err != nil {
			log.Printf("Error deleting alert: %v", err)
		} else {
			fmt.Println("Alert deleted successfully")
		}
	}

	// Example: Search streams
	fmt.Println("\n=== Searching Streams ===")
	results, err := client.SearchStreams(ctx, "gaming setup", 1, 5, "week")
	if err != nil {
		log.Printf("Error searching streams: %v", err)
	} else {
		fmt.Printf("Found %d results\n", results.Pagination.TotalItems)
		for _, result := range results.Results {
			fmt.Printf("  - %s by %s\n", result.Title, result.UserDisplayName)
			for _, highlight := range result.Highlights {
				fmt.Printf("    %s\n", highlight)
			}
		}
	}

	// Example: Get a streamer
	fmt.Println("\n=== Getting Streamer ===")
	streamer, err := client.GetStreamer(ctx, "streamer_xyz789")
	if err != nil {
		var apiErr *corestream.APIError
		if errors.As(err, &apiErr) {
			if corestream.IsNotFound(err) {
				fmt.Println("Streamer not found")
			} else {
				fmt.Printf("API error: %s\n", apiErr.Message)
			}
		} else {
			log.Printf("Error getting streamer: %v", err)
		}
	} else {
		fmt.Printf("Streamer: %s (%s)\n", streamer.DisplayName, streamer.Login)
	}
}
