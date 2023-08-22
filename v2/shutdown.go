package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// Shutdown shut down the node. Return a GRPC error if the shutdown failed.
func (c *Client) Shutdown(ctx context.Context) (err error) {
	_, err = c.GrpcClient.Shutdown(ctx, new(pb.Empty))
	if err != nil {
		return err
	}

	return nil
}
