package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// BanPeer ban the given peer. Returns a GRPC error if the action failed.
func (c *Client) BanPeer(ctx context.Context, req *pb.PeerToBan) (err error) {
	_, err = c.grpcClient.BanPeer(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
