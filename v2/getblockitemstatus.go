package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockItemStatus get the status of and information about a specific block item (transaction).
func (c *Client) GetBlockItemStatus(ctx context.Context, req TransactionHash) (_ *pb.BlockItemStatus, err error) {
	blockItemStatus, err := c.GrpcClient.GetBlockItemStatus(ctx, &pb.TransactionHash{
		Value: req.Value[:],
	})
	if err != nil {
		return &pb.BlockItemStatus{}, err
	}

	return blockItemStatus, nil
}
