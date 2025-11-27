package corestream

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestListAlerts(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v2/alerts" {
			t.Errorf("expected path /v2/alerts, got %s", r.URL.Path)
		}

		// Check pagination params
		page := r.URL.Query().Get("page")
		pageSize := r.URL.Query().Get("page_size")
		if page != "1" {
			t.Errorf("expected page=1, got %s", page)
		}
		if pageSize != "20" {
			t.Errorf("expected page_size=20, got %s", pageSize)
		}

		resp := ListAlertsResponse{
			Alerts: []Alert{
				{
					ID:        "alert_123",
					Name:      "Test Alert",
					Phrases:   []string{"test phrase"},
					IsActive:  true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
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
	result, err := client.ListAlerts(ctx, 1, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Alerts) != 1 {
		t.Errorf("expected 1 alert, got %d", len(result.Alerts))
	}
	if result.Alerts[0].ID != "alert_123" {
		t.Errorf("expected alert ID 'alert_123', got %s", result.Alerts[0].ID)
	}
}

func TestCreateAlert(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v2/alerts" {
			t.Errorf("expected path /v2/alerts, got %s", r.URL.Path)
		}

		var req CreateAlertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Name != "My Alert" {
			t.Errorf("expected name 'My Alert', got %s", req.Name)
		}

		w.WriteHeader(http.StatusCreated)
		alert := Alert{
			ID:        "alert_new",
			Name:      req.Name,
			Phrases:   req.Phrases,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		json.NewEncoder(w).Encode(alert)
	})
	defer server.Close()

	ctx := context.Background()
	isActive := true
	result, err := client.CreateAlert(ctx, &CreateAlertRequest{
		Name:     "My Alert",
		Phrases:  []string{"phrase1", "phrase2"},
		IsActive: &isActive,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "alert_new" {
		t.Errorf("expected alert ID 'alert_new', got %s", result.ID)
	}
}

func TestGetAlert(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v2/alerts/alert_123" {
				t.Errorf("expected path /v2/alerts/alert_123, got %s", r.URL.Path)
			}

			alert := Alert{
				ID:        "alert_123",
				Name:      "Test Alert",
				Phrases:   []string{"test"},
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			json.NewEncoder(w).Encode(alert)
		})
		defer server.Close()

		ctx := context.Background()
		result, err := client.GetAlert(ctx, "alert_123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "alert_123" {
			t.Errorf("expected alert ID 'alert_123', got %s", result.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"not_found","message":"Alert not found"}}`))
		})
		defer server.Close()

		ctx := context.Background()
		_, err := client.GetAlert(ctx, "nonexistent")

		if err == nil {
			t.Fatal("expected error")
		}
		if !IsNotFound(err) {
			t.Errorf("expected not found error, got %v", err)
		}
	})
}

func TestUpdateAlert(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/v2/alerts/alert_123" {
			t.Errorf("expected path /v2/alerts/alert_123, got %s", r.URL.Path)
		}

		var req UpdateAlertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		alert := Alert{
			ID:        "alert_123",
			Name:      *req.Name,
			Phrases:   req.Phrases,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		json.NewEncoder(w).Encode(alert)
	})
	defer server.Close()

	ctx := context.Background()
	newName := "Updated Alert"
	result, err := client.UpdateAlert(ctx, "alert_123", &UpdateAlertRequest{
		Name:    &newName,
		Phrases: []string{"new phrase"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated Alert" {
		t.Errorf("expected name 'Updated Alert', got %s", result.Name)
	}
}

func TestDeleteAlert(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			if r.URL.Path != "/v2/alerts/alert_123" {
				t.Errorf("expected path /v2/alerts/alert_123, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		defer server.Close()

		ctx := context.Background()
		err := client.DeleteAlert(ctx, "alert_123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"not_found","message":"Alert not found"}}`))
		})
		defer server.Close()

		ctx := context.Background()
		err := client.DeleteAlert(ctx, "nonexistent")

		if err == nil {
			t.Fatal("expected error")
		}
		if !IsNotFound(err) {
			t.Errorf("expected not found error, got %v", err)
		}
	})
}

func TestGetAlertNotifications(t *testing.T) {
	client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v2/alerts/alert_123/notifications" {
			t.Errorf("expected path /v2/alerts/alert_123/notifications, got %s", r.URL.Path)
		}

		resp := ListNotificationsResponse{
			Notifications: []Notification{
				{
					ID:            "notif_456",
					AlertID:       "alert_123",
					AlertName:     "Test Alert",
					MatchedPhrase: "test phrase",
					Context:       "...context...",
					StreamSource:  "Twitch",
					StreamTitle:   "Test Stream",
					Timestamp:     time.Now(),
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
	result, err := client.GetAlertNotifications(ctx, "alert_123", 1, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Notifications) != 1 {
		t.Errorf("expected 1 notification, got %d", len(result.Notifications))
	}
	if result.Notifications[0].ID != "notif_456" {
		t.Errorf("expected notification ID 'notif_456', got %s", result.Notifications[0].ID)
	}
}
