package main

import (
	"context"
	"fmt"
	"log"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// This example retrieves and prints projected earliest wintime of a baker.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("Failed to instantiate client, err: %v", err)
	}

	// sending empty context, can also use any other context instead.
	resp, err := client.GetBakerEarliestWinTime(context.TODO(), &pb.BakerId{
		Value: 1,
	})
	if err != nil {
		log.Fatalf("failed to get wintime, err: %v", err)
	}

	fmt.Println("wintime: ", resp)
}
