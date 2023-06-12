package v2

import (
	"context"
)

// DumpStop stop dumping packages. Only enabled if the node was built with the `network_dump` feature.
// Returns a GRPC error if the network dump failed to be stopped.
func (c *Client) DumpStop(ctx context.Context) (err error) {
	_, err = c.grpcClient.DumpStop(ctx, new(Empty))
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}
