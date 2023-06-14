package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetPoolInfo get information about a given pool at the end of a given block.
func (c *Client) GetPoolInfo(ctx context.Context, req *pb.PoolInfoRequest) (_ *pb.PoolInfoResponse, err error) {
	poolInfo, err := c.grpcClient.GetPoolInfo(ctx, req)
	if err != nil {
		return &pb.PoolInfoResponse{}, Error.Wrap(err)
	}

	return poolInfo, nil
}
