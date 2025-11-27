package corestream

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// generateSignature creates a valid HMAC-SHA256 signature for testing.
func generateSignature(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{"id":"test"}`)

	t.Run("valid signature", func(t *testing.T) {
		signature := generateSignature(body, secret)
		if !VerifyWebhookSignature(body, signature, secret) {
			t.Error("expected signature to be valid")
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		if VerifyWebhookSignature(body, "invalid-signature", secret) {
			t.Error("expected signature to be invalid")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		signature := generateSignature(body, secret)
		if VerifyWebhookSignature(body, signature, "wrong-secret") {
			t.Error("expected signature to be invalid with wrong secret")
		}
	})

	t.Run("tampered body", func(t *testing.T) {
		signature := generateSignature(body, secret)
		tamperedBody := []byte(`{"id":"tampered"}`)
		if VerifyWebhookSignature(tamperedBody, signature, secret) {
			t.Error("expected signature to be invalid for tampered body")
		}
	})

	t.Run("invalid hex signature", func(t *testing.T) {
		if VerifyWebhookSignature(body, "not-hex", secret) {
			t.Error("expected invalid hex to fail")
		}
	})
}

func TestParseWebhookNotification(t *testing.T) {
	t.Run("valid payload", func(t *testing.T) {
		payload := WebhookNotification{
			ID:            "notif_123",
			AlertID:       "alert_456",
			StreamID:      "stream_789",
			MatchedPhrase: "test phrase",
			ContextText:   "...context around test phrase...",
			Timestamp:     time.Now(),
		}
		body, _ := json.Marshal(payload)

		result, err := ParseWebhookNotification(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "notif_123" {
			t.Errorf("expected ID 'notif_123', got %s", result.ID)
		}
		if result.MatchedPhrase != "test phrase" {
			t.Errorf("expected matched phrase 'test phrase', got %s", result.MatchedPhrase)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := ParseWebhookNotification([]byte(`{invalid json`))
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("with full transcript", func(t *testing.T) {
		payload := WebhookNotification{
			ID:             "notif_123",
			AlertID:        "alert_456",
			MatchedPhrase:  "test",
			Timestamp:      time.Now(),
			FullTranscript: "This is the full transcript of the stream...",
		}
		body, _ := json.Marshal(payload)

		result, err := ParseWebhookNotification(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.FullTranscript == "" {
			t.Error("expected full transcript to be set")
		}
	})
}

func TestWebhookReceiver_ServeHTTP(t *testing.T) {
	secret := "test-secret"

	t.Run("valid request", func(t *testing.T) {
		handlerCalled := false
		var receivedNotification *WebhookNotification

		receiver := NewWebhookReceiver(secret, func(n *WebhookNotification) error {
			handlerCalled = true
			receivedNotification = n
			return nil
		})

		payload := WebhookNotification{
			ID:            "notif_123",
			AlertID:       "alert_456",
			MatchedPhrase: "test phrase",
			Timestamp:     time.Now(),
		}
		body, _ := json.Marshal(payload)
		signature := generateSignature(body, secret)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(SignatureHeader, signature)
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
		if !handlerCalled {
			t.Error("handler was not called")
		}
		if receivedNotification.ID != "notif_123" {
			t.Errorf("expected notification ID 'notif_123', got %s", receivedNotification.ID)
		}
	})

	t.Run("missing signature", func(t *testing.T) {
		receiver := NewWebhookReceiver(secret, func(n *WebhookNotification) error {
			return nil
		})

		body := []byte(`{"id":"test"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		receiver := NewWebhookReceiver(secret, func(n *WebhookNotification) error {
			return nil
		})

		body := []byte(`{"id":"test"}`)
		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(SignatureHeader, "invalid-signature")
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		receiver := NewWebhookReceiver(secret, func(n *WebhookNotification) error {
			return nil
		})

		req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", rec.Code)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		receiver := NewWebhookReceiver(secret, func(n *WebhookNotification) error {
			return nil
		})

		body := []byte(`{invalid json}`)
		signature := generateSignature(body, secret)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(SignatureHeader, signature)
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("handler error", func(t *testing.T) {
		receiver := NewWebhookReceiver(secret, func(n *WebhookNotification) error {
			return http.ErrAbortHandler // Any error
		})

		payload := WebhookNotification{
			ID:        "notif_123",
			AlertID:   "alert_456",
			Timestamp: time.Now(),
		}
		body, _ := json.Marshal(payload)
		signature := generateSignature(body, secret)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(SignatureHeader, signature)
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}
	})

	t.Run("verification disabled - no signature", func(t *testing.T) {
		handlerCalled := false
		receiver := NewWebhookReceiver("", func(n *WebhookNotification) error {
			handlerCalled = true
			return nil
		}, WithoutSignatureVerification())

		payload := WebhookNotification{
			ID:            "notif_123",
			AlertID:       "alert_456",
			MatchedPhrase: "test phrase",
			Timestamp:     time.Now(),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
		if !handlerCalled {
			t.Error("handler was not called")
		}
	})

	t.Run("verification disabled - invalid signature ignored", func(t *testing.T) {
		handlerCalled := false
		receiver := NewWebhookReceiver("", func(n *WebhookNotification) error {
			handlerCalled = true
			return nil
		}, WithoutSignatureVerification())

		payload := WebhookNotification{
			ID:            "notif_123",
			AlertID:       "alert_456",
			MatchedPhrase: "test phrase",
			Timestamp:     time.Now(),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
		req.Header.Set(SignatureHeader, "invalid-signature")
		rec := httptest.NewRecorder()

		receiver.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
		if !handlerCalled {
			t.Error("handler was not called")
		}
	})
}
