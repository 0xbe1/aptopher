// Package aptos provides a Go SDK for the Aptos blockchain.
//
// This SDK allows you to interact with the Aptos blockchain through its REST API,
// including reading account data, executing view functions, building transactions,
// and submitting signed transactions.
//
// # Quick Start
//
// Create a client connected to a network:
//
//	client, err := aptos.NewClient(aptos.MainnetConfig)
//	// or aptos.TestnetConfig, aptos.DevnetConfig, aptos.LocalnetConfig
//
// Query account information:
//
//	account, err := client.GetAccount(ctx, address)
//	fmt.Println(account.Data.SequenceNumber)
//
// Execute a view function:
//
//	result, err := client.View(ctx, aptos.ViewRequest{
//	    Function: "0x1::coin::balance",
//	    TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
//	    Arguments: []interface{}{address.String()},
//	})
//
// Submit a transaction:
//
//	account, _ := aptos.AccountFromEd25519Seed(privateKey)
//	rawTxn, _ := client.BuildTransaction(ctx, account.Address, payload)
//	signedTxn, _ := account.SignTransaction(rawTxn)
//	txnBytes, _ := signedTxn.Bytes()
//	pending, _ := client.SubmitTransaction(ctx, txnBytes)
//	txn, _ := client.WaitForTransactionByHash(ctx, pending.Data.Hash)
//
// # Package Structure
//
// The SDK is organized as follows:
//
//   - aptos: Main package with Client, Account, and core types
//   - aptos/bcs: Binary Canonical Serialization for transaction encoding
//   - aptos/crypto: Cryptographic primitives (Ed25519, Secp256k1)
//   - aptos/examples: Runnable examples
//
// # Response Metadata
//
// All API responses are wrapped in Response[T] which includes both the data
// and metadata from Aptos API headers:
//
//	type Response[T any] struct {
//	    Data     T
//	    Metadata ResponseMetadata
//	}
//
//	type ResponseMetadata struct {
//	    ChainID       uint8
//	    LedgerVersion uint64
//	    Epoch         uint64
//	    BlockHeight   uint64
//	    // ... other fields
//	}
//
// # Error Handling
//
// API errors are returned as *APIError and can be checked using errors.Is:
//
//	_, err := client.GetAccount(ctx, address)
//	if errors.Is(err, aptos.ErrAccountNotFound) {
//	    // Handle missing account
//	}
//
// # Transaction Building
//
// Build, sign, and submit transactions step by step:
//
//	rawTxn, err := client.BuildTransaction(ctx, sender, payload)
//	signedTxn, err := account.SignTransaction(rawTxn)
//	txnBytes, err := signedTxn.Bytes()
//	pending, err := client.SubmitTransaction(ctx, txnBytes)
package aptos
