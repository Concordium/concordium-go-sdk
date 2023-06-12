package v2

import (
	"context"
)

// GetPeersInfo get a list of the peers that the node is connected to / and assoicated network related information for each peer.
func (c *Client) GetPeersInfo(ctx context.Context) (_ *PeersInfo, err error) {
	peersInfo, err := c.grpcClient.GetPeersInfo(ctx, new(Empty))
	if err != nil {
		return &PeersInfo{}, Error.Wrap(err)
	}

	return peersInfo, nil
}
