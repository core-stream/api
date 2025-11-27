package corestream

import (
	"context"
	"fmt"
	"net/http"
)

// CreateWebhook creates a webhook for an alert.
func (c *Client) CreateWebhook(ctx context.Context, alertID string, req *CreateWebhookRequest) (*Webhook, error) {
	path := fmt.Sprintf("/v2/alerts/%s/webhook", alertID)
	var webhook Webhook
	if err := c.request(ctx, http.MethodPost, path, nil, req, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// GetWebhook retrieves the webhook configuration for an alert.
func (c *Client) GetWebhook(ctx context.Context, alertID string) (*Webhook, error) {
	path := fmt.Sprintf("/v2/alerts/%s/webhook", alertID)
	var webhook Webhook
	if err := c.request(ctx, http.MethodGet, path, nil, nil, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// UpdateWebhook updates the webhook configuration for an alert.
func (c *Client) UpdateWebhook(ctx context.Context, alertID string, req *UpdateWebhookRequest) (*Webhook, error) {
	path := fmt.Sprintf("/v2/alerts/%s/webhook", alertID)
	var webhook Webhook
	if err := c.request(ctx, http.MethodPut, path, nil, req, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// DeleteWebhook removes the webhook configuration from an alert.
func (c *Client) DeleteWebhook(ctx context.Context, alertID string) error {
	path := fmt.Sprintf("/v2/alerts/%s/webhook", alertID)
	return c.request(ctx, http.MethodDelete, path, nil, nil, nil)
}

// TestWebhook sends a test webhook notification.
// If req is nil, tests the saved webhook configuration.
// If req is provided, tests with the specified URL/secret.
func (c *Client) TestWebhook(ctx context.Context, alertID string, req *TestWebhookRequest) error {
	path := fmt.Sprintf("/v2/alerts/%s/webhook/test", alertID)
	return c.request(ctx, http.MethodPost, path, nil, req, nil)
}
