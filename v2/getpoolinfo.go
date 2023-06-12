package v2

import (
	"context"
)

// GetPoolInfo get information about a given pool at the end of a given block.
func (c *Client) GetPoolInfo(ctx context.Context, req *PoolInfoRequest) (_ *PoolInfoResponse, err error) {
	poolInfo, err := c.grpcClient.GetPoolInfo(ctx, req)
	if err != nil {
		return &PoolInfoResponse{}, Error.Wrap(err)
	}

	return poolInfo, nil
}
