package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// Shutdown shut down the node. Return a GRPC error if the shutdown failed.
func (c *Client) Shutdown(ctx context.Context) (err error) {
	_, err = c.grpcClient.Shutdown(ctx, new(pb.Empty))
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}
