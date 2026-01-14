package playvideo

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// DefaultWebhookTolerance is the default maximum age of a webhook in seconds
const DefaultWebhookTolerance = 300

// VerifyWebhookSignature verifies a webhook signature
//
// Parameters:
//   - payload: The raw webhook payload (as bytes or string)
//   - signature: The X-PlayVideo-Signature header value
//   - timestamp: The X-PlayVideo-Timestamp header value
//   - secret: Your webhook secret (whsec_xxx)
//   - tolerance: Maximum age of the webhook in seconds (use 0 for default of 300s)
//
// Returns nil if valid, WebhookSignatureError if invalid
func VerifyWebhookSignature(payload []byte, signature, timestamp, secret string, tolerance int) error {
	if signature == "" || timestamp == "" || secret == "" {
		return NewWebhookSignatureError("Missing required parameters for signature verification")
	}

	if tolerance == 0 {
		tolerance = DefaultWebhookTolerance
	}

	// Parse timestamp
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return NewWebhookSignatureError("Invalid timestamp")
	}

	// Check timestamp tolerance
	now := time.Now().UnixMilli()
	age := math.Abs(float64(now-ts)) / 1000
	if age > float64(tolerance) {
		return NewWebhookSignatureError(fmt.Sprintf("Webhook timestamp too old (%.0fs > %ds)", age, tolerance))
	}

	// Extract signature value (format: sha256=xxx)
	sigParts := strings.SplitN(signature, "=", 2)
	if len(sigParts) != 2 || sigParts[0] != "sha256" {
		return NewWebhookSignatureError("Invalid signature format")
	}
	expectedSig := sigParts[1]

	// Compute signature
	signedPayload := fmt.Sprintf("%s.%s", timestamp, string(payload))
	computedSig := computeHMACSHA256(signedPayload, secret)

	// Constant-time comparison
	if !secureCompare(expectedSig, computedSig) {
		return NewWebhookSignatureError("Signature mismatch")
	}

	return nil
}

// ConstructEvent verifies and parses a webhook event
//
// Parameters:
//   - payload: The raw webhook payload
//   - signature: The X-PlayVideo-Signature header value
//   - timestamp: The X-PlayVideo-Timestamp header value
//   - secret: Your webhook secret (whsec_xxx)
//
// Returns the parsed WebhookPayload if valid, error if invalid
func ConstructEvent(payload []byte, signature, timestamp, secret string) (*WebhookPayload, error) {
	if err := VerifyWebhookSignature(payload, signature, timestamp, secret, 0); err != nil {
		return nil, err
	}

	var event WebhookPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, NewWebhookSignatureError(fmt.Sprintf("Failed to parse webhook payload: %v", err))
	}

	return &event, nil
}

// computeHMACSHA256 computes an HMAC-SHA256 signature
func computeHMACSHA256(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// secureCompare performs constant-time string comparison
func secureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
