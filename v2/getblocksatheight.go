package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetBlocksAtHeight get a list of live blocks at a given height.
func (c *Client) GetBlocksAtHeight(ctx context.Context, req *pb.BlocksAtHeightRequest) (_ []BlockHash, err error) {
	blockAtHeight, err := c.GrpcClient.GetBlocksAtHeight(ctx, req)
	if err != nil {
		return []BlockHash{}, err
	}

	var blockHashes []BlockHash
	for i := 0; i < len(blockAtHeight.Blocks); i++ {
		var blockHash BlockHash
		copy(blockHash[:], blockAtHeight.Blocks[i].Value)
		blockHashes = append(blockHashes, blockHash)
	}

	return blockHashes, nil
}
