package concordium

const (
	NodeInfoNotInCommittee               = 0
	NodeInfoAddedButNotActiveInCommittee = 1
	NodeInfoAddedButWrongKeys            = 2
	NodeInfoActiveInCommittee            = 3
)

type NodeInfoIsInBakingCommittee int

type NodeId string

type NodeInfo struct {
	NodeId                      NodeId
	CurrentLocaltime            int64
	PeerType                    PeerType
	ConsensusBakerRunning       bool
	ConsensusRunning            bool
	ConsensusType               ConsensusType
	ConsensusBakerCommittee     NodeInfoIsInBakingCommittee
	ConsensusFinalizerCommittee bool
	BakerId                     BakerId
}
