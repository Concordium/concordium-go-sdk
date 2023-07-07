package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetInstanceState get the exact state of a specific contract instance, streamed as a list of key-value pairs.
// The list is streamed in lexicographic order of keys.
func (c *Client) GetInstanceState(ctx context.Context, blockHash isBlockHashInput, address ContractAddress) (_ []*pb.InstanceStateKVPair, err error) {
	stream, err := c.GrpcClient.GetInstanceState(ctx, &pb.InstanceInfoRequest{
		BlockHash: convertBlockHashInput(blockHash),
		Address: &pb.ContractAddress{
			Index:    address.Index,
			Subindex: address.Subindex,
		},
	})
	if err != nil {
		return nil, err
	}

	var instanceStateKVPairs []*pb.InstanceStateKVPair

	for err == nil {
		instanceStateKVPair, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		instanceStateKVPairs = append(instanceStateKVPairs, instanceStateKVPair)
	}

	return instanceStateKVPairs, nil
}
