package v2

import (
	"context"
)

// GetInstanceInfo get info about a smart contract instance as it appears at the end of the given block.
func (c *Client) GetInstanceInfo(ctx context.Context, req *InstanceInfoRequest) (_ *InstanceInfo, err error) {
	instanceInfo, err := c.grpcClient.GetInstanceInfo(ctx, req)
	if err != nil {
		return &InstanceInfo{}, Error.Wrap(err)
	}

	return instanceInfo, nil
}
