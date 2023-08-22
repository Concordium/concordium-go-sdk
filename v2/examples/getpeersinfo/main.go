package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Concordium/concordium-go-sdk/v2"
)

// in this example we receive and print peer info.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// sending empty context, can also use any other context instead.
	resp, err := client.GetPeersInfo(context.TODO())
	if err != nil {
		log.Fatalf("failed to get peers info, err: %v", err)
	}

	fmt.Println("peers info: ", resp.String())
}
