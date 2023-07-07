package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetInstanceList get a list of addresses for all smart contract instances.
// The stream will end when all instances that exist in the state
// at the end of the given block has been returned.
func (c *Client) GetInstanceList(ctx context.Context, req isBlockHashInput) (_ []*ContractAddress, err error) {
	stream, err := c.GrpcClient.GetInstanceList(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var contractAddresses []*pb.ContractAddress

	for err == nil {
		contractAddress, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		contractAddresses = append(contractAddresses, contractAddress)
	}

	var result []*ContractAddress

	for i := 0; i < len(contractAddresses); i++ {
		result = append(result, &ContractAddress{
			Index:    contractAddresses[i].Index,
			Subindex: contractAddresses[i].Subindex,
		})
	}

	return result, nil
}
