package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetAccountList retrieve the list of accounts that exist at the end of the given block.
func (c *Client) GetAccountList(ctx context.Context, req isBlockHashInput) (_ []*AccountAddress, err error) {
	stream, err := c.GrpcClient.GetAccountList(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var accounts []*pb.AccountAddress

	for err == nil {
		account, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		accounts = append(accounts, account)
	}

	var result []*AccountAddress

	for i := 0; i < len(accounts); i++ {
		var accountAddress AccountAddress
		copy(accountAddress.Value[:], accounts[i].Value)
		result = append(result, &accountAddress)
	}

	return result, nil
}
