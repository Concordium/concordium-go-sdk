package concordium

const (
	PeerElementUpToDate   PeerElementCatchupStatus = 0
	PeerElementPending    PeerElementCatchupStatus = 1
	PeerElementCatchingUp PeerElementCatchupStatus = 2

	PeerTypeNode         PeerType = "Node"
	PeerTypeBootstrapper PeerType = "Bootstrapper"
)

type PeerElementCatchupStatus int

type PeerType string

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
