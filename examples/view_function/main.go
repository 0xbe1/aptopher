// Package main demonstrates calling view functions on Aptos.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	aptos "github.com/0xbe1/lets-go-aptos"
)

func main() {
	// Create a client connected to mainnet
	client, err := aptos.NewClient(aptos.MainnetConfig)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Get the coin balance of an account using view function
	address := aptos.AccountOne
	fmt.Printf("Querying APT balance for: %s\n", address.ShortString())

	result, err := client.View(ctx, aptos.ViewRequest{
		Function:      "0x1::coin::balance",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		Arguments:     []interface{}{address.String()},
	})
	if err != nil {
		log.Fatalf("Failed to call view function: %v", err)
	}

	if len(result.Data) > 0 {
		var balance string
		if err := json.Unmarshal(result.Data[0], &balance); err != nil {
			log.Fatalf("Failed to parse balance: %v", err)
		}
		fmt.Printf("APT Balance: %s (in octas)\n", balance)
	}

	// Example 2: Check if an account exists
	fmt.Println("\nChecking if account 0x1 exists...")
	existsResult, err := client.View(ctx, aptos.ViewRequest{
		Function:      "0x1::account::exists_at",
		TypeArguments: []string{},
		Arguments:     []interface{}{address.String()},
	})
	if err != nil {
		log.Fatalf("Failed to call view function: %v", err)
	}

	if len(existsResult.Data) > 0 {
		var exists bool
		if err := json.Unmarshal(existsResult.Data[0], &exists); err != nil {
			log.Fatalf("Failed to parse result: %v", err)
		}
		fmt.Printf("Account exists: %v\n", exists)
	}

	// Example 3: Get supply of APT coin
	fmt.Println("\nQuerying APT coin supply...")
	supplyResult, err := client.View(ctx, aptos.ViewRequest{
		Function:      "0x1::coin::supply",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		Arguments:     []interface{}{},
	})
	if err != nil {
		log.Fatalf("Failed to call view function: %v", err)
	}

	if len(supplyResult.Data) > 0 {
		fmt.Printf("APT Supply: %s\n", string(supplyResult.Data[0]))
	}
}
