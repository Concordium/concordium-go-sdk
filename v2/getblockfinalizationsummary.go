package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetBlockFinalizationSummary get the summary of the finalization data in a given block.
func (c *Client) GetBlockFinalizationSummary(ctx context.Context, req isBlockHashInput) (_ *pb.BlockFinalizationSummary, err error) {
	blockFinalizationSummary, err := c.GrpcClient.GetBlockFinalizationSummary(ctx, convertBlockHashInput(req))
	if err != nil {
		return &pb.BlockFinalizationSummary{}, err
	}

	return blockFinalizationSummary, nil
}
