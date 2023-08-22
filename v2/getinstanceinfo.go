package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetInstanceInfo get info about a smart contract instance as it appears at the end of the given block.
func (c *Client) GetInstanceInfo(ctx context.Context, blockHash isBlockHashInput, address ContractAddress) (_ *pb.InstanceInfo, err error) {
	instanceInfo, err := c.GrpcClient.GetInstanceInfo(ctx, &pb.InstanceInfoRequest{
		BlockHash: convertBlockHashInput(blockHash),
		Address: &pb.ContractAddress{
			Index:    address.Index,
			Subindex: address.Subindex,
		},
	})
	if err != nil {
		return &pb.InstanceInfo{}, err
	}

	return instanceInfo, nil
}
