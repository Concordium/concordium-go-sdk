package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetIdentityProviders get the identity providers registered as of the end of a given block.
// The stream will end when all the identity providers have been returned.
func (c *Client) GetIdentityProviders(ctx context.Context, req isBlockHashInput) (_ pb.Queries_GetIdentityProvidersClient, err error) {
	stream, err := c.GrpcClient.GetIdentityProviders(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
