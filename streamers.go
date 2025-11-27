package corestream

import (
	"context"
	"fmt"
	"net/http"
)

// GetStreamer retrieves detailed information about a specific streamer.
func (c *Client) GetStreamer(ctx context.Context, streamerID string) (*Streamer, error) {
	path := fmt.Sprintf("/v2/streamers/%s", streamerID)
	var streamer Streamer
	if err := c.request(ctx, http.MethodGet, path, nil, nil, &streamer); err != nil {
		return nil, err
	}
	return &streamer, nil
}
