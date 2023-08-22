package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// DumpStart start dumping packages into the specified file.
// Only enabled if the node was built with the `network_dump` feature.
// Returns a GRPC error if the network dump failed to start.
func (c *Client) DumpStart(ctx context.Context, req *pb.DumpRequest) (err error) {
	_, err = c.GrpcClient.DumpStart(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
