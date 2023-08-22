package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetBlocks returns a stream of blocks that arrive from the time
// the query is made onward. This can be used to listen for incoming blocks.
func (c *Client) GetBlocks(ctx context.Context) (_ pb.Queries_GetBlocksClient, err error) {
	stream, err := c.GrpcClient.GetBlocks(ctx, new(pb.Empty))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
