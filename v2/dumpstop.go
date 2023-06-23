package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// DumpStop stop dumping packages. Only enabled if the node was built with the `network_dump` feature.
// Returns a GRPC error if the network dump failed to be stopped.
func (c *Client) DumpStop(ctx context.Context) (err error) {
	_, err = c.grpcClient.DumpStop(ctx, new(pb.Empty))
	if err != nil {
		return err
	}

	return nil
}
