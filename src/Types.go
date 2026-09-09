package keymint

import "fmt"

// NewCustomer represents the structure for creating a new customer when creating a license key.
type NewCustomer struct {
	// Name is the name of the new customer.
	Name string `json:"name"`
	// Email is the optional email of the new customer.
	Email *string `json:"email,omitempty"`
}

// KeyFormat represents key format options for custom license key shapes.
type KeyFormat struct {
	// Sections is the optional number of sections (1-10).
	Sections *int `json:"sections,omitempty"`
	// SectionLength is the optional length of each section (1-32).
	SectionLength *int `json:"sectionLength,omitempty"`
	// Separator is the optional separator character (max 3 chars).
	Separator *string `json:"separator,omitempty"`
	// Charset is the optional custom character set (no whitespace).
	Charset *string `json:"charset,omitempty"`
	// Prefix is the optional prefix prepended to the key (max 16 chars).
	Prefix *string `json:"prefix,omitempty"`
	// Suffix is the optional suffix appended to the key (max 16 chars).
	Suffix *string `json:"suffix,omitempty"`
	// Case is the optional character case: "upper", "lower", or "mixed".
	Case *string `json:"case,omitempty"`
}

// CreateKeyParams represents parameters for the createKey API endpoint.
type CreateKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// MaxActivations is the optional maximum number of times the key can be activated.
	MaxActivations *string `json:"maxActivations,omitempty"`
	// ExpiryDate is the optional expiration date of the key in ISO 8601 format.
	ExpiryDate *string `json:"expiryDate,omitempty"`
	// CustomerID is the optional ID of an existing customer to associate with the key.
	CustomerID *string `json:"customerId,omitempty"`
	// VersionID is the optional ID of a specific product version to associate with the key.
	VersionID *string `json:"versionId,omitempty"`
	// Metadata is an optional custom dictionary payload to attach to the license key.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// NewCustomer is an optional object to create and associate a new customer with the key.
	NewCustomer *NewCustomer `json:"newCustomer,omitempty"`
	// AllowedHosts is an optional list of machine IDs authorized to use this license.
	AllowedHosts []string `json:"allowedHosts,omitempty"`
	// Format is an optional custom key format.
	Format *KeyFormat `json:"format,omitempty"`
	// AmountKeys is the optional number of keys to generate at once (bulk creation).
	AmountKeys *string `json:"amountKeys,omitempty"`
	// LicenseType is the optional license type: "node-locked" or "floating" (defaults to "node-locked").
	LicenseType *string `json:"licenseType,omitempty"`
	// MaxConcurrentSessions is the optional max concurrent floating sessions.
	MaxConcurrentSessions *int `json:"maxConcurrentSessions,omitempty"`
	// HeartbeatInterval is the optional floating heartbeat interval in seconds (min 60).
	HeartbeatInterval *int `json:"heartbeatInterval,omitempty"`
	// SessionLeaseDuration is the optional floating session lease duration in seconds (min 300).
	SessionLeaseDuration *int `json:"sessionLeaseDuration,omitempty"`
}

// CreateKeyResponse represents response structure for a successful createKey API call.
type CreateKeyResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// Key is the generated license key for single-key creation.
	Key string `json:"key"`
	// Keys holds generated license keys for bulk creation.
	Keys []string `json:"keys,omitempty"`
}

// ApiErrorDetails carries the nested error envelope from the current Keymint API.
type ApiErrorDetails struct {
	// Code is the string error code (e.g., "RATE_LIMIT_EXCEEDED").
	Code *string `json:"code,omitempty"`
	// Message is the detailed error message.
	Message *string `json:"message,omitempty"`
	// Details holds optional extra error context.
	Details interface{} `json:"details,omitempty"`
}

// ApiError represents standard error response structure from the KeyMint API.
type ApiError struct {
	// Message is a descriptive error message.
	Message string `json:"message"`
	// Code is the API specific error code.
	Code int `json:"code"`
	// Status is the optional HTTP status code.
	Status *int `json:"status,omitempty"`
	// ErrorDetail holds the nested error envelope when present.
	ErrorDetail *ApiErrorDetails `json:"error,omitempty"`
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
	ProductID string `json:"productId"`
	// LicenseKey is the license key to activate.
	LicenseKey string `json:"licenseKey"`
	// HostID is an optional unique identifier for the device.
	HostID *string `json:"hostId,omitempty"`
	// DeviceTag is an optional user-friendly name for the device.
	DeviceTag *string `json:"deviceTag,omitempty"`
	// Licensee is an optional customer name and email to set during activation.
	Licensee *ActivationLicensee `json:"licensee,omitempty"`
	// Version is an optional product version string (max 32 chars).
	Version *string `json:"version,omitempty"`
}

