package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

// in this example we receive and print all block transaction events in base58 format.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// TODO: add events to specific block and query it.
	blockItemSummaries, err := client.GetBlockTransactionEvents(context.TODO(), v2.BlockHashInputBest{})
	if err != nil {
		log.Fatalf("failed to get block transaction events, err: %v", err)
	}

	// print all events.
	for i := 0; i < len(blockItemSummaries); i++ {
		fmt.Println("event: ", blockItemSummaries[i].String())
	}
}
