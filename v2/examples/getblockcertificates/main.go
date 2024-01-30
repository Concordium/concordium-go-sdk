package main

import (
	"context"
	"fmt"
	"log"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
)

// This example retrieves and prints the blockcertificates of a non-genesis block.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("Failed to instantiate client, err: %v", err)
	}

	// sending empty context, can also use any other context instead.
	resp, err := client.GetBlockCertificates(context.TODO(), v2.BlockHashInputBest{})
	if err != nil {
		log.Fatalf("failed to get certificates, err: %v", err)
	}

	fmt.Println("certificates: ", resp)
}
