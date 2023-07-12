package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBannedPeers get a list of banned peers.
func (c *Client) GetBannedPeers(ctx context.Context) (_ *pb.BannedPeers, err error) {
	bannedPeers, err := c.GrpcClient.GetBannedPeers(ctx, new(pb.Empty))
	if err != nil {
		return &pb.BannedPeers{}, err
	}

	return bannedPeers, nil
}
