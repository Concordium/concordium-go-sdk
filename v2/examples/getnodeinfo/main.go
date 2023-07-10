package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

// in this example we receive and print node info.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// sending empty context, can also use any other context instead.
	resp, err := client.GetNodeInfo(context.TODO())
	if err != nil {
		log.Fatalf("failed to get node info, err: %v", err)
	}

	fmt.Println("node info: ", resp.String())
}
