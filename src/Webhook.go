package keymint

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// VerifyWebhookSignature verifies a webhook payload signature received from Keymint.
// payload: The raw request body as string.
// header: The value of the "Keymint-Signature" header.
// secret: The webhook endpoint's signing secret.
// tolerance: Time tolerance duration (e.g. 5 * time.Minute) to prevent replay attacks. Set to 0 to use default (5 minutes).
// Returns nil if verification is successful, or an error if verification fails.
func VerifyWebhookSignature(payload string, header string, secret string, tolerance time.Duration) error {
	if header == "" {
		return fmt.Errorf("missing signature header")
	}
	if secret == "" {
		return fmt.Errorf("missing signing secret")
	}

	// Parse header (e.g., t=1719374021,v1=signature)
	var timestampStr string
	var signature string
	parts := strings.Split(header, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			if kv[0] == "t" {
				timestampStr = kv[1]
			} else if kv[0] == "v1" {
				signature = kv[1]
			}
		}
	}

	if timestampStr == "" || signature == "" {
		return fmt.Errorf("invalid signature header format")
	}

	// Check timestamp validity
	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	if tolerance <= 0 {
		tolerance = 5 * time.Minute
	}

	eventTime := time.Unix(timestampInt, 0)
	diff := time.Since(eventTime)
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		return fmt.Errorf("timestamp is outside tolerance limit")
	}

	// Verify HMAC signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestampStr + "." + payload))
	expectedMac := mac.Sum(nil)

	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding")
	}

	if !hmac.Equal(sigBytes, expectedMac) {
		return fmt.Errorf("signatures do not match")
	}

	return nil
}
