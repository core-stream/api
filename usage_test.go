package corestream

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGetMonthlyUsage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v2/usage/monthly" {
				t.Errorf("expected path /v2/usage/monthly, got %s", r.URL.Path)
			}

			resp := MonthlyUsageResponse{
				BillingSummary: BillingSummary{
					UserID:             "user_123",
					BillingPeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					BillingPeriodEnd:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
					TotalRequests:      315750,
					IncludedRequests:   300000,
					BillableRequests:   15750,
					SubscriptionTier:   "enterprise",
				},
				Subscription: Subscription{
					Status: "active",
					Tier:   "enterprise",
				},
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		ctx := context.Background()
		result, err := client.GetMonthlyUsage(ctx)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.BillingSummary.TotalRequests != 315750 {
			t.Errorf("expected 315750 total requests, got %d", result.BillingSummary.TotalRequests)
		}
		if result.BillingSummary.BillableRequests != 15750 {
			t.Errorf("expected 15750 billable requests, got %d", result.BillingSummary.BillableRequests)
		}
		if result.Subscription.Tier != "enterprise" {
			t.Errorf("expected tier 'enterprise', got %s", result.Subscription.Tier)
		}
	})

	t.Run("forbidden for non-enterprise", func(t *testing.T) {
		client, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":{"code":"forbidden","message":"This endpoint is only available for Enterprise tier users"}}`))
		})
		defer server.Close()

		ctx := context.Background()
		_, err := client.GetMonthlyUsage(ctx)

		if err == nil {
			t.Fatal("expected error")
		}
		if !IsForbidden(err) {
			t.Errorf("expected forbidden error, got %v", err)
		}
	})
}
