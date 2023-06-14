package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetBannedPeers get a list of banned peers.
func (c *Client) GetBannedPeers(ctx context.Context) (_ *pb.BannedPeers, err error) {
	bannedPeers, err := c.grpcClient.GetBannedPeers(ctx, new(pb.Empty))
	if err != nil {
		return &pb.BannedPeers{}, Error.Wrap(err)
	}

	return bannedPeers, nil
}
