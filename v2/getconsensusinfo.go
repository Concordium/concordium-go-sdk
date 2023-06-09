package v2

import (
	"context"
)

// GetConsensusInfo get information about the current state of consensus.
func (c *Client) GetConsensusInfo(ctx context.Context, req *AccountAddress) (_ *ConsensusInfo, err error) {
	consensusInfo, err := c.grpcClient.GetConsensusInfo(ctx, new(Empty))
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return consensusInfo, nil
}
