package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// InstanceStateLookup get the value at a specific key of a contract state.
// In contrast to `GetInstanceState` this is more efficient,
// but requires the user to know the specific key to look for.
func (c *Client) InstanceStateLookup(ctx context.Context, req *pb.InstanceStateLookupRequest) (_ *pb.InstanceStateValueAtKey, err error) {
	instanceStateKeyValue, err := c.GrpcClient.InstanceStateLookup(ctx, req)
	if err != nil {
		return nil, err
	}

	return instanceStateKeyValue, nil
}
