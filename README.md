# KeyMint Go

A professional, production-ready SDK for integrating with the KeyMint API in Go. Provides robust access to all major KeyMint features, with idiomatic Go error handling.

## Features
- **Idiomatic Go**: Clean, type-safe API using a `Client` struct and `keymint.New()` constructor.
- **Comprehensive**: Complete API coverage for all KeyMint endpoints.
- **Consistent error handling**: All API errors are returned as `*ApiError`.
- **Security**: Credentials are always loaded from environment variables.

## Installation
Add the SDK to your project:

```bash
go get github.com/keymint-dev/Keymint-Go
```

## Usage

```go
import (
    "os"
    "github.com/keymint-dev/Keymint-Go"
)

accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
productId := os.Getenv("KEYMINT_PRODUCT_ID")

client, err := keymint.New(accessToken, "")
if err != nil {
    // handle error
}

key, err := client.CreateKey(keymint.CreateKeyParams{ProductID: productId})
if err != nil {
    // handle error
}
```

## Error Handling
All SDK methods return a result and an error. If the error is from the API, it will be of type `*ApiError` with `Message`, `Code`, and `Status` fields.

## API Methods

All methods return a result struct and an error.

### License Key Management

| Method           | Description                                     |
|------------------|-------------------------------------------------|
| `CreateKey`      | Creates a new license key.                      |
| `ActivateKey`    | Activates a license key for a device.           |
| `DeactivateKey`  | Deactivates a device from a license key.        |
| `GetKey`         | Retrieves detailed information about a key.     |
| `BlockKey`       | Blocks a license key.                           |
| `UnblockKey`     | Unblocks a previously blocked license key.      |

### Customer Management

| Method                  | Description                                      |
|-------------------------|--------------------------------------------------|
| `CreateCustomer`        | Creates a new customer.                          |
| `GetAllCustomers`       | Retrieves all customers.                         |
| `GetCustomerById`       | Gets a specific customer by ID.                  |
| `GetCustomerWithKeys`   | Gets a customer along with their license keys.   |
| `UpdateCustomer`        | Updates customer information.                    |
| `ToggleCustomerStatus`  | Toggles customer active status.                  |
| `DeleteCustomer`        | Permanently deletes a customer and their keys.   |

For detailed parameter and response types, see the [KeyMint API docs](https://docs.keymint.dev) or use your IDE's autocomplete.

## License
MIT

## Support
For help, see [KeyMint API docs](https://docs.keymint.dev) or open an issue.
