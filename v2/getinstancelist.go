package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetInstanceList get a list of addresses for all smart contract instances.
// The stream will end when all instances that exist in the state
// at the end of the given block has been returned.
func (c *Client) GetInstanceList(ctx context.Context, req *BlockHashInput) (_ pb.Queries_GetInstanceListClient, err error) {
	stream, err := c.grpcClient.GetInstanceList(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
