package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// in this example we receive and print all block transaction events in base58 format.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	eventsStream, err := client.GetBlockTransactionEvents(context.TODO(), &pb.BlockHashInput{
		BlockHashInput: &pb.BlockHashInput_AbsoluteHeight{
			AbsoluteHeight: &pb.AbsoluteBlockHeight{
				Value: uint64(3101973),
			}}})
	if err != nil {
		log.Fatalf("failed to get block transaction events, err: %v", err)
	}

	var events []*pb.BlockItemSummary

	// calling Recv method of received stream until we get EOF, to be sure we collected all tx events.
	for err == nil {
		event, err := eventsStream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			log.Fatalf("failed to receive event, err: %v", err)
		}

		events = append(events, event)
	}

	// print all events.
	for i := 0; i < len(events); i++ {
		fmt.Println("event: ", events[i].String())
	}
}
