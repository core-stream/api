package corestream

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestListStreams(t *testing.T) {
	t.Run("basic list", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v2/streams" {
				t.Errorf("expected path /v2/streams, got %s", r.URL.Path)
			}

			resp := ListStreamsResponse{
				Streams: []Stream{
					{
						ID:         "stream_abc",
						StreamerID: "streamer_xyz",
						Title:      "Test Stream",
						StartedAt:  time.Now(),
						CreatedAt:  time.Now(),
					},
				},
				Pagination: Pagination{
					Page:       1,
					PageSize:   20,
					TotalItems: 1,
					TotalPages: 1,
				},
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		ctx := context.Background()
		result, err := client.ListStreams(ctx, 1, 20, "")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Streams) != 1 {
			t.Errorf("expected 1 stream, got %d", len(result.Streams))
		}
	})

	t.Run("with streamer_id filter", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			streamerID := r.URL.Query().Get("streamer_id")
			if streamerID != "streamer_xyz" {
				t.Errorf("expected streamer_id=streamer_xyz, got %s", streamerID)
			}

			resp := ListStreamsResponse{
				Streams: []Stream{},
				Pagination: Pagination{
					Page:       1,
					PageSize:   20,
					TotalItems: 0,
					TotalPages: 0,
				},
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		ctx := context.Background()
		_, err := client.ListStreams(ctx, 1, 20, "streamer_xyz")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestSearchStreams(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v2/streams/search" {
			t.Errorf("expected path /v2/streams/search, got %s", r.URL.Path)
		}

		query := r.URL.Query().Get("q")
		if query != "gaming setup" {
			t.Errorf("expected q='gaming setup', got %s", query)
		}

		timeRange := r.URL.Query().Get("time_range")
		if timeRange != "week" {
			t.Errorf("expected time_range=week, got %s", timeRange)
		}

		resp := SearchStreamsResponse{
			Results: []SearchResult{
				{
					StreamID:        "stream_abc",
					StreamerID:      "streamer_xyz",
					Title:           "Gaming Setup Tour",
					UserDisplayName: "TestStreamer",
					Highlights:      []string{"check out my <em>gaming setup</em>"},
					CreatedAt:       time.Now(),
				},
			},
			Pagination: Pagination{
				Page:       1,
				PageSize:   20,
				TotalItems: 1,
				TotalPages: 1,
			},
		}
		json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.SearchStreams(ctx, "gaming setup", 1, 20, "week")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(result.Results))
	}
	if result.Results[0].Title != "Gaming Setup Tour" {
		t.Errorf("expected title 'Gaming Setup Tour', got %s", result.Results[0].Title)
	}
}

func TestGetStream(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v2/streams/stream_abc" {
				t.Errorf("expected path /v2/streams/stream_abc, got %s", r.URL.Path)
			}

			resp := GetStreamResponse{
				Stream: Stream{
					ID:              "stream_abc",
					StreamerID:      "streamer_xyz",
					Title:           "Test Stream",
					VodURL:          "https://twitch.tv/videos/123",
					StartedAt:       time.Now(),
					DurationSeconds: 3600,
					CreatedAt:       time.Now(),
				},
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		ctx := context.Background()
		result, err := client.GetStream(ctx, "stream_abc")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "stream_abc" {
			t.Errorf("expected stream ID 'stream_abc', got %s", result.ID)
		}
		if result.DurationSeconds != 3600 {
			t.Errorf("expected duration 3600, got %d", result.DurationSeconds)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"not_found","message":"Stream not found"}}`))
		})
		defer server.Close()

		ctx := context.Background()
		_, err := client.GetStream(ctx, "nonexistent")

		if err == nil {
			t.Fatal("expected error")
		}
		if !IsNotFound(err) {
			t.Errorf("expected not found error, got %v", err)
		}
	})
}

func TestGetStreamTranscript(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v2/streams/stream_abc/transcript" {
			t.Errorf("expected path /v2/streams/stream_abc/transcript, got %s", r.URL.Path)
		}

		resp := TranscriptResponse{
			Segments: []TranscriptSegment{
				{Start: 0.0, End: 3.5, Text: "Hello everyone!"},
				{Start: 3.5, End: 7.0, Text: "Welcome to the stream."},
				{Start: 7.0, End: 12.0, Text: "Today we're going to talk about coding."},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	ctx := context.Background()
	result, err := client.GetStreamTranscript(ctx, "stream_abc")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Segments) != 3 {
		t.Errorf("expected 3 segments, got %d", len(result.Segments))
	}
	if result.Segments[0].Text != "Hello everyone!" {
		t.Errorf("expected first segment text 'Hello everyone!', got %s", result.Segments[0].Text)
	}
}
