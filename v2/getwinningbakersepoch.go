package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetWinningBakersEpoch retrieves a list of bakers that won the lottery in a particular historical epoch
// (i.e. the last finalized block in a later epoch). This lists the winners for each round in the epoch,
// starting from the round after the last block in the previous epoch, running to the round before the first
// block in the next epoch. It also indicates if a block in each round was included in the finalized chain.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetWinningBakersEpoch(ctx context.Context, req isEpochRequest) (_ pb.Queries_GetWinningBakersEpochClient, err error) {
	stream, err := c.GrpcClient.GetWinningBakersEpoch(ctx, convertEpochRequest(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
