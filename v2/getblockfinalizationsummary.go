package v2

import (
	"context"
)

// GetBlockFinalizationSummary get the summary of the finalization data in a given block.
func (c *Client) GetBlockFinalizationSummary(ctx context.Context, req *BlockHashInput) (_ *BlockFinalizationSummary, err error) {
	blockFinalizationSummary, err := c.grpcClient.GetBlockFinalizationSummary(ctx, req)
	if err != nil {
		return &BlockFinalizationSummary{}, Error.Wrap(err)
	}

	return blockFinalizationSummary, nil
}
