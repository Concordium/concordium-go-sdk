package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetModuleList get a list of all smart contract modules. The stream will end when
// all modules that exist in the state at the end of the given block have been returned.
func (c *Client) GetModuleList(ctx context.Context, req *pb.BlockHashInput) (_ pb.Queries_GetModuleListClient, err error) {
	stream, err := c.grpcClient.GetModuleList(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
