package concordium

const (
	NodeInfoNotInCommittee               = 0
	NodeInfoAddedButNotActiveInCommittee = 1
	NodeInfoAddedButWrongKeys            = 2
	NodeInfoActiveInCommittee            = 3
)

type NodeInfoIsInBakingCommittee int

type NodeId string

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
	ConsensusBakerCommittee NodeInfoIsInBakingCommittee
	// Whether the node is part of the finalization committee.
	ConsensusFinalizerCommittee bool
	// The baker id. This will be 0 if the node is not a baker.
	ConsensusBakerId BakerId
}
