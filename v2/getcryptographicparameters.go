package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetCryptographicParameters get the cryptographic parameters in a given block.
func (c *Client) GetCryptographicParameters(ctx context.Context, req isBlockHashInput) (_ *pb.CryptographicParameters, err error) {
	cryptographicParameters, err := c.GrpcClient.GetCryptographicParameters(ctx, convertBlockHashInput(req))
	if err != nil {
		return &pb.CryptographicParameters{}, err
	}

	return cryptographicParameters, nil
}
