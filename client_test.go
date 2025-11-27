package corestream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupTestServer creates a mock server and client for testing.
func setupTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	client, err := NewClient("test-token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	return client, server
}

func TestNewClient(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		client, err := NewClient("my-token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("client should not be nil")
		}
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := NewClient("")
		if err == nil {
			t.Fatal("expected error for empty token")
		}
	})
}

func TestNewClient_WithBaseURL(t *testing.T) {
	client, err := NewClient("token", WithBaseURL("https://custom.api.com"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.baseURL.String() != "https://custom.api.com" {
		t.Errorf("expected base URL https://custom.api.com, got %s", client.baseURL.String())
	}
}

func TestNewClient_WithBaseURL_Invalid(t *testing.T) {
	_, err := NewClient("token", WithBaseURL("://invalid"))
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	customClient := &http.Client{}
	client, err := NewClient("token", WithHTTPClient(customClient))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.httpClient != customClient {
		t.Error("HTTP client was not set correctly")
	}
}

func TestNewClient_WithHTTPClient_Nil(t *testing.T) {
	_, err := NewClient("token", WithHTTPClient(nil))
	if err == nil {
		t.Fatal("expected error for nil HTTP client")
	}
}

func TestClient_AuthorizationHeader(t *testing.T) {
	var receivedAuth string
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	ctx := context.Background()
	client.GetStreamer(ctx, "test-id")

	expected := "Bearer test-token"
	if receivedAuth != expected {
		t.Errorf("expected Authorization header %q, got %q", expected, receivedAuth)
	}
}

func TestClient_UserAgent(t *testing.T) {
	var receivedUA string
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	ctx := context.Background()
	client.GetStreamer(ctx, "test-id")

	if receivedUA != userAgent {
		t.Errorf("expected User-Agent %q, got %q", userAgent, receivedUA)
	}
}

func TestClient_ErrorResponses(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		checkFunc  func(error) bool
	}{
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"error":{"code":"unauthorized","message":"Invalid token"}}`,
			checkFunc:  IsUnauthorized,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			body:       `{"error":{"code":"not_found","message":"Resource not found"}}`,
			checkFunc:  IsNotFound,
		},
		{
			name:       "429 Rate Limited",
			statusCode: http.StatusTooManyRequests,
			body:       `{"error":{"code":"rate_limit_exceeded","message":"Too many requests"}}`,
			checkFunc:  IsRateLimited,
		},
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			body:       `{"error":{"code":"forbidden","message":"Access denied"}}`,
			checkFunc:  IsForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			})
			defer server.Close()

			ctx := context.Background()
			_, err := client.GetStreamer(ctx, "test-id")

			if err == nil {
				t.Fatal("expected error")
			}
			if !tt.checkFunc(err) {
				t.Errorf("expected %s error, got %v", tt.name, err)
			}
		})
	}
}

func TestClient_ErrorResponse_WithDetails(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid_request","message":"Invalid parameters"}}`))
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.GetStreamer(ctx, "test-id")

	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != 400 {
		t.Errorf("expected status code 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Code != "invalid_request" {
		t.Errorf("expected code 'invalid_request', got %q", apiErr.Code)
	}
	if apiErr.Message != "Invalid parameters" {
		t.Errorf("expected message 'Invalid parameters', got %q", apiErr.Message)
	}
}
