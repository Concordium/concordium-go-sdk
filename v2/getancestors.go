package v2

import (
	"context"
)

// GetAncestors get a stream of ancestors for the provided block.
// Starting with the provided block itself, moving backwards until
// no more ancestors or the requested number of ancestors has been returned.
func (c *Client) GetAncestors(ctx context.Context, req *AncestorsRequest) (_ Queries_GetAncestorsClient, err error) {
	stream, err := c.grpcClient.GetAncestors(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
