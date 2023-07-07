package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetAnonymityRevokers get the anonymity revokers registered as of the end of a given block.
// The stream will end when all the anonymity revokers have been returned.
func (c *Client) GetAnonymityRevokers(ctx context.Context, req isBlockHashInput) (_ []*pb.ArInfo, err error) {
	stream, err := c.GrpcClient.GetAnonymityRevokers(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var arInfos []*pb.ArInfo

	for err == nil {
		arInfo, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		arInfos = append(arInfos, arInfo)
	}

	return arInfos, nil
}
