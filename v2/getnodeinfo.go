package v2

import (
	"context"
)

// GetNodeInfo get information about the node.
// The `NodeInfo` includes information of * Meta information such as the, version of the node,
// type of the node, uptime and the local time of the node.
// * NetworkInfo which yields data such as the node id, packets sent/received, average bytes per second sent/received.
// * ConsensusInfo. The `ConsensusInfo` returned depends on if the node supports the protocol on chain and whether the node is configured as a baker or not.
func (c *Client) GetNodeInfo(ctx context.Context) (_ *NodeInfo, err error) {
	nodeInfo, err := c.grpcClient.GetNodeInfo(ctx, new(Empty))
	if err != nil {
		return &NodeInfo{}, Error.Wrap(err)
	}

	return nodeInfo, nil
}
