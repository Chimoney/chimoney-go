# Chimoney Go SDK

A Go client for the [Chimoney API](https://chimoney.io/).

## Installation

```bash
go get github.com/chimoney/chimoney-go
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/chimoney/chimoney-go"
)

func main() {
    client := chimoney.New(
        chimoney.WithAPIKey("your-api-key"),
        chimoney.WithSandbox(true), // Use sandbox environment
    )

    ctx := context.Background()
    
    // Get supported assets
    assets, err := client.Info.GetSupportedAssets(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Supported assets: %+v\n", assets)
}
```

## Features

- Simple, idiomatic Go API
- Supports all Chimoney API endpoints
- Configurable HTTP client
- Sandbox mode support
- Context support for cancellation and timeouts

## Modules

- **Account**: Account verification and management
- **Info**: System information and supported assets
- **MobileMoney**: Mobile money payments and transactions
- **Payouts**: Handle various payout methods
  - Airtime
  - Bank transfers
  - Chimoney transfers
  - Gift cards
  - Crypto payments
- **Redeem**: Redeem and verify Chimoney transactions
- **SubAccount**: Manage sub-accounts
- **Wallet**: Wallet operations and transfers

## Examples

### Mobile Money Payment
```go
momoReq := &mobilemoney.PaymentRequest{
    Amount:      10,
    Currency:    "NGN",
    PhoneNumber: "+2348123456789",
    FullName:    "John Doe",
    Country:     "Nigeria",
    Email:       "john@example.com",
    TxRef:       "tx123",
}
resp, err := client.MobileMoney.MakePayment(ctx, momoReq)
```

### Airtime Payout
```go
airtimes := []payouts.AirtimePayload{
    {
        CountryToSend: "Nigeria",
        PhoneNumber:   "+2348123456789",
        ValueInUSD:    3,
    },
}
resp, err := client.Payouts.Airtime(ctx, airtimes, "")
```

### Wallet Operations
```go
// Get wallet balance
balance, err := client.Wallet.GetBalance(ctx, "")

// List wallets
wallets, err := client.Wallet.List(ctx, "")

// Transfer funds
transfer, err := client.Wallet.Transfer(ctx, "receiver123", "wallet_type")
```

## Testing

The SDK includes comprehensive unit tests. To run all tests:

```bash
go test -v ./test/...
```

### Test Coverage

Currently implemented test coverage:

✅ **Account Module**
- GetAllTransactions
- GetTransactionByID
- GetTransactionsByIssueID
- Transfer
- DeleteUnpaidTransaction

✅ **Info Module**
- GetSupportedAssets
- GetAirtimeCountries
- GetBanks
- GetLocalAmountInUSD
- GetMobileMoneyCodes
- GetUSDInLocalAmount



### Sub-Account Management
```go
// Create sub-account
subAccReq := &subaccount.CreateRequest{
    Name:  "Test Account",
    Email: "test@example.com",
}
subAcc, err := client.SubAccount.Create(ctx, subAccReq)

// List sub-accounts
subAccounts, err := client.SubAccount.List(ctx)
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
