package concordium

type PeerElementCatchupStatus int32

const (
	// PeerElementUpToDate means that the peer does not have any data unknown to us. If we receive a message from the
	// peer that refers to unknown data (e.g., an unknown block) the peer is marked as pending.
	PeerElementUpToDate PeerElementCatchupStatus = 0
	// PeerElementPending means that the peer might have some data unknown to us. A peer can be in this state either because
	// it sent a message that refers to data unknown to us, or before we have established a baseline with it.
	// The latter happens during node startup, as well as upon protocol updates until the initial catchup handshake
	// completes.
	PeerElementPending PeerElementCatchupStatus = 1
	// PeerElementCatchingUp means that the node is currently catching up by requesting blocks from this peer.
	// There will be at most one peer with this status at a time.
	// Once the peer has responded to the request, its status will be changed to:
	// - 'UPTODATE' if the peer has no more data that is not known to us
	// - 'PENDING' if the node has more data that is unknown to us.
	PeerElementCatchingUp PeerElementCatchupStatus = 2
)

type PeerType string

const (
	PeerTypeNode         PeerType = "Node"
	PeerTypeBootstrapper PeerType = "Bootstrapper"
)

// PeerElement is a peer node.
type PeerElement struct {
	// The id of the node.
	NodeId NodeId
	// The IP of the node.
	Ip string
	// The port of the node.
	Port uint32
	// The current status of the peer.
	CatchupStatus PeerElementCatchupStatus
}

// PeerList contains a list of peers.
type PeerList struct {
	// The type of the queried node.
	Type PeerType
	// A list of peers.
	Peers []*PeerElement
}

type PeerStatsElement struct {
	// The node id.
	NodeId NodeId
	// The number of messages sent to the peer.
	PacketsSent uint64
	// The number of messages received from the peer.
	PacketsReceived uint64
	// The connection latency (i.e., ping time) in milliseconds.
	Latency uint64
}

// PeerStats contains information about a peer.
type PeerStats struct {
	// A list of stats for the peers.
	Peers []*PeerStatsElement
	// Average outbound throughput in bytes per second.
	AvgBpsIn uint64
	// Average inbound throughput in bytes per second.
	AvgBpsOut uint64
}
