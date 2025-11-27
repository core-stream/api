package corestream

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestCreateWebhook(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v2/alerts/alert_123/webhook" {
			t.Errorf("expected path /v2/alerts/alert_123/webhook, got %s", r.URL.Path)
		}

		var req CreateWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.URL != "https://example.com/webhook" {
			t.Errorf("expected URL 'https://example.com/webhook', got %s", req.URL)
		}

		w.WriteHeader(http.StatusCreated)
		webhook := Webhook{
			ID:                    "webhook_789",
			AlertID:               "alert_123",
			URL:                   req.URL,
			IsActive:              true,
			IncludeFullTranscript: false,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}
		json.NewEncoder(w).Encode(webhook)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.CreateWebhook(ctx, "alert_123", &CreateWebhookRequest{
		URL:    "https://example.com/webhook",
		Secret: "my-secret",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "webhook_789" {
		t.Errorf("expected webhook ID 'webhook_789', got %s", result.ID)
	}
}

func TestGetWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v2/alerts/alert_123/webhook" {
				t.Errorf("expected path /v2/alerts/alert_123/webhook, got %s", r.URL.Path)
			}

			webhook := Webhook{
				ID:                    "webhook_789",
				AlertID:               "alert_123",
				URL:                   "https://example.com/webhook",
				IsActive:              true,
				IncludeFullTranscript: false,
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			}
			json.NewEncoder(w).Encode(webhook)
		})
		defer server.Close()

		ctx := context.Background()
		result, err := client.GetWebhook(ctx, "alert_123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "webhook_789" {
			t.Errorf("expected webhook ID 'webhook_789', got %s", result.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"not_found","message":"Webhook not found"}}`))
		})
		defer server.Close()

		ctx := context.Background()
		_, err := client.GetWebhook(ctx, "alert_123")

		if err == nil {
			t.Fatal("expected error")
		}
		if !IsNotFound(err) {
			t.Errorf("expected not found error, got %v", err)
		}
	})
}

func TestUpdateWebhook(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/v2/alerts/alert_123/webhook" {
			t.Errorf("expected path /v2/alerts/alert_123/webhook, got %s", r.URL.Path)
		}

		var req UpdateWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		webhook := Webhook{
			ID:                    "webhook_789",
			AlertID:               "alert_123",
			URL:                   req.URL,
			IsActive:              req.IsActive,
			IncludeFullTranscript: req.IncludeFullTranscript,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}
		json.NewEncoder(w).Encode(webhook)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.UpdateWebhook(ctx, "alert_123", &UpdateWebhookRequest{
		URL:                   "https://new-url.com/webhook",
		IsActive:              true,
		IncludeFullTranscript: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.URL != "https://new-url.com/webhook" {
		t.Errorf("expected URL 'https://new-url.com/webhook', got %s", result.URL)
	}
	if !result.IncludeFullTranscript {
		t.Error("expected IncludeFullTranscript to be true")
	}
}

func TestDeleteWebhook(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v2/alerts/alert_123/webhook" {
			t.Errorf("expected path /v2/alerts/alert_123/webhook, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	ctx := context.Background()
	err := client.DeleteWebhook(ctx, "alert_123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTestWebhook(t *testing.T) {
	t.Run("test saved webhook", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/v2/alerts/alert_123/webhook/test" {
				t.Errorf("expected path /v2/alerts/alert_123/webhook/test, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Test webhook delivered successfully"}`))
		})
		defer server.Close()

		ctx := context.Background()
		err := client.TestWebhook(ctx, "alert_123", nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("test custom webhook", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			var req TestWebhookRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}

			if req.URL != "https://test-url.com/webhook" {
				t.Errorf("expected URL 'https://test-url.com/webhook', got %s", req.URL)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Test webhook delivered successfully"}`))
		})
		defer server.Close()

		ctx := context.Background()
		err := client.TestWebhook(ctx, "alert_123", &TestWebhookRequest{
			URL:    "https://test-url.com/webhook",
			Secret: "test-secret",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
