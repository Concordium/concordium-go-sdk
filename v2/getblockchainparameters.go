package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlockChainParameters get the values of chain parameters in effect in the given block.
func (c *Client) GetBlockChainParameters(ctx context.Context, req *pb.BlockHashInput) (_ *pb.ChainParameters, err error) {
	chainParameters, err := c.grpcClient.GetBlockChainParameters(ctx, req)
	if err != nil {
		return &pb.ChainParameters{}, err
	}

	return chainParameters, nil
}
