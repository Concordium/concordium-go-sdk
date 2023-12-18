package main

import (
	"context"
	"fmt"
	"io"
	"log"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
)

// This example retrieves and prints the bakers that won the lottery in a particular historical epoch.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("Failed to instantiate client, err: %v", err)
	}

	// sending empty context, can also use any other context instead.
	stream, err := client.GetWinningBakersEpoch(context.TODO(), v2.EpochRequestRelativeEpoch{
		GenesisIndex: v2.GenesisIndex{Value: 3},
		Epoch:        v2.Epoch{Value: 5},
	})
	if err != nil {
		log.Fatalf("failed to get winning bakers, err: %v", err)
	}

	for err == nil {
		winningBaker, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// All WinningBakers recieved.
				break
			}
			log.Fatalf("Could not receive winning baker, err: %v", err)
		}
		fmt.Println("Winning baker: ", winningBaker)
	}
}
