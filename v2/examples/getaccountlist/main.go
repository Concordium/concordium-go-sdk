package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

// in this example we receive and print all accounts in base58 format that we received from specific block but absolute height.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// sending empty context, can also use any other context instead, and best block to receive stream of accounts belong to it.
	accounts, err := client.GetAccountList(context.TODO(), v2.BlockHashInputBest{})
	if err != nil {
		log.Fatalf("failed to get accounts stream, err: %v", err)
	}

	// print all accounts as base58.
	for i := 0; i < len(accounts); i++ {
		fmt.Println("account: ", accounts[i].ToBase58())
	}
}
