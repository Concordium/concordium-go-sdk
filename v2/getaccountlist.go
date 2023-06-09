package v2

import (
	"context"
)

// GetAccountList retrieve the list of accounts that exist at the end of the given block.
func (c *Client) GetAccountList(ctx context.Context, req *BlockHashInput) (_ Queries_GetAccountListClient, err error) {
	stream, err := c.grpcClient.GetAccountList(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
