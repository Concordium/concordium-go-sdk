package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetElectionInfo get information related to the baker election for a particular block.
func (c *Client) GetElectionInfo(ctx context.Context, req *pb.BlockHashInput) (_ *pb.ElectionInfo, err error) {
	electionInfo, err := c.grpcClient.GetElectionInfo(ctx, req)
	if err != nil {
		return &pb.ElectionInfo{}, err
	}

	return electionInfo, nil
}
