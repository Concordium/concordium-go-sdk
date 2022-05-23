package concordium

const (
	PeerElementUpToDate   PeerElementCatchupStatus = 0
	PeerElementPending    PeerElementCatchupStatus = 1
	PeerElementCatchingUp PeerElementCatchupStatus = 2
)

type PeerElementCatchupStatus int

type PeerType string

type PeerElement struct {
	NodeId        NodeId
	Ip            string
	Port          int
	CatchupStatus PeerElementCatchupStatus
}

type PeerList struct {
	Type  PeerType
	Peers []*PeerElement
}

type PeerStatsElement struct {
	NodeId          NodeId
	PacketsSent     int
	PacketsReceived int
	Latency         int
}

type PeerStats struct {
	Peers     []*PeerStatsElement
	AvgBpsIn  int
	AvgBpsOut int
}
