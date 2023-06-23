package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetAnonymityRevokers get the anonymity revokers registered as of the end of a given block.
// The stream will end when all the anonymity revokers have been returned.
func (c *Client) GetAnonymityRevokers(ctx context.Context, req *pb.BlockHashInput) (_ pb.Queries_GetAnonymityRevokersClient, err error) {
	stream, err := c.grpcClient.GetAnonymityRevokers(ctx, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
