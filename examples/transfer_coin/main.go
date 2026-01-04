// Package main demonstrates transferring APT coins on Aptos devnet.
//
// Prerequisites:
// 1. Generate a new account or use an existing one
// 2. Fund the account using the devnet faucet: https://aptos.dev/network/faucet
// 3. Set the APTOS_PRIVATE_KEY environment variable
package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	aptos "github.com/0xbe1/lets-go-aptos"
)

func main() {
	// Get private key from environment
	privateKeyHex := os.Getenv("APTOS_PRIVATE_KEY")
	if privateKeyHex == "" {
		// Generate a new account for demonstration
		fmt.Println("No APTOS_PRIVATE_KEY set. Generating a new account...")
		fmt.Println("To use this account, fund it via the devnet faucet and set the private key.")
		fmt.Println()

		account, err := aptos.NewEd25519Account()
		if err != nil {
			log.Fatalf("Failed to generate account: %v", err)
		}

		fmt.Printf("Address: %s\n", account.Address.String())
		fmt.Printf("Auth Key: %s\n", bytesToHex(account.AuthKey()))
		fmt.Println()
		fmt.Println("Fund this account at: https://aptos.dev/network/faucet")
		fmt.Println("Then set APTOS_PRIVATE_KEY and run again.")
		return
	}

	// Parse private key
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// Create account from private key
	account, err := aptos.AccountFromEd25519Seed(privateKey)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}

	fmt.Printf("Sender Address: %s\n", account.Address.String())

	// Create a client connected to devnet
	client, err := aptos.NewClient(aptos.DevnetConfig)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Check sender balance
	senderAccount, err := client.GetAccount(ctx, account.Address)
	if err != nil {
		log.Fatalf("Failed to get sender account: %v", err)
	}
	fmt.Printf("Sender Sequence Number: %s\n", senderAccount.Data.SequenceNumber)

	// Generate a random recipient for demonstration
	recipient, err := aptos.NewEd25519Account()
	if err != nil {
		log.Fatalf("Failed to generate recipient: %v", err)
	}
	fmt.Printf("Recipient Address: %s\n", recipient.Address.String())

	// Transfer 1000 octas (0.00001 APT)
	amount := uint64(1000)
	fmt.Printf("Transferring %d octas...\n", amount)

	// Create the transfer payload
	payload := aptos.TransactionPayload{
		Payload: &aptos.EntryFunction{
			Module: aptos.ModuleId{
				Address: aptos.AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer",
			TypeArgs: nil,
			Args: [][]byte{
				mustSerializeAddress(recipient.Address),
				serializeU64(amount),
			},
		},
	}

	// Build the transaction
	rawTxn, err := client.BuildTransaction(ctx, account.Address, payload)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v", err)
	}

	// Sign the transaction
	signedTxn, err := account.SignTransaction(rawTxn)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Get signed transaction bytes
	txnBytes, err := signedTxn.Bytes()
	if err != nil {
		log.Fatalf("Failed to serialize transaction: %v", err)
	}

	// Submit the transaction
	pending, err := client.SubmitTransaction(ctx, txnBytes)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("Transaction submitted: %s\n", pending.Data.Hash)

	// Wait for transaction to be committed
	fmt.Println("Waiting for transaction...")
	txn, err := client.WaitForTransactionByHash(ctx, pending.Data.Hash)
	if err != nil {
		log.Fatalf("Failed waiting for transaction: %v", err)
	}

	if txn.Data.Success {
		fmt.Printf("Transaction successful! Version: %s\n", txn.Data.Version)
		fmt.Printf("Gas used: %s\n", txn.Data.GasUsed)
	} else {
		fmt.Printf("Transaction failed: %s\n", txn.Data.VMStatus)
	}
}

func bytesToHex(b [32]byte) string {
	return "0x" + hex.EncodeToString(b[:])
}

func mustSerializeAddress(addr aptos.AccountAddress) []byte {
	return addr[:]
}

func serializeU64(v uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return buf
}
