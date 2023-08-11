package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

// in this example we receive and print all accounts in base58 format that we received from specific block but absolute height.
func main() {

	creds, err := v2.HostTLSRoots()

	if err != nil {
		return
	}

	client, err := v2.NewClient(v2.Config{NodeAddress: "grpc.testnet.concordium.com:20000", TlsCredentials: creds})

	// sending empty context, can also use any other context instead, best block can be changed with any other type that implement isBlockHashInput.
	accounts, err := client.GetAccountList(context.TODO(), v2.BlockHashInputBest{})
	if err != nil {
		log.Fatalf("failed to get accounts stream, err: %v", err)
	}

	// print all accounts as base58.
	for i := 0; i < len(accounts); i++ {
		fmt.Println("account: ", accounts[i].ToBase58())
	}
}
