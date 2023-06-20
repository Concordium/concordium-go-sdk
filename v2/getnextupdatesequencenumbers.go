package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetNextUpdateSequenceNumbers get next available sequence numbers for updating chain parameters after a given block.
func (c *Client) GetNextUpdateSequenceNumbers(ctx context.Context, req *pb.BlockHashInput) (_ *pb.NextUpdateSequenceNumbers, err error) {
	nextUpdateSequenceNumbers, err := c.grpcClient.GetNextUpdateSequenceNumbers(ctx, req)
	if err != nil {
		return &pb.NextUpdateSequenceNumbers{}, err
	}

	return nextUpdateSequenceNumbers, nil
}
