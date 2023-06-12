package v2

import (
	"context"
)

// DumpStart start dumping packages into the specified file.
// Only enabled if the node was built with the `network_dump` feature.
// Returns a GRPC error if the network dump failed to start.
func (c *Client) DumpStart(ctx context.Context, req *DumpRequest) (err error) {
	_, err = c.grpcClient.DumpStart(ctx, req)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}
