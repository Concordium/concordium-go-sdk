package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// in this example we receive stream of events, receive all data and print.
// it is possible to receive 0 events on block hash input best
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// TODO: create some static events on specific block and hardcode.
	blockTxEventsStream, err := client.GetBlockTransactionEvents(context.TODO(), v2.BlockHashInputBest{})
	if err != nil {
		log.Fatalf("failed to get block transaction events, err: %v", err)
	}

	var totalSummaries []*pb.BlockItemSummary

	for err == nil {
		blockTxEvent, err := blockTxEventsStream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			log.Fatalf("could not receive tx event, err: %v", err)
		}

		totalSummaries = append(totalSummaries, blockTxEvent)
	}

	// print all events.
	for i := 0; i < len(totalSummaries); i++ {
		fmt.Println("event: ", totalSummaries[i].String())
	}
}
