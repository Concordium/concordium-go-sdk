package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetPassiveDelegatorsRewardPeriod get the fixed passive delegators for the reward period of the given block.
// In contracts to the `GetPassiveDelegators` which returns delegators registered for the given block,
// this endpoint returns the fixed delegators contributing stake in the reward period containing the given block.
// The stream will end when all the delegators has been returned.
func (c *Client) GetPassiveDelegatorsRewardPeriod(ctx context.Context, req *pb.BlockHashInput) (_ pb.Queries_GetPassiveDelegatorsRewardPeriodClient, err error) {
	stream, err := c.grpcClient.GetPassiveDelegatorsRewardPeriod(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