// ActivationLicensee represents customer info set during activation.
type ActivationLicensee struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ActivateKeyResponse represents response structure for a successful activateKey API call.
type ActivateKeyResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// Message is the activation status message (e.g., "License valid").
	Message string `json:"message"`
	// LicenseeName is the optional name of the licensee.
	LicenseeName *string `json:"licenseeName,omitempty"`
	// LicenseeEmail is the optional email of the licensee.
	LicenseeEmail *string `json:"licenseeEmail,omitempty"`
	// Metadata is the optional custom dictionary attached to the license key.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// VersionID is the optional associated product version ID.
	VersionID *string `json:"versionId,omitempty"`
	// Version is the optional detailed product version information.
	Version map[string]interface{} `json:"version,omitempty"`
	// AllowedHosts is an optional list of authorized machine IDs.
	AllowedHosts []string `json:"allowedHosts,omitempty"`
}

// DeactivateKeyParams represents parameters for the deactivateKey API endpoint.
type DeactivateKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key to deactivate.
	LicenseKey string `json:"licenseKey"`
	// HostID is an optional unique identifier of the device to deactivate. If omitted, all devices are deactivated.
	HostID *string `json:"hostId,omitempty"`
}

// DeactivateKeyResponse represents response structure for a successful deactivateKey API call.
type DeactivateKeyResponse struct {
	// Message is the confirmation message (e.g., "Device deactivated").
	Message string `json:"message"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// DevicesRemoved is the number of device activations removed.
	DevicesRemoved *int `json:"devicesRemoved,omitempty"`
}

// DeviceDetails represents device details included in the GetKeyResponse.
type DeviceDetails struct {
	// HostID is the updated field name.
	HostID string `json:"hostId"`
	// DeviceTag is the updated field name.
	DeviceTag *string `json:"deviceTag,omitempty"`
	// IPAddress is the updated field name.
	IPAddress *string `json:"ipAddress,omitempty"`
	// ActivationTime is the updated field name.
	ActivationTime string `json:"activationTime"`
}

// LicenseDetails represents license details included in the GetKeyResponse.
type LicenseDetails struct {
	// ID is the license ID.
	ID string `json:"id"`
	// Key is the license key.
	Key string `json:"key"`
	// ProductID is the updated field name.
	ProductID string `json:"productId"`
	// MaxActivations is the updated field name.
	MaxActivations int `json:"maxActivations"`
	// Activations is the number of times the license has been activated.
	Activations int `json:"activations"`
	// Devices is the list of devices associated with the license.
	Devices []DeviceDetails `json:"devices"`
	// Activated indicates if the license is activated.
	Activated bool `json:"activated"`
	// ExpirationDate is the updated field name.
	ExpirationDate *string `json:"expirationDate,omitempty"`
	// Metadata is the optional custom dictionary attached to the license key.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// VersionID is the optional associated product version ID.
	VersionID *string `json:"versionId,omitempty"`
	// Version is the optional detailed product version information.
	Version map[string]interface{} `json:"version,omitempty"`
	// AllowedHosts is the optional list of authorized machine IDs.
	AllowedHosts []string `json:"allowedHosts,omitempty"`
}

// CustomerDetails represents customer details included in the GetKeyResponse.
type CustomerDetails struct {
	// ID is the customer ID.
	ID string `json:"id"`
	// Name is the optional updated customer name.
	Name *string `json:"name,omitempty"`
	// Email is the optional updated customer email.
	Email *string `json:"email,omitempty"`
	// Active indicates if the customer is active.
	Active bool `json:"active"`
}

// GetKeyParams represents parameters for the getKey API endpoint.
type GetKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
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
		License LicenseDetails `json:"license"`
		// Customer contains the optional customer details.
		Customer *CustomerDetails `json:"customer,omitempty"`
	} `json:"data"`
}

// BlockKeyParams represents parameters for the blockKey API endpoint.
type BlockKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key to block.
	LicenseKey string `json:"licenseKey"`
}

// BlockKeyResponse represents response structure for a successful blockKey API call.
type BlockKeyResponse struct {
	// Message is the confirmation message (e.g., "Key blocked").
	Message string `json:"message"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
}

// UnblockKeyParams represents parameters for the unblockKey API endpoint.
type UnblockKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key to unblock.
	LicenseKey string `json:"licenseKey"`
}

