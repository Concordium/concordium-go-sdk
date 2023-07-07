package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// InvokeInstance run the smart contract entrypoint in a given context and in the state at the end of the given block.
func (c *Client) InvokeInstance(ctx context.Context, req *pb.InvokeInstanceRequest) (_ *pb.InvokeInstanceResponse, err error) {
	invokeInstanceResponse, err := c.GrpcClient.InvokeInstance(ctx, req)
	if err != nil {
		return &pb.InvokeInstanceResponse{}, err
	}

	return invokeInstanceResponse, nil
}
