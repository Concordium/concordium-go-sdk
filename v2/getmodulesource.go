package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetModuleSource get the source of a smart contract module.
func (c *Client) GetModuleSource(ctx context.Context, req *pb.ModuleSourceRequest) (_ *pb.VersionedModuleSource, err error) {
	source, err := c.grpcClient.GetModuleSource(ctx, req)
	if err != nil {
		return &pb.VersionedModuleSource{}, err
	}

	return source, nil
}
