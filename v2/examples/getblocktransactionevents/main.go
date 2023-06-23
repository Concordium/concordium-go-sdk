package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	txEventsStream, err := client.GetBlockTransactionEvents(context.TODO(), &pb.BlockHashInput{
		BlockHashInput: &pb.BlockHashInput_AbsoluteHeight{
			AbsoluteHeight: &pb.AbsoluteBlockHeight{
				Value: uint64(3101973),
			}}})
	if err != nil {
		log.Fatalf("failed to get block transaction events, err: %v", err)
	}

	blockItemSummary, err := txEventsStream.Recv()
	if err != nil {
		log.Fatalf("failed to receive block item summary, err: %v", err)
	}

	fmt.Println("account: ", blockItemSummary.String())
}
