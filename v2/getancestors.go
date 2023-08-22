package v2

import (
	"context"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// GetAncestors get a stream of ancestors for the provided block.
// Starting with the provided block itself, moving backwards until
// no more ancestors or the requested number of ancestors has been returned.
func (c *Client) GetAncestors(ctx context.Context, amount uint64, b isBlockHashInput) (_ []BlockHash, err error) {
	stream, err := c.GrpcClient.GetAncestors(ctx, &pb.AncestorsRequest{
		BlockHash: convertBlockHashInput(b),
		Amount:    amount,
	})
	if err != nil {
		return nil, err
	}

	var ancestorsBlockHashes []*pb.BlockHash

	for err == nil {
		ancestorBlockHash, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		ancestorsBlockHashes = append(ancestorsBlockHashes, ancestorBlockHash)
	}

	var result []BlockHash

	for i := 0; i < len(ancestorsBlockHashes); i++ {
		var blockHash BlockHash
		copy(blockHash.Value[:], ancestorsBlockHashes[i].Value)
		result = append(result, blockHash)
	}

	return result, nil
}
