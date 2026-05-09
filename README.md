# Keymint Go

A professional, production-ready SDK for integrating with the Keymint API in Go. Provides robust access to all major Keymint features, with idiomatic Go error handling.

## Features
- **Idiomatic Go**: Clean, type-safe API using a `Client` struct and `keymint.New()` constructor.
- **Comprehensive**: Complete API coverage for all Keymint endpoints.
- **Consistent error handling**: All API errors are returned as `*ApiError`.
- **Security**: Credentials are always loaded from environment variables.

## Installation
Add the SDK to your project:

```bash
go get github.com/keymint-dev/keymint-go
```

## Usage

```go
package main

import (
    "os"
    "fmt"
    "github.com/keymint-dev/keymint-go"
)

func main() {
    accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
    productId := os.Getenv("KEYMINT_PRODUCT_ID")

    client, err := keymint.New(accessToken, "")
    if err != nil {
        panic(err)
    }

    // Create a key with authorized hosts
    key, err := client.CreateKey(keymint.CreateKeyParams{
        ProductID:    productId,
        AllowedHosts: []string{"machine-a"},
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Created Key: %s\n", key.Key)
}
```

## Error Handling
All SDK methods return a result and an error. If the error is from the API, it will be of type `*ApiError` with `Message`, `Code`, and `Status` fields.

## API Methods

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

## License
MIT

## Support
For help, see [Keymint API docs](https://docs.keymint.dev) or open an issue.
