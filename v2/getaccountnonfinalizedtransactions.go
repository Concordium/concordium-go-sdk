package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetAccountNonFinalizedTransactions get a list of non-finalized transaction hashes for a given account.
// This endpoint is not expected to return a large amount of data in most cases, but in bad network conditions it might.
// The stream will end when all the non-finalized transaction hashes have been returned.
func (c *Client) GetAccountNonFinalizedTransactions(ctx context.Context, req *AccountAddress) (_ []*TransactionHash, err error) {
	stream, err := c.GrpcClient.GetAccountNonFinalizedTransactions(ctx, &pb.AccountAddress{Value: req.Value[:]})
	if err != nil {
		return nil, err
	}

	var transactionHashes []*pb.TransactionHash

	for {
		transactionHash, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		transactionHashes = append(transactionHashes, transactionHash)
	}

	var result []*TransactionHash

	for i := 0; i < len(transactionHashes); i++ {
		var txHash TransactionHash
		copy(txHash.Value[:], transactionHashes[i].Value)
		result = append(result, &txHash)
	}

	return result, nil
}
