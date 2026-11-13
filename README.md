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

Keymint API keys have a scope. Use two different keys, in two different places:

- A **`client`** key, bound to one product, goes in the app you ship. It can
  only activate and deactivate devices and run floating-license sessions.
- An **`admin`** key stays on your server. It creates, updates, signs, blocks
  and deletes keys and manages customers. **Never ship an admin key in an
  app.**

Anything compiled into an app can be extracted from it, so treat any key you
ship as public.

### In your app (client key)

Get this machine's installation ID, then activate. Calling `ActivateKey` again
with the same `HostID` reuses that device's slot, so you can call it at
start-up to check the license is still valid.

```go
package main

import (
    "fmt"

    "github.com/keymint-dev/keymint-go"
)

const (
    // A client-scoped API key bound to this product. Never put an admin key here.
    clientAPIKey = "YOUR_PRODUCT_BOUND_CLIENT_KEY"
    productId    = "YOUR_PRODUCT_ID"
)

func main() {
    client, err := keymint.New(clientAPIKey, "")
    if err != nil {
        panic(err)
    }

    // 1. Get this installation's ID (see Machine Identity below)
    hostId, err := keymint.GetOrCreateInstallationID("")
    if err != nil {
        panic(err)
    }

    licenseKey := "the key your user entered"

    // 2. Activate (or re-check) the license on this machine
    res, err := client.ActivateKey(keymint.ActivateKeyParams{
        ProductID:  productId,
        LicenseKey: licenseKey,
        HostID:     keymint.String(hostId),
    })
    if err != nil {
        fmt.Printf("Not activated: %v\n", err)
        return
    }
    fmt.Println("Activated:", res.Message)

    // 3. Later, free this machine's slot (for example when the user moves to
    //    another computer). Client keys must pass HostID.
    if _, err := client.DeactivateKey(keymint.DeactivateKeyParams{
        ProductID:  productId,
        LicenseKey: licenseKey,
        HostID:     keymint.String(hostId),
    }); err != nil {
        fmt.Printf("Deactivation failed: %v\n", err)
    }
}
```

`FloatingCheckout`, `FloatingHeartbeat` and `FloatingCheckin` also work with a
client key.

### On your server (admin key)

Creating, updating, signing, blocking and deleting keys, and managing
customers, need an admin key. Run this code on infrastructure you control and
load the key from your server's environment or secret store.

```go
package main

import (
    "fmt"
    "os"

    "github.com/keymint-dev/keymint-go"
)

func main() {
    // Admin key: server only. Never ship this key in an app.
    client, err := keymint.New(os.Getenv("KEYMINT_ADMIN_API_KEY"), "")
    if err != nil {
        panic(err)
    }
    productId := os.Getenv("KEYMINT_PRODUCT_ID")

    // Create a customer
    customer, err := client.CreateCustomer(keymint.CreateCustomerParams{
        Name:  "Ada Lovelace",
        Email: "ada@example.com",
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Println("Created customer:", customer.Message)

    // Create a key. AllowedHosts is optional: to lock the key to one
    // machine, pass the host ID your app reported to your server.
    key, err := client.CreateKey(keymint.CreateKeyParams{
        ProductID: productId,
        // AllowedHosts: []string{hostIdFromYourApp},
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Created Key: %s\n", key.Key)

    // Block a key, and unblock it again
    if _, err := client.BlockKey(keymint.BlockKeyParams{ProductID: productId, LicenseKey: key.Key}); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    if _, err := client.UnblockKey(keymint.UnblockKeyParams{ProductID: productId, LicenseKey: key.Key}); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

Customer management (`CreateCustomer`, `UpdateCustomer`, `ToggleCustomerStatus`,
`DeleteCustomer`) also needs an admin key. The read calls (`GetKey`,
`GetAllCustomers`, `GetCustomerById`, `GetCustomerWithKeys`) need at least a
`read-only` key, which a `client` key is not, so keep those on your server too.

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
   modal, or call `SignKey` **on your server**: it needs an admin key. While
   the machine is online, your app sends its host ID to your server, which
   signs the license and returns the file:

   ```go
   // On your server, with an admin-key client:
   signed, err := adminClient.SignKey(keymint.SignKeyParams{
       ProductID:  productId,
       LicenseKey: licenseKey,
       HostID:     hostIdFromYourApp, // the machine code the license is bound to
   })
   // signed.File holds the .lic contents; send it to the app to save.
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
Keymint provides utilities to identify an installation for node-locking:

