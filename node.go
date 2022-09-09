package concordium

// NodeId is a node identifier.
type NodeId string

type NodeInfoBakingCommittee int32

const (
	// NodeInfoNotInCommittee means that the node is not the baking committee.
	NodeInfoNotInCommittee NodeInfoBakingCommittee = 0
	// NodeInfoAddedButNotActiveInCommittee means that the node has baker keys, but the
	// account is not currently a baker (and possibly never will be).
	NodeInfoAddedButNotActiveInCommittee NodeInfoBakingCommittee = 1
	// NodeInfoAddedButWrongKeys means that the node has baker keys, but they don't match
	// the current keys on the baker account.
	NodeInfoAddedButWrongKeys NodeInfoBakingCommittee = 2
	// NodeInfoActiveInCommittee means that the node has valid baker keys and is active in
	// the baker committee.
	NodeInfoActiveInCommittee NodeInfoBakingCommittee = 3
)

// NodeInfo contains information about the node.
type NodeInfo struct {
	// The unique node identifier.
	NodeId NodeId
	// The local time of the node represented as a unix timestamp in seconds.
	CurrentLocaltime uint64
	// The node type.
	PeerType PeerType
	// Whether the node is a baker.
	ConsensusBakerRunning bool
	// Whether consensus is running. This is only false if the protocol was updated
	// to a version which the node software does not support.
	ConsensusRunning bool
	// The node consensus type.
	ConsensusType ConsensusType
	// The baking status of the node.
	ConsensusBakerCommittee NodeInfoBakingCommittee
	// Whether the node is part of the finalization committee.
	ConsensusFinalizerCommittee bool
	// The baker id. This will be 0 if the node is not a baker.
	ConsensusBakerId BakerId
}
