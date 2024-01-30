package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetBlockCertificates returns the quorum certificate, a timeout certificate (if present) and epoch finalization certificate (if present)
// for a non-genesis block.
//
// This endpoint is only supported for protocol version 6 and onwards.
func (c *Client) GetBlockCertificates(ctx context.Context, req isBlockHashInput) (_ *pb.BlockCertificates, err error) {
	certificates, err := c.GrpcClient.GetBlockCertificates(ctx, convertBlockHashInput(req))
	if err != nil {
		return &pb.BlockCertificates{}, err
	}

	return certificates, nil
}
