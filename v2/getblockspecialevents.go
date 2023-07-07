package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockSpecialEvents get a list of transaction events in a given block.
// The stream will end when all the transaction events for a given block have been returned.
func (c *Client) GetBlockSpecialEvents(ctx context.Context, req isBlockHashInput) (_ []*pb.BlockSpecialEvent, err error) {
	stream, err := c.GrpcClient.GetBlockSpecialEvents(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var blockSpecialEvents []*pb.BlockSpecialEvent

	for err == nil {
		blockSpecialEvent, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		blockSpecialEvents = append(blockSpecialEvents, blockSpecialEvent)
	}

	return blockSpecialEvents, nil
}
