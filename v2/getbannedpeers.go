package v2

import (
	"context"
)

// GetBannedPeers get a list of banned peers.
func (c *Client) GetBannedPeers(ctx context.Context) (_ *BannedPeers, err error) {
	bannedPeers, err := c.grpcClient.GetBannedPeers(ctx, new(Empty))
	if err != nil {
		return &BannedPeers{}, Error.Wrap(err)
	}

	return bannedPeers, nil
}
