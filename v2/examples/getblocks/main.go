package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// in this example we receive and print 3 blocks in base58 format.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// sending empty context, can also use any other context instead since we don't need any specific params to be passed with it.
	blocksStream, err := client.GetBlocks(context.TODO())
	if err != nil {
		log.Fatalf("failed to get stream of blocks, err: %v", err)
	}

	var blocks []*pb.ArrivedBlockInfo

	// calling Recv method of received stream 3 times, since it will endlessly receive upcoming blocks (not until EOF as some other methods).
	for i := 0; i < 3; i++ {
		block, err := blocksStream.Recv()
		if err != nil {
			log.Fatalf("failed to receive block, err: %v", err)
		}

		blocks = append(blocks, block)
	}

	// print all blocks.
	for i := 0; i < len(blocks); i++ {
		fmt.Println("arrived block info: ", blocks[i].String())
	}
}
