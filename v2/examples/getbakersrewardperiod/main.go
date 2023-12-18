package main

import (
	"context"
	"fmt"
	"io"
	"log"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
)

// This example retrieves and prints the info of the bakers in the reward period of a block.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("Failed to instantiate client, err: %v", err)
	}

	// sending empty context, can also use any other context instead.
	stream, err := client.GetBakersRewardPeriod(context.TODO(), v2.BlockHashInputBest{})
	if err != nil {
		log.Fatalf("failed to get BakerRewardPeriodInfos, err: %v", err)
	}

	for err == nil {
		bakerRewardPeriodInfo, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// All BakerRewardPeriodInfo recieved.
				break
			}
			log.Fatalf("Could not receive BakerRewardPeriodInfo, err: %v", err)
		}
		fmt.Println("BakerRewardPeriodInfo: ", bakerRewardPeriodInfo)
	}
}
