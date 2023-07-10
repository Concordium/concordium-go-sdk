package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetConsensusInfo get information about the current state of consensus.
func (c *Client) GetConsensusInfo(ctx context.Context) (_ *pb.ConsensusInfo, err error) {
	consensusInfo, err := c.GrpcClient.GetConsensusInfo(ctx, new(pb.Empty))
	if err != nil {
		return &pb.ConsensusInfo{}, err
	}

	return consensusInfo, nil
}
