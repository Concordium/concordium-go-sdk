package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetBakerEarliestWinTime retrieves the projected earliest time at which a particular baker will be required to bake a block.
//
//   - If the baker is not a baker for the current reward period, this returns a timestamp at the start of the next reward period.
//   - If the baker is a baker for the current reward period, the earliest win time is projected from the current round forwward,
//     assuming that each round after the last finalixed round will take the minimum block time. (If blocks take longer, or timeouts occur,
//     the actual time may be later, and the reported time in subsequent queries may reflect this.)
//   - At the end of an epoch (or if the baker is not projected to bake before the end of the epoch)
//     the earliest win time for a (current) baker will be projected as the start of the next epoch.
//     This is because the seed for the leader election is updated at the epoch boundary,
//     and so the winners cannot be predicted beyond that.
//
// Note that in some circumstances the returned timestamp can be in the past, especially at the end of an epoch.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetBakerEarliestWinTime(ctx context.Context, req *pb.BakerId) (_ Timestamp, err error) {
	timestamp, err := c.GrpcClient.GetBakerEarliestWinTime(ctx, req)
	if err != nil {
		return Timestamp{}, err
	}

	return parseTimestamp(timestamp), nil
}
