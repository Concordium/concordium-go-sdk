package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	resp, err := client.GetConsensusInfo(context.TODO())
	if err != nil {
		log.Fatalf("failed to get consensus info, err: %v", err)
	}

	fmt.Println("consensus info: ", resp.String())
}
