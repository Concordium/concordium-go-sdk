package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetBlockTransactionEvents returns stream of transaction events in a given block.
// The stream will end when all the transaction events for a given block have been returned.
func (c *Client) GetBlockTransactionEvents(ctx context.Context, req isBlockHashInput) (_ pb.Queries_GetBlockTransactionEventsClient, err error) {
	stream, err := c.GrpcClient.GetBlockTransactionEvents(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
