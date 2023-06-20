package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetPoolDelegators get the registered delegators of a given pool at the end of a given block.
// In contrast to the `GetPoolDelegatorsRewardPeriod` which returns delegators that are fixed
// for the reward period of the block, this endpoint returns the list of delegators
// that are registered in the block. Any changes to delegators are immediately visible in this list.
// The stream will end when all the delegators has been returned.
func (c *Client) GetPoolDelegators(ctx context.Context, req *pb.GetPoolDelegatorsRequest) (_ pb.Queries_GetPoolDelegatorsClient, err error) {
	stream, err := c.grpcClient.GetPoolDelegators(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
