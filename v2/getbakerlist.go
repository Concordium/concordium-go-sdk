package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetBakerList get all the bakers at the end of the given block.
func (c *Client) GetBakerList(ctx context.Context, req *pb.BlockHashInput) (_ pb.Queries_GetBakerListClient, err error) {
	stream, err := c.grpcClient.GetBakerList(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
