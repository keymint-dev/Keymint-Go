package keymint
//go test -v to test ./...

import (
	keymint "KeymintGoSdk/src"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCustomer_IntegrationStyle(t *testing.T) {
	fmt.Println("KeyMint SDK Test")
	fmt.Println("================")

	// Create a test server that simulates the KeyMint API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/customer":
			// Handle customer creation
			var params keymint.CreateCustomerParams
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(keymint.ApiError{
					Message: "Invalid request body",
					Code:    400,
				})
				return
			}

			// Simulate successful customer creation
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(keymint.CreateCustomerResponse{
				Action:  "createCustomer",
				Status:  true,
				Message: "Customer created successfully",
				Data: struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Email string `json:"email"`
				}{
					ID:    "12345",
					Name:  params.Name,
					Email: params.Email,
				},
				Code: 0,
			})

		case "/key":
			// Handle key creation (for completeness)
			var params keymint.CreateKeyParams
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(keymint.ApiError{
					Message: "Invalid request body",
					Code:    400,
				})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(keymint.CreateKeyResponse{
				Code: 0,
				Key:  "LICENSE-KEY-12345",
			})

		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(keymint.ApiError{
				Message: "Endpoint not found",
				Code:    404,
			})
		}
	}))
	defer server.Close()

	// Initialize SDK
	sdk, err := keymint.NewSDK("at_CNNOGtg3ZqclCRXK0QQZutYyTyWwSSLflwAJx82WCPU28b13e3aeeb54fac9f43fae61d47dde29", server.URL)
	if err != nil {
		t.Fatalf("Failed to initialize SDK: %v", err)
	}
	fmt.Println("✓ SDK initialized successfully")

	productId := "f6736ea441ac4ca3a52a84"

	// Create customer
	customerParams := keymint.CreateCustomerParams{
		Name:  "Go customer",
		Email: "customer@go.com",
	}

	fmt.Println("Creating customer...")
	customer, err := sdk.CreateCustomer(customerParams)
	if err != nil {
		t.Fatalf("✗ Customer creation failed: %v", err)
	}
	fmt.Println("✓ Customer created successfully")
	fmt.Printf("Customer ID: %s\n", customer.Data.ID)

	// Use the customer ID for key creation
	customerId := customer.Data.ID

	// Option 1: Create key with existing customer ID
	if customerId != "" {
		fmt.Printf("Creating key with customer ID: %s\n", customerId)
		
		keyParams := keymint.CreateKeyParams{
			ProductID:     productId,
			CustomerID:    &customerId,
			MaxActivations: func() *string { s := "5"; return &s }(),
		}

		key, err := sdk.CreateKey(keyParams)
		if err != nil {
			fmt.Printf("✗ Key creation with customer ID failed: %v\n", err)
		} else {
			fmt.Println("✓ Key created successfully with existing customer")
			fmt.Printf("License Key: %s\n", key.Key)
			fmt.Printf("Response Code: %d\n", key.Code)
		}
	}

	// Option 2: Create key with new customer
	fmt.Println("\nTrying to create key with new customer...")

	newCustomer := keymint.NewCustomer{
		Name:  "Go Customer 2",
		Email: func() *string { s := "customer2@go.com"; return &s }(),
	}

	keyParams2 := keymint.CreateKeyParams{
		ProductID:     productId,
		NewCustomer:   &newCustomer,
		MaxActivations: func() *string { s := "3"; return &s }(),
	}

	// Debug: Show what we're sending
	fmt.Println("Key params being sent:")
	fmt.Printf("  productId: %s\n", keyParams2.ProductID)
	fmt.Printf("  maxActivations: %s\n", *keyParams2.MaxActivations)
	if keyParams2.NewCustomer != nil {
		fmt.Printf("  newCustomer.name: %s\n", keyParams2.NewCustomer.Name)
		if keyParams2.NewCustomer.Email != nil {
			fmt.Printf("  newCustomer.email: %s\n", *keyParams2.NewCustomer.Email)
		}
	}

	key2, err := sdk.CreateKey(keyParams2)
	if err != nil {
		fmt.Printf("✗ Key creation with new customer failed: %v\n", err)
		// Check if it's an ApiError
		if apiErr, ok := err.(*keymint.ApiError); ok {
			fmt.Printf("Status Code: %v\n", apiErr.Status)
			fmt.Printf("API Code: %d\n", apiErr.Code)
		}
	} else {
		fmt.Println("✓ Key created successfully with new customer")
		fmt.Printf("License Key: %s\n", key2.Key)
		fmt.Printf("Response Code: %d\n", key2.Code)
	}

	fmt.Println("Test completed!")
}

// Additional test for error scenarios
func TestCreateCustomer_ErrorCases(t *testing.T) {
	// Test server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(keymint.ApiError{
			Message: "Email already exists",
			Code:    1001,
		})
	}))
	defer server.Close()

	sdk, err := keymint.NewSDK("test-token", server.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Test duplicate email error
	_, err = sdk.CreateCustomer(keymint.CreateCustomerParams{
		Name:  "Test Customer",
		Email: "duplicate@mail.com",
	})

	if err == nil {
		t.Error("Expected an error for duplicate email, but got none")
	}

	// Verify it's an ApiError
	if apiErr, ok := err.(*keymint.ApiError); ok {
		if apiErr.Message != "Email already exists" {
			t.Errorf("Expected error message 'Email already exists', got '%s'", apiErr.Message)
		}
		if apiErr.Code != 1001 {
			t.Errorf("Expected error code 1001, got %d", apiErr.Code)
		}
	} else {
		t.Errorf("Expected ApiError, got %T", err)
	}
}

// Test for invalid token
func TestCreateCustomer_InvalidToken(t *testing.T) {
	_, err := keymint.NewSDK("", "https://api.keymint.dev")
	if err == nil {
		t.Error("Expected error for empty token, but got none")
	}

	expectedError := "access token is required to initialize the SDK"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}