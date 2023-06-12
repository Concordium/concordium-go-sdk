package v2

import (
	"context"
)

// GetNextUpdateSequenceNumbers get next available sequence numbers for updating chain parameters after a given block.
func (c *Client) GetNextUpdateSequenceNumbers(ctx context.Context, req *BlockHashInput) (_ *NextUpdateSequenceNumbers, err error) {
	nextUpdateSequenceNumbers, err := c.grpcClient.GetNextUpdateSequenceNumbers(ctx, req)
	if err != nil {
		return &NextUpdateSequenceNumbers{}, Error.Wrap(err)
	}

	return nextUpdateSequenceNumbers, nil
}
