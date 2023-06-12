package v2

import (
	"context"
)

// PeerDisconnect disconnect from the peer and remove them from the given addresses list if they are on it.
// Return if the request was processed successfully. Otherwise return a GRPC error.
func (c *Client) PeerDisconnect(ctx context.Context, req *IpSocketAddress) (err error) {
	_, err = c.grpcClient.PeerDisconnect(ctx, req)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}
