package main

import (
	"context"
	"fmt"
	"log"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
)

// in this example we receive and print consensus info.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// sending empty context, can also use any other context instead.
	resp, err := client.GetConsensusInfo(context.TODO())
	if err != nil {
		log.Fatalf("failed to get consensus info, err: %v", err)
	}

	fmt.Println("consensus info: ", resp.String())
}
