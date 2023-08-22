package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/pb"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/send"
)

type updateParams struct {
	SenderPubKey    string `env:"SENDER_PUBLIC_KEY_IN_BASE58"`
	SenderSignKey   string `env:"SENDER_SIGN_KEY_IN_HEX"`
	SenderVerifyKey string `env:"SENDER_VERIFY_KEY_IN_HEX"`

	ContractIndex    uint64 `env:"CONTRACT_INDEX"`
	ContractSubIndex uint64 `env:"CONTRACT_SUB_INDEX"`
	ReceiveName      string `env:"RECEIVE_NAME"`
	Parameter        string `env:"PARAMETER_IN_HEX"`
	AmountOfCCD      uint64 `env:"AMOUNT_OF_CCD"`
	MaxEnergyAmount  uint64 `env:"MAX_ENERGY_AMOUNT"`
}

// for this example we need to fill .env file (located in the same directory with this main.go file) with values,
// described in struct updateParams above.
// sends block item with update contract payload + verifies that transaction exists in block.
// to run this example write in terminal in directory of main.go file "go run main.go".
// signKey and verifyKey are keys you can download after creating your wallet and copy from file, pubKey you can check in wallet accounts.
func main() {
	if err := godotenv.Overload(); err != nil {
		log.Fatalf("failed to load update params, err: %v", err)
	}

	params := new(updateParams)
	if err := env.Parse(params); err != nil {
		log.Fatalf("failed to parse update params, err: %v", err)
	}

	privateKey, err := hex.DecodeString(params.SenderSignKey + params.SenderVerifyKey)
	if err != nil {
		log.Fatalf("failed to decode private key, err: %v", err)
	}

	client, err := v2.NewClient(v2.Config{
		NodeAddress: "node.testnet.concordium.com:20000",
	})
	if err != nil {
		log.Fatalf("failed to create new client, err: %v", err)
	}

	sender, err := v2.AccountAddressFromString(params.SenderPubKey)
	if err != nil {
		log.Fatalf("failed to decode sender, err: %v", err)
	}

	ctx := context.Background()
	sequenceNumber, err := client.GetNextAccountSequenceNumber(ctx, &sender)
	if err != nil {
		log.Fatalf("failed to get next sender sequnce number, err: %v", err)
	}

	nonce := v2.SequenceNumber{
		Value: sequenceNumber.SequenceNumber.Value,
	}
	expiry := v2.TransactionTime{
		Value: uint64(time.Now().UTC().Add(time.Hour).Unix()),
	}
	amount := v2.Amount{
		Value: params.AmountOfCCD,
	}
	maxEnergyValue := v2.Energy{
		Value: params.MaxEnergyAmount,
	}

	keyPair, err := v2.NewKeyPairFromSignKeyAndVerifyKey(privateKey[:32], privateKey[32:])
	if err != nil {
		log.Fatalf("failed to create key pair, err: %v", err)
	}

	senderWallet := v2.NewWalletAccount(sender, *keyPair)

	parameter, err := hex.DecodeString(params.Parameter)
	if err != nil {
		log.Fatalf("failed to update contract parameter value, err: %v", err)
	}

	updateContractPayload := v2.UpdateContractPayload{
		Amount: &amount,
		Address: &v2.ContractAddress{
			Index:    params.ContractIndex,
			Subindex: params.ContractSubIndex,
		},
		ReceiveName: &v2.ReceiveName{Value: params.ReceiveName},
		Parameter:   &v2.Parameter{Value: parameter},
	}
	accountTx, err := send.UpdateContract(
		senderWallet,
		*senderWallet.Address,
		nonce,
		expiry,
		updateContractPayload,
		maxEnergyValue,
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
		for _, item := range items {
			if item.Hash.Value == txHash.Value {
				return
			}
		}

		log.Fatalf("tx hash not match expected")
	}
}
