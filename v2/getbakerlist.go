package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBakerList get all the bakers at the end of the given block.
func (c *Client) GetBakerList(ctx context.Context, req isBlockHashInput) (_ []*pb.BakerId, err error) {
	stream, err := c.GrpcClient.GetBakerList(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var backerIDs []*pb.BakerId

	for err == nil {
		backerID, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		backerIDs = append(backerIDs, backerID)
	}

	return backerIDs, nil
}
