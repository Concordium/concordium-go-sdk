package v2

import (
	"context"
)

// GetBranches get the current branches of blocks starting from and including the last finalized block.
func (c *Client) GetBranches(ctx context.Context) (_ *Branch, err error) {
	branch, err := c.grpcClient.GetBranches(ctx, new(Empty))
	if err != nil {
		return &Branch{}, Error.Wrap(err)
	}

	return branch, nil
}
