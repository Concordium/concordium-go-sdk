package v2

import (
	"context"
)

// GetBlockChainParameters get the values of chain parameters in effect in the given block.
func (c *Client) GetBlockChainParameters(ctx context.Context, req *BlockHashInput) (_ *ChainParameters, err error) {
	chainParameters, err := c.grpcClient.GetBlockChainParameters(ctx, req)
	if err != nil {
		return &ChainParameters{}, Error.Wrap(err)
	}

	return chainParameters, nil
}
