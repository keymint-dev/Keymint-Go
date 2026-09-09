package keymint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the main entry point for the KeyMint API client
// Client provides methods to interact with the KeyMint API for license and customer management.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a new KeyMint API client instance.
// apiKey: Your Keymint API key (required).
// baseURL: Optional API base URL (defaults to https://api.keymint.dev).
// Returns a new Client instance or an error if apiKey is missing.
func New(apiKey string, baseURL string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required to initialize the client")
	}

	if baseURL == "" {
		baseURL = "https://api.keymint.dev"
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// parseApiError decodes an error response body, surfacing the nested error
// envelope message when the API provides one.
func parseApiError(body []byte, statusCode int) *ApiError {
	var apiErr ApiError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
		if apiErr.ErrorDetail != nil && apiErr.ErrorDetail.Message != nil && *apiErr.ErrorDetail.Message != "" {
			apiErr.Message = *apiErr.ErrorDetail.Message
		}
		apiErr.Status = &statusCode
		return &apiErr
	}
	return &ApiError{
		Message: fmt.Sprintf("API error: %s", string(body)),
		Code:    -1,
		Status:  &statusCode,
	}
}

// handleRequest is a generic method to handle POST/PUT requests.
// method: HTTP method (POST/PUT).
// endpoint: API endpoint.
// params: Request body parameters.
// result: Pointer to the result struct to unmarshal response into.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns an error if the request fails or the API returns an error.
func (c *Client) handleRequest(method, endpoint string, params interface{}, result interface{}, opts ...*RequestOptions) error {
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

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	if len(opts) > 0 && opts[0] != nil && opts[0].IdempotencyKey != "" {
		req.Header.Set("Idempotency-Key", opts[0].IdempotencyKey)
	}

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
		return parseApiError(body, resp.StatusCode)
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

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
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
		return parseApiError(body, resp.StatusCode)
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
// opts: Optional request configurations (e.g. idempotency keys).
// Returns an error if the request fails or the API returns an error.
func (c *Client) handleDeleteRequest(endpoint string, queryParams map[string]string, result interface{}, opts ...*RequestOptions) error {
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

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	if len(opts) > 0 && opts[0] != nil && opts[0].IdempotencyKey != "" {
		req.Header.Set("Idempotency-Key", opts[0].IdempotencyKey)
	}

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
		return parseApiError(body, resp.StatusCode)
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

// handlePatchRequest is a generic method to handle PATCH requests.
func (c *Client) handlePatchRequest(endpoint string, params interface{}, result interface{}, opts ...*RequestOptions) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to marshal request: %v", err),
			Code:    -1,
		}
	}

	req, err := http.NewRequest("PATCH", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return &ApiError{
			Message: fmt.Sprintf("failed to create request: %v", err),
			Code:    -1,
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	if len(opts) > 0 && opts[0] != nil && opts[0].IdempotencyKey != "" {
		req.Header.Set("Idempotency-Key", opts[0].IdempotencyKey)
	}

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
		return parseApiError(body, resp.StatusCode)
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
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the created key information or an error.
func (c *Client) CreateKey(params CreateKeyParams, opts ...*RequestOptions) (*CreateKeyResponse, error) {
	var result CreateKeyResponse
	err := c.handleRequest("POST", "/key", params, &result, opts...)
	return &result, err
}

// ActivateKey activates a license key for a specific device.
//
// If HostID is omitted, the activation uses the shared "N/A" hostless slot —
// all hostless activations share one reusable device slot per license.
// Pass the returned "N/A" host ID back unchanged for re-activation and
// deactivation.
// params: Parameters for activating the key.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the activation status or an error.
func (c *Client) ActivateKey(params ActivateKeyParams, opts ...*RequestOptions) (*ActivateKeyResponse, error) {
	var result ActivateKeyResponse
	err := c.handleRequest("POST", "/key/activate", params, &result, opts...)
	return &result, err
}

// DeactivateKey deactivates a device from a license key.
// params: Parameters for deactivating the key.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the deactivation confirmation or an error.
func (c *Client) DeactivateKey(params DeactivateKeyParams, opts ...*RequestOptions) (*DeactivateKeyResponse, error) {
	var result DeactivateKeyResponse
	err := c.handleRequest("POST", "/key/deactivate", params, &result, opts...)
	return &result, err
}

// FloatingCheckout checks out a floating license seat.
// params: Parameters for checking out the license.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the checkout response or an error.
func (c *Client) FloatingCheckout(params FloatingCheckoutParams, opts ...*RequestOptions) (*FloatingCheckoutResponse, error) {
	var result FloatingCheckoutResponse
	err := c.handleRequest("POST", "/key/checkout", params, &result, opts...)
	return &result, err
}

// FloatingHeartbeat sends a heartbeat to keep a floating license session alive.
// params: Parameters for the heartbeat.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the heartbeat response or an error.
func (c *Client) FloatingHeartbeat(params FloatingHeartbeatParams, opts ...*RequestOptions) (*FloatingHeartbeatResponse, error) {
	var result FloatingHeartbeatResponse
	err := c.handleRequest("POST", "/key/heartbeat", params, &result, opts...)
	return &result, err
}

// FloatingCheckin checks in a floating license session, releasing the seat.
// params: Parameters for checking in the license.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the checkin response or an error.
func (c *Client) FloatingCheckin(params FloatingCheckinParams, opts ...*RequestOptions) (*FloatingCheckinResponse, error) {
	var result FloatingCheckinResponse
	err := c.handleRequest("POST", "/key/checkin", params, &result, opts...)
	return &result, err
}

// GetKey retrieves detailed information about a specific license key.
// params: Parameters for fetching the key details.
// The license key travels in the x-license-key header (never the query
// string) so it stays out of logs, history, and referrers.
// Returns the license key details or an error.
func (c *Client) GetKey(params GetKeyParams) (*GetKeyResponse, error) {
	var result GetKeyResponse
	req, err := http.NewRequest("GET", c.baseURL+"/key?productId="+url.QueryEscape(params.ProductID), nil)
	if err != nil {
		return nil, &ApiError{
			Message: fmt.Sprintf("failed to create request: %v", err),
			Code:    -1,
		}
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("x-license-key", params.LicenseKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &ApiError{
			Message: fmt.Sprintf("request failed: %v", err),
			Code:    -1,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ApiError{
			Message: fmt.Sprintf("failed to read response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	if resp.StatusCode >= 400 {
		return nil, parseApiError(body, resp.StatusCode)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &ApiError{
			Message: fmt.Sprintf("failed to unmarshal response: %v", err),
			Code:    -1,
			Status:  &resp.StatusCode,
		}
	}

	return &result, nil
}

// BlockKey blocks a specific license key.
// params: Parameters for blocking the key.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the block confirmation or an error.
func (c *Client) BlockKey(params BlockKeyParams, opts ...*RequestOptions) (*BlockKeyResponse, error) {
	var result BlockKeyResponse
	err := c.handleRequest("POST", "/key/block", params, &result, opts...)
	return &result, err
}

// UnblockKey unblocks a previously blocked license key.
// params: Parameters for unblocking the key.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the unblock confirmation or an error.
func (c *Client) UnblockKey(params UnblockKeyParams, opts ...*RequestOptions) (*UnblockKeyResponse, error) {
	var result UnblockKeyResponse
	err := c.handleRequest("POST", "/key/unblock", params, &result, opts...)
	return &result, err
}

// UpdateKey updates an existing license key via PATCH /api/key.
func (c *Client) UpdateKey(params UpdateKeyParams, opts ...*RequestOptions) (*UpdateKeyResponse, error) {
	var result UpdateKeyResponse
	err := c.handlePatchRequest("/key", params, &result, opts...)
	return &result, err
}

// SignKey signs a license key for offline (air-gapped) validation via POST /api/key/sign.
func (c *Client) SignKey(params SignKeyParams, opts ...*RequestOptions) (*SignKeyResponse, error) {
	var result SignKeyResponse
	err := c.handleRequest("POST", "/key/sign", params, &result, opts...)
	return &result, err
}

// CreateCustomer creates a new customer.
// params: Parameters for creating the customer.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the created customer information or an error.
func (c *Client) CreateCustomer(params CreateCustomerParams, opts ...*RequestOptions) (*CreateCustomerResponse, error) {
	var result CreateCustomerResponse
	err := c.handleRequest("POST", "/customer", params, &result, opts...)
	return &result, err
}

// GetAllCustomers retrieves all customers.
// params: Optional parameters for pagination and filtering.
// Returns a list of all customers or an error.
func (c *Client) GetAllCustomers(params GetAllCustomersParams) (*GetAllCustomersResponse, error) {
	var result GetAllCustomersResponse
	queryParams := make(map[string]string)
	if params.Page != nil {
		queryParams["page"] = fmt.Sprintf("%d", *params.Page)
	}
	if params.Limit != nil {
		queryParams["limit"] = fmt.Sprintf("%d", *params.Limit)
	}
	if params.Email != nil {
		queryParams["email"] = *params.Email
	}

	err := c.handleGetRequest("/customer", queryParams, &result)
	return &result, err
}

// GetCustomerWithKeys retrieves license keys belonging to a specific customer.
// Returns a flat list of license keys.
func (c *Client) GetCustomerWithKeys(params GetCustomerWithKeysParams) ([]CustomerLicenseKey, error) {
	var result []CustomerLicenseKey
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := c.handleGetRequest("/customer/keys", queryParams, &result)
	return result, err
}

// UpdateCustomer updates an existing customer.
// params: Parameters for updating the customer.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the update confirmation or an error.
func (c *Client) UpdateCustomer(params UpdateCustomerParams, opts ...*RequestOptions) (*UpdateCustomerResponse, error) {
	var result UpdateCustomerResponse
	err := c.handleRequest("PUT", "/customer/by-id", params, &result, opts...)
	return &result, err
}

// DeleteCustomer deletes a customer and all associated license keys permanently.
// params: Parameters containing the customer ID.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the deletion confirmation or an error.
func (c *Client) DeleteCustomer(params DeleteCustomerParams, opts ...*RequestOptions) (*DeleteCustomerResponse, error) {
	var result DeleteCustomerResponse
	queryParams := map[string]string{
		"customerId": params.CustomerID,
	}
	err := c.handleDeleteRequest("/customer/by-id", queryParams, &result, opts...)
	return &result, err
}

// ToggleCustomerStatus toggles the status of a customer (active/inactive).
// params: Parameters containing the customer ID.
// opts: Optional request configurations (e.g. idempotency keys).
// Returns the status toggle confirmation or an error.
func (c *Client) ToggleCustomerStatus(params ToggleCustomerStatusParams, opts ...*RequestOptions) (*ToggleCustomerStatusResponse, error) {
	var result ToggleCustomerStatusResponse
	endpoint := fmt.Sprintf("/customer/disable?customerId=%s", params.CustomerID)
	err := c.handleRequest("POST", endpoint, struct{}{}, &result, opts...)
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
