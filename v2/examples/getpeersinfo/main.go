package main

import (
	"context"
	"fmt"
	"log"

	"concordium-go-sdk/v2"
)

func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	resp, err := client.GetPeersInfo(context.TODO())
	if err != nil {
		log.Fatalf("failed to get peers info, err: %v", err)
	}

	fmt.Println("peers info: ", resp.String())
}
