Got it 👍 — I’ll adapt the docs for the **Go SDK** instead of NodeJS/TypeScript. Here’s the updated version:

---

# KeyMint Go SDK

Welcome to the official **KeyMint SDK for Go**!
This library provides a simple and convenient way to interact with the KeyMint API, allowing you to manage license keys and customers for your applications with ease.

## ✨ Features

* **Simple & Intuitive**: Clean Go API that feels natural to use.
* **Comprehensive**: Covers all the essential KeyMint API endpoints.
* **Well-Documented**: Clear examples and descriptions for each method.
* **Error Handling**: Standardized `ApiError` responses for easier debugging.

---

## 🚀 Quick Start

Here’s a complete example showing how to create and activate a license key using the SDK:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/your-org/keymint-go-sdk/keymint"
)

func main() {
    accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
    productId := os.Getenv("KEYMINT_PRODUCT_ID")

    if accessToken == "" || productId == "" {
        log.Fatal("Please set KEYMINT_ACCESS_TOKEN and KEYMINT_PRODUCT_ID environment variables")
    }

    sdk, err := keymint.NewSDK(accessToken, "")
    if err != nil {
        log.Fatalf("Failed to initialize SDK: %v", err)
    }

    // 1. Create a new license key
    createResp, err := sdk.CreateKey(keymint.CreateKeyParams{
        ProductID:      productId,
        MaxActivations: "5",
    })
    if err != nil {
        log.Fatalf("Failed to create key: %v", err)
    }

    licenseKey := createResp.Key
    fmt.Println("Key created:", licenseKey)

    // 2. Activate the license key
    activateResp, err := sdk.ActivateKey(keymint.ActivateKeyParams{
        ProductID:  productId,
        LicenseKey: licenseKey,
        HostID:     "UNIQUE_DEVICE_ID",
    })
    if err != nil {
        log.Fatalf("Failed to activate key: %v", err)
    }

    fmt.Println("Key activated:", activateResp.Message)
}
```

---

## 📦 Installation

```bash
go get github.com/your-org/keymint-go-sdk
```

---

## 🛠️ Usage

### Initialization

Import the SDK and initialize it with your access token.
You can find your access token in your [KeyMint dashboard](https://app.keymint.dev/dashboard/developer/access-tokens).

```go
sdk, err := keymint.NewSDK(os.Getenv("KEYMINT_ACCESS_TOKEN"), "")
if err != nil {
    log.Fatal(err)
}
```

---

### API Methods

#### License Key Management

| Method          | Description                                 |
| --------------- | ------------------------------------------- |
| `CreateKey`     | Creates a new license key.                  |
| `ActivateKey`   | Activates a license key for a device.       |
| `DeactivateKey` | Deactivates a device from a license key.    |
| `GetKey`        | Retrieves detailed information about a key. |
| `BlockKey`      | Blocks a license key.                       |
| `UnblockKey`    | Unblocks a previously blocked license key.  |

#### Customer Management

| Method                 | Description                                    |
| ---------------------- | ---------------------------------------------- |
| `CreateCustomer`       | Creates a new customer.                        |
| `GetAllCustomers`      | Retrieves all customers.                       |
| `GetCustomerById`      | Gets a specific customer by ID.                |
| `GetCustomerWithKeys`  | Gets a customer along with their license keys. |
| `UpdateCustomer`       | Updates an existing customer's information.    |
| `ToggleCustomerStatus` | Toggles a customer's active status.            |
| `DeleteCustomer`       | Permanently deletes a customer and their keys. |

---

## 🚨 Error Handling

All errors are returned as `error`, but if they come from the API, they will be of type `*ApiError`.

```go
resp, err := sdk.GetKey(keymint.GetKeyParams{
    ProductID:  "prod_123",
    LicenseKey: "ABC-DEF-GHI",
})
if err != nil {
    if apiErr, ok := err.(*keymint.ApiError); ok {
        fmt.Println("API Error:", apiErr.Message)
        fmt.Println("Status:", *apiErr.Status)
        fmt.Println("Code:", apiErr.Code)
    } else {
        fmt.Println("Unexpected error:", err)
    }
}
```

---

## 📋 Examples

### Customer Management

```go
// Create a new customer
customer, _ := sdk.CreateCustomer(keymint.CreateCustomerParams{
    Name:  "John Doe",
    Email: "john@example.com",
})

