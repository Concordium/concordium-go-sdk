package concordium

import (
	"testing"
	"time"
)

const testdataConsensusStatus = "testdata/consensus_status.json"

var testConsensusStatus = &ConsensusStatus{
	BestBlock:                "8671be6e747c09b9ee694c33f909bcd04e5a3497efcf5858a71e27b42d569e39",
	GenesisBlock:             "b6078154d6717e909ce0da4a45a25151b592824f31624b755900a74429e3073d",
	GenesisTime:              testTimeMustParse(time.RFC3339, "2021-05-07T12:00:00Z"),
	SlotDuration:             250,
	EpochDuration:            3600000,
	LastFinalizedBlock:       "8671be6e747c09b9ee694c33f909bcd04e5a3497efcf5858a71e27b42d569e39",
	BestBlockHeight:          2965925,
	LastFinalizedBlockHeight: 2965925,
	BlocksReceivedCount:      1200833,
	BlockLastReceivedTime:    testTimeMustParse(time.RFC3339, "2022-04-23T11:38:05.519319626Z"),
	BlockReceiveLatencyEMA:   0.39017566218808836,
	BlockReceiveLatencyEMSD:  0.14564296734552024,
	BlockReceivePeriodEMA:    12.264950053211724,
	BlockReceivePeriodEMSD:   14.771515854036915,
	BlocksVerifiedCount:      1200833,
	BlockLastArrivedTime:     testTimeMustParse(time.RFC3339, "2022-04-23T11:38:05.52783592Z"),
	BlockArriveLatencyEMA:    0.39862098584103156,
	BlockArriveLatencyEMSD:   0.14664278987775503,
	BlockArrivePeriodEMA:     12.264957938805841,
	BlockArrivePeriodEMSD:    14.772069494891108,
	TransactionsPerBlockEMA:  2.7856107500049184e-3,
	TransactionsPerBlockEMSD: 5.2705323485886274e-2,
	FinalizationCount:        1022079,
	LastFinalizedTime:        testTimeMustParse(time.RFC3339, "2022-04-23T11:38:06.860688322Z"),
	FinalizationPeriodEMA:    13.033186677708576,
	FinalizationPeriodEMSD:   14.942897961412053,
	ProtocolVersion:          3,
	GenesisIndex:             2,
	CurrentEraGenesisBlock:   "396e1ac6b3cd4aef76fbd463275f270dbacf70294d3e33eee8e28cfb51aa1625",
	CurrentEraGenesisTime:    testTimeMustParse(time.RFC3339, "2021-12-06T11:00:03Z"),
}

func TestConsensusStatus_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, testConsensusStatus, testdataConsensusStatus)
}

func TestConsensusStatus_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, &ConsensusStatus{}, testConsensusStatus, testdataConsensusStatus)
}
