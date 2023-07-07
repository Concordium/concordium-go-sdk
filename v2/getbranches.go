package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBranches get the current branches of blocks starting from and including the last finalized block.
func (c *Client) GetBranches(ctx context.Context) (_ *pb.Branch, err error) {
	branch, err := c.GrpcClient.GetBranches(ctx, new(pb.Empty))
	if err != nil {
		return &pb.Branch{}, err
	}

	return branch, nil
}
