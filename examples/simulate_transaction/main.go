// Package main demonstrates simulating transactions to estimate gas.
package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"

	aptos "github.com/0xbe1/aptopher"
)

func main() {
	// Create a client connected to devnet
	client, err := aptos.NewClient(aptos.DevnetConfig)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Generate a temporary account for simulation
	account, err := aptos.NewEd25519Account()
	if err != nil {
		log.Fatalf("Failed to generate account: %v", err)
	}

	fmt.Printf("Simulating with account: %s\n", account.Address.String())

	// Generate a random recipient
	recipient, err := aptos.NewEd25519Account()
	if err != nil {
		log.Fatalf("Failed to generate recipient: %v", err)
	}

	ctx := context.Background()

	// Create a transfer payload
	payload := aptos.TransactionPayload{
		Payload: &aptos.EntryFunction{
			Module: aptos.ModuleId{
				Address: aptos.AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer",
			TypeArgs: nil,
			Args: [][]byte{
				recipient.Address[:],
				serializeU64(1000),
			},
		},
	}

	// Build a raw transaction
	// Note: For simulation, the sequence number doesn't need to be correct
	rawTxn := &aptos.RawTransaction{
		Sender:                  account.Address,
		SequenceNumber:          0, // Doesn't matter for simulation
		Payload:                 payload,
		MaxGasAmount:            200000,
		GasUnitPrice:            100,
		ExpirationTimestampSecs: 9999999999, // Far future
		ChainID:                 4,          // Devnet
	}

	// Create a fake signature for simulation (all zeros)
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

	txnBytes, err := fakeSignedTxn.Bytes()
	if err != nil {
		log.Fatalf("Failed to serialize transaction: %v", err)
	}

	// Simulate the transaction with gas estimation
	fmt.Println("Simulating APT transfer transaction...")
	result, err := client.SimulateTransaction(ctx, txnBytes,
		aptos.WithEstimateMaxGasAmount(),
		aptos.WithEstimateGasUnitPrice(),
	)
	if err != nil {
		log.Fatalf("Failed to simulate transaction: %v", err)
	}

	if len(result.Data) > 0 {
		simResult := result.Data[0]
		fmt.Printf("\nSimulation Result:\n")
		fmt.Printf("  Success: %v\n", simResult.Success)
		fmt.Printf("  Gas Used: %s\n", simResult.GasUsed)
		fmt.Printf("  Max Gas Amount: %s\n", simResult.MaxGasAmount)
		fmt.Printf("  Gas Unit Price: %s\n", simResult.GasUnitPrice)
		fmt.Printf("  VM Status: %s\n", simResult.VMStatus)

		if len(simResult.Events) > 0 {
			fmt.Printf("  Events (%d):\n", len(simResult.Events))
			for i, event := range simResult.Events {
				fmt.Printf("    %d: %s\n", i+1, event.Type)
			}
		}
	}

	// Get current gas price estimate
	gasEstimate, err := client.EstimateGasPrice(ctx)
	if err != nil {
		log.Printf("Warning: Failed to get gas estimate: %v", err)
	} else {
		fmt.Printf("\nCurrent Gas Price Estimates:\n")
		fmt.Printf("  Normal: %d\n", gasEstimate.Data.GasEstimate)
		fmt.Printf("  Prioritized: %d\n", gasEstimate.Data.PrioritizedGasEstimate)
		fmt.Printf("  Deprioritized: %d\n", gasEstimate.Data.DeprioritizedGasEstimate)
	}
}

func serializeU64(v uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return buf
}
