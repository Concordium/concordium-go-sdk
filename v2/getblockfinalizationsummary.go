package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockFinalizationSummary get the summary of the finalization data in a given block.
func (c *Client) GetBlockFinalizationSummary(ctx context.Context, req *pb.BlockHashInput) (_ *pb.BlockFinalizationSummary, err error) {
	blockFinalizationSummary, err := c.grpcClient.GetBlockFinalizationSummary(ctx, req)
	if err != nil {
		return &pb.BlockFinalizationSummary{}, err
	}

	return blockFinalizationSummary, nil
}
