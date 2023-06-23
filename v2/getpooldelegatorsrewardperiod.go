package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetPoolDelegatorsRewardPeriod get the fixed delegators of a given pool for the reward period of the given block.
// In contracts to the `GetPoolDelegators` which returns delegators registered for the given block,
// this endpoint returns the fixed delegators contributing stake in the reward period containing the given block.
// The stream will end when all the delegators has been returned.
func (c *Client) GetPoolDelegatorsRewardPeriod(ctx context.Context, req *pb.GetPoolDelegatorsRequest) (_ pb.Queries_GetPoolDelegatorsRewardPeriodClient, err error) {
	stream, err := c.grpcClient.GetPoolDelegatorsRewardPeriod(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