- `keymint.GetOrCreateInstallationID(storagePath)`: **Recommended** as the
  `HostID`. On first call it combines a random UUID, the `GetMachineID()`
  fingerprint (empty if none could be read) and the current time, and saves
  that value to `storagePath` (default `~/.keymint/installation-id`, or the
  OS temp directory if the home directory can't be found). It returns the
  SHA-256 hash of the saved value, so later calls return the same ID across
  restarts and app updates, and hardware changes don't affect it. It is
  random, not derived from the hardware: a new ID is created if the file is
  deleted, after an OS reinstall, or for another OS user (each has their own
  home directory). If the file can't be written, no error is returned and
  each call produces a new ID.
- `keymint.GetMachineID()`: A best-effort SHA-256 fingerprint of the first
  usable value among the BIOS/hardware UUID, the OS machine ID (on macOS, the
  hardware serial number) and the first active network interface's MAC
  address. It can change after hardware changes or an OS reinstall, can
  collide on virtual machines, and is empty if nothing could be read.

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
| `GetCustomerWithKeys`   | Lists the license keys belonging to a customer.  |
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

## Error Handling

API errors are returned as `*keymint.ApiError`. Branch on `ErrorCode`, a stable string such as `license_invalid_or_forbidden` or `rate_limited`, not on `Message`, which is written for people and may be reworded. See the [error code reference](https://keymint.dev/docs/api-reference/error-codes) for the full list.

```go
res, err := client.ActivateKey(keymint.ActivateKeyParams{
    ProductID:  productId,
    LicenseKey: licenseKey,
    HostID:     keymint.String(hostId),
})
var apiErr *keymint.ApiError
switch {
case err == nil:
    fmt.Println("Activated:", res.Message)
case errors.As(err, &apiErr):
    switch apiErr.ErrorCode {
    case "license_invalid_or_forbidden":
        // Rejected: ask the user to re-enter the key or contact support.
    case "activation_limit_reached":
        // Ask the user to deactivate another device.
    case "rate_limited", "service_unavailable":
        // No answer yet: keep the last known state and retry later.
    default:
        // Status is nil when no HTTP response arrived, so don't dereference it blindly.
        fmt.Printf("Keymint error %s (code %d): %s\n", apiErr.ErrorCode, apiErr.Code, apiErr.Message)
    }
}
```

`ErrorCode` is empty when the response has none, for example a non-JSON body from a proxy, where `Code` is `-1`. `Code` (numeric) and `ErrorDetail` / `ErrorText` (the legacy `error` field) are still set as before.

## Idempotency

Mutating SDK methods accept an idempotency key so you can safely retry a request after a network drop. Pass a pointer to a `keymint.RequestOptions` struct as the optional variadic argument:

```go
// On your server: CreateKey needs an admin key.
key, err := client.CreateKey(keymint.CreateKeyParams{
    ProductID: productId,
}, &keymint.RequestOptions{
    IdempotencyKey: "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
})
```

Use a fresh random UUID for each operation and reuse it only to retry that same operation. Never derive it from the license key or host ID: a replay returns the original response, which can include a newly created license key.

The floating-license methods (`FloatingCheckout`, `FloatingHeartbeat`, `FloatingCheckin`) are the exception: the server ignores `Idempotency-Key` on them, so nothing is stored or replayed and a retry is a new request. Each heartbeat and check-in spends the single-use nonce it carries. If a checkout response is lost, the seat stays reserved until its lease expires.

## License
MIT

## Support
For help, see [Keymint API docs](https://keymint.dev/docs) or open an issue.
