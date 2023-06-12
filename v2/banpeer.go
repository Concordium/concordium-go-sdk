package v2

import (
	"context"
)

// BanPeer ban the given peer. Returns a GRPC error if the action failed.
func (c *Client) BanPeer(ctx context.Context, req *PeerToBan) (err error) {
	_, err = c.grpcClient.BanPeer(ctx, req)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}
