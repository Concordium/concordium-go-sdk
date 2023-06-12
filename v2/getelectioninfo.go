package v2

import (
	"context"
)

// GetElectionInfo get information related to the baker election for a particular block.
func (c *Client) GetElectionInfo(ctx context.Context, req *BlockHashInput) (_ *ElectionInfo, err error) {
	electionInfo, err := c.grpcClient.GetElectionInfo(ctx, req)
	if err != nil {
		return &ElectionInfo{}, Error.Wrap(err)
	}

	return electionInfo, nil
}
