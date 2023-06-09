package v2

import (
	"context"
)

// InstanceStateLookup get the value at a specific key of a contract state.
// In contrast to `GetInstanceState` this is more efficient,
// but requires the user to know the specific key to look for.
func (c *Client) InstanceStateLookup(ctx context.Context, req *InstanceStateLookupRequest) (_ *InstanceStateValueAtKey, err error) {
	instanceStateKeyValue, err := c.grpcClient.InstanceStateLookup(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return instanceStateKeyValue, nil
}