// UnblockKeyResponse represents response structure for a successful unblockKey API call.
type UnblockKeyResponse struct {
	// Message is the confirmation message (e.g., "Key unblocked").
	Message string `json:"message"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
}

// CreateCustomerParams represents parameters for the createCustomer API endpoint.
type CreateCustomerParams struct {
	// Name is the required customer name.
	Name string `json:"name"`
	// Email is the required customer email.
	Email string `json:"email"`
}

// CreateCustomerResponse represents response structure for a successful createCustomer API call.
type CreateCustomerResponse struct {
	// ID is the customer ID.
	ID string `json:"id"`
	// Action is the action performed (e.g., "createCustomer").
	Action string `json:"action"`
	// Status indicates the success status.
	Status bool `json:"status"`
	// Message is the success message.
	Message string `json:"message"`
	// Data contains the created customer details.
	Data struct {
		// ID is the customer ID.
		ID string `json:"id"`
		// Name is the customer name.
		Name string `json:"name"`
		// Email is the customer email.
		Email string `json:"email"`
	} `json:"data"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
}

// Customer represents customer information in the getAllCustomers response.
type Customer struct {
	// ID is the customer ID.
	ID string `json:"id"`
	// Name is the customer name.
	Name string `json:"name"`
	// Email is the customer email.
	Email string `json:"email"`
	// Active indicates if the customer is active.
	Active bool `json:"active"`
	// CreatedAt is the timestamp when the customer was created.
	CreatedAt string `json:"createdAt"`
	// UpdatedAt is the timestamp when the customer was last updated.
	UpdatedAt string `json:"updatedAt"`
	// CreatedBy is the identifier of the user who created the customer.
	// The REST API sanitizes this field, so it is usually empty.
	CreatedBy string `json:"createdBy"`
}

// GetAllCustomersParams represents parameters for the getAllCustomers API endpoint.
type GetAllCustomersParams struct {
	// Page is the optional page number.
	Page *int `json:"page,omitempty"`
	// Limit is the optional number of items per page.
	Limit *int `json:"limit,omitempty"`
	// Email is the optional filter by email.
	Email *string `json:"email,omitempty"`
}

// PaginationMeta represents pagination metadata included in list responses.
type PaginationMeta struct {
	// Total is the total number of items.
	Total int `json:"total"`
	// Page is the current page number.
	Page int `json:"page"`
	// Limit is the number of items per page.
	Limit int `json:"limit"`
	// TotalPages is the total number of pages.
	TotalPages int `json:"totalPages"`
}

