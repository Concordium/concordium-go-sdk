package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockTransactionEvents get a list of transaction events in a given block.
// The stream will end when all the transaction events for a given block have been returned.
func (c *Client) GetBlockTransactionEvents(ctx context.Context, req isBlockHashInput) (_ []*pb.BlockItemSummary, err error) {
	stream, err := c.GrpcClient.GetBlockTransactionEvents(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var blockItemSummaries []*pb.BlockItemSummary

	for err == nil {
		blockItemSummary, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		blockItemSummaries = append(blockItemSummaries, blockItemSummary)
	}

	return blockItemSummaries, nil
}
