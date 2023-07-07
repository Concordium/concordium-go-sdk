package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetIdentityProviders get the identity providers registered as of the end of a given block.
// The stream will end when all the identity providers have been returned.
func (c *Client) GetIdentityProviders(ctx context.Context, req isBlockHashInput) (_ []*pb.IpInfo, err error) {
	stream, err := c.GrpcClient.GetIdentityProviders(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var ipInfos []*pb.IpInfo

	for err == nil {
		ipInfo, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		ipInfos = append(ipInfos, ipInfo)
	}

	return ipInfos, nil
}
