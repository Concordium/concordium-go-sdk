package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// SendBlockItem send a block item. A block item is either an `AccountTransaction`,
// which is a transaction signed and paid for by an account, a `CredentialDeployment`,
// which creates a new account, or `UpdateInstruction`, which is an instruction to change
// some parameters of the chain. Update instructions can only be sent by the governance committee.
// Returns a hash of the block item, which can be used with `GetBlockItemStatus`.
func (c *Client) SendBlockItem(ctx context.Context, req *pb.SendBlockItemRequest) (_ *pb.TransactionHash, err error) {
	txHash, err := c.grpcClient.SendBlockItem(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return txHash, nil
}
