package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	respStream, err := client.GetBlocks(context.TODO())
	if err != nil {
		log.Fatalf("failed to get stream of blocks, err: %v", err)
	}

	block, err := respStream.Recv()
	if err != nil {
		log.Fatalf("failed to recieve block, err: %v", err)
	}

	fmt.Println("block: ", block.String())
}
