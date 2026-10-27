# Keymint Go

License key validation, activation, and management for Go. Supports node-locking, offline licensing, and hardware fingerprinting. Production-ready SDK for integrating with the Keymint API with idiomatic Go error handling.

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

### Upgrading from v1.5.0 or earlier

Import the module root, `github.com/keymint-dev/keymint-go`. The old
`github.com/keymint-dev/keymint-go/src` path still compiles but is deprecated
and will be removed in a future major version.

Optional fields take pointers; use the helpers `keymint.String`, `keymint.Int`
and `keymint.Bool`, for example `HostID: keymint.String(hostId)`.

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

## Showing The Expiry Date

`ActivateKey` and `FloatingCheckout` return `ExpirationDate`: when the license
expires, as an ISO 8601 string, or `nil` if it never expires. On
`FloatingCheckoutResponse` this is the license expiry; `ExpiresAt` is the
floating seat lease.

```go
res, err := client.ActivateKey(keymint.ActivateKeyParams{
    ProductID:  productId,
    LicenseKey: licenseKey,
    HostID:     keymint.String(hostId),
})
if err != nil {
    // handle error
}

if res.ExpirationDate == nil {
    fmt.Println("License never expires")
} else if exp, err := time.Parse(time.RFC3339, *res.ExpirationDate); err == nil {
    daysLeft := int(math.Ceil(time.Until(exp).Hours() / 24))
    if daysLeft < 0 {
        daysLeft = 0
    }
    fmt.Printf("Expires %s (%d days left)\n", exp.Format("2006-01-02"), daysLeft)
}
```

Store the last `ExpirationDate` your app received. An expired license fails
activation with the same generic error as an invalid key, so the stored date
lets your app tell the user the license expired on that date rather than
calling the key invalid.

## Offline Licenses

An offline license is a signed `.lic` file your app can check with no network
call. `VerifyOfflineLicense` verifies it locally with the standard library
(Ed25519); no extra dependencies.

1. **Get the `.lic` file.** Sign it from the dashboard's Offline License
   modal, or call `SignKey` while the machine is online:

   ```go
   signed, err := client.SignKey(keymint.SignKeyParams{
       ProductID:  productId,
       LicenseKey: licenseKey,
       HostID:     hostId, // the machine code the license is bound to
   })
   // signed.File holds the .lic contents; save it next to your app.
   ```

2. **Bundle your public key.** Copy the PEM from **Developer → Public Keys**
   in the dashboard and ship it with your app (it is public; never ship the
   API key).

3. **Verify on start-up**, passing the product ID and this machine's code:

   ```go
   //go:embed keymint-public-key.pem
   var publicKeyPEM string

   lic, err := keymint.VerifyOfflineLicense(licFileContents, keymint.OfflineVerifyOptions{
       PublicKeys:  []string{publicKeyPEM},
       ProductID:   productId,
       MachineCode: hostId,
       Now:         trustedNow(), // optional, see below
   })
   var licErr *keymint.OfflineLicenseError
   switch {
   case err == nil:
       fmt.Printf("Licensed %s until %s\n", lic.ProductID, lic.ExpiresAt.Format("2006-01-02"))
   case errors.As(err, &licErr):
       switch licErr.Code {
       case keymint.OfflineExpired:
           // Ask the user to go online and re-sign (SignKey) to renew.
       case keymint.OfflineWrongMachine:
           // The file was signed for another machine.
       case keymint.OfflineWrongProduct:
           // The file is for a different product.
       default:
           // OfflineMalformed, OfflineUnsupportedAlgorithm, OfflineUnknownKey,
           // OfflineBadSignature, OfflineWrongType: not a valid license file.
       }
   case errors.Is(err, keymint.ErrInvalidOfflineOptions):
       // Programming error: missing/invalid public key, empty ProductID or
       // MachineCode, or ClockSkew outside 0..1h.
       panic(err)
   }
   ```

`ProductID` is required because one signing key covers every product in your
organization: without the check, a license for one product would unlock
another. Comparisons are exact (no trimming or case folding).

`license` may be the `.lic` file contents or the bare signed token. Pass more
than one key in `PublicKeys` to accept licenses across a key rotation.
`ClockSkew` (default 60s, max 1h; pass a pointer to 0 for none) is the leeway
allowed after expiry.

**Clock rollback.** The SDK can't tell if the user set the system clock back.
Store the latest time your app has trusted (for example after each successful
check or online call) and pass the later of `time.Now()` and that stored time as `Now`.

An offline license stays valid until it expires; blocking the key online does
not revoke a file already issued.

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
| `UpdateKey`      | Updates an existing license key.                |
| `BlockKey`       | Blocks a license key.                           |
| `UnblockKey`     | Unblocks a previously blocked license key.      |
| `DeleteKey`      | Permanently deletes a key (`BlockKey` is reversible). |
| `SignKey`        | Signs a key for offline (air-gapped) validation.|
| `FloatingCheckout` | Checks out a floating license seat.            |
| `FloatingHeartbeat`| Sends a heartbeat to keep a session alive.     |
| `FloatingCheckin`  | Checks in a session, releasing the seat.       |

### Customer Management

| Method                  | Description                                      |
|-------------------------|--------------------------------------------------|
| `CreateCustomer`        | Creates a new customer.                          |
| `GetAllCustomers`       | Retrieves all customers.                         |
| `GetCustomerById`       | Gets a specific customer by ID.                  |
| `GetCustomerWithKeys`   | Gets a customer along with their license keys.   |
| `UpdateCustomer`        | Updates customer information.                    |
| `ToggleCustomerStatus`  | Toggles customer active status.                  |
| `DeleteCustomer`        | Deletes a customer. Their license keys are kept; use `DeleteKey` to remove them. |

### Webhook Verification

| Method                  | Description                                      |
|-------------------------|--------------------------------------------------|
| `VerifyWebhookSignature`| Verifies the signature of a webhook request payload. |

### Offline Verification

| Function                | Description                                      |
|-------------------------|--------------------------------------------------|
| `VerifyOfflineLicense`  | Verifies a signed `.lic` file locally, with no network call. |
| `OfflineKeyID`          | Returns the key ID Keymint computes for a public key PEM. |

## Idempotency

All mutating SDK methods support idempotency keys to safely retry requests in case of network drops. Pass a pointer to a `keymint.RequestOptions` struct as the optional variadic argument:

```go
key, err := client.CreateKey(keymint.CreateKeyParams{
    ProductID: productId,
}, &keymint.RequestOptions{
    IdempotencyKey: "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
})
```

Use a fresh random UUID for each operation and reuse it only to retry that same operation. Never derive it from the license key or host ID: a replay returns the original response, including floating-session secrets.

## License
MIT

## Support
For help, see [Keymint API docs](https://docs.keymint.dev) or open an issue.
