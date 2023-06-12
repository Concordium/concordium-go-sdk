package v2

import (
	"context"
)

// GetPassiveDelegators get the registered passive delegators at the end of a given block.
// In contrast to the `GetPassiveDelegatorsRewardPeriod` which returns delegators that are fixed for the reward period of the block,
// this endpoint returns the list of delegators that are registered in the block.
// Any changes to delegators are immediately visible in this list. The stream will end when all the delegators has been returned.
func (c *Client) GetPassiveDelegators(ctx context.Context, req *BlockHashInput) (_ Queries_GetPassiveDelegatorsClient, err error) {
	stream, err := c.grpcClient.GetPassiveDelegators(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
