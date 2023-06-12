package v2

import (
	"context"
)

// GetPassiveDelegationInfo get information about the passive delegators at the end of a given block.
func (c *Client) GetPassiveDelegationInfo(ctx context.Context, req *BlockHashInput) (_ *PassiveDelegationInfo, err error) {
	passiveDelegationInfo, err := c.grpcClient.GetPassiveDelegationInfo(ctx, req)
	if err != nil {
		return &PassiveDelegationInfo{}, Error.Wrap(err)
	}

	return passiveDelegationInfo, nil
}
