package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	corestream "github.com/core-stream/api"
)

func main() {
	secret := os.Getenv("CORESTREAM_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("CORESTREAM_WEBHOOK_SECRET environment variable is required")
	}

	// Create a webhook receiver with a handler function
	receiver := corestream.NewWebhookReceiver(secret, handleWebhook)

	// Register the webhook endpoint
	http.Handle("/webhooks/corestream", receiver)

	// Start the server
	addr := ":8080"
	log.Printf("Webhook server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleWebhook(notification *corestream.WebhookNotification) error {
	log.Printf("Received webhook notification:")
	log.Printf("  ID: %s", notification.ID)
	log.Printf("  Alert ID: %s", notification.AlertID)
	log.Printf("  Matched Phrase: %s", notification.MatchedPhrase)
	log.Printf("  Timestamp: %s", notification.Timestamp)

	if notification.StreamID != "" {
		log.Printf("  Stream ID: %s", notification.StreamID)
	}
	if notification.StreamerID != "" {
		log.Printf("  Streamer ID: %s", notification.StreamerID)
	}
	if notification.ContextText != "" {
		log.Printf("  Context: %s", notification.ContextText)
	}
	if notification.FullTranscript != "" {
		log.Printf("  Full Transcript: %s...", truncate(notification.FullTranscript, 100))
	}

	// You can add your custom logic here, for example:
	// - Save to database
	// - Send to Slack/Discord
	// - Trigger other workflows

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// Example of manual webhook handling (without WebhookReceiver)
func manualWebhookHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("CORESTREAM_WEBHOOK_SECRET")

	// Read the request body
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get the signature from header
	signature := r.Header.Get(corestream.SignatureHeader)
	if signature == "" {
		http.Error(w, "missing signature", http.StatusUnauthorized)
		return
	}

	// Verify the signature
	if !corestream.VerifyWebhookSignature(body, signature, secret) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse the notification
	notification, err := corestream.ParseWebhookNotification(body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Process the notification
	log.Printf("Received notification for alert %s: %s", notification.AlertID, notification.MatchedPhrase)

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
