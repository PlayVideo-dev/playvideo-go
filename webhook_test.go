package playvideo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// Helper to compute signature for testing
func computeTestSignature(payload string, timestamp int64, secret string) string {
	signedPayload := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "whsec_test_secret_key"
	payload := `{"event":"video.completed","data":{"videoId":"vid1"}}`

	t.Run("valid signature", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()
		sig := computeTestSignature(payload, timestamp, secret)

		err := VerifyWebhookSignature(
			[]byte(payload),
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", timestamp),
			secret,
			300,
		)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing signature", func(t *testing.T) {
		err := VerifyWebhookSignature(
			[]byte(payload),
			"",
			"123456",
			secret,
			0,
		)

		if err == nil {
			t.Error("expected error for missing signature")
		}
		if !IsWebhookSignatureError(err) {
			t.Errorf("expected WebhookSignatureError, got %T", err)
		}
	})

	t.Run("missing timestamp", func(t *testing.T) {
		err := VerifyWebhookSignature(
			[]byte(payload),
			"sha256=abc",
			"",
			secret,
			0,
		)

		if err == nil {
			t.Error("expected error for missing timestamp")
		}
	})

	t.Run("missing secret", func(t *testing.T) {
		err := VerifyWebhookSignature(
			[]byte(payload),
			"sha256=abc",
			"123456",
			"",
			0,
		)

		if err == nil {
			t.Error("expected error for missing secret")
		}
	})

	t.Run("invalid timestamp", func(t *testing.T) {
		err := VerifyWebhookSignature(
			[]byte(payload),
			"sha256=abc",
			"not-a-number",
			secret,
			0,
		)

		if err == nil {
			t.Error("expected error for invalid timestamp")
		}
	})

	t.Run("old timestamp", func(t *testing.T) {
		oldTimestamp := time.Now().Add(-400 * time.Second).UnixMilli()
		sig := computeTestSignature(payload, oldTimestamp, secret)

		err := VerifyWebhookSignature(
			[]byte(payload),
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", oldTimestamp),
			secret,
			300,
		)

		if err == nil {
			t.Error("expected error for old timestamp")
		}
	})

	t.Run("custom tolerance", func(t *testing.T) {
		oldTimestamp := time.Now().Add(-500 * time.Second).UnixMilli()
		sig := computeTestSignature(payload, oldTimestamp, secret)

		// Should fail with 5 minute tolerance
		err := VerifyWebhookSignature(
			[]byte(payload),
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", oldTimestamp),
			secret,
			300,
		)

		if err == nil {
			t.Error("expected error with short tolerance")
		}

		// Should pass with 10 minute tolerance
		err = VerifyWebhookSignature(
			[]byte(payload),
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", oldTimestamp),
			secret,
			600,
		)

		if err != nil {
			t.Errorf("unexpected error with long tolerance: %v", err)
		}
	})

	t.Run("invalid signature format", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()

		err := VerifyWebhookSignature(
			[]byte(payload),
			"invalid_format",
			fmt.Sprintf("%d", timestamp),
			secret,
			0,
		)

		if err == nil {
			t.Error("expected error for invalid format")
		}
	})

	t.Run("non-sha256 algorithm", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()

		err := VerifyWebhookSignature(
			[]byte(payload),
			"md5=abc123",
			fmt.Sprintf("%d", timestamp),
			secret,
			0,
		)

		if err == nil {
			t.Error("expected error for non-sha256 algorithm")
		}
	})

	t.Run("signature mismatch", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()

		err := VerifyWebhookSignature(
			[]byte(payload),
			"sha256=invalid_signature",
			fmt.Sprintf("%d", timestamp),
			secret,
			0,
		)

		if err == nil {
			t.Error("expected error for signature mismatch")
		}
	})

	t.Run("tampered payload", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()
		sig := computeTestSignature(payload, timestamp, secret)

		tamperedPayload := `{"event":"video.completed","data":{"videoId":"vid2"}}`

		err := VerifyWebhookSignature(
			[]byte(tamperedPayload),
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", timestamp),
			secret,
			0,
		)

		if err == nil {
			t.Error("expected error for tampered payload")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()
		sig := computeTestSignature(payload, timestamp, secret)

		err := VerifyWebhookSignature(
			[]byte(payload),
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", timestamp),
			"wrong_secret",
			0,
		)

		if err == nil {
			t.Error("expected error for wrong secret")
		}
	})
}

func TestConstructEvent(t *testing.T) {
	secret := "whsec_test_secret_key"

	t.Run("valid event", func(t *testing.T) {
		eventData := map[string]interface{}{
			"event":     "video.completed",
			"timestamp": time.Now().UnixMilli(),
			"data":      map[string]interface{}{"videoId": "vid123"},
		}
		payload, _ := json.Marshal(eventData)
		timestamp := time.Now().UnixMilli()
		sig := computeTestSignature(string(payload), timestamp, secret)

		event, err := ConstructEvent(
			payload,
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", timestamp),
			secret,
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if event.Event != "video.completed" {
			t.Errorf("expected event 'video.completed', got %s", event.Event)
		}

		if event.Data["videoId"] != "vid123" {
			t.Errorf("expected videoId 'vid123', got %v", event.Data["videoId"])
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		timestamp := time.Now().UnixMilli()

		_, err := ConstructEvent(
			[]byte(`{"event":"video.completed"}`),
			"sha256=invalid",
			fmt.Sprintf("%d", timestamp),
			secret,
		)

		if err == nil {
			t.Error("expected error for invalid signature")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		payload := "not valid json"
		timestamp := time.Now().UnixMilli()
		sig := computeTestSignature(payload, timestamp, secret)

		_, err := ConstructEvent(
			[]byte(payload),
			fmt.Sprintf("sha256=%s", sig),
			fmt.Sprintf("%d", timestamp),
			secret,
		)

		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})
}

func TestWebhookSignatureError(t *testing.T) {
	t.Run("error message", func(t *testing.T) {
		err := NewWebhookSignatureError("Test message")
		if err.Message != "Test message" {
			t.Errorf("expected message 'Test message', got %s", err.Message)
		}
	})

	t.Run("default message", func(t *testing.T) {
		err := NewWebhookSignatureError("")
		if err.Message != "Invalid webhook signature" {
			t.Errorf("expected default message, got %s", err.Message)
		}
	})

	t.Run("implements error interface", func(t *testing.T) {
		var err error = NewWebhookSignatureError("Test")
		if err.Error() == "" {
			t.Error("expected non-empty error string")
		}
	})
}
