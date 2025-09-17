package keymint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the main entry point for the KeyMint API client
// Client provides methods to interact with the KeyMint API for license and customer management.
type Client struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

// New creates a new KeyMint API client instance.
// accessToken: Your KeyMint API access token (required).
// baseURL: Optional API base URL (defaults to https://api.keymint.dev).
// Returns a new Client instance or an error if accessToken is missing.
func New(accessToken string, baseURL string) (*Client, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required to initialize the client")
	}

	if baseURL == "" {
		baseURL = "https://api.keymint.dev"
	}

	return &Client{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// handleRequest is a generic method to handle POST/PUT requests.
// method: HTTP method (POST/PUT).
// endpoint: API endpoint.
// params: Request body parameters.
// result: Pointer to the result struct to unmarshal response into.
// Returns an error if the request fails or the API returns an error.
func (c *Client) handleRequest(method, endpoint string, params interface{}, result interface{}) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to marshal request: %v", err),
			Code:    -1,
		}
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to create request: %v", err),
			Code:    -1,
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
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

// handleGetRequest is a generic method to handle GET requests.
// endpoint: API endpoint.
// queryParams: Query parameters as a map.
// result: Pointer to the result struct to unmarshal response into.
// Returns an error if the request fails or the API returns an error.
func (c *Client) handleGetRequest(endpoint string, queryParams map[string]string, result interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
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

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
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

// handleDeleteRequest is a generic method to handle DELETE requests.
// endpoint: API endpoint.
// queryParams: Query parameters as a map.
// result: Pointer to the result struct to unmarshal response into.
// Returns an error if the request fails or the API returns an error.
func (c *Client) handleDeleteRequest(endpoint string, queryParams map[string]string, result interface{}) error {
	req, err := http.NewRequest("DELETE", c.baseURL+endpoint, nil)
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

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
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

// CreateKey creates a new license key.
// params: Parameters for creating the key.
// Returns the created key information or an error.
func (c *Client) CreateKey(params CreateKeyParams) (*CreateKeyResponse, error) {
	var result CreateKeyResponse
	err := c.handleRequest("POST", "/key", params, &result)
	return &result, err
}

// ActivateKey activates a license key for a specific device.
// params: Parameters for activating the key.
// Returns the activation status or an error.
func (c *Client) ActivateKey(params ActivateKeyParams) (*ActivateKeyResponse, error) {
	var result ActivateKeyResponse
	err := c.handleRequest("POST", "/key/activate", params, &result)
	return &result, err
}

// DeactivateKey deactivates a device from a license key.
// params: Parameters for deactivating the key.
// Returns the deactivation confirmation or an error.
func (c *Client) DeactivateKey(params DeactivateKeyParams) (*DeactivateKeyResponse, error) {
	var result DeactivateKeyResponse
	err := c.handleRequest("POST", "/key/deactivate", params, &result)
	return &result, err
}

// GetKey retrieves detailed information about a specific license key.
// params: Parameters for fetching the key details.
// Returns the license key details or an error.
func (c *Client) GetKey(params GetKeyParams) (*GetKeyResponse, error) {
	var result GetKeyResponse
	queryParams := map[string]string{
		"productId":  params.ProductID,
		"licenseKey": params.LicenseKey,
	}
	err := c.handleGetRequest("/key", queryParams, &result)
	return &result, err
}

// BlockKey blocks a specific license key.
// params: Parameters for blocking the key.
// Returns the block confirmation or an error.
func (c *Client) BlockKey(params BlockKeyParams) (*BlockKeyResponse, error) {
	var result BlockKeyResponse
	err := c.handleRequest("POST", "/key/block", params, &result)
	return &result, err
}

// UnblockKey unblocks a previously blocked license key.
// params: Parameters for unblocking the key.
// Returns the unblock confirmation or an error.
func (c *Client) UnblockKey(params UnblockKeyParams) (*UnblockKeyResponse, error) {
	var result UnblockKeyResponse
	err := c.handleRequest("POST", "/key/unblock", params, &result)
	return &result, err
}

// CreateCustomer creates a new customer.
// params: Parameters for creating the customer.
// Returns the created customer information or an error.
func (c *Client) CreateCustomer(params CreateCustomerParams) (*CreateCustomerResponse, error) {
	var result CreateCustomerResponse
	err := c.handleRequest("POST", "/customer", params, &result)
	return &result, err
}

// GetAllCustomers retrieves all customers.
// Returns a list of all customers or an error.
func (c *Client) GetAllCustomers() (*GetAllCustomersResponse, error) {
	var result GetAllCustomersResponse
	err := c.handleGetRequest("/customer", nil, &result)
	return &result, err
}

// GetCustomerWithKeys retrieves a customer along with their associated license keys.
// params: Parameters containing the customer ID.
// Returns the customer information with associated license keys or an error.
func (c *Client) GetCustomerWithKeys(params GetCustomerWithKeysParams) (*GetCustomerWithKeysResponse, error) {
	var result GetCustomerWithKeysResponse
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := c.handleGetRequest("/customer/keys", queryParams, &result)
	return &result, err
}

// UpdateCustomer updates an existing customer.
// params: Parameters for updating the customer.
// Returns the update confirmation or an error.
func (c *Client) UpdateCustomer(params UpdateCustomerParams) (*UpdateCustomerResponse, error) {
	var result UpdateCustomerResponse
	err := c.handleRequest("PUT", "/customer/by-id", params, &result)
	return &result, err
}

// DeleteCustomer deletes a customer and all associated license keys permanently.
// params: Parameters containing the customer ID.
// Returns the deletion confirmation or an error.
func (c *Client) DeleteCustomer(params DeleteCustomerParams) (*DeleteCustomerResponse, error) {
	var result DeleteCustomerResponse
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := c.handleDeleteRequest("/customer/by-id", queryParams, &result)
	return &result, err
}

// ToggleCustomerStatus toggles the status of a customer (active/inactive).
// params: Parameters containing the customer ID.
// Returns the status toggle confirmation or an error.
func (c *Client) ToggleCustomerStatus(params ToggleCustomerStatusParams) (*ToggleCustomerStatusResponse, error) {
	var result ToggleCustomerStatusResponse
	err := c.handleRequest("POST", "/customer/disable", params, &result)
	return &result, err
}

// GetCustomerById retrieves detailed information about a specific customer by ID.
// params: Parameters containing the customer ID.
// Returns the customer information or an error.
func (c *Client) GetCustomerById(params GetCustomerByIdParams) (*GetCustomerByIdResponse, error) {
	var result GetCustomerByIdResponse
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := c.handleGetRequest("/customer/by-id", queryParams, &result)
	return &result, err
}
