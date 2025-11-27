package corestream

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
)

const (
	// SignatureHeader is the HTTP header containing the HMAC signature.
	SignatureHeader = "X-Webhook-Signature"

	// MaxWebhookBodySize limits the webhook body to prevent DoS (1 MB).
	MaxWebhookBodySize = 1 << 20
)

// WebhookHandler is a function that processes validated webhook notifications.
type WebhookHandler func(notification *WebhookNotification) error

// WebhookReceiver handles incoming webhooks with signature verification.
// It implements http.Handler for easy integration with HTTP servers.
type WebhookReceiver struct {
	secret      []byte
	handler     WebhookHandler
	maxBodySize int64
}

// NewWebhookReceiver creates a new webhook receiver.
// The secret is used for HMAC-SHA256 signature verification.
// The handler is called for each validated webhook notification.
func NewWebhookReceiver(secret string, handler WebhookHandler) *WebhookReceiver {
	maxBodySize := int64(MaxWebhookBodySize)
	return &WebhookReceiver{
		secret:      []byte(secret),
		handler:     handler,
		maxBodySize: maxBodySize,
	}
}

// ServeHTTP implements http.Handler.
func (r *WebhookReceiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	signature := req.Header.Get(SignatureHeader)
	if signature == "" {
		http.Error(w, ErrMissingSignature.Error(), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(io.LimitReader(req.Body, r.maxBodySize))
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	if !verifySignature(body, signature, r.secret) {
		http.Error(w, ErrInvalidSignature.Error(), http.StatusUnauthorized)
		return
	}

	notification, err := ParseWebhookNotification(body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := r.handler(notification); err != nil {
		http.Error(w, "handler error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// VerifyWebhookSignature verifies the HMAC-SHA256 signature of a webhook payload.
// This is useful for manual webhook handling outside of WebhookReceiver.
func VerifyWebhookSignature(body []byte, signature, secret string) bool {
	return verifySignature(body, signature, []byte(secret))
}

func verifySignature(body []byte, signature string, secret []byte) bool {
	expectedSig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	computedSig := mac.Sum(nil)

	return hmac.Equal(expectedSig, computedSig)
}

// ParseWebhookNotification parses a webhook payload into a WebhookNotification.
// This is useful for manual webhook handling outside of WebhookReceiver.
func ParseWebhookNotification(body []byte) (*WebhookNotification, error) {
	var notification WebhookNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		return nil, err
	}
	return &notification, nil
}
