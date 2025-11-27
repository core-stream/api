package corestream

import (
	"net/http"
	"time"
)

// HTTPClient interface allows for custom HTTP client implementations and testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Pagination contains pagination information for list responses.
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// Alert represents an alert configuration.
type Alert struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phrases   []string  `json:"phrases"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAlertRequest is the request body for creating an alert.
type CreateAlertRequest struct {
	Name     string   `json:"name"`
	Phrases  []string `json:"phrases"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// UpdateAlertRequest is the request body for updating an alert.
type UpdateAlertRequest struct {
	Name     *string  `json:"name,omitempty"`
	Phrases  []string `json:"phrases,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// ListAlertsResponse is the response for listing alerts.
type ListAlertsResponse struct {
	Alerts     []Alert    `json:"alerts"`
	Pagination Pagination `json:"pagination"`
}

// Notification represents an alert notification.
type Notification struct {
	ID            string    `json:"id"`
	AlertID       string    `json:"alert_id"`
	AlertName     string    `json:"alert_name"`
	MatchedPhrase string    `json:"matched_phrase"`
	Context       string    `json:"context"`
	StreamSource  string    `json:"stream_source"`
	StreamTitle   string    `json:"stream_title"`
	Timestamp     time.Time `json:"timestamp"`
	TranscriptURL string    `json:"transcript_url,omitempty"`
}

// ListNotificationsResponse is the response for listing alert notifications.
type ListNotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
	Pagination    Pagination     `json:"pagination"`
}

// Webhook represents a webhook configuration.
type Webhook struct {
	ID                    string    `json:"id"`
	AlertID               string    `json:"alert_id"`
	URL                   string    `json:"url"`
	Secret                string    `json:"secret,omitempty"`
	IsActive              bool      `json:"is_active"`
	IncludeFullTranscript bool      `json:"include_full_transcript"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// CreateWebhookRequest is the request body for creating a webhook.
type CreateWebhookRequest struct {
	URL                   string `json:"url"`
	Secret                string `json:"secret,omitempty"`
	IsActive              *bool  `json:"is_active,omitempty"`
	IncludeFullTranscript *bool  `json:"include_full_transcript,omitempty"`
}

// UpdateWebhookRequest is the request body for updating a webhook.
type UpdateWebhookRequest struct {
	URL                   string `json:"url"`
	Secret                string `json:"secret,omitempty"`
	IsActive              bool   `json:"is_active"`
	IncludeFullTranscript bool   `json:"include_full_transcript"`
}

// TestWebhookRequest is the request body for testing a webhook.
type TestWebhookRequest struct {
	URL                   string `json:"url,omitempty"`
	Secret                string `json:"secret,omitempty"`
	IncludeFullTranscript *bool  `json:"include_full_transcript,omitempty"`
}

// Stream represents a stream.
type Stream struct {
	ID              string    `json:"id"`
	StreamerID      string    `json:"streamer_id"`
	TwitchID        string    `json:"twitch_id,omitempty"`
	Title           string    `json:"title,omitempty"`
	VodID           string    `json:"vod_id,omitempty"`
	VodURL          string    `json:"vod_url,omitempty"`
	StartedAt       time.Time `json:"started_at"`
	DurationSeconds int       `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// ListStreamsResponse is the response for listing streams.
type ListStreamsResponse struct {
	Streams  []Stream `json:"streams"`
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}

// GetStreamResponse wraps a single stream response.
type GetStreamResponse struct {
	Stream Stream `json:"stream"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	StreamID        string    `json:"stream_id"`
	StreamerID      string    `json:"streamer_id"`
	Title           string    `json:"title"`
	UserDisplayName string    `json:"user_display_name"`
	Highlights      []string  `json:"highlights"`
	CreatedAt       time.Time `json:"created_at"`
}

// SearchStreamsResponse is the response for searching streams.
type SearchStreamsResponse struct {
	Results  []SearchResult `json:"results"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// TranscriptSegment represents a single transcript segment.
type TranscriptSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

// TranscriptResponse is the response for getting a stream transcript.
type TranscriptResponse struct {
	Segments []TranscriptSegment `json:"segments"`
}

// Streamer represents a streamer profile.
type Streamer struct {
	ID              string    `json:"id"`
	TwitchID        string    `json:"twitch_id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	Type            string    `json:"type,omitempty"`
	BroadcasterType string    `json:"broadcaster_type,omitempty"`
	Description     string    `json:"description,omitempty"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	OfflineImageURL string    `json:"offline_image_url,omitempty"`
	ViewCount       int       `json:"view_count,omitempty"`
	Followers       int       `json:"followers,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	FetchedAt       time.Time `json:"fetched_at"`
}

// BillingSummary contains billing information for Enterprise users.
type BillingSummary struct {
	UserID             string    `json:"user_id"`
	BillingPeriodStart time.Time `json:"billing_period_start"`
	BillingPeriodEnd   time.Time `json:"billing_period_end"`
	TotalRequests      int       `json:"total_requests"`
	IncludedRequests   int       `json:"included_requests"`
	BillableRequests   int       `json:"billable_requests"`
	SubscriptionTier   string    `json:"subscription_tier"`
}

// Subscription contains subscription status information.
type Subscription struct {
	Status string `json:"status"`
	Tier   string `json:"tier"`
}

// MonthlyUsageResponse is the response for getting monthly usage.
type MonthlyUsageResponse struct {
	BillingSummary BillingSummary `json:"billing_summary"`
	Subscription   Subscription   `json:"subscription"`
}

// WebhookNotification is the payload received from core.stream webhooks.
type WebhookNotification struct {
	ID             string    `json:"id"`
	AlertID        string    `json:"alert_id"`
	StreamID       string    `json:"stream_id,omitempty"`
	StreamerID     string    `json:"streamer_id,omitempty"`
	MatchedPhrase  string    `json:"matched_phrase"`
	ContextText    string    `json:"context_text,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	FullTranscript string    `json:"full_transcript,omitempty"`
}
