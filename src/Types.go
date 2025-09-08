package keymint

import "fmt"

// NewCustomer represents the structure for creating a new customer
// when creating a license key.
type NewCustomer struct {
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"` // Optional: Email of the new customer
}

// CreateKeyParams represents parameters for the createKey API endpoint.
type CreateKeyParams struct {
	ProductID     string       `json:"productId"`               // Required: The unique identifier of the product.
	MaxActivations *string      `json:"maxActivations,omitempty"` // Optional: The maximum number of times the key can be activated.
	ExpiryDate    *string      `json:"expiryDate,omitempty"`    // Optional: The expiration date of the key in ISO 8601 format.
	CustomerID    *string      `json:"customerId,omitempty"`    // Optional: The ID of an existing customer to associate with the key.
	NewCustomer   *NewCustomer `json:"newCustomer,omitempty"`   // Optional: An object to create and associate a new customer with the key.
}

// CreateKeyResponse represents response structure for a successful createKey API call.
type CreateKeyResponse struct {
	Code int    `json:"code"` // API response code (e.g., 0 for success)
	Key  string `json:"key"`  // The generated license key
}

// ApiError represents standard error response structure from the KeyMint API.
type ApiError struct {
	Message string `json:"message"`        // Descriptive error message
	Code    int    `json:"code"`           // API specific error code
	Status  *int   `json:"status,omitempty"` // HTTP status code, optional
}

func (e *ApiError) Error() string {
	if e.Status != nil {
		return fmt.Sprintf("KeyMint API Error (code: %d, status: %d): %s", e.Code, *e.Status, e.Message)
	}
	return fmt.Sprintf("KeyMint API Error (code: %d): %s", e.Code, e.Message)
}

// ActivateKeyParams represents parameters for the activateKey API endpoint.
type ActivateKeyParams struct {
	ProductID   string  `json:"productId"`             // Required: The unique identifier of the product.
	LicenseKey  string  `json:"licenseKey"`            // Required: The license key to activate.
	HostID      *string `json:"hostId,omitempty"`      // Optional: A unique identifier for the device.
	DeviceTag   *string `json:"deviceTag,omitempty"`   // Optional: A user-friendly name for the device.
}

// ActivateKeyResponse represents response structure for a successful activateKey API call.
type ActivateKeyResponse struct {
	Code          int     `json:"code"`                    // API response code (e.g., 0 for success)
	Message       string  `json:"message"`                 // Activation status message (e.g., "License valid")
	LicenseeName  *string `json:"licenseeName,omitempty"`  // Optional: Name of the licensee
	LicenseeEmail *string `json:"licenseeEmail,omitempty"` // Optional: Email of the licensee
}

// DeactivateKeyParams represents parameters for the deactivateKey API endpoint.
type DeactivateKeyParams struct {
	ProductID  string  `json:"productId"`            // Required: The unique identifier of the product.
	LicenseKey string  `json:"licenseKey"`           // Required: The license key to deactivate.
	HostID     *string `json:"hostId,omitempty"`     // Optional: The unique identifier of the device to deactivate. If omitted, all devices are deactivated.
}

// DeactivateKeyResponse represents response structure for a successful deactivateKey API call.
type DeactivateKeyResponse struct {
	Message string `json:"message"` // Confirmation message (e.g., "Device deactivated")
	Code    int    `json:"code"`    // API response code (e.g., 0 for success)
}

// DeviceDetails represents device details included in the GetKeyResponse.
type DeviceDetails struct {
	HostID         string  `json:"hostId"`                   // Updated field name
	DeviceTag      *string `json:"deviceTag,omitempty"`      // Updated field name  
	IPAddress      *string `json:"ipAddress,omitempty"`      // Updated field name
	ActivationTime string  `json:"activationTime"`           // Updated field name
}

// LicenseDetails represents license details included in the GetKeyResponse.
type LicenseDetails struct {
	ID             string          `json:"id"`
	Key            string          `json:"key"`
	ProductID      string          `json:"productId"`          // Updated field name
	MaxActivations int             `json:"maxActivations"`     // Updated field name
	Activations    int             `json:"activations"`
	Devices        []DeviceDetails `json:"devices"`
	Activated      bool            `json:"activated"`
	ExpirationDate *string         `json:"expirationDate,omitempty"` // Updated field name
}

// CustomerDetails represents customer details included in the GetKeyResponse.
type CustomerDetails struct {
	ID     string  `json:"id"`
	Name   *string `json:"name,omitempty"`   // Optional
	Email  *string `json:"email,omitempty"`  // Optional
	Active bool    `json:"active"`
}

// GetKeyParams represents parameters for the getKey API endpoint.
type GetKeyParams struct {
	ProductID  string `json:"productId"`  // Required: The unique identifier of the product.
	LicenseKey string `json:"licenseKey"` // Required: The license key to retrieve.
}

// GetKeyResponse represents response structure for a successful getKey API call.
type GetKeyResponse struct {
	Code int `json:"code"` // API response code (e.g., 0 for success)
	Data struct {
		License  LicenseDetails  `json:"license"`
		Customer *CustomerDetails `json:"customer,omitempty"` // Optional, customer data might not be present
	} `json:"data"`
}

// BlockKeyParams represents parameters for the blockKey API endpoint.
type BlockKeyParams struct {
	ProductID  string `json:"productId"`  // Required: The unique identifier of the product.
	LicenseKey string `json:"licenseKey"` // Required: The license key to block.
}

