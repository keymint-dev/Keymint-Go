package keymint

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// OfflineLicenseErrorCode is the reason an offline license failed verification.
type OfflineLicenseErrorCode string

// Reason codes returned in OfflineLicenseError.Code. Checks run in the order
// listed; the first one that fails decides the code.
const (
	// OfflineMalformed: the input is not a .lic file or a well-formed token,
	// or the signed claims are missing exp or sub.
	OfflineMalformed OfflineLicenseErrorCode = "malformed"
	// OfflineUnsupportedAlgorithm: the token header's alg is not EdDSA, or it
	// carries a crit header.
	OfflineUnsupportedAlgorithm OfflineLicenseErrorCode = "unsupported_algorithm"
	// OfflineUnknownKey: the token was not issued under any of the trusted
	// public keys.
	OfflineUnknownKey OfflineLicenseErrorCode = "unknown_key"
	// OfflineBadSignature: the signature does not verify (the file was
	// altered or signed by a different key).
	OfflineBadSignature OfflineLicenseErrorCode = "bad_signature"
	// OfflineWrongType: the token is not an offline license.
	OfflineWrongType OfflineLicenseErrorCode = "wrong_type"
	// OfflineWrongProduct: the license was issued for a different product.
	OfflineWrongProduct OfflineLicenseErrorCode = "wrong_product"
	// OfflineWrongMachine: the license was issued for a different machine.
	OfflineWrongMachine OfflineLicenseErrorCode = "wrong_machine"
	// OfflineExpired: the license has expired; sign a new one online.
	OfflineExpired OfflineLicenseErrorCode = "expired"
)

// OfflineLicenseError is returned by VerifyOfflineLicense when the license
// itself fails verification. Use errors.As to read the Code. The message never
// contains the license key.
type OfflineLicenseError struct {
	Code    OfflineLicenseErrorCode
	Message string
}

func (e *OfflineLicenseError) Error() string {
	return fmt.Sprintf("offline license %s: %s", e.Code, e.Message)
}

// ErrInvalidOfflineOptions is wrapped by the error VerifyOfflineLicense
// returns when OfflineVerifyOptions are invalid (no public keys, a key that is
// not an Ed25519 SPKI PEM, an empty ProductID or MachineCode, or a ClockSkew
// outside 0..1h). This is a programming error, not a verification result.
var ErrInvalidOfflineOptions = errors.New("keymint: invalid offline verify options")

// OfflineVerifyOptions configures VerifyOfflineLicense.
type OfflineVerifyOptions struct {
	// PublicKeys are the trusted SPKI PEM public keys, from Developer → Public
	// Keys in the dashboard. At least one is required.
	PublicKeys []string
	// ProductID must equal the product the license was signed for. Required:
	// one signing key covers every product in an organization.
	ProductID string
	// MachineCode must equal the machine code (host ID) the license was
	// signed for. Required.
	MachineCode string
	// Now is the time to check expiry against. Zero means time.Now(). To
	// guard against clock rollback, pass the later of the system time and a
	// last-trusted time your app stores.
	Now time.Time
	// ClockSkew is the leeway allowed after the expiry time. Nil means 60s;
	// an explicit 0 means none. Must be between 0 and 1h.
	ClockSkew *time.Duration
}

// OfflineLicense is a verified offline license.
type OfflineLicense struct {
	// LicenseKey is the license key the file was signed for.
	LicenseKey  string
	ProductID   string
	MachineCode string
	// KeyID identifies the public key that verified the signature.
	KeyID string
	// IssuedAt is when the license was signed, or nil if the token has no iat.
	IssuedAt *time.Time
	// ExpiresAt is when the license expires.
	ExpiresAt time.Time
}

const (
	defaultOfflineClockSkew = 60 * time.Second
	maxOfflineClockSkew     = time.Hour
)

var offlineSegment = regexp.MustCompile(`^[A-Za-z0-9_-]*$`)

// OfflineKeyID returns the key ID Keymint uses for an Ed25519 SPKI PEM public
// key: base64url (no padding) of SHA-256 over the DER-encoded key.
func OfflineKeyID(publicKeyPEM string) (string, error) {
	_, kid, err := parseOfflinePublicKey(publicKeyPEM)
	return kid, err
}

func parseOfflinePublicKey(publicKeyPEM string) (ed25519.PublicKey, string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, "", fmt.Errorf("%w: public key is not a PEM block", ErrInvalidOfflineOptions)
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, "", fmt.Errorf("%w: public key is not an SPKI public key", ErrInvalidOfflineOptions)
	}
	key, ok := parsed.(ed25519.PublicKey)
	if !ok {
		return nil, "", fmt.Errorf("%w: public key is not an Ed25519 key", ErrInvalidOfflineOptions)
	}
	sum := sha256.Sum256(block.Bytes)
	return key, base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

// decodeOfflineSegment decodes a base64url JSON segment that must be an
// object. Numbers are kept as json.Number so their type can be checked.
func decodeOfflineSegment(segment string) (map[string]any, bool) {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return nil, false
	}
	return decodeOfflineObject(raw)
}

func decodeOfflineObject(raw []byte) (map[string]any, bool) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var value any
	if err := dec.Decode(&value); err != nil {
		return nil, false
	}
	// Reject anything but whitespace after the value, as JSON.parse does.
	if len(bytes.Trim(raw[dec.InputOffset():], " \t\r\n")) != 0 {
		return nil, false
	}
	obj, ok := value.(map[string]any)
	return obj, ok
}

// offlineNumber returns v as a finite float64 if it is a JSON number.
func offlineNumber(v any) (float64, bool) {
	n, ok := v.(json.Number)
	if !ok {
		return 0, false
	}
	f, err := strconv.ParseFloat(string(n), 64)
	if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
		return 0, false
	}
	return f, true
}

