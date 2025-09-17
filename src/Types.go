package keymint

import "fmt"

// NewCustomer represents the structure for creating a new customer when creating a license key.
type NewCustomer struct {
	// Name is the name of the new customer.
	Name  string  `json:"name"`
	// Email is the optional email of the new customer.
	Email *string `json:"email,omitempty"`
}

// CreateKeyParams represents parameters for the createKey API endpoint.
type CreateKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID     string       `json:"productId"`
	// MaxActivations is the optional maximum number of times the key can be activated.
	MaxActivations *string      `json:"maxActivations,omitempty"`
	// ExpiryDate is the optional expiration date of the key in ISO 8601 format.
	ExpiryDate    *string      `json:"expiryDate,omitempty"`
	// CustomerID is the optional ID of an existing customer to associate with the key.
	CustomerID    *string      `json:"customerId,omitempty"`
	// NewCustomer is an optional object to create and associate a new customer with the key.
	NewCustomer   *NewCustomer `json:"newCustomer,omitempty"`
}

// CreateKeyResponse represents response structure for a successful createKey API call.
type CreateKeyResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int    `json:"code"`
	// Key is the generated license key.
	Key  string `json:"key"`
}

// ApiError represents standard error response structure from the KeyMint API.
type ApiError struct {
	// Message is a descriptive error message.
	Message string `json:"message"`
	// Code is the API specific error code.
	Code    int    `json:"code"`
	// Status is the optional HTTP status code.
	Status  *int   `json:"status,omitempty"`
}

// Error implements the error interface for ApiError.
func (e *ApiError) Error() string {
	if e.Status != nil {
		return fmt.Sprintf("KeyMint API Error (code: %d, status: %d): %s", e.Code, *e.Status, e.Message)
	}
	return fmt.Sprintf("KeyMint API Error (code: %d): %s", e.Code, e.Message)
}

// ActivateKeyParams represents parameters for the activateKey API endpoint.
type ActivateKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID   string  `json:"productId"`
	// LicenseKey is the license key to activate.
	LicenseKey  string  `json:"licenseKey"`
	// HostID is an optional unique identifier for the device.
	HostID      *string `json:"hostId,omitempty"`
	// DeviceTag is an optional user-friendly name for the device.
	DeviceTag   *string `json:"deviceTag,omitempty"`
}

// ActivateKeyResponse represents response structure for a successful activateKey API call.
type ActivateKeyResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code          int     `json:"code"`
	// Message is the activation status message (e.g., "License valid").
	Message       string  `json:"message"`
	// LicenseeName is the optional name of the licensee.
	LicenseeName  *string `json:"licenseeName,omitempty"`
	// LicenseeEmail is the optional email of the licensee.
	LicenseeEmail *string `json:"licenseeEmail,omitempty"`
}

// DeactivateKeyParams represents parameters for the deactivateKey API endpoint.
type DeactivateKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID  string  `json:"productId"`
	// LicenseKey is the license key to deactivate.
	LicenseKey string  `json:"licenseKey"`
	// HostID is an optional unique identifier of the device to deactivate. If omitted, all devices are deactivated.
	HostID     *string `json:"hostId,omitempty"`
}

// DeactivateKeyResponse represents response structure for a successful deactivateKey API call.
type DeactivateKeyResponse struct {
	// Message is the confirmation message (e.g., "Device deactivated").
	Message string `json:"message"`
	// Code is the API response code (e.g., 0 for success).
	Code    int    `json:"code"`
}

// DeviceDetails represents device details included in the GetKeyResponse.
type DeviceDetails struct {
	// HostID is the updated field name.
	HostID         string  `json:"hostId"`
	// DeviceTag is the updated field name.
	DeviceTag      *string `json:"deviceTag,omitempty"`
	// IPAddress is the updated field name.
	IPAddress      *string `json:"ipAddress,omitempty"`
	// ActivationTime is the updated field name.
	ActivationTime string  `json:"activationTime"`
}

