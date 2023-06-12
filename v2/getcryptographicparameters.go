package v2

import (
	"context"
)

// GetCryptographicParameters get the cryptographic parameters in a given block.
func (c *Client) GetCryptographicParameters(ctx context.Context, req *BlockHashInput) (_ *CryptographicParameters, err error) {
	cryptographicParameters, err := c.grpcClient.GetCryptographicParameters(ctx, req)
	if err != nil {
		return &CryptographicParameters{}, Error.Wrap(err)
	}

	return cryptographicParameters, nil
}
