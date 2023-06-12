package v2

import (
	"context"
)

// GetTokenomicsInfo get information about tokenomics at the end of a given block.
func (c *Client) GetTokenomicsInfo(ctx context.Context, req *BlockHashInput) (_ *TokenomicsInfo, err error) {
	tokenomicsInfo, err := c.grpcClient.GetTokenomicsInfo(ctx, req)
	if err != nil {
		return &TokenomicsInfo{}, Error.Wrap(err)
	}

	return tokenomicsInfo, nil
}
