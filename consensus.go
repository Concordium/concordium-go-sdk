package concordium

import "time"

const (
	// ConsensusTypeActive indicates that the node has baker credentials and
	// can thus potentially participate in baking and finalization.
	ConsensusTypeActive ConsensusType = "Active"

	// ConsensusTypePassive indicates that the node has no baker credentials
	// is thus only an observer of the consensus protocol.
	ConsensusTypePassive ConsensusType = "Passive"
)

// ConsensusType is the node consensus type.
type ConsensusType string

// ConsensusStatus is summary of the current state of consensus.
type ConsensusStatus struct {
	// Hash of the current best block. The best block is a protocol defined block that
	// the node must use a parent block to build the chain on. Note that this is subjective,
	// in the sense that it is only the best block among the blocks the node knows about.
	BestBlock BlockHash `json:"bestBlock"`
	// Hash of the genesis block.
	GenesisBlock BlockHash `json:"genesisBlock"`
	// Slot time of the genesis block.
	GenesisTime time.Time `json:"genesisTime"`
	// Duration of a slot.
	SlotDuration uint64 `json:"slotDuration"`
	// Duration of an epoch.
	EpochDuration int64 `json:"epochDuration"`
	// Hash of the last, i.e., most recent, finalized block.
	LastFinalizedBlock BlockHash `json:"lastFinalizedBlock"`
	// Height of the best block. See ConsensusStatus.BestBlock.
	BestBlockHeight BlockHeight `json:"bestBlockHeight"`
	// Height of the last finalized block. Genesis block has height 0.
	LastFinalizedBlockHeight BlockHeight `json:"lastFinalizedBlockHeight"`
	// The number of blocks that have been received.
	BlocksReceivedCount uint64 `json:"blocksReceivedCount"`
	// The time (local time of the node) that a block was last received.
	BlockLastReceivedTime *time.Time `json:"blockLastReceivedTime"`
	// Exponential moving average of block receive latency (in seconds), i.e. the
	// time between a block's nominal slot time, and the time at which is received.
	BlockReceiveLatencyEMA float64 `json:"blockReceiveLatencyEMA"`
	// Exponential moving average standard deviation of block receive latency (in seconds),
	// i.e. the time between a block's nominal slot time, and the time at which is received.
	BlockReceiveLatencyEMSD float64 `json:"blockReceiveLatencyEMSD"`
	// Exponential moving average of the time between receiving blocks (in seconds).
	BlockReceivePeriodEMA *float64 `json:"blockReceivePeriodEMA"`
	// Exponential moving average standard deviation of the time between receiving
	// blocks (in seconds).
	BlockReceivePeriodEMSD *float64 `json:"blockReceivePeriodEMSD"`
	// Number of blocks that arrived, i.e., were added to the tree. Note that in some cases
	// this can be more than [ConsensusInfo::blocks_received_count] since blocks that the
	// node itself produces count towards this, but are not received.
	BlocksVerifiedCount uint64 `json:"blocksVerifiedCount"`
	// The time (local time of the node) that a block last arrived, i.e., was verified and
	// added to the node's tree.
	BlockLastArrivedTime *time.Time `json:"blockLastArrivedTime"`
	// The exponential moving average of the time between a block's nominal slot time,
	// and the time at which it is verified.
	BlockArriveLatencyEMA float64 `json:"blockArriveLatencyEMA"`
	// The exponential moving average standard deviation of the time between a block's
	// nominal slot time, and the time at which it is verified.
	BlockArriveLatencyEMSD float64 `json:"blockArriveLatencyEMSD"`
	// Exponential moving average of the time between receiving blocks (in seconds).
	BlockArrivePeriodEMA *float64 `json:"blockArrivePeriodEMA"`
	// Exponential moving average standard deviation of the time between blocks being verified.
	BlockArrivePeriodEMSD *float64 `json:"blockArrivePeriodEMSD"`
	// Exponential moving average of the number of transactions per block.
	TransactionsPerBlockEMA float64 `json:"transactionsPerBlockEMA"`
	// Exponential moving average standard deviation of the number of transactions per block.
	TransactionsPerBlockEMSD float64 `json:"transactionsPerBlockEMSD"`
	// The number of completed finalizations.
	FinalizationCount uint64 `json:"finalizationCount"`
	// Time at which a block last became finalized. Note that this is the local time of the
	// node at the time the block was finalized.
	LastFinalizedTime *time.Time `json:"lastFinalizedTime"`
	// Exponential moving average of the time between finalizations. Will be `None` if
	// there are no finalizations yet since the node start.
	FinalizationPeriodEMA *float64 `json:"finalizationPeriodEMA"`
	// Exponential moving average standard deviation of the time between finalizations.
	// Will be `None` if there are no finalizations yet since the node start.
	FinalizationPeriodEMSD *float64 `json:"finalizationPeriodEMSD"`
	// Currently active protocol version.
	ProtocolVersion uint64 `json:"protocolVersion"`
	// The number of chain restarts via a protocol update. An effected protocol update
	// instruction might not change the protocol version specified in the previous
	// field, but it always increments the genesis index.
	GenesisIndex uint32 `json:"genesisIndex"`
	// Block hash of the genesis block of current era, i.e., since the last protocol
	// update. Initially this is equal to ConsensusStatus.GenesisBlock.
	CurrentEraGenesisBlock BlockHash `json:"currentEraGenesisBlock"`
	// Time when the current era started.
	CurrentEraGenesisTime time.Time `json:"currentEraGenesisTime"`
}
