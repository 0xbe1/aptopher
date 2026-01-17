# Aptopher

A clean, simple Go SDK for the [Aptos](https://aptos.dev) blockchain.

## Features

- **Full REST API coverage** - All Aptos node API endpoints
- **BCS serialization** - Binary Canonical Serialization for transactions
- **Multiple signature schemes** - Ed25519 and Secp256k1 support
- **Minimal dependencies** - Only `golang.org/x/crypto` and `secp256k1`
- **Simple API** - Clean, idiomatic Go interface
- **Response metadata** - Access to chain ID, ledger version, epoch from headers

## Installation

```bash
go get github.com/0xbe1/aptopher
```

Requires Go 1.21 or later.

## Quick Start

### Connect to a Network

```go
import aptos "github.com/0xbe1/aptopher"

// Mainnet
client, err := aptos.NewClient(aptos.MainnetConfig)

// Testnet
client, err := aptos.NewClient(aptos.TestnetConfig)

// Devnet
client, err := aptos.NewClient(aptos.DevnetConfig)

// Custom endpoint
client, err := aptos.NewClient(aptos.ClientConfig{
    NodeURL: "https://your-node.example.com/v1",
})
```

### Query Account Information

```go
ctx := context.Background()

// Parse an address
address, err := aptos.ParseAccountAddress("0x1")

// Get account info
account, err := client.GetAccount(ctx, address)
fmt.Println("Sequence Number:", account.Data.SequenceNumber)
fmt.Println("Ledger Version:", account.Metadata.LedgerVersion)

// Get account resources
resources, err := client.GetAccountResources(ctx, address)
for _, r := range resources.Data {
    fmt.Println(r.Type)
}

// Get APT balance
balance, err := client.GetAccountBalance(ctx, address, "0x1::aptos_coin::AptosCoin")
```

### Execute View Functions

```go
result, err := client.View(ctx, aptos.ViewRequest{
    Function:      "0x1::coin::balance",
    TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
    Arguments:     []interface{}{address.String()},
})

var balance string
json.Unmarshal(result.Data[0], &balance)
fmt.Println("Balance:", balance)
```

### Submit a Transaction

```go
// Create account from private key (32-byte Ed25519 seed)
account, err := aptos.AccountFromEd25519Seed(privateKeyBytes)

// Create a payload (example: APT transfer)
recipient, _ := aptos.ParseAccountAddress("0x123...")
payload := aptos.TransactionPayload{
    Payload: &aptos.EntryFunction{
        Module: aptos.ModuleId{
            Address: aptos.AccountOne,
            Name:    "aptos_account",
        },
        Function: "transfer",
        TypeArgs: nil,
        Args: [][]byte{
            recipient[:],           // recipient address
            serializeU64(100_000_000), // amount in octas (1 APT)
        },
    },
}

// Build the transaction
rawTxn, err := client.BuildTransaction(ctx, account.Address, payload)

// Sign the transaction
signedTxn, err := account.SignTransaction(rawTxn)
txnBytes, err := signedTxn.Bytes()

// Submit and wait for confirmation
pending, err := client.SubmitTransaction(ctx, txnBytes)
txn, err := client.WaitForTransactionByHash(ctx, pending.Data.Hash)
if txn.Data.Success {
    fmt.Println("Transaction successful!")
}
```

### Orderless Transactions

Orderless transactions use a replay protection nonce instead of a sequence number, allowing multiple transactions to be signed and submitted in any order. This is useful for multi-agent scenarios or when transaction ordering doesn't matter.

```go
import "crypto/rand"

// Generate a random nonce for replay protection
var nonceBuf [8]byte
rand.Read(nonceBuf[:])
nonce := binary.LittleEndian.Uint64(nonceBuf[:])

// Build an orderless transaction
rawTxn, err := client.BuildTransaction(ctx, account.Address, payload,
    aptos.WithReplayProtectionNonce(nonce),
)

// Sign and submit as usual
signedTxn, err := account.SignTransaction(rawTxn)
txnBytes, err := signedTxn.Bytes()
pending, err := client.SubmitTransaction(ctx, txnBytes)
```

**Note:** Orderless transactions have a maximum expiration time of 60 seconds.

### Simulate Transactions

```go
// Build a raw transaction
rawTxn, err := client.BuildTransaction(ctx, account.Address, payload)

// Create a fake signature for simulation
fakeSignedTxn := &aptos.SignedTransaction{
    RawTxn: rawTxn,
    Authenticator: aptos.TransactionAuthenticator{
        Variant: aptos.TransactionAuthenticatorSingleSender,
        Auth: &aptos.AccountAuthenticatorSingleKey{
            PublicKey: aptos.AnyPublicKey{
                Variant:   account.Signer.Scheme(),
                PublicKey: account.Signer.PublicKey(),
            },
            Signature: aptos.AnySignature{
                Variant:   account.Signer.Scheme(),
                Signature: make([]byte, 64), // Zero signature for simulation
            },
        },
    },
}
txnBytes, _ := fakeSignedTxn.Bytes()

// Simulate to estimate gas
result, err := client.SimulateTransaction(ctx, txnBytes,
    aptos.WithEstimateMaxGasAmount(),
    aptos.WithEstimateGasUnitPrice(),
)

fmt.Println("Gas used:", result.Data[0].GasUsed)
fmt.Println("Success:", result.Data[0].Success)
```

## API Reference

### Client Methods

#### General
- `GetLedgerInfo(ctx)` - Get current ledger state
- `GetNodeInfo(ctx)` - Get node information
- `HealthCheck(ctx)` - Check node health
- `EstimateGasPrice(ctx)` - Get gas price estimates

#### Accounts
- `GetAccount(ctx, address)` - Get account info (sequence number, auth key)
- `GetAccountResources(ctx, address)` - List all resources
- `GetAccountResourcesBCS(ctx, address)` - List all resources (BCS format)
- `GetAccountResource(ctx, address, resourceType)` - Get specific resource
- `GetAccountResourceBCS(ctx, address, resourceType)` - Get specific resource (BCS format)
- `GetAccountModules(ctx, address)` - List all modules
- `GetAccountModulesBCS(ctx, address)` - List all modules (BCS format)
- `GetAccountModule(ctx, address, moduleName)` - Get specific module
- `GetAccountModuleBCS(ctx, address, moduleName)` - Get specific module (BCS format)
- `GetAccountBalance(ctx, address, assetType)` - Get coin balance

#### Transactions
- `GetTransactions(ctx)` - List transactions
- `GetTransactionByHash(ctx, hash)` - Get by hash
- `GetTransactionByVersion(ctx, version)` - Get by version
- `GetAccountTransactions(ctx, address)` - Get account's transactions
- `SubmitTransaction(ctx, signedTxnBytes)` - Submit signed transaction
- `SimulateTransaction(ctx, signedTxnBytes)` - Simulate transaction
- `WaitForTransactionByHash(ctx, hash)` - Wait for confirmation (long-polling)
- `PollForTransaction(ctx, hash, interval)` - Poll for confirmation

#### Blocks
- `GetBlockByHeight(ctx, height, withTxns)` - Get block by height
- `GetBlockByVersion(ctx, version, withTxns)` - Get block by version

#### Events
- `GetEventsByCreationNumber(ctx, address, creationNum)` - Get events
- `GetEventsByEventHandle(ctx, address, handle, field)` - Get events by handle

#### Tables
- `GetTableItem(ctx, tableHandle, request)` - Get table item
- `GetTableItemBCS(ctx, tableHandle, request)` - Get table item (BCS format)
- `GetRawTableItem(ctx, tableHandle, request)` - Get raw table item

#### View Functions
- `View(ctx, request)` - Execute view function
- `ViewBCS(ctx, request)` - Execute view function (BCS format)

#### Transaction Building
- `BuildTransaction(ctx, sender, payload)` - Build raw transaction
- `BuildSignAndSubmitTransaction(ctx, account, payload)` - Build, sign, and submit in one call

### Request Options

```go
// Pagination
client.GetAccountResources(ctx, address,
    aptos.WithStart(0),
    aptos.WithLimit(100),
)

// Historical state
client.GetAccount(ctx, address,
    aptos.WithLedgerVersion(12345678),
)

// Transaction building
client.BuildTransaction(ctx, sender, payload,
    aptos.WithMaxGasAmount(50000),
    aptos.WithGasUnitPrice(100),
    aptos.WithSequenceNumber(5),
)

// Orderless transaction (uses nonce instead of sequence number)
client.BuildTransaction(ctx, sender, payload,
    aptos.WithReplayProtectionNonce(12345),
)
```

### Response Metadata

All API responses include metadata from Aptos headers:

```go
resp, err := client.GetAccount(ctx, address)

fmt.Println(resp.Data.SequenceNumber)     // Response data
fmt.Println(resp.Metadata.ChainID)        // Chain ID (1=mainnet, 2=testnet)
fmt.Println(resp.Metadata.LedgerVersion)  // Current ledger version
fmt.Println(resp.Metadata.Epoch)          // Current epoch
fmt.Println(resp.Metadata.BlockHeight)    // Current block height
```

### Error Handling

```go
import "errors"

account, err := client.GetAccount(ctx, address)
if err != nil {
    if errors.Is(err, aptos.ErrAccountNotFound) {
        // Account doesn't exist
    }
    if errors.Is(err, aptos.ErrResourceNotFound) {
        // Resource not found
    }

    // Check specific API error
    var apiErr *aptos.APIError
    if errors.As(err, &apiErr) {
        fmt.Println("Error code:", apiErr.ErrorCode)
        fmt.Println("Message:", apiErr.Message)
    }
}
```

### Typed Argument Builders

Use these helpers to build BCS-encoded entry function arguments:

```go
// Primitives
aptos.BoolArg(true)
aptos.U8Arg(255)
aptos.U16Arg(65535)
aptos.U32Arg(4294967295)
aptos.U64Arg(1000000)
aptos.U128Arg(big.NewInt(1000000))
aptos.U256Arg(big.NewInt(1000000))

// Address and strings
aptos.AddressArg(address)
aptos.StringArg("hello")
aptos.BytesArg([]byte{0x01, 0x02})
aptos.ObjectArg(objectAddress)  // Same as AddressArg

// Vectors
aptos.VectorU8Arg([]byte{1, 2, 3})
aptos.VectorU64Arg([]uint64{100, 200})
aptos.VectorAddressArg([]aptos.AccountAddress{addr1, addr2})
aptos.VectorStringArg([]string{"a", "b"})

// Options (pass nil for None)
aptos.OptionU64Arg(&value)      // Some(value)
aptos.OptionU64Arg(nil)         // None
aptos.OptionAddressArg(&addr)
aptos.OptionStringArg(&str)

// Combine into entry function args
payload := aptos.TransactionPayload{
    Payload: &aptos.EntryFunction{
        Module:   aptos.ModuleId{Address: aptos.AccountOne, Name: "aptos_account"},
        Function: "transfer",
        Args:     aptos.EntryFunctionArgs(aptos.AddressArg(recipient), aptos.U64Arg(amount)),
    },
}
```

## Examples

See the [examples](./examples) directory for complete, runnable examples:

- **[query_account](./examples/query_account)** - Query account data from mainnet
- **[view_function](./examples/view_function)** - Execute view functions
- **[transfer_coin](./examples/transfer_coin)** - Transfer APT on devnet
- **[simulate_transaction](./examples/simulate_transaction)** - Simulate and estimate gas

Run an example:

```bash
go run ./examples/query_account
```

## Package Structure

```
github.com/0xbe1/aptopher/
├── bcs/                    # Binary Canonical Serialization
│   ├── serializer.go       # BCS encoding
│   ├── deserializer.go     # BCS decoding
│   └── interfaces.go       # Marshaler/Unmarshaler interfaces
├── crypto/                 # Cryptographic operations
│   ├── ed25519.go          # Ed25519 signing
│   ├── secp256k1.go        # Secp256k1 ECDSA
│   ├── hash.go             # SHA3-256 hashing
│   └── signer.go           # Signer interface
├── examples/               # Runnable examples
├── internal/hex/           # Hex encoding utilities
├── client.go               # Main Client type
├── account_address.go      # AccountAddress type
├── move_types.go           # TypeTag, StructTag, U128, U256
├── transaction.go          # Transaction types
├── transaction_payload.go  # EntryFunction, Script payloads
├── raw_transaction.go      # RawTransaction for signing
├── signed_transaction.go   # SignedTransaction for submission
└── ...
```

## Dependencies

This SDK has minimal dependencies:

- `golang.org/x/crypto` - SHA3-256, Ed25519
- `github.com/decred/dcrd/dcrec/secp256k1/v4` - Secp256k1 ECDSA

## Acknowledgments

This SDK implements the [Aptos Node API specification](https://api.mainnet.aptoslabs.com/v1/spec) and was built with [Claude](https://claude.ai), using the official [aptos-go-sdk](https://github.com/aptos-labs/aptos-go-sdk) as reference.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
