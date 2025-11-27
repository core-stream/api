package corestream

import (
	"errors"
	"fmt"
)

// APIError represents an error response from the core.stream API.
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("corestream: %s (status %d, code %s)", e.Message, e.StatusCode, e.Code)
	}
	return fmt.Sprintf("corestream: request failed with status %d", e.StatusCode)
}

// Webhook signature errors.
var (
	ErrMissingSignature = errors.New("corestream: missing webhook signature")
	ErrInvalidSignature = errors.New("corestream: invalid webhook signature")
)

// IsNotFound returns true if the error is a 404 Not Found response.
func IsNotFound(err error) bool {
	return isStatusCode(err, 404)
}

// IsUnauthorized returns true if the error is a 401 Unauthorized response.
func IsUnauthorized(err error) bool {
	return isStatusCode(err, 401)
}

// IsRateLimited returns true if the error is a 429 Rate Limit Exceeded response.
func IsRateLimited(err error) bool {
	return isStatusCode(err, 429)
}

// IsForbidden returns true if the error is a 403 Forbidden response.
func IsForbidden(err error) bool {
	return isStatusCode(err, 403)
}

func isStatusCode(err error, statusCode int) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == statusCode
	}
	return false
}
