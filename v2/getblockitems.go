package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetBlockItems get the items of a block.
func (c *Client) GetBlockItems(ctx context.Context, req *pb.BlockHashInput) (_ pb.Queries_GetBlockItemsClient, err error) {
	stream, err := c.grpcClient.GetBlockItems(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
