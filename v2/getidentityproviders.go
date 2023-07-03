package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetIdentityProviders get the identity providers registered as of the end of a given block.
// The stream will end when all the identity providers have been returned.
func (c *Client) GetIdentityProviders(ctx context.Context, req *pb.BlockHashInput) (_ pb.Queries_GetIdentityProvidersClient, err error) {
	stream, err := c.grpcClient.GetIdentityProviders(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
