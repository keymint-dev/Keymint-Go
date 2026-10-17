// Package keymint (import path github.com/keymint-dev/keymint-go/src) is
// deprecated: import github.com/keymint-dev/keymint-go instead. This package
// only forwards to the module root so existing /src imports keep compiling.
package keymint

import km "github.com/keymint-dev/keymint-go"

// Types.
type (
	Client                       = km.Client
	NewCustomer                  = km.NewCustomer
	KeyFormat                    = km.KeyFormat
	CreateKeyParams              = km.CreateKeyParams
	CreateKeyResponse            = km.CreateKeyResponse
	ApiErrorDetails              = km.ApiErrorDetails
	ApiError                     = km.ApiError
	ActivateKeyParams            = km.ActivateKeyParams
	ActivationLicensee           = km.ActivationLicensee
	ActivateKeyResponse          = km.ActivateKeyResponse
	DeactivateKeyParams          = km.DeactivateKeyParams
	DeactivateKeyResponse        = km.DeactivateKeyResponse
	DeviceDetails                = km.DeviceDetails
	LicenseDetails               = km.LicenseDetails
	CustomerDetails              = km.CustomerDetails
	GetKeyParams                 = km.GetKeyParams
	GetKeyResponse               = km.GetKeyResponse
	BlockKeyParams               = km.BlockKeyParams
	BlockKeyResponse             = km.BlockKeyResponse
	UnblockKeyParams             = km.UnblockKeyParams
	UnblockKeyResponse           = km.UnblockKeyResponse
	CreateCustomerParams         = km.CreateCustomerParams
	CreateCustomerResponse       = km.CreateCustomerResponse
	Customer                     = km.Customer
	GetAllCustomersParams        = km.GetAllCustomersParams
	PaginationMeta               = km.PaginationMeta
	GetAllCustomersResponse      = km.GetAllCustomersResponse
	GetCustomerWithKeysParams    = km.GetCustomerWithKeysParams
	CustomerLicenseKey           = km.CustomerLicenseKey
	GetCustomerWithKeysResponse  = km.GetCustomerWithKeysResponse
	UpdateCustomerParams         = km.UpdateCustomerParams
	UpdateCustomerResponse       = km.UpdateCustomerResponse
	ToggleCustomerStatusParams   = km.ToggleCustomerStatusParams
	ToggleCustomerStatusResponse = km.ToggleCustomerStatusResponse
	GetCustomerByIdParams        = km.GetCustomerByIdParams
	GetCustomerByIdResponse      = km.GetCustomerByIdResponse
	DeleteCustomerParams         = km.DeleteCustomerParams
	DeleteCustomerResponse       = km.DeleteCustomerResponse
	FloatingCheckoutParams       = km.FloatingCheckoutParams
	FloatingCheckoutResponse     = km.FloatingCheckoutResponse
	FloatingHeartbeatParams      = km.FloatingHeartbeatParams
	FloatingHeartbeatResponse    = km.FloatingHeartbeatResponse
	FloatingCheckinParams        = km.FloatingCheckinParams
	FloatingCheckinResponse      = km.FloatingCheckinResponse
	UpdateKeyParams              = km.UpdateKeyParams
	UpdateKeyResponse            = km.UpdateKeyResponse
	SignKeyParams                = km.SignKeyParams
	SignKeyResponse              = km.SignKeyResponse
	RequestOptions               = km.RequestOptions
)

// Functions.
var (
	GetMachineID              = km.GetMachineID
	GetOrCreateInstallationID = km.GetOrCreateInstallationID
	GenerateSessionSignature  = km.GenerateSessionSignature
	New                       = km.New
	VerifyWebhookSignature    = km.VerifyWebhookSignature
	String                    = km.String
	Int                       = km.Int
	Bool                      = km.Bool
)
