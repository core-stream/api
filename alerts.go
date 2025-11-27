package corestream

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListAlerts returns all alerts for the authenticated user.
func (c *Client) ListAlerts(ctx context.Context, page, pageSize int) (*ListAlertsResponse, error) {
	query := url.Values{}
	if page > 0 {
		query.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		query.Set("page_size", strconv.Itoa(pageSize))
	}

	var resp ListAlertsResponse
	if err := c.request(ctx, http.MethodGet, "/v2/alerts", query, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateAlert creates a new alert.
func (c *Client) CreateAlert(ctx context.Context, req *CreateAlertRequest) (*Alert, error) {
	var alert Alert
	if err := c.request(ctx, http.MethodPost, "/v2/alerts", nil, req, &alert); err != nil {
		return nil, err
	}
	return &alert, nil
}

// GetAlert retrieves a specific alert by ID.
func (c *Client) GetAlert(ctx context.Context, alertID string) (*Alert, error) {
	path := fmt.Sprintf("/v2/alerts/%s", alertID)
	var alert Alert
	if err := c.request(ctx, http.MethodGet, path, nil, nil, &alert); err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateAlert updates an existing alert.
func (c *Client) UpdateAlert(ctx context.Context, alertID string, req *UpdateAlertRequest) (*Alert, error) {
	path := fmt.Sprintf("/v2/alerts/%s", alertID)
	var alert Alert
	if err := c.request(ctx, http.MethodPut, path, nil, req, &alert); err != nil {
		return nil, err
	}
	return &alert, nil
}

// DeleteAlert permanently deletes an alert.
func (c *Client) DeleteAlert(ctx context.Context, alertID string) error {
	path := fmt.Sprintf("/v2/alerts/%s", alertID)
	return c.request(ctx, http.MethodDelete, path, nil, nil, nil)
}

// GetAlertNotifications retrieves notifications for a specific alert.
func (c *Client) GetAlertNotifications(ctx context.Context, alertID string, page, pageSize int) (*ListNotificationsResponse, error) {
	path := fmt.Sprintf("/v2/alerts/%s/notifications", alertID)

	query := url.Values{}
	if page > 0 {
		query.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		query.Set("page_size", strconv.Itoa(pageSize))
	}

	var resp ListNotificationsResponse
	if err := c.request(ctx, http.MethodGet, path, query, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
