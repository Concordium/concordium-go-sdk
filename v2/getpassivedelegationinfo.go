package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetPassiveDelegationInfo get information about the passive delegators at the end of a given block.
func (c *Client) GetPassiveDelegationInfo(ctx context.Context, req *pb.BlockHashInput) (_ *pb.PassiveDelegationInfo, err error) {
	passiveDelegationInfo, err := c.grpcClient.GetPassiveDelegationInfo(ctx, req)
	if err != nil {
		return &pb.PassiveDelegationInfo{}, err
	}

	return passiveDelegationInfo, nil
}
