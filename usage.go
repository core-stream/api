package corestream

import (
	"context"
	"net/http"
)

// GetMonthlyUsage retrieves monthly API usage with billing information.
// This endpoint is only available for Enterprise tier users.
func (c *Client) GetMonthlyUsage(ctx context.Context) (*MonthlyUsageResponse, error) {
	var resp MonthlyUsageResponse
	if err := c.request(ctx, http.MethodGet, "/v2/usage/monthly", nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
