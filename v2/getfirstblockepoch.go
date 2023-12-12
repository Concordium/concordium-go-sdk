package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetFirstBlockEpoch retrieves the block hash of the first finalized block in a specified epoch.
//
// The following errors are possible:
//   - `NOT_FOUND` if the query specifies an unknown block.
//   - `UNAVAILABLE` if the query is for an epoch that is not finalized in the current genesis index, or is for a future genesis index.
//   - `INVALID_ARGUMENT` if the query is for an epoch with no finalized blocks for a past genesis index.
//   - `INVALID_ARGUMENT` if the input `EpochRequest` is malformed.
//   - `UNIMPLEMENTED` if the endpoint is disabled on the node.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetFirstBlockEpoch(ctx context.Context, req isEpochRequest) (_ *pb.BlockHash, err error) {
	resp, err := c.GrpcClient.GetFirstBlockEpoch(ctx, convertEpochRequest(req))
	if err != nil {
		return nil, err
	}

	return resp, nil
}
