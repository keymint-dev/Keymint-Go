package keymint
//go test -v to test ./...

import (
	keymint "KeymintGoSdk/src"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	integrationLicenseKey string
	integrationCustomerID string
	testInitOnce sync.Once
)

// Helper to create a unique customer and license key for all tests
func integrationSetup(t *testing.T) (customerID, licenseKey string) {
	testInitOnce.Do(func() {
		accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
		productId := os.Getenv("KEYMINT_PRODUCT_ID")
		if accessToken == "" || productId == "" {
			t.Skip("Set KEYMINT_ACCESS_TOKEN and KEYMINT_PRODUCT_ID in your environment to run integration tests.")
		}
		client, err := keymint.New(accessToken, "")
		if err != nil {
			t.Fatalf("Failed to initialize client: %v", err)
		}
		// Generate unique email
		rand.Seed(time.Now().UnixNano())
		email := fmt.Sprintf("integration-customer-%d@go.com", rand.Intn(1e9))
		params := keymint.CreateCustomerParams{Name: "Go Integration Customer", Email: email}
		resp, err := client.CreateCustomer(params)
		if err != nil {
			t.Fatalf("CreateCustomer failed: %v", err)
		}
		integrationCustomerID = resp.Data.ID
		// Create key for this customer
		keyParams := keymint.CreateKeyParams{ProductID: productId, CustomerID: &integrationCustomerID}
		keyResp, err := client.CreateKey(keyParams)
		if err != nil {
			t.Fatalf("CreateKey failed: %v", err)
		}
		integrationLicenseKey = keyResp.Key
	})
	if integrationCustomerID == "" || integrationLicenseKey == "" {
		t.Fatal("Failed to generate integration customer and license key")
	}
	return integrationCustomerID, integrationLicenseKey
}

func TestIntegration_CreateCustomer(t *testing.T) {
	_, _ = integrationSetup(t)
	// If setup passes, customer creation is tested.
}

func TestIntegration_CreateKey(t *testing.T) {
	customerID, _ := integrationSetup(t)
	accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
	productId := os.Getenv("KEYMINT_PRODUCT_ID")
	client, err := keymint.New(accessToken, "")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	keyParams := keymint.CreateKeyParams{ProductID: productId, CustomerID: &customerID}
	key, err := client.CreateKey(keyParams)
	if err != nil {
		t.Fatalf("CreateKey failed: %v", err)
	}
	if key.Key == "" {
		t.Error("Expected a license key, got empty string")
	}
}

func TestIntegration_GetAllCustomers(t *testing.T) {
	accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
	client, err := keymint.New(accessToken, "")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	_, err = client.GetAllCustomers()
	if err != nil {
		t.Fatalf("GetAllCustomers failed: %v", err)
	}
}

func TestIntegration_ActivateKey(t *testing.T) {
	_, licenseKey := integrationSetup(t)
	accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
	productId := os.Getenv("KEYMINT_PRODUCT_ID")
	client, err := keymint.New(accessToken, "")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	_, err = client.ActivateKey(keymint.ActivateKeyParams{ProductID: productId, LicenseKey: licenseKey})
	if err != nil {
		t.Fatalf("ActivateKey failed: %v", err)
	}
}

func TestIntegration_DeactivateKey(t *testing.T) {
	_, licenseKey := integrationSetup(t)
	accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
	productId := os.Getenv("KEYMINT_PRODUCT_ID")
	client, err := keymint.New(accessToken, "")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	_, err = client.DeactivateKey(keymint.DeactivateKeyParams{ProductID: productId, LicenseKey: licenseKey})
	if err != nil {
		t.Fatalf("DeactivateKey failed: %v", err)
	}
}

func TestIntegration_GetKey(t *testing.T) {
	_, licenseKey := integrationSetup(t)
	accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
	productId := os.Getenv("KEYMINT_PRODUCT_ID")
	client, err := keymint.New(accessToken, "")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	_, err = client.GetKey(keymint.GetKeyParams{ProductID: productId, LicenseKey: licenseKey})
	if err != nil {
		t.Fatalf("GetKey failed: %v", err)
	}
}

// Add similar integration tests for other endpoints as needed.