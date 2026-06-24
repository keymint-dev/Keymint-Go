# Keymint Go

A professional, production-ready SDK for integrating with the Keymint API in Go. Provides robust access to all major Keymint features, with idiomatic Go error handling.

## Features
- **Idiomatic Go**: Clean, type-safe API using a `Client` struct and `keymint.New()` constructor.
- **Comprehensive**: Complete API coverage for all Keymint endpoints.
- **Consistent error handling**: All API errors are returned as `*ApiError`.
- **Machine Identity**: Built-in utilities for hardware fingerprinting and stable installation IDs.

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
    apiKey := os.Getenv("KEYMINT_API_KEY")
    productId := os.Getenv("KEYMINT_PRODUCT_ID")

    client, err := keymint.New(apiKey, "")
    if err != nil {
        panic(err)
    }

    // 1. Get a stable, unique ID for this machine
    hostId, err := keymint.GetOrCreateInstallationID("")
    if err != nil {
        panic(err)
    }

    // 2. Create a key authorized only for this machine
    key, err := client.CreateKey(keymint.CreateKeyParams{
        ProductID:    productId,
        AllowedHosts: []string{hostId},
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Created Key: %s\n", key.Key)
}
```

## Machine Identity
Keymint provides utilities to uniquely identify machines for node-locking:

- `keymint.GetOrCreateInstallationID(storagePath)`: **Recommended.** Generates a stable UUID anchored to hardware and persists it to `~/.keymint/installation-id`.
- `keymint.GetMachineID()`: Generates a SHA-256 fingerprint based on BIOS UUID, OS machine ID, and MAC address.

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
