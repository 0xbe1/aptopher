// Package main demonstrates querying account information from Aptos.
package main

import (
	"context"
	"fmt"
	"log"

	aptos "github.com/0xbe1/aptopher"
)

func main() {
	// Create a client connected to mainnet
	client, err := aptos.NewClient(aptos.MainnetConfig)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Query the core framework account (0x1)
	address := aptos.AccountOne
	fmt.Printf("Querying account: %s\n", address.ShortString())

	ctx := context.Background()

	// Get account info
	account, err := client.GetAccount(ctx, address)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}

	fmt.Printf("Sequence Number: %s\n", account.Data.SequenceNumber)
	fmt.Printf("Authentication Key: %s\n", account.Data.AuthenticationKey)
	fmt.Printf("Ledger Version: %d\n", account.Metadata.LedgerVersion)
	fmt.Printf("Chain ID: %d\n", account.Metadata.ChainID)

	// Get account resources (limit to first 5 for brevity)
	resources, err := client.GetAccountResources(ctx, address, aptos.WithLimit(5))
	if err != nil {
		log.Fatalf("Failed to get resources: %v", err)
	}

	fmt.Printf("\nFirst %d resources:\n", len(resources.Data))
	for _, r := range resources.Data {
		fmt.Printf("  - %s\n", r.Type)
	}

	// Get ledger info
	ledgerInfo, err := client.GetLedgerInfo(ctx)
	if err != nil {
		log.Fatalf("Failed to get ledger info: %v", err)
	}

	fmt.Printf("\nLedger Info:\n")
	fmt.Printf("  Chain ID: %d\n", ledgerInfo.Data.ChainID)
	fmt.Printf("  Epoch: %s\n", ledgerInfo.Data.Epoch)
	fmt.Printf("  Block Height: %s\n", ledgerInfo.Data.BlockHeight)
	fmt.Printf("  Ledger Version: %s\n", ledgerInfo.Data.LedgerVersion)
}
