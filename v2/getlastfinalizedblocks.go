package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetFinalizedBlocks return a stream of blocks that are finalized from the time the query is made onward.
// This can be used to listen for newly finalized blocks. Note that there is no guarantee that blocks
// will not be skipped if the client is too slow in processing the stream, however blocks will always be sent by increasing block height.
func (c *Client) GetFinalizedBlocks(ctx context.Context) (_ pb.Queries_GetFinalizedBlocksClient, err error) {
	stream, err := c.grpcClient.GetFinalizedBlocks(ctx, new(pb.Empty))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