// Get all customers
customers, _ := sdk.GetAllCustomers()

// Get a specific customer by ID
customerById, _ := sdk.GetCustomerById(keymint.GetCustomerByIdParams{
    CustomerID: "customer_123",
})

// Get customer with their license keys
customerWithKeys, _ := sdk.GetCustomerWithKeys(keymint.GetCustomerWithKeysParams{
    CustomerID: customer.Data.ID,
})

// Update customer
updatedCustomer, _ := sdk.UpdateCustomer(keymint.UpdateCustomerParams{
    CustomerID: customer.Data.ID,
    Name:       "John Smith",
    Email:      "john.smith@example.com",
})

// Toggle customer status
toggleResp, _ := sdk.ToggleCustomerStatus(keymint.ToggleCustomerStatusParams{
    CustomerID: customer.Data.ID,
})

// Delete customer permanently
deleteResp, _ := sdk.DeleteCustomer(keymint.DeleteCustomerParams{
    CustomerID: customer.Data.ID,
})
```

---

### Creating a License Key with a New Customer

```go
licenseResp, _ := sdk.CreateKey(keymint.CreateKeyParams{
    ProductID:      os.Getenv("KEYMINT_PRODUCT_ID"),
    MaxActivations: "3",
    NewCustomer: &keymint.NewCustomerParams{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    },
})
```

---

## 🔒 Security Best Practices

**Never hardcode your access tokens!** Always use environment variables.

1. **Create a `.env` file**:

   ```bash
   KEYMINT_ACCESS_TOKEN=your_actual_token_here
   KEYMINT_PRODUCT_ID=your_product_id_here
   ```

2. **Load environment variables**:

   ```bash
   export $(cat .env | xargs)
   ```

3. **Use them in your Go code**:

   ```go
   accessToken := os.Getenv("KEYMINT_ACCESS_TOKEN")
   productId := os.Getenv("KEYMINT_PRODUCT_ID")
   ```

⚠️ **Important**: Never commit `.env` files to version control.

---

## 📚 API Reference

### `NewSDK(accessToken string, baseUrl string) (*SDK, error)`

| Parameter     | Type     | Description                                                      |
| ------------- | -------- | ---------------------------------------------------------------- |
| `accessToken` | `string` | **Required.** Your KeyMint API access token.                     |
| `baseUrl`     | `string` | *Optional.* API base URL. Defaults to `https://api.keymint.dev`. |

### License Key Methods

* `CreateKey(params CreateKeyParams) (*CreateKeyResponse, error)`
* `ActivateKey(params ActivateKeyParams) (*ActivateKeyResponse, error)`
* `DeactivateKey(params DeactivateKeyParams) (*DeactivateKeyResponse, error)`
* `GetKey(params GetKeyParams) (*GetKeyResponse, error)`
* `BlockKey(params BlockKeyParams) (*BlockKeyResponse, error)`
* `UnblockKey(params UnblockKeyParams) (*UnblockKeyResponse, error)`

### Customer Methods

* `CreateCustomer(params CreateCustomerParams) (*CreateCustomerResponse, error)`
* `GetAllCustomers() (*GetAllCustomersResponse, error)`
* `GetCustomerById(params GetCustomerByIdParams) (*GetCustomerByIdResponse, error)`
* `GetCustomerWithKeys(params GetCustomerWithKeysParams) (*GetCustomerWithKeysResponse, error)`
* `UpdateCustomer(params UpdateCustomerParams) (*UpdateCustomerResponse, error)`
* `ToggleCustomerStatus(params ToggleCustomerStatusParams) (*ToggleCustomerStatusResponse, error)`
* `DeleteCustomer(params DeleteCustomerParams) (*DeleteCustomerResponse, error)`

---

## 📜 License

This SDK is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Would you like me to also **add GoDoc-style comments** (`// CreateKey ...`) to your SDK methods so that this documentation auto-generates in `pkg.go.dev`? That way the docs are always in sync with the code.
