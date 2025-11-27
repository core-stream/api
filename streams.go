package corestream

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListStreams returns a paginated list of streams.
// Use streamerID to filter streams by a specific streamer (optional, pass empty string to skip).
func (c *Client) ListStreams(ctx context.Context, page, pageSize int, streamerID string) (*ListStreamsResponse, error) {
	query := url.Values{}
	if page > 0 {
		query.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		query.Set("page_size", strconv.Itoa(pageSize))
	}
	if streamerID != "" {
		query.Set("streamer_id", streamerID)
	}

	var resp ListStreamsResponse
	if err := c.request(ctx, http.MethodGet, "/v2/streams", query, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchStreams searches for streams by keywords or phrases in their transcripts.
// The query supports individual words and "quoted phrases" for exact matches.
// timeRange can be "today", "week", or "month" (defaults to "today" if empty).
func (c *Client) SearchStreams(ctx context.Context, query string, page, pageSize int, timeRange string) (*SearchStreamsResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		params.Set("page_size", strconv.Itoa(pageSize))
	}
	if timeRange != "" {
		params.Set("time_range", timeRange)
	}

	var resp SearchStreamsResponse
	if err := c.request(ctx, http.MethodGet, "/v2/streams/search", params, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStream retrieves detailed information about a specific stream.
func (c *Client) GetStream(ctx context.Context, streamID string) (*Stream, error) {
	path := fmt.Sprintf("/v2/streams/%s", streamID)
	var resp GetStreamResponse
	if err := c.request(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Stream, nil
}

// GetStreamTranscript retrieves the full transcript for a specific stream.
func (c *Client) GetStreamTranscript(ctx context.Context, streamID string) (*TranscriptResponse, error) {
	path := fmt.Sprintf("/v2/streams/%s/transcript", streamID)
	var resp TranscriptResponse
	if err := c.request(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
