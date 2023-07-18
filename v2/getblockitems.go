package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockItems get the items of a block.
func (c *Client) GetBlockItems(ctx context.Context, req isBlockHashInput) (_ []*BlockItem, err error) {
	stream, err := c.GrpcClient.GetBlockItems(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var blockItems []*pb.BlockItem

	for err == nil {
		blockItem, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		blockItems = append(blockItems, blockItem)
	}

	return ConvertBlockItems(blockItems), nil
}