func unixSecondsToTime(sec float64) time.Time {
	// Clamp to what int64 seconds can hold; signed licenses are capped at
	// 10 years, so this only matters for absurd values.
	const limit = float64(math.MaxInt64 / 2)
	sec = math.Max(-limit, math.Min(limit, sec))
	whole, frac := math.Modf(sec)
	return time.Unix(int64(whole), int64(frac*1e9))
}

// VerifyOfflineLicense verifies an offline license locally, with no network
// call. license is the contents of the .lic file from SignKey or the
// dashboard, or the bare signed token. It returns the verified license, an
// *OfflineLicenseError if the license fails verification, or an error
// wrapping ErrInvalidOfflineOptions if opts are invalid.
func VerifyOfflineLicense(license string, opts OfflineVerifyOptions) (*OfflineLicense, error) {
	if len(opts.PublicKeys) == 0 {
		return nil, fmt.Errorf("%w: PublicKeys must contain at least one key", ErrInvalidOfflineOptions)
	}
	if opts.ProductID == "" {
		return nil, fmt.Errorf("%w: ProductID is required", ErrInvalidOfflineOptions)
	}
	if opts.MachineCode == "" {
		return nil, fmt.Errorf("%w: MachineCode is required", ErrInvalidOfflineOptions)
	}
	clockSkew := defaultOfflineClockSkew
	if opts.ClockSkew != nil {
		clockSkew = *opts.ClockSkew
	}
	if clockSkew < 0 || clockSkew > maxOfflineClockSkew {
		return nil, fmt.Errorf("%w: ClockSkew must be between 0 and 1h", ErrInvalidOfflineOptions)
	}
	keys := make(map[string]ed25519.PublicKey, len(opts.PublicKeys))
	for _, p := range opts.PublicKeys {
		key, kid, err := parseOfflinePublicKey(p)
		if err != nil {
			return nil, err
		}
		keys[kid] = key
	}
	fail := func(code OfflineLicenseErrorCode, msg string) (*OfflineLicense, error) {
		return nil, &OfflineLicenseError{Code: code, Message: msg}
	}

	// 1. Parse: a .lic file (JSON) or a bare token.
	token := strings.TrimSpace(license)
	if strings.HasPrefix(token, "{") {
		file, ok := decodeOfflineObject([]byte(token))
		if !ok {
			return fail(OfflineMalformed, "license file is not a JSON object")
		}
		signed, ok := file["signedKey"].(string)
		if !ok {
			return fail(OfflineMalformed, "license file has no signedKey")
		}
		token = signed
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fail(OfflineMalformed, "token must have three segments")
	}
	for _, p := range parts {
		if !offlineSegment.MatchString(p) {
			return fail(OfflineMalformed, "token segment is not base64url")
		}
	}
	headerB64, payloadB64, signatureB64 := parts[0], parts[1], parts[2]
	header, okHeader := decodeOfflineSegment(headerB64)
	payload, okPayload := decodeOfflineSegment(payloadB64)
	if !okHeader || !okPayload {
		return fail(OfflineMalformed, "token header or payload is not a JSON object")
	}

	// 2. Algorithm: exactly EdDSA, no critical extensions.
	if alg, _ := header["alg"].(string); alg != "EdDSA" {
		return fail(OfflineUnsupportedAlgorithm, "token algorithm must be EdDSA")
	}
	if _, hasCrit := header["crit"]; hasCrit {
		return fail(OfflineUnsupportedAlgorithm, "token has unsupported critical header")
	}

	// 3. Key: the kid must name a trusted key.
	kid, _ := header["kid"].(string)
	key, ok := keys[kid]
	if !ok {
		return fail(OfflineUnknownKey, "token was not signed by a trusted key")
	}

	// 4. Signature over the segments exactly as received.
	signature, err := base64.RawURLEncoding.DecodeString(signatureB64)
	if err != nil || len(signature) != ed25519.SignatureSize ||
		!ed25519.Verify(key, []byte(headerB64+"."+payloadB64), signature) {
		return fail(OfflineBadSignature, "signature is invalid")
	}

	// 5-7. Claims, compared exactly.
	if typ, _ := payload["type"].(string); typ != "offline" {
		return fail(OfflineWrongType, "token is not an offline license")
	}
	if productID, ok := payload["productId"].(string); !ok || productID != opts.ProductID {
		return fail(OfflineWrongProduct, "license is for a different product")
	}
	if machineCode, ok := payload["machineCode"].(string); !ok || machineCode != opts.MachineCode {
		return fail(OfflineWrongMachine, "license is for a different machine")
	}

	// 8. Shape of the remaining claims, then 9. expiry.
	exp, ok := offlineNumber(payload["exp"])
	if !ok {
		return fail(OfflineMalformed, "license has no numeric exp")
	}
	sub, ok := payload["sub"].(string)
	if !ok || sub == "" {
		return fail(OfflineMalformed, "license has no sub")
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}
	nowSeconds := float64(now.Unix()) + float64(now.Nanosecond())/1e9
	if !(nowSeconds < exp+clockSkew.Seconds()) {
		return fail(OfflineExpired, "license has expired")
	}

	result := &OfflineLicense{
		LicenseKey:  sub,
		ProductID:   opts.ProductID,
		MachineCode: opts.MachineCode,
		KeyID:       kid,
		ExpiresAt:   unixSecondsToTime(exp),
	}
	if iat, ok := offlineNumber(payload["iat"]); ok {
		issuedAt := unixSecondsToTime(iat)
		result.IssuedAt = &issuedAt
	}
	return result, nil
}
