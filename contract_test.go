package keymint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	c, err := New("test-key", srv.URL)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	return c, srv.Close
}

func TestSignKeyDeserializesFileStringResponse(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/key/sign" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"file": `{"signedDate":"2026-09-09","signedKey":"jwt.here","keyId":"kid","publicKeyFingerprint":"kid"}`,
		})
	})
	defer done()

	res, err := c.SignKey(SignKeyParams{ProductID: "p", LicenseKey: "k", HostID: "h"})
	if err != nil {
		t.Fatalf("SignKey failed: %v", err)
	}
	file, ok := res.File.(string)
	if !ok || !strings.Contains(file, "signedKey") || !strings.Contains(file, "keyId") {
		t.Fatalf("unexpected file payload: %#v", res.File)
	}
}

func TestGetKeySendsLicenseInHeaderInsteadOfURL(t *testing.T) {
	var gotHeader, gotQuery string
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("x-license-key")
		gotQuery = r.URL.Query().Get("licenseKey")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 0,
			"data": map[string]interface{}{
				"license": map[string]interface{}{
					"id": "l1", "key": "AAAA", "productId": "p",
					"maxActivations": 2, "activations": 0,
					"devices": []interface{}{}, "activated": false,
				},
			},
		})
	})
	defer done()

	res, err := c.GetKey(GetKeyParams{ProductID: "p", LicenseKey: "AAAA"})
	if err != nil {
		t.Fatalf("GetKey failed: %v", err)
	}
	if gotHeader != "AAAA" {
		t.Fatalf("license key not sent in header, got %q", gotHeader)
	}
	if gotQuery != "" {
		t.Fatalf("license key leaked into query string: %q", gotQuery)
	}
	if res.Data.License.ProductID != "p" {
		t.Fatalf("unexpected product %q", res.Data.License.ProductID)
	}
}

func TestActivateKeyDeserializesMetadataAndVersion(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 0, "message": "License valid",
			"metadata": map[string]interface{}{
				"hostId":   "h1",
				"features": map[string]interface{}{"Pro Mode": true},
			},
			"versionId": "v1",
			"version":   map[string]interface{}{"version": "1.2.3"},
		})
	})
	defer done()

	host := "h1"
	res, err := c.ActivateKey(ActivateKeyParams{ProductID: "p", LicenseKey: "k", HostID: &host})
	if err != nil {
		t.Fatalf("ActivateKey failed: %v", err)
	}
	if res.Metadata["hostId"] != "h1" {
		t.Fatalf("metadata not parsed: %#v", res.Metadata)
	}
	if res.Version["version"] != "1.2.3" {
		t.Fatalf("version not parsed: %#v", res.Version)
	}
}

func TestCreateKeyDeserializesBulkKeys(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 0, "key": "AAAA", "keys": []string{"AAAA", "BBBB"},
		})
	})
	defer done()

	res, err := c.CreateKey(CreateKeyParams{ProductID: "p"})
	if err != nil {
		t.Fatalf("CreateKey failed: %v", err)
	}
	if res.Key != "AAAA" || len(res.Keys) != 2 || res.Keys[1] != "BBBB" {
		t.Fatalf("bulk keys not parsed: %#v", res)
	}
}

func TestFloatingCheckoutSerializesRenewalCredentials(t *testing.T) {
	var body map[string]interface{}
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 0, "message": "ok", "sessionId": "s",
			"sessionSecret": "sec", "nextNonce": "n",
			"expiresAt": "x", "heartbeatInterval": 60,
		})
	})
	defer done()

	ts := interface{}("nonce-1")
	sig := "sig-1"
	_, err := c.FloatingCheckout(FloatingCheckoutParams{
		ProductID: "p", LicenseKey: "k", HostID: "h",
		Timestamp: ts, Signature: &sig,
	})
	if err != nil {
		t.Fatalf("FloatingCheckout failed: %v", err)
	}
	if body["timestamp"] != "nonce-1" || body["signature"] != "sig-1" {
		t.Fatalf("renewal credentials not serialized: %#v", body)
	}
}

func TestApiErrorsExposeNestedMessage(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message":"Request failed","code":1,"error":{"code":"RATE_LIMIT_EXCEEDED","message":"Too many requests. Please try again later."}}`))
	})
	defer done()

	_, err := c.CreateKey(CreateKeyParams{ProductID: "p"})
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if !strings.Contains(apiErr.Message, "Too many requests") {
		t.Fatalf("nested message not surfaced: %q", apiErr.Message)
	}
	if apiErr.Status == nil || *apiErr.Status != http.StatusTooManyRequests {
		t.Fatalf("status not captured: %#v", apiErr.Status)
	}
}

