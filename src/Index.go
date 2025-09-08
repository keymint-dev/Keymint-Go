package keymint


import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SDK is the main entry point for the KeyMint API client
type SDK struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

// NewSDK creates a new KeyMint SDK instance
func NewSDK(accessToken string, baseURL string) (*SDK, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required to initialize the SDK")
	}

	if baseURL == "" {
		baseURL = "https://api.keymint.dev"
	}

	return &SDK{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// handleRequest is a generic method to handle POST/PUT requests
func (s *SDK) handleRequest(method, endpoint string, params interface{}, result interface{}) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to marshal request: %v", err),
			Code:    -1,
		}
	}

	req, err := http.NewRequest(method, s.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to create request: %v", err),
			Code:    -1,
		}
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("request failed: %v", err),
			Code:    -1,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to read response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	if resp.StatusCode >= 400 {
		var apiErr ApiError
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = &resp.StatusCode
			return &apiErr
		}
		return &ApiError{
			Message: fmt.Sprintf("API error: %s", string(body)),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	if err := json.Unmarshal(body, result); err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to unmarshal response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	return nil
}

// handleGetRequest is a generic method to handle GET requests
func (s *SDK) handleGetRequest(endpoint string, queryParams map[string]string, result interface{}) error {
	req, err := http.NewRequest("GET", s.baseURL+endpoint, nil)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to create request: %v", err),
			Code:    -1,
		}
	}

	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("request failed: %v", err),
			Code:    -1,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to read response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	if resp.StatusCode >= 400 {
		var apiErr ApiError
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = &resp.StatusCode
			return &apiErr
		}
		return &ApiError{
			Message: fmt.Sprintf("API error: %s", string(body)),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	if err := json.Unmarshal(body, result); err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to unmarshal response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	return nil
}

// handleDeleteRequest is a generic method to handle DELETE requests
func (s *SDK) handleDeleteRequest(endpoint string, queryParams map[string]string, result interface{}) error {
	req, err := http.NewRequest("DELETE", s.baseURL+endpoint, nil)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to create request: %v", err),
			Code:    -1,
		}
	}

	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("request failed: %v", err),
			Code:    -1,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to read response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	if resp.StatusCode >= 400 {
		var apiErr ApiError
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = &resp.StatusCode
			return &apiErr
		}
		return &ApiError{
			Message: fmt.Sprintf("API error: %s", string(body)),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	if err := json.Unmarshal(body, result); err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to unmarshal response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	return nil
}

// CreateKey creates a new license key
func (s *SDK) CreateKey(params CreateKeyParams) (*CreateKeyResponse, error) {
	var result CreateKeyResponse
	err := s.handleRequest("POST", "/key", params, &result)
	return &result, err
}

// ActivateKey activates a license key for a specific device
func (s *SDK) ActivateKey(params ActivateKeyParams) (*ActivateKeyResponse, error) {
	var result ActivateKeyResponse
	err := s.handleRequest("POST", "/key/activate", params, &result)
	return &result, err
}

// DeactivateKey deactivates a device from a license key
func (s *SDK) DeactivateKey(params DeactivateKeyParams) (*DeactivateKeyResponse, error) {
	var result DeactivateKeyResponse
	err := s.handleRequest("POST", "/key/deactivate", params, &result)
	return &result, err
}

// GetKey retrieves detailed information about a specific license key
func (s *SDK) GetKey(params GetKeyParams) (*GetKeyResponse, error) {
	var result GetKeyResponse
	queryParams := map[string]string{
		"productId":  params.ProductID,
		"licenseKey": params.LicenseKey,
	}
	err := s.handleGetRequest("/key", queryParams, &result)
	return &result, err
}

// BlockKey blocks a specific license key
func (s *SDK) BlockKey(params BlockKeyParams) (*BlockKeyResponse, error) {
	var result BlockKeyResponse
	err := s.handleRequest("POST", "/key/block", params, &result)
	return &result, err
}

// UnblockKey unblocks a previously blocked license key
func (s *SDK) UnblockKey(params UnblockKeyParams) (*UnblockKeyResponse, error) {
	var result UnblockKeyResponse
	err := s.handleRequest("POST", "/key/unblock", params, &result)
	return &result, err
}

// CreateCustomer creates a new customer
func (s *SDK) CreateCustomer(params CreateCustomerParams) (*CreateCustomerResponse, error) {
	var result CreateCustomerResponse
	err := s.handleRequest("POST", "/customer", params, &result)
	return &result, err
}

// GetAllCustomers retrieves all customers
func (s *SDK) GetAllCustomers() (*GetAllCustomersResponse, error) {
	var result GetAllCustomersResponse
	err := s.handleGetRequest("/customer", nil, &result)
	return &result, err
}

// GetCustomerWithKeys retrieves a customer along with their associated license keys
func (s *SDK) GetCustomerWithKeys(params GetCustomerWithKeysParams) (*GetCustomerWithKeysResponse, error) {
	var result GetCustomerWithKeysResponse
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := s.handleGetRequest("/customer/keys", queryParams, &result)
	return &result, err
}

// UpdateCustomer updates an existing customer
func (s *SDK) UpdateCustomer(params UpdateCustomerParams) (*UpdateCustomerResponse, error) {
	var result UpdateCustomerResponse
	err := s.handleRequest("PUT", "/customer/by-id", params, &result)
	return &result, err
}

// DeleteCustomer deletes a customer and all associated license keys permanently
func (s *SDK) DeleteCustomer(params DeleteCustomerParams) (*DeleteCustomerResponse, error) {
	var result DeleteCustomerResponse
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := s.handleDeleteRequest("/customer/by-id", queryParams, &result)
	return &result, err
}

// ToggleCustomerStatus toggles the status of a customer (active/inactive)
func (s *SDK) ToggleCustomerStatus(params ToggleCustomerStatusParams) (*ToggleCustomerStatusResponse, error) {
	var result ToggleCustomerStatusResponse
	err := s.handleRequest("POST", "/customer/disable", params, &result)
	return &result, err
}

// GetCustomerById retrieves detailed information about a specific customer by ID
func (s *SDK) GetCustomerById(params GetCustomerByIdParams) (*GetCustomerByIdResponse, error) {
	var result GetCustomerByIdResponse
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := s.handleGetRequest("/customer/by-id", queryParams, &result)
	return &result, err
}
