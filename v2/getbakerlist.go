package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetBakerList get all the bakers at the end of the given block.
func (c *Client) GetBakerList(ctx context.Context, req isBlockHashInput) (_ pb.Queries_GetBakerListClient, err error) {
	stream, err := c.GrpcClient.GetBakerList(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
