package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetAnonymityRevokers get the anonymity revokers registered as of the end of a given block.
// The stream will end when all the anonymity revokers have been returned.
func (c *Client) GetAnonymityRevokers(ctx context.Context, req isBlockHashInput) (_ pb.Queries_GetAnonymityRevokersClient, err error) {
	stream, err := c.GrpcClient.GetAnonymityRevokers(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
