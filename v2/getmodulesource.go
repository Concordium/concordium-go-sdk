package v2

import (
	"context"
)

// GetModuleSource get the source of a smart contract module.
func (c *Client) GetModuleSource(ctx context.Context, req *ModuleSourceRequest) (_ *VersionedModuleSource, err error) {
	source, err := c.grpcClient.GetModuleSource(ctx, req)
	if err != nil {
		return &VersionedModuleSource{}, Error.Wrap(err)
	}

	return source, nil
}
