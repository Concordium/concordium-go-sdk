package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetBranches get the current branches of blocks starting from and including the last finalized block.
func (c *Client) GetBranches(ctx context.Context) (_ *pb.Branch, err error) {
	branch, err := c.grpcClient.GetBranches(ctx, new(pb.Empty))
	if err != nil {
		return &pb.Branch{}, Error.Wrap(err)
	}

	return branch, nil
}
