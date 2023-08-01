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

	nonce := v2.SequenceNumber{
		Value: 1,
	}
	expiry := v2.TransactionTime{
		Value: uint64(time.Now().Add(time.Hour).UTC().Unix()),
	}
	amount := v2.Amount{
		Value: 10,
	}

	signer := v2.NewSimpleSigner(privateKey)

	senderAccount, err := v2.AccountAddressFromBytes(sender)
	if err != nil {
		log.Fatalf("failed to receive account from bytes, err: %v", err)
	}

	receiverAccount, err := v2.AccountAddressFromBytes(receiver)
	if err != nil {
		log.Fatalf("failed to receive account from bytes, err: %v", err)
	}

	accountTx, err := send.Transfer(
		signer,
		senderAccount,
		nonce,
		expiry,
		receiverAccount,
		amount,
	)
	if err != nil {
		log.Fatalf("failed to costruct account transfer, err: %v", err)
	}

	signature := &pb.AccountTransactionSignature{Signatures: map[uint32]*pb.AccountSignatureMap{0: {Signatures: map[uint32]*pb.Signature{
		0: {Value: accountTx.Signature.Signatures[0].Signatures[0].Value},
	}}}}

	txHash, err := client.SendBlockItem(context.Background(), &pb.SendBlockItemRequest{
		BlockItem: &pb.SendBlockItemRequest_AccountTransaction{
			AccountTransaction: &pb.AccountTransaction{
				Signature: signature, Header: &pb.AccountTransactionHeader{
					Sender: &pb.AccountAddress{
						Value: accountTx.Header.Sender.Value[:],
					}, SequenceNumber: &pb.SequenceNumber{
						Value: accountTx.Header.SequenceNumber.Value,
					},
					EnergyAmount: &pb.Energy{
						Value: accountTx.Header.EnergyAmount.Value,
					}, Expiry: &pb.TransactionTime{
						Value: accountTx.Header.Expiry.Value,
					},
				},
				Payload: &pb.AccountTransactionPayload{
					Payload: &pb.AccountTransactionPayload_RawPayload{
						RawPayload: accountTx.Payload.Payload.Encode().Value,
					}},
			},
		}})
	if err != nil {
		log.Fatalf("failed to send block item, err: %v", err)
	}
	fmt.Println("transaction hash: ", txHash.Hex())

	// wait till transaction status turns Finalized
	time.Sleep(25 * time.Second)

	status, err := client.GetBlockItemStatus(context.Background(), *txHash)
	if err != nil {
		log.Fatalf("failed to get block item status, err: %v", err)
	}

	switch v := status.Status.(type) {
	case *pb.BlockItemStatus_Finalized_:
		var bh [32]byte
		copy(bh[:], v.Finalized.Outcome.BlockHash.Value)

		// verify that transaction exists on block
		items, err := client.GetBlockItems(context.Background(), v2.BlockHashInputGiven{
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
