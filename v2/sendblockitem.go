package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// SendBlockItem send a block item. A block item is either an `AccountTransaction`,
// which is a transaction signed and paid for by an account, a `CredentialDeployment`,
// which creates a new account, or `UpdateInstruction`, which is an instruction to change
// some parameters of the chain. Update instructions can only be sent by the governance committee.
// Returns a hash of the block item, which can be used with `GetBlockItemStatus`.
func (c *Client) SendBlockItem(ctx context.Context, req *pb.SendBlockItemRequest) (_ *TransactionHash, err error) {
	txHash, err := c.GrpcClient.SendBlockItem(ctx, req)
	if err != nil {
		return &TransactionHash{}, err
	}

	var res *TransactionHash
	copy(res.Value[:], txHash.Value)

	return res, nil
}
