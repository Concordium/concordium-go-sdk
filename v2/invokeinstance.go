package v2

import (
	"context"
)

// InvokeInstance run the smart contract entrypoint in a given context and in the state at the end of the given block.
func (c *Client) InvokeInstance(ctx context.Context, req *InvokeInstanceRequest) (_ *InvokeInstanceResponse, err error) {
	invokeInstanceResponse, err := c.grpcClient.InvokeInstance(ctx, req)
	if err != nil {
		return &InvokeInstanceResponse{}, Error.Wrap(err)
	}

	return invokeInstanceResponse, nil
}
