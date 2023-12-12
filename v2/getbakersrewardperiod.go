package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetBakersRewardPeriod retrieves all bakers in the reward period of a block.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetBakersRewardPeriod(ctx context.Context, req isBlockHashInput) (_ pb.Queries_GetBakersRewardPeriodClient, err error) {
	stream, err := c.GrpcClient.GetBakersRewardPeriod(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
