package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetAncestors get a stream of ancestors for the provided block.
// Starting with the provided block itself, moving backwards until
// no more ancestors or the requested number of ancestors has been returned.
func (c *Client) GetAncestors(ctx context.Context, req *pb.AncestorsRequest) (_ pb.Queries_GetAncestorsClient, err error) {
	stream, err := c.grpcClient.GetAncestors(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
