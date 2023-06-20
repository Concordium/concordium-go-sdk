package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetCryptographicParameters get the cryptographic parameters in a given block.
func (c *Client) GetCryptographicParameters(ctx context.Context, req *pb.BlockHashInput) (_ *pb.CryptographicParameters, err error) {
	cryptographicParameters, err := c.grpcClient.GetCryptographicParameters(ctx, req)
	if err != nil {
		return &pb.CryptographicParameters{}, err
	}

	return cryptographicParameters, nil
}
