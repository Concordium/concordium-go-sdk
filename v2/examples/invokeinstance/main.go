package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

type invokeInstanceParams struct {
	SenderPubKey string `env:"SENDER_PUBLIC_KEY_IN_BASE58"`

	ContractIndex    uint64 `env:"CONTRACT_INDEX"`
	ContractSubIndex uint64 `env:"CONTRACT_SUB_INDEX"`
	ReceiveName      string `env:"RECEIVE_NAME"`
	Parameter        string `env:"PARAMETER_IN_HEX"`
	AmountOfCCD      uint64 `env:"AMOUNT_OF_CCD"`
	MaxEnergyAmount  uint64 `env:"MAX_ENERGY_AMOUNT"`
}

// for this example we need to fill .env file (located in the same directory with this main.go file) with values,
// described in struct invokeInstanceParams above.
// invokes instance + returns success or failed result.
// to run this example write in terminal in directory of main.go file "go run main.go".
func main() {
	if err := godotenv.Overload(); err != nil {
		log.Fatalf("failed to load update params, err: %v", err)
	}

	params := new(invokeInstanceParams)
	if err := env.Parse(params); err != nil {
		log.Fatalf("failed to parse update params, err: %v", err)
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
	amount := v2.Amount{
		Value: params.AmountOfCCD,
	}
	maxEnergyValue := v2.Energy{
		Value: params.MaxEnergyAmount,
	}

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

	resp, err := client.InvokeInstance(ctx, updateContractPayload, v2.BlockHashInputLastFinal{}, maxEnergyValue, &sender)
	if err != nil {
		log.Fatalf("failed to invoke instance, err: %v", err)
	}

	success := resp.GetSuccess()
	if success == nil {
		failure := resp.GetFailure()
		if failure == nil {
			fmt.Println("couldn't estimate transaction")
			return
		}

		fmt.Println("Failure! Reason: ", failure.GetReason().String())
		return
	}

	fmt.Println("Success!")
	fmt.Println("Used energy: ", success.UsedEnergy.Value)
	fmt.Println("Returned value: ", success.ReturnValue)
}