// BlockKeyResponse represents response structure for a successful blockKey API call.
type BlockKeyResponse struct {
	Message string `json:"message"` // Confirmation message (e.g., "Key blocked")
	Code    int    `json:"code"`    // API response code (e.g., 0 for success)
}

// UnblockKeyParams represents parameters for the unblockKey API endpoint.
type UnblockKeyParams struct {
	ProductID  string `json:"productId"`  // Required: The unique identifier of the product.
	LicenseKey string `json:"licenseKey"` // Required: The license key to unblock.
}

// UnblockKeyResponse represents response structure for a successful unblockKey API call.
type UnblockKeyResponse struct {
	Message string `json:"message"` // Confirmation message (e.g., "Key unblocked")
	Code    int    `json:"code"`    // API response code (e.g., 0 for success)
}

// CreateCustomerParams represents parameters for the createCustomer API endpoint.
type CreateCustomerParams struct {
	Name  string `json:"name"`  // Required: Customer name
	Email string `json:"email"` // Required: Customer email
}

// CreateCustomerResponse represents response structure for a successful createCustomer API call.
type CreateCustomerResponse struct {
	ID      string `json:"id"`
	Action  string `json:"action"`  // Action performed (e.g., "createCustomer")
	Status  bool   `json:"status"`  // Success status
	Message string `json:"message"` // Success message
	Data    struct {
		ID    string `json:"id"`    // Customer ID
		Name  string `json:"name"`  // Customer name
		Email string `json:"email"` // Customer email
	} `json:"data"`
	Code int `json:"code"` // API response code (e.g., 0 for success)
}

// Customer represents customer information in the getAllCustomers response.
type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
}

// GetAllCustomersResponse represents response structure for a successful getAllCustomers API call.
type GetAllCustomersResponse struct {
	Action string     `json:"action"`   // Action performed (e.g., "getCustomers")
	Status bool       `json:"status"`   // Success status
	Data   []Customer `json:"data"`     // Array of customer objects
	Code   int        `json:"code"`     // API response code (e.g., 0 for success)
}

// GetCustomerWithKeysParams represents parameters for the getCustomerWithKeys API endpoint.
type GetCustomerWithKeysParams struct {
	CustomerID string `json:"customerId"` // Required: The customer ID
}

// CustomerLicenseKey represents license key information in customer with keys response.
type CustomerLicenseKey struct {
	ID             string  `json:"id"`
	Key            string  `json:"key"`
	ProductID      string  `json:"productId"`
	MaxActivations int     `json:"maxActivations"`
	Activations    int     `json:"activations"`
	Activated      bool    `json:"activated"`
	ExpirationDate *string `json:"expirationDate,omitempty"`
}

// GetCustomerWithKeysResponse represents response structure for a successful getCustomerWithKeys API call.
type GetCustomerWithKeysResponse struct {
	Action string `json:"action"`
	Status bool   `json:"status"`
	Data   struct {
		Customer    Customer            `json:"customer"`
		LicenseKeys []CustomerLicenseKey `json:"licenseKeys"`
	} `json:"data"`
	Code int `json:"code"`
}

// UpdateCustomerParams represents parameters for the updateCustomer API endpoint.
type UpdateCustomerParams struct {
	CustomerID string  `json:"customerId"`          // Required: The customer ID
	Name       *string `json:"name,omitempty"`      // Optional: Updated customer name
	Email      *string `json:"email,omitempty"`     // Optional: Updated customer email
	Active     *bool   `json:"active,omitempty"`    // Optional: Customer active status
}

// UpdateCustomerResponse represents response structure for a successful updateCustomer API call.
type UpdateCustomerResponse struct {
	Action  string   `json:"action"`
	Status  bool     `json:"status"`
	Message string   `json:"message"`
	Data    Customer `json:"data"`
	Code    int      `json:"code"`
}

// ToggleCustomerStatusParams represents parameters for the toggleCustomerStatus API endpoint.
type ToggleCustomerStatusParams struct {
	CustomerID string `json:"customerId"` // Required: The customer ID
}

// ToggleCustomerStatusResponse represents response structure for a successful toggleCustomerStatus API call.
type ToggleCustomerStatusResponse struct {
	Action  string `json:"action"`  // Action performed (e.g., "toggleActive")
	Status  bool   `json:"status"`  // Success status
	Message string `json:"message"` // Status message (e.g., "Customer disabled")
	Code    int    `json:"code"`    // API response code
}

// GetCustomerByIdParams represents parameters for the getCustomerById API endpoint.
type GetCustomerByIdParams struct {
	CustomerID string `json:"customerId"` // Required: The customer ID
}

// GetCustomerByIdResponse represents response structure for a successful getCustomerById API call.
type GetCustomerByIdResponse struct {
	Action string     `json:"action"`      // Action performed (e.g., "getCustomerById")
	Status bool       `json:"status"`      // Success status
	Data   []Customer `json:"data"`        // Array containing the customer object
	Code   int        `json:"code"`        // API response code
}

// DeleteCustomerParams represents parameters for the deleteCustomer API endpoint.
type DeleteCustomerParams struct {
	CustomerID string `json:"customerId"` // Required: The customer ID
}

// DeleteCustomerResponse represents response structure for a successful deleteCustomer API call.
type DeleteCustomerResponse struct {
	Action  string `json:"action"`  // Action performed (e.g., "deleteCustomer")
	Status  bool   `json:"status"`  // Success status
	Message string `json:"message"` // Status message (e.g., "Customer deleted")
	Code    int    `json:"code"`    // API response code
}