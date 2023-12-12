package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetFirstBlockEpoch retrieves the block hash of the first finalized block in a specified epoch.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetFirstBlockEpoch(ctx context.Context, req isEpochRequest) (_ *pb.BlockHash, err error) {
	resp, err := c.GrpcClient.GetFirstBlockEpoch(ctx, convertEpochRequest(req))
	if err != nil {
		return nil, err
	}

	return resp, nil
}
