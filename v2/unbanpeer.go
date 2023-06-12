package v2

import (
	"context"
)

// UnbanPeer unban the banned peer. Returns a GRPC error if the action failed.
func (c *Client) UnbanPeer(ctx context.Context, req *BannedPeer) (err error) {
	_, err = c.grpcClient.UnbanPeer(ctx, req)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}
