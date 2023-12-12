package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetWinningBakersEpoch retrieves the list of bakers that won the lottery in a particular historical epoch
// (i.e. the last finalized block is in a later epoch). This lists the winners for each round in the
// epoch, starting from the round after the last block in the previous epoch, running to
// the round before the first block in the next epoch. It also indicates if a block in each
// round was included in the finalized chain.
//
// The following error cases are possible:
//   - `NOT_FOUND` if the query specifies an unknown block.
//   - `UNAVAILABLE` if the query is for an epoch that is not finalized in the current genesis index, or is for a future genesis index.
//   - `INVALID_ARGUMENT` if the query is for an epoch that is not finalized for a past genesis index.
//   - `INVALID_ARGUMENT` if the query is for a genesis index at consensus version 0.
//   - `INVALID_ARGUMENT` if the input `EpochRequest` is malformed.
//   - `UNIMPLEMENTED` if the endpoint is disabled on the node.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetWinningBakersEpoch(ctx context.Context, req isEpochRequest) (_ pb.Queries_GetWinningBakersEpochClient, err error) {
	stream, err := c.GrpcClient.GetWinningBakersEpoch(ctx, convertEpochRequest(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
