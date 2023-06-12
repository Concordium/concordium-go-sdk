package v2

import (
	"context"
)

// GetAccountNonFinalizedTransactions get a list of non-finalized transaction hashes for a given account.
// This endpoint is not expected to return a large amount of data in most cases, but in bad network conditions it might.
// The stream will end when all the non-finalized transaction hashes have been returned.
func (c *Client) GetAccountNonFinalizedTransactions(ctx context.Context, req *AccountAddress) (_ Queries_GetAccountNonFinalizedTransactionsClient, err error) {
	stream, err := c.grpcClient.GetAccountNonFinalizedTransactions(ctx, req)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return stream, nil
}
