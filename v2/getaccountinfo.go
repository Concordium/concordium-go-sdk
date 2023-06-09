package v2

import (
	"context"
)

// GetAccountInfo retrieve the information about the given account in the given block.
func (c *Client) GetAccountInfo(ctx context.Context, req *AccountInfoRequest) (_ *AccountInfo, err error) {
	accountInfo, err := c.grpcClient.GetAccountInfo(ctx, req)
	if err != nil {
		return &AccountInfo{}, Error.Wrap(err)
	}

	return accountInfo, nil
}
