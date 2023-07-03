package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetTokenomicsInfo get information about tokenomics at the end of a given block.
func (c *Client) GetTokenomicsInfo(ctx context.Context, req *pb.BlockHashInput) (_ *pb.TokenomicsInfo, err error) {
	tokenomicsInfo, err := c.grpcClient.GetTokenomicsInfo(ctx, req)
	if err != nil {
		return &pb.TokenomicsInfo{}, err
	}

	return tokenomicsInfo, nil
}
