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

	respStream, err := client.GetAccountList(context.TODO(), &pb.BlockHashInput{
		BlockHashInput: &pb.BlockHashInput_AbsoluteHeight{
			AbsoluteHeight: &pb.AbsoluteBlockHeight{
				Value: uint64(3101973),
			}}})
	if err != nil {
		log.Fatalf("failed to get accounts stream, err: %v", err)
	}

	account, err := respStream.Recv()
	if err != nil {
		log.Fatalf("failed to receive account, err: %v", err)
	}

	fmt.Println("account: ", account.String())
}
