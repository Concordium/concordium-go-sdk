package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// UnbanPeer unban the banned peer. Returns a GRPC error if the action failed.
func (c *Client) UnbanPeer(ctx context.Context, req *pb.BannedPeer) (err error) {
	_, err = c.GrpcClient.UnbanPeer(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
