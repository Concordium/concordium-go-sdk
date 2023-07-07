package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetElectionInfo get information related to the baker election for a particular block.
func (c *Client) GetElectionInfo(ctx context.Context, req isBlockHashInput) (_ *pb.ElectionInfo, err error) {
	electionInfo, err := c.GrpcClient.GetElectionInfo(ctx, convertBlockHashInput(req))
	if err != nil {
		return &pb.ElectionInfo{}, err
	}

	return electionInfo, nil
}