// LicenseDetails represents license details included in the GetKeyResponse.
type LicenseDetails struct {
	// ID is the license ID.
	ID             string          `json:"id"`
	// Key is the license key.
	Key            string          `json:"key"`
	// ProductID is the updated field name.
	ProductID      string          `json:"productId"`
	// MaxActivations is the updated field name.
	MaxActivations int             `json:"maxActivations"`
	// Activations is the number of times the license has been activated.
	Activations    int             `json:"activations"`
	// Devices is the list of devices associated with the license.
	Devices        []DeviceDetails `json:"devices"`
	// Activated indicates if the license is activated.
	Activated      bool            `json:"activated"`
	// ExpirationDate is the updated field name.
	ExpirationDate *string         `json:"expirationDate,omitempty"`
}

// CustomerDetails represents customer details included in the GetKeyResponse.
type CustomerDetails struct {
	// ID is the customer ID.
	ID     string  `json:"id"`
	// Name is the optional updated customer name.
	Name   *string `json:"name,omitempty"`
	// Email is the optional updated customer email.
	Email  *string `json:"email,omitempty"`
	// Active indicates if the customer is active.
	Active bool    `json:"active"`
}

// GetKeyParams represents parameters for the getKey API endpoint.
type GetKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID  string `json:"productId"`
	// LicenseKey is the license key to retrieve.
	LicenseKey string `json:"licenseKey"`
}

// GetKeyResponse represents response structure for a successful getKey API call.
type GetKeyResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// Data contains the license and optional customer details.
	Data struct {
		// License contains the license details.
		License  LicenseDetails  `json:"license"`
		// Customer contains the optional customer details.
		Customer *CustomerDetails `json:"customer,omitempty"`
	} `json:"data"`
}

// BlockKeyParams represents parameters for the blockKey API endpoint.
type BlockKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID  string `json:"productId"`
	// LicenseKey is the license key to block.
	LicenseKey string `json:"licenseKey"`
}

// BlockKeyResponse represents response structure for a successful blockKey API call.
type BlockKeyResponse struct {
	// Message is the confirmation message (e.g., "Key blocked").
	Message string `json:"message"`
	// Code is the API response code (e.g., 0 for success).
	Code    int    `json:"code"`
}

// UnblockKeyParams represents parameters for the unblockKey API endpoint.
type UnblockKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID  string `json:"productId"`
	// LicenseKey is the license key to unblock.
	LicenseKey string `json:"licenseKey"`
}

// UnblockKeyResponse represents response structure for a successful unblockKey API call.
type UnblockKeyResponse struct {
	// Message is the confirmation message (e.g., "Key unblocked").
	Message string `json:"message"`
	// Code is the API response code (e.g., 0 for success).
	Code    int    `json:"code"`
}

// CreateCustomerParams represents parameters for the createCustomer API endpoint.
type CreateCustomerParams struct {
	// Name is the required customer name.
	Name  string `json:"name"`
	// Email is the required customer email.
	Email string `json:"email"`
}

