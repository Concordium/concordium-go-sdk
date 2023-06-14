package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetInstanceState get the exact state of a specific contract instance, streamed as a list of key-value pairs.
// The list is streamed in lexicographic order of keys.
func (c *Client) GetInstanceState(ctx context.Context, req *pb.InstanceInfoRequest) (_ pb.Queries_GetInstanceStateClient, err error) {
	stream, err := c.grpcClient.GetInstanceState(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
