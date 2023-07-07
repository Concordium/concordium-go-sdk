package v2

import (
	"context"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// GetModuleList get a list of all smart contract modules. The stream will end when
// all modules that exist in the state at the end of the given block have been returned.
func (c *Client) GetModuleList(ctx context.Context, req isBlockHashInput) (_ []*ModuleRef, err error) {
	stream, err := c.GrpcClient.GetModuleList(ctx, convertBlockHashInput(req))
	if err != nil {
		return nil, err
	}

	var moduleRefs []*pb.ModuleRef

	for err == nil {
		moduleRef, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, err
		}

		moduleRefs = append(moduleRefs, moduleRef)
	}

	var result []*ModuleRef

	for i := 0; i < len(moduleRefs); i++ {
		var m ModuleRef
		copy(m[:], moduleRefs[i].Value)
		result = append(result, &m)
	}

	return result, nil
}
