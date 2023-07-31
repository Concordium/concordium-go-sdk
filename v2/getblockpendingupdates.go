package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockPendingUpdates get the pending updates to chain parameters at the end of a given block.
// The stream will end when all the pending updates for a given block have been returned.
func (c *Client) GetBlockPendingUpdates(ctx context.Context, req isBlockHashInput) (_ pb.Queries_GetBlockPendingUpdatesClient, err error) {
	stream, err := c.GrpcClient.GetBlockPendingUpdates(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
