package keymint

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func liveClient(t *testing.T, keyEnv string) *Client {
	t.Helper()
	key := os.Getenv(keyEnv)
	base := os.Getenv("KEYMINT_TEST_BASE_URL")
	if base == "" {
		base = "https://api.keymint.dev"
	}
	if key == "" || os.Getenv("KEYMINT_TEST_PRODUCT_ID") == "" {
		t.Skip("live credentials not set (KEYMINT_TEST_*); skipping live suite")
	}
	c, err := New(key, base)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	return c
}

func isRateLimited(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Too many requests")
}

// live calls fn up to 4 times on 429s (shared test workspaces are rate-limited).
func live[T any](t *testing.T, what string, fn func() (T, error)) T {
	t.Helper()
	var err error
	var out T
	for attempt := 1; attempt <= 4; attempt++ {
		out, err = fn()
		if err == nil {
			return out
		}
		if !isRateLimited(err) || attempt == 4 {
			t.Fatalf("%s failed: %v", what, err)
		}
		time.Sleep(10 * time.Second)
	}
	t.Fatalf("%s failed: %v", what, err)
	var zero T
	return zero
}

func strptr(s string) *string { return &s }

func TestLiveLicenseWorkflowsEndToEnd(t *testing.T) {
	admin := liveClient(t, "KEYMINT_TEST_ADMIN_API_KEY")
	readOnly := liveClient(t, "KEYMINT_TEST_READONLY_API_KEY")
	client := liveClient(t, "KEYMINT_TEST_CLIENT_API_KEY")
	productID := os.Getenv("KEYMINT_TEST_PRODUCT_ID")
	runID := strings.ReplaceAll(uuid.New().String(), "-", "")
	nodeHost := "ci-go-node-" + runID
	var nodeKey, floatingKey, customerID string

	defer func() {
		if nodeKey != "" {
			_, _ = admin.BlockKey(BlockKeyParams{ProductID: productID, LicenseKey: nodeKey})
		}
		if floatingKey != "" {
			_, _ = admin.BlockKey(BlockKeyParams{ProductID: productID, LicenseKey: floatingKey})
		}
		if customerID != "" {
			_, _ = admin.DeleteCustomer(DeleteCustomerParams{CustomerID: customerID})
		}
	}()

	created := live(t, "create node key", func() (*CreateKeyResponse, error) {
		return admin.CreateKey(CreateKeyParams{
			ProductID:      productID,
			MaxActivations: strptr("2"),
			Metadata:       map[string]interface{}{"purpose": "go-sdk-live-test", "runId": runID},
		})
	})
	nodeKey = created.Key
	if nodeKey == "" {
		t.Fatal("empty license key returned")
	}

	lookup := live(t, "readonly lookup", func() (*GetKeyResponse, error) {
		return readOnly.GetKey(GetKeyParams{ProductID: productID, LicenseKey: nodeKey})
	})
	if lookup.Data.License.ProductID != productID {
		t.Fatalf("product mismatch: %q", lookup.Data.License.ProductID)
	}

	activation := live(t, "activate", func() (*ActivateKeyResponse, error) {
		return client.ActivateKey(ActivateKeyParams{
			ProductID: productID, LicenseKey: nodeKey,
			HostID: &nodeHost, DeviceTag: strptr("Go SDK CI"),
		})
	})
	if activation.Code != 0 {
		t.Fatalf("activation code %d", activation.Code)
	}
	if host, _ := activation.Metadata["hostId"].(string); host != nodeHost {
		t.Fatalf("hostId mismatch: %#v", activation.Metadata["hostId"])
	}

	deactivated := live(t, "deactivate", func() (*DeactivateKeyResponse, error) {
		return client.DeactivateKey(DeactivateKeyParams{
			ProductID: productID, LicenseKey: nodeKey, HostID: &nodeHost,
		})
	})
	if deactivated.DevicesRemoved == nil || *deactivated.DevicesRemoved != 1 {
		t.Fatalf("unexpected devicesRemoved: %#v", deactivated.DevicesRemoved)
	}

	live(t, "update", func() (*UpdateKeyResponse, error) {
		return admin.UpdateKey(UpdateKeyParams{
			ProductID: productID, LicenseKey: nodeKey, MaxActivations: 3,
		})
	})
	live(t, "block", func() (*BlockKeyResponse, error) {
		return admin.BlockKey(BlockKeyParams{ProductID: productID, LicenseKey: nodeKey})
	})
	live(t, "unblock", func() (*UnblockKeyResponse, error) {
		return admin.UnblockKey(UnblockKeyParams{ProductID: productID, LicenseKey: nodeKey})
	})

	maxSessions, hb, lease := 1, 60, 300
	floating := live(t, "create floating key", func() (*CreateKeyResponse, error) {
		return admin.CreateKey(CreateKeyParams{
			ProductID: productID, LicenseType: strptr("floating"),
			MaxConcurrentSessions: &maxSessions, HeartbeatInterval: &hb,
			SessionLeaseDuration: &lease,
			Metadata:             map[string]interface{}{"purpose": "go-sdk-floating-live-test", "runId": runID},
		})
	})
	floatingKey = floating.Key

	checkout := live(t, "floating checkout", func() (*FloatingCheckoutResponse, error) {
		return client.FloatingCheckout(FloatingCheckoutParams{
			ProductID: productID, LicenseKey: floatingKey,
			HostID: "ci-go-floating-" + runID, DeviceTag: strptr("Go SDK CI floating"),
		})
	})
	if checkout.SessionSecret == "" {
		t.Fatal("empty session secret")
	}

	heartbeat := live(t, "floating heartbeat", func() (*FloatingHeartbeatResponse, error) {
		return client.FloatingHeartbeat(FloatingHeartbeatParams{
			ProductID: productID, LicenseKey: floatingKey,
			SessionID: checkout.SessionID, Timestamp: checkout.NextNonce,
			Signature: GenerateSessionSignature(checkout.SessionID, checkout.NextNonce, checkout.SessionSecret),
		})
	})

	live(t, "floating checkin", func() (*FloatingCheckinResponse, error) {
		return client.FloatingCheckin(FloatingCheckinParams{
			ProductID: productID, LicenseKey: floatingKey,
			SessionID: checkout.SessionID, Timestamp: heartbeat.NextNonce,
			Signature: GenerateSessionSignature(checkout.SessionID, heartbeat.NextNonce, checkout.SessionSecret),
		})
	})

	ttl := 300
	signed := live(t, "sign offline", func() (*SignKeyResponse, error) {
		return admin.SignKey(SignKeyParams{
			ProductID: productID, LicenseKey: nodeKey, HostID: nodeHost, TTL: &ttl,
		})
	})
	fileStr, ok := signed.File.(string)
	if !ok || !strings.Contains(fileStr, "signedKey") || !strings.Contains(fileStr, "keyId") {
		t.Fatalf("unexpected signed file: %#v", signed.File)
	}

	email := "ci-go-" + runID + "@example.com"
	createdCustomer := live(t, "create customer", func() (*CreateCustomerResponse, error) {
		return admin.CreateCustomer(CreateCustomerParams{Name: "Go SDK CI", Email: email})
	})
	customerID = createdCustomer.Data.ID

	fetched := live(t, "get customer", func() (*GetCustomerByIdResponse, error) {
		return admin.GetCustomerById(GetCustomerByIdParams{CustomerID: customerID})
	})
	found := false
	for _, c := range fetched.Data {
		if c.ID == customerID {
			found = true
		}
	}
	if !found {
		t.Fatal("created customer missing from get-by-id")
	}

	updated := live(t, "update customer", func() (*UpdateCustomerResponse, error) {
		return admin.UpdateCustomer(UpdateCustomerParams{CustomerID: customerID, Name: strptr("Go SDK CI Updated")})
	})
	if updated.Data.Name != "Go SDK CI Updated" {
		t.Fatalf("unexpected updated customer: %#v", updated.Data)
	}

	all := live(t, "list customers", func() (*GetAllCustomersResponse, error) {
		return admin.GetAllCustomers(GetAllCustomersParams{})
	})
	found = false
	for _, c := range all.Data {
		if c.ID == customerID {
			found = true
		}
	}
	if !found {
		t.Fatal("created customer missing from list")
	}

	withKeys := live(t, "customer keys", func() ([]CustomerLicenseKey, error) {
		return admin.GetCustomerWithKeys(GetCustomerWithKeysParams{CustomerID: customerID})
	})
	_ = withKeys

	toggledOff := live(t, "disable customer", func() (*ToggleCustomerStatusResponse, error) {
		return admin.ToggleCustomerStatus(ToggleCustomerStatusParams{CustomerID: customerID})
	})
	if !toggledOff.Status {
		t.Fatal("disable reported failure")
	}
	toggledOn := live(t, "enable customer", func() (*ToggleCustomerStatusResponse, error) {
		return admin.ToggleCustomerStatus(ToggleCustomerStatusParams{CustomerID: customerID})
	})
	if !toggledOn.Status {
		t.Fatal("enable reported failure")
	}

	deleted := live(t, "delete customer", func() (*DeleteCustomerResponse, error) {
		return admin.DeleteCustomer(DeleteCustomerParams{CustomerID: customerID})
	})
	if !deleted.Status {
		t.Fatal("delete reported failure")
	}
	customerID = ""
}
