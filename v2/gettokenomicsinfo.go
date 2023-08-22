package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetTokenomicsInfo get information about tokenomics at the end of a given block.
func (c *Client) GetTokenomicsInfo(ctx context.Context, req isBlockHashInput) (_ *pb.TokenomicsInfo, err error) {
	tokenomicsInfo, err := c.GrpcClient.GetTokenomicsInfo(ctx, convertBlockHashInput(req))
	if err != nil {
		return &pb.TokenomicsInfo{}, err
	}

	return tokenomicsInfo, nil
}
