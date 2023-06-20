package v2

import (
	"context"

	"concordium-go-sdk/v2/pb"
)

// GetAccountTransactionSignHash get the hash to be signed for an account transaction.
// The hash returned should be signed and the signatures included as an AccountTransactionSignature when calling `SendBlockItem`.
// This is provided as a convenience to support cases where the right SDK is not available for interacting with the node.
// If an SDK is available then it is strongly recommended to compute this hash off-line using it.
// That reduces the trust in the node, removes networking failure modes, and will perform better.
func (c *Client) GetAccountTransactionSignHash(ctx context.Context, req *pb.PreAccountTransaction) (_ *pb.AccountTransactionSignHash, err error) {
	accountTransactionSignHash, err := c.grpcClient.GetAccountTransactionSignHash(ctx, req)
	if err != nil {
		return &pb.AccountTransactionSignHash{}, err
	}

	return accountTransactionSignHash, nil
}
