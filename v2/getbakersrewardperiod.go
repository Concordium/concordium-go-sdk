package v2

import (
	"context"
)

// GetBakersRewardPeriod retrieves all bakers in the reward period of a block.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetBakersRewardPeriod(ctx context.Context, req isBlockHashInput) (_ BakerRewardPeriodInfoStream, err error) {
	stream, err := c.GrpcClient.GetBakersRewardPeriod(ctx, convertBlockHashInput(req))
	if err != nil {
		return BakerRewardPeriodInfoStream{}, err
	}

	return BakerRewardPeriodInfoStream{stream: stream}, nil
}
