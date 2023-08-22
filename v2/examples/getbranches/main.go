package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Concordium/concordium-go-sdk/v2"
)

// in this example we receive and print blocks branching of this block.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// sending empty context, can also use any other context instead.
	resp, err := client.GetBranches(context.TODO())
	if err != nil {
		log.Fatalf("failed to get branch, err: %v", err)
	}

	fmt.Println("branch: ", resp.String())
}