// expirationCases covers the three shapes the API can send: a date, an
// explicit null (never expires), and an older server that omits the field.
var expirationCases = []struct {
	name string
	body string
	want *string
}{
	{"date", `"expirationDate":"2027-03-01T00:00:00.000Z",`, String("2027-03-01T00:00:00.000Z")},
	{"null", `"expirationDate":null,`, nil},
	{"absent", ``, nil},
}

func checkExpiration(t *testing.T, got, want *string) {
	t.Helper()
	switch {
	case want == nil && got != nil:
		t.Fatalf("expected nil ExpirationDate, got %q", *got)
	case want != nil && got == nil:
		t.Fatalf("expected ExpirationDate %q, got nil", *want)
	case want != nil && *got != *want:
		t.Fatalf("expected ExpirationDate %q, got %q", *want, *got)
	}
}

func TestActivateKeyDeserializesExpirationDate(t *testing.T) {
	for _, tc := range expirationCases {
		t.Run(tc.name, func(t *testing.T) {
			c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"code":0,"message":"License valid",` + tc.body + `"metadata":{"hostId":"h1"}}`))
			})
			defer done()

			res, err := c.ActivateKey(ActivateKeyParams{ProductID: "p", LicenseKey: "k", HostID: String("h1")})
			if err != nil {
				t.Fatalf("ActivateKey failed: %v", err)
			}
			checkExpiration(t, res.ExpirationDate, tc.want)
		})
	}
}

func TestFloatingCheckoutDeserializesExpirationDate(t *testing.T) {
	for _, tc := range expirationCases {
		t.Run(tc.name, func(t *testing.T) {
			c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"code":0,"message":"ok","sessionId":"s","sessionSecret":"sec","nextNonce":"n",` +
					`"expiresAt":"2026-10-26T12:00:00.000Z","heartbeatInterval":60,` + tc.body + `"currentSessions":1,"maxSessions":5}`))
			})
			defer done()

			res, err := c.FloatingCheckout(FloatingCheckoutParams{ProductID: "p", LicenseKey: "k", HostID: "h"})
			if err != nil {
				t.Fatalf("FloatingCheckout failed: %v", err)
			}
			if res.ExpiresAt != "2026-10-26T12:00:00.000Z" {
				t.Fatalf("seat lease ExpiresAt changed: %q", res.ExpiresAt)
			}
			checkExpiration(t, res.ExpirationDate, tc.want)
		})
	}
}

func TestDeleteKeySendsLicenseInBodyInsteadOfURL(t *testing.T) {
	var gotMethod, gotPath, gotRawQuery, gotURI, gotAuth, gotContentType string
	var body map[string]interface{}
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotRawQuery = r.URL.RawQuery
		gotURI = r.RequestURI
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"Key deleted","code":0}`))
	})
	defer done()

	res, err := c.DeleteKey(DeleteKeyParams{ProductID: "p", LicenseKey: "AAAA-BBBB"})
	if err != nil {
		t.Fatalf("DeleteKey failed: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Fatalf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/key" {
		t.Fatalf("unexpected path %s", gotPath)
	}
	if gotRawQuery != "" || strings.Contains(gotURI, "AAAA-BBBB") {
		t.Fatalf("license key leaked into URL: %q", gotURI)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("unexpected Authorization header %q", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Fatalf("unexpected Content-Type %q", gotContentType)
	}
	if body["productId"] != "p" || body["licenseKey"] != "AAAA-BBBB" {
		t.Fatalf("unexpected body: %#v", body)
	}
	if res.Message != "Key deleted" || res.Code != 0 {
		t.Fatalf("unexpected response: %#v", res)
	}
}

func TestDeleteKeyReturnsNotFoundError(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Invalid license key","code":1}`))
	})
	defer done()

	_, err := c.DeleteKey(DeleteKeyParams{ProductID: "p", LicenseKey: "missing"})
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T", err)
	}
	if apiErr.Message != "Invalid license key" || apiErr.Code != 1 {
		t.Fatalf("unexpected error: %#v", apiErr)
	}
	if apiErr.Status == nil || *apiErr.Status != http.StatusNotFound {
		t.Fatalf("status not captured: %#v", apiErr.Status)
	}
}
