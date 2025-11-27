package corestream

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGetStreamer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v2/streamers/streamer_xyz" {
				t.Errorf("expected path /v2/streamers/streamer_xyz, got %s", r.URL.Path)
			}

			streamer := Streamer{
				ID:              "streamer_xyz",
				TwitchID:        "123456789",
				Login:           "teststreamer",
				DisplayName:     "TestStreamer",
				BroadcasterType: "partner",
				Description:     "A test streamer",
				ViewCount:       1000000,
				Followers:       50000,
				FetchedAt:       time.Now(),
			}
			json.NewEncoder(w).Encode(streamer)
		})
		defer server.Close()

		ctx := context.Background()
		result, err := client.GetStreamer(ctx, "streamer_xyz")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "streamer_xyz" {
			t.Errorf("expected streamer ID 'streamer_xyz', got %s", result.ID)
		}
		if result.DisplayName != "TestStreamer" {
			t.Errorf("expected display name 'TestStreamer', got %s", result.DisplayName)
		}
		if result.Followers != 50000 {
			t.Errorf("expected 50000 followers, got %d", result.Followers)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"not_found","message":"Streamer not found"}}`))
		})
		defer server.Close()

		ctx := context.Background()
		_, err := client.GetStreamer(ctx, "nonexistent")

		if err == nil {
			t.Fatal("expected error")
		}
		if !IsNotFound(err) {
			t.Errorf("expected not found error, got %v", err)
		}
	})
}
