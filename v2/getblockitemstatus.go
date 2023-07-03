package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockItemStatus get the status of and information about a specific block item (transaction).
func (c *Client) GetBlockItemStatus(ctx context.Context, req *pb.TransactionHash) (_ *pb.BlockItemStatus, err error) {
	blockItemStatus, err := c.grpcClient.GetBlockItemStatus(ctx, req)
	if err != nil {
		return &pb.BlockItemStatus{}, err
	}

	return blockItemStatus, nil
}
