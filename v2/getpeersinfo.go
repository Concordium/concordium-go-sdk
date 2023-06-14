package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetPeersInfo get a list of the peers that the node is connected to / and assoicated network related information for each peer.
func (c *Client) GetPeersInfo(ctx context.Context) (_ *pb.PeersInfo, err error) {
	peersInfo, err := c.grpcClient.GetPeersInfo(ctx, new(pb.Empty))
	if err != nil {
		return &pb.PeersInfo{}, Error.Wrap(err)
	}

	return peersInfo, nil
}
