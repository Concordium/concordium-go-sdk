package v2

import (
	"context"
)

// GetAnonymityRevokers get the anonymity revokers registered as of the end of a given block.
// The stream will end when all the anonymity revokers have been returned.
func (c *Client) GetAnonymityRevokers(ctx context.Context, req *BlockHashInput) (_ Queries_GetAnonymityRevokersClient, err error) {
	stream, err := c.grpcClient.GetAnonymityRevokers(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
