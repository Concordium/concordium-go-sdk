package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetNextUpdateSequenceNumbers get next available sequence numbers for updating chain parameters after a given block.
func (c *Client) GetNextUpdateSequenceNumbers(ctx context.Context, req isBlockHashInput) (_ *pb.NextUpdateSequenceNumbers, err error) {
	nextUpdateSequenceNumbers, err := c.GrpcClient.GetNextUpdateSequenceNumbers(ctx, convertBlockHashInput(req))
	if err != nil {
		return &pb.NextUpdateSequenceNumbers{}, err
	}

	return nextUpdateSequenceNumbers, nil
}