// GetAllCustomersResponse represents response structure for a successful getAllCustomers API call.
type GetAllCustomersResponse struct {
	// Action is the action performed (e.g., "getCustomers").
	Action string `json:"action"`
	// Status indicates the success status.
	Status bool `json:"status"`
	// Data is the array of customer objects.
	Data []Customer `json:"data"`
	// Meta contains pagination metadata.
	Meta *PaginationMeta `json:"meta,omitempty"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
}

// GetCustomerWithKeysParams represents parameters for the getCustomerWithKeys API endpoint.
type GetCustomerWithKeysParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// CustomerLicenseKey represents license key information in customer with keys response.
type CustomerLicenseKey struct {
	// ID is the license key ID.
	ID string `json:"id"`
	// Key is the license key.
	Key string `json:"key"`
	// ProductID is the product ID associated with the license key.
	ProductID string `json:"productId"`
	// MaxActivations is the maximum number of activations for the license key.
	MaxActivations int `json:"maxActivations"`
	// Activations is the number of times the license key has been activated.
	Activations int `json:"activations"`
	// Activated indicates if the license key is activated.
	Activated bool `json:"activated"`
	// ExpirationDate is the expiration date of the license key.
	ExpirationDate *string `json:"expirationDate,omitempty"`
	// Metadata is the optional custom dictionary attached to the license key.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// VersionID is the optional associated product version ID.
	VersionID *string `json:"versionId,omitempty"`
	// AllowedHosts is the optional list of authorized machine IDs.
	AllowedHosts []string `json:"allowedHosts,omitempty"`
}

// GetCustomerWithKeysResponse represents response structure for a successful getCustomerWithKeys API call.
// Returns a flat list of license keys for the customer.
type GetCustomerWithKeysResponse []CustomerLicenseKey

// UpdateCustomerParams represents parameters for the updateCustomer API endpoint.
type UpdateCustomerParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
	// Name is the optional updated customer name.
	Name *string `json:"name,omitempty"`
	// Email is the optional updated customer email.
	Email *string `json:"email,omitempty"`
}

// UpdateCustomerResponse represents response structure for a successful updateCustomer API call.
type UpdateCustomerResponse struct {
	// Action is the action performed (e.g., "updateCustomer").
	Action string `json:"action"`
	// Status indicates the success status.
	Status bool `json:"status"`
	// Message is the status message.
	Message string `json:"message"`
	// Data contains the updated customer details.
	Data Customer `json:"data"`
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
}

// ToggleCustomerStatusParams represents parameters for the toggleCustomerStatus API endpoint.
type ToggleCustomerStatusParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// ToggleCustomerStatusResponse represents response structure for a successful toggleCustomerStatus API call.
type ToggleCustomerStatusResponse struct {
	// Action is the action performed (e.g., "toggleActive").
	Action string `json:"action"`
	// Status indicates the success status.
	Status bool `json:"status"`
	// Message is the status message when present (absent on current API envelope).
	Message *string `json:"message,omitempty"`
	// Code is the API response code when present.
	Code *int `json:"code,omitempty"`
	// CustomerName is the customer name as returned by the API.
	CustomerName *string `json:"customerName,omitempty"`
	// Active is the resulting active flag as returned by the API.
	Active *bool `json:"active,omitempty"`
}

// GetCustomerByIdParams represents parameters for the getCustomerById API endpoint.
type GetCustomerByIdParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// GetCustomerByIdResponse represents response structure for a successful getCustomerById API call.
type GetCustomerByIdResponse struct {
	// Action is the action performed (e.g., "getCustomerById").
	Action string `json:"action"`
	// Status indicates the success status.
	Status bool `json:"status"`
	// Data is the array containing the customer object.
	Data []Customer `json:"data"`
	// Code is the API response code.
	Code int `json:"code"`
}

// DeleteCustomerParams represents parameters for the deleteCustomer API endpoint.
type DeleteCustomerParams struct {
	// CustomerID is the required customer ID.
	CustomerID string `json:"customerId"`
}

// DeleteCustomerResponse represents response structure for a successful deleteCustomer API call.
type DeleteCustomerResponse struct {
	// Action is the action performed (e.g., "deleteCustomer").
	Action string `json:"action"`
	// Status indicates the success status.
	Status bool `json:"status"`
	// Message is the status message (e.g., "Customer deleted").
	Message string `json:"message"`
	// Code is the API response code.
	Code int `json:"code"`
}

// FloatingCheckoutParams represents parameters for the floating license checkout API endpoint.
type FloatingCheckoutParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key.
	LicenseKey string `json:"licenseKey"`
	// HostID is the unique hardware identifier of the device.
	HostID string `json:"hostId"`
	// DeviceTag is an optional friendly name for the device.
	DeviceTag *string `json:"deviceTag,omitempty"`
	// UserIdentifier is an optional user identifier.
	UserIdentifier *string `json:"userIdentifier,omitempty"`
	// Timestamp carries the session's current nextNonce when re-checking out
	// an existing active hostId (session extension proof of possession).
	Timestamp interface{} `json:"timestamp,omitempty"`
	// Signature is HMAC-SHA256(sessionSecret, "sessionId:timestamp"), required
	// for session extension (see GenerateSessionSignature).
	Signature *string `json:"signature,omitempty"`
}

// FloatingCheckoutResponse represents response structure for a successful floating license checkout API call.
type FloatingCheckoutResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// Message is the activation status message.
	Message string `json:"message"`
	// SessionID is the unique session ID.
	SessionID string `json:"sessionId"`
	// SessionSecret is the temporary session secret key.
	SessionSecret string `json:"sessionSecret"`
	// NextNonce is the rotating nonce string to use for the next request.
	NextNonce string `json:"nextNonce"`
	// ExpiresAt is the expiration time of the session in ISO 8601 format.
	ExpiresAt string `json:"expiresAt"`
	// HeartbeatInterval is the interval (in seconds) the client must heartbeat within.
	HeartbeatInterval int `json:"heartbeatInterval"`
	// Metadata is the optional custom dictionary attached to the license key.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// CurrentSessions is the number of active sessions for the license key.
	CurrentSessions *int `json:"currentSessions,omitempty"`
	// MaxSessions is the maximum concurrent sessions allowed for the license key.
	MaxSessions *int `json:"maxSessions,omitempty"`
	// LicenseeName is the optional name of the customer licensee.
	LicenseeName *string `json:"licenseeName,omitempty"`
	// LicenseeEmail is the optional email of the customer licensee.
	LicenseeEmail *string `json:"licenseeEmail,omitempty"`
}

// FloatingHeartbeatParams represents parameters for the floating license heartbeat API endpoint.
type FloatingHeartbeatParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key.
	LicenseKey string `json:"licenseKey"`
	// SessionID is the unique session ID.
	SessionID string `json:"sessionId"`
	// Timestamp is the rotating nonce (nextNonce) received from the previous response.
	Timestamp interface{} `json:"timestamp"`
	// Signature is the HMAC-SHA256 signature generated using the sessionSecret over the payload 'sessionId:nonce'.
	Signature string `json:"signature"`
}

// FloatingHeartbeatResponse represents response structure for a successful floating license heartbeat API call.
type FloatingHeartbeatResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// Message is the response status message.
	Message string `json:"message"`
	// ExpiresAt is the extended expiration time of the session in ISO 8601 format.
	ExpiresAt string `json:"expiresAt"`
	// NextNonce is the newly rotated nonce to be used for the next subsequent heartbeat.
	NextNonce string `json:"nextNonce"`
}

// FloatingCheckinParams represents parameters for the floating license checkin API endpoint.
type FloatingCheckinParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key.
	LicenseKey string `json:"licenseKey"`
	// SessionID is the unique session ID.
	SessionID string `json:"sessionId"`
	// Timestamp is the rotating nonce (nextNonce) received from the previous response.
	Timestamp interface{} `json:"timestamp"`
	// Signature is the HMAC-SHA256 signature generated using the sessionSecret over the payload 'sessionId:nonce'.
	Signature string `json:"signature"`
}

// FloatingCheckinResponse represents response structure for a successful floating license checkin API call.
type FloatingCheckinResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// Message is the confirmation message.
	Message string `json:"message"`
}

// UpdateKeyParams represents parameters for the updateKey API endpoint (PATCH /api/key).
type UpdateKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key to update.
	LicenseKey string `json:"licenseKey"`
	// MaxActivations is the optional new max activations count (string or number).
	MaxActivations interface{} `json:"maxActivations,omitempty"`
	// ExpiryDate is the optional new expiration date in ISO 8601 format.
	ExpiryDate *string `json:"expiryDate,omitempty"`
	// CustomerID is the optional new customer ID to associate.
	CustomerID *string `json:"customerId,omitempty"`
	// NewCustomer is an optional object to create and associate a new customer.
	NewCustomer *NewCustomer `json:"newCustomer,omitempty"`
	// Metadata is the optional updated custom metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// VersionID is the optional updated product version ID.
	VersionID *string `json:"versionId,omitempty"`
	// AllowedHosts is the optional updated list of authorized machine IDs.
	AllowedHosts []string `json:"allowedHosts,omitempty"`
	// LicenseType is the optional license type: "node-locked" or "floating".
	LicenseType *string `json:"licenseType,omitempty"`
	// MaxConcurrentSessions is the optional updated max concurrent sessions.
	MaxConcurrentSessions *int `json:"maxConcurrentSessions,omitempty"`
	// HeartbeatInterval is the optional updated heartbeat interval in seconds (min 60).
	HeartbeatInterval *int `json:"heartbeatInterval,omitempty"`
	// SessionLeaseDuration is the optional updated session lease duration in seconds (min 300).
	SessionLeaseDuration *int `json:"sessionLeaseDuration,omitempty"`
}

// UpdateKeyResponse represents response structure for a successful updateKey API call.
type UpdateKeyResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// Message is the confirmation message.
	Message string `json:"message"`
	// AffectedCount is the number of keys affected.
	AffectedCount *int `json:"affectedCount,omitempty"`
}

// SignKeyParams represents parameters for the signKey API endpoint (POST /api/key/sign).
type SignKeyParams struct {
	// ProductID is the unique identifier of the product.
	ProductID string `json:"productId"`
	// LicenseKey is the license key to sign.
	LicenseKey string `json:"licenseKey"`
	// HostID is the required machine code to bind the offline license to.
	HostID string `json:"hostId"`
	// TTL is the optional time-to-live in seconds (min 60).
	TTL *int `json:"ttl,omitempty"`
}

// SignKeyResponse represents response structure for a successful signKey API call.
type SignKeyResponse struct {
	// Code is the API response code (e.g., 0 for success).
	Code int `json:"code"`
	// File contains the signed license file (signedKey, keyId, publicKeyFingerprint) — may be a JSON string.
	File interface{} `json:"file"`
}

// RequestOptions contains optional parameters for Keymint API requests (e.g. idempotency keys).
type RequestOptions struct {
	IdempotencyKey string
}
