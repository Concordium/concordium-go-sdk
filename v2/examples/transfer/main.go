package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/send"
	"github.com/btcsuite/btcutil/base58"
)

// for this example we need to pass signKey, verifyKey and 2 wallet pubkeys - sender and receiver, sender need to have funds to transfer. note: accounts must be on test network.
// sending block item + verify that transaction exists in block.
// to run this example write in terminal in directory of main.go file "go run main.go *signKey* *verifyKey* *senderPubKey* *receiverPubKey*".
// signKey and verifyKey are keys you can download after creating your wallet and copy from file, pubkey you can check in wallet accounts.
// for repeating test nonce (SequenceNumber) need to be incremented by 1 each run.
func main() {
	args := os.Args
	if len(args) < 5 {
		log.Fatalf("require sender and receiver addresses")
	}

	signKey := args[1]
	verifyKey := args[2]
	senderPubKey := args[3]
	receiverPubKey := args[4]

	privateKey, err := hex.DecodeString(signKey + verifyKey)
	if err != nil {
		log.Fatalf("failed to decode private key, err: %v", err)
	}

	client, err := v2.NewClient(v2.Config{
		NodeAddress: "node.testnet.concordium.com:20000",
	})
	if err != nil {
		log.Fatalf("failed to create new client, err: %v", err)
	}

	sender, _, err := base58.CheckDecode(senderPubKey)
	if err != nil {
		log.Fatalf("failed to decode sender, err: %v", err)
	}

	receiver, _, err := base58.CheckDecode(receiverPubKey)
	if err != nil {
		log.Fatalf("failed to decode receiver, err: %v", err)
	}

	ctx := context.Background()
	senderAddr := v2.AccountAddressFromBytes(sender)
	sequenceNumber, err := client.GetNextAccountSequenceNumber(ctx, &senderAddr)
	if err != nil {
		log.Fatalf("failed to get next sender sequnce number, err: %v", err)
	}

	nonce := v2.SequenceNumber{
		Value: sequenceNumber.SequenceNumber.Value,
	}
	expiry := v2.TransactionTime{
		Value: uint64(time.Now().Add(time.Hour).UTC().Unix()),
	}
	amount := v2.Amount{
		Value: 10,
	}

	signer := v2.NewSimpleSigner(privateKey)

	accountTx, err := send.Transfer(
		signer,
		senderAddr,
		nonce,
		expiry,
		v2.AccountAddressFromBytes(receiver),
		amount,
	)
	if err != nil {
		log.Fatalf("failed to costruct account transfer, err: %v", err)
	}

	txHash, err := accountTx.Send(ctx, client)
	if err != nil {
		log.Fatalf("failed to send block item, err: %v", err)
	}
	fmt.Println("transaction hash: ", txHash.Hex())

	// wait till transaction status turns Finalized
	time.Sleep(25 * time.Second)

	status, err := client.GetBlockItemStatus(ctx, *txHash)
	if err != nil {
		log.Fatalf("failed to get block item status, err: %v", err)
	}

	switch v := status.Status.(type) {
	case *pb.BlockItemStatus_Finalized_:
		var bh [32]byte
		copy(bh[:], v.Finalized.Outcome.BlockHash.Value)

		// verify that transaction exists on block
		items, err := client.GetBlockItems(ctx, v2.BlockHashInputGiven{
			Given: v2.BlockHash{
				Value: bh,
			}})
		if err != nil {
			log.Fatalf("failed to get block item, err: %v", err)
		}
		fmt.Println("block item hash:", items[0].Hash.Hex())

		// compare transaction hash value
		if items[0].Hash.Value != txHash.Value {
			log.Fatalf("tx hash not match expected")
		}
	}
}