// CreateCustomerResponse represents response structure for a successful createCustomer API call.
type CreateCustomerResponse struct {
	// ID is the customer ID.
	ID      string `json:"id"`
	// Action is the action performed (e.g., "createCustomer").
	Action  string `json:"action"`
	// Status indicates the success status.
	Status  bool   `json:"status"`
	// Message is the success message.
	Message string `json:"message"`
	// Data contains the created customer details.
	Data    struct {
		// ID is the customer ID.
		ID    string `json:"id"`
		// Name is the customer name.
		Name  string `json:"name"`
		// Email is the customer email.
		Email string `json:"email"`
	} `json:"data"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
}

// Customer represents customer information in the getAllCustomers response.
type Customer struct {
	// ID is the customer ID.
	ID        string `json:"id"`
	// Name is the customer name.
	Name      string `json:"name"`
	// Email is the customer email.
	Email     string `json:"email"`
	// Active indicates if the customer is active.
	Active    bool   `json:"active"`
	// CreatedAt is the timestamp when the customer was created.
	CreatedAt string `json:"createdAt"`
	// UpdatedAt is the timestamp when the customer was last updated.
	UpdatedAt string `json:"updatedAt"`
	// CreatedBy is the identifier of the user who created the customer.
	CreatedBy string `json:"createdBy"`
}

// GetAllCustomersResponse represents response structure for a successful getAllCustomers API call.
type GetAllCustomersResponse struct {
	// Action is the action performed (e.g., "getCustomers").
	Action string     `json:"action"`
	// Status indicates the success status.
	Status bool       `json:"status"`
	// Data is the array of customer objects.
	Data   []Customer `json:"data"`
	// Code is the API response code (e.g., 0 for success).
	Code   int        `json:"code"`
}

// GetCustomerWithKeysParams represents parameters for the getCustomerWithKeys API endpoint.
type GetCustomerWithKeysParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// CustomerLicenseKey represents license key information in customer with keys response.
type CustomerLicenseKey struct {
	// ID is the license key ID.
	ID             string  `json:"id"`
	// Key is the license key.
	Key            string  `json:"key"`
	// ProductID is the product ID associated with the license key.
	ProductID      string  `json:"productId"`
	// MaxActivations is the maximum number of activations for the license key.
	MaxActivations int     `json:"maxActivations"`
	// Activations is the number of times the license key has been activated.
	Activations    int     `json:"activations"`
	// Activated indicates if the license key is activated.
	Activated      bool    `json:"activated"`
	// ExpirationDate is the expiration date of the license key.
	ExpirationDate *string `json:"expirationDate,omitempty"`
}

// GetCustomerWithKeysResponse represents response structure for a successful getCustomerWithKeys API call.
type GetCustomerWithKeysResponse struct {
	// Action is the action performed (e.g., "getCustomerWithKeys").
	Action string `json:"action"`
	// Status indicates the success status.
	Status bool   `json:"status"`
	// Data contains the customer and license keys information.
	Data   struct {
		// Customer contains the customer details.
		Customer    Customer            `json:"customer"`
		// LicenseKeys contains the array of license keys associated with the customer.
		LicenseKeys []CustomerLicenseKey `json:"licenseKeys"`
	} `json:"data"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
}

// UpdateCustomerParams represents parameters for the updateCustomer API endpoint.
type UpdateCustomerParams struct {
	// CustomerID is the required customer ID.
	CustomerID string  `json:"customerId"`
	// Name is the optional updated customer name.
	Name       *string `json:"name,omitempty"`
	// Email is the optional updated customer email.
	Email      *string `json:"email,omitempty"`
	// Active is the optional customer active status.
	Active     *bool   `json:"active,omitempty"`
}

// UpdateCustomerResponse represents response structure for a successful updateCustomer API call.
type UpdateCustomerResponse struct {
	// Action is the action performed (e.g., "updateCustomer").
	Action  string   `json:"action"`
	// Status indicates the success status.
	Status  bool     `json:"status"`
	// Message is the status message.
	Message string   `json:"message"`
	// Data contains the updated customer details.
	Data    Customer `json:"data"`
	// Code is the API response code (e.g., 0 for success).
	Code    int      `json:"code"`
}

// ToggleCustomerStatusParams represents parameters for the toggleCustomerStatus API endpoint.
type ToggleCustomerStatusParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// ToggleCustomerStatusResponse represents response structure for a successful toggleCustomerStatus API call.
type ToggleCustomerStatusResponse struct {
	// Action is the action performed (e.g., "toggleActive").
	Action  string `json:"action"`
	// Status indicates the success status.
	Status  bool   `json:"status"`
	// Message is the status message (e.g., "Customer disabled").
	Message string `json:"message"`
	// Code is the API response code.
	Code    int    `json:"code"`
}

// GetCustomerByIdParams represents parameters for the getCustomerById API endpoint.
type GetCustomerByIdParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// GetCustomerByIdResponse represents response structure for a successful getCustomerById API call.
type GetCustomerByIdResponse struct {
	// Action is the action performed (e.g., "getCustomerById").
	Action string     `json:"action"`
	// Status indicates the success status.
	Status bool       `json:"status"`
	// Data is the array containing the customer object.
	Data   []Customer `json:"data"`
	// Code is the API response code.
	Code   int        `json:"code"`
}

// DeleteCustomerParams represents parameters for the deleteCustomer API endpoint.
type DeleteCustomerParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// DeleteCustomerResponse represents response structure for a successful deleteCustomer API call.
type DeleteCustomerResponse struct {
	// Action is the action performed (e.g., "deleteCustomer").
	Action  string `json:"action"`
	// Status indicates the success status.
	Status  bool   `json:"status"`
	// Message is the status message (e.g., "Customer deleted").
	Message string `json:"message"`
	// Code is the API response code.
	Code    int    `json:"code"`
}