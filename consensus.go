package concordium

import "time"

type ConsensusType string

type ConsensusStatus struct {
	BestBlock                BlockHash   `json:"bestBlock"`
	GenesisBlock             BlockHash   `json:"genesisBlock"`
	GenesisTime              time.Time   `json:"genesisTime"`
	SlotDuration             int64       `json:"slotDuration"`
	EpochDuration            int64       `json:"epochDuration"`
	LastFinalizedBlock       BlockHash   `json:"lastFinalizedBlock"`
	BestBlockHeight          BlockHeight `json:"bestBlockHeight"`
	LastFinalizedBlockHeight BlockHeight `json:"lastFinalizedBlockHeight"`
	BlocksReceivedCount      int         `json:"blocksReceivedCount"`
	BlockLastReceivedTime    time.Time   `json:"blockLastReceivedTime"`
	BlockReceiveLatencyEMA   float64     `json:"blockReceiveLatencyEMA"`
	BlockReceiveLatencyEMSD  float64     `json:"blockReceiveLatencyEMSD"`
	BlockReceivePeriodEMA    float64     `json:"blockReceivePeriodEMA"`
	BlockReceivePeriodEMSD   float64     `json:"blockReceivePeriodEMSD"`
	BlocksVerifiedCount      int         `json:"blocksVerifiedCount"`
	BlockLastArrivedTime     time.Time   `json:"blockLastArrivedTime"`
	BlockArriveLatencyEMA    float64     `json:"blockArriveLatencyEMA"`
	BlockArriveLatencyEMSD   float64     `json:"blockArriveLatencyEMSD"`
	BlockArrivePeriodEMA     float64     `json:"blockArrivePeriodEMA"`
	BlockArrivePeriodEMSD    float64     `json:"blockArrivePeriodEMSD"`
	TransactionsPerBlockEMA  float64     `json:"transactionsPerBlockEMA"`
	TransactionsPerBlockEMSD float64     `json:"transactionsPerBlockEMSD"`
	FinalizationCount        int         `json:"finalizationCount"`
	LastFinalizedTime        time.Time   `json:"lastFinalizedTime"`
	FinalizationPeriodEMA    float64     `json:"finalizationPeriodEMA"`
	FinalizationPeriodEMSD   float64     `json:"finalizationPeriodEMSD"`
	ProtocolVersion          int         `json:"protocolVersion"`
	GenesisIndex             int         `json:"genesisIndex"`
	CurrentEraGenesisBlock   BlockHash   `json:"currentEraGenesisBlock"`
	CurrentEraGenesisTime    time.Time   `json:"currentEraGenesisTime"`
}
