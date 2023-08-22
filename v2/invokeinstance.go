package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// InvokeInstance run the smart contract entrypoint in a given context and in the state at the end of the given block.
func (c *Client) InvokeInstance(ctx context.Context, payload UpdateContractPayload, input isBlockHashInput, energy Energy, address isAddress) (_ *pb.InvokeInstanceResponse, err error) {
	var pbAddress pb.Address

	switch k := address.(type) {
	case *AccountAddress:
		accountAddress := make([]byte, AccountAddressLength)
		copy(accountAddress, k.Value[:])
		pbAddress.Type = &pb.Address_Account{
			Account: &pb.AccountAddress{
				Value: accountAddress,
			},
		}
	case *ContractAddress:
		pbAddress.Type = &pb.Address_Contract{
			Contract: &pb.ContractAddress{
				Index:    k.Index,
				Subindex: k.Subindex,
			},
		}
	}

	invokeInstanceResponse, err := c.GrpcClient.InvokeInstance(ctx, &pb.InvokeInstanceRequest{
		BlockHash: convertBlockHashInput(input),
		Invoker:   &pbAddress,
		Instance: &pb.ContractAddress{
			Index:    payload.Address.Index,
			Subindex: payload.Address.Subindex,
		},
		Amount: &pb.Amount{
			Value: payload.Amount.Value,
		},
		Entrypoint: &pb.ReceiveName{
			Value: payload.ReceiveName.Value,
		},
		Parameter: &pb.Parameter{
			Value: payload.Parameter.Value,
		},
		Energy: &pb.Energy{
			Value: energy.Value,
		},
	})
	if err != nil {
		return &pb.InvokeInstanceResponse{}, err
	}

	return invokeInstanceResponse, nil
}
