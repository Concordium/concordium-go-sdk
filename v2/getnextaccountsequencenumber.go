package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetNextAccountSequenceNumber get the best guess as to what the next account sequence number should be.
// If all account transactions are finalized then this information is reliable.
// Otherwise, this is the best guess, assuming all other transactions will be committed to blocks and eventually finalized.
func (c *Client) GetNextAccountSequenceNumber(ctx context.Context, req *pb.AccountAddress) (_ *pb.NextAccountSequenceNumber, err error) {
	sequenceNumber, err := c.grpcClient.GetNextAccountSequenceNumber(ctx, req)
	if err != nil {
		return &pb.NextAccountSequenceNumber{}, err
	}

	return sequenceNumber, nil
}
