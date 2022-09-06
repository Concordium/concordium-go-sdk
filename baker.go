package concordium

type BakerId uint64

// BakerInfo is the state of an individual baker.
type BakerInfo struct {
	// ID of the baker. Matches their account index.
	BakerId BakerId `json:"bakerId"`
	// Address of the account this baker is associated with.
	BakerAccount AccountAddress `json:"bakerAccount"`
	// The lottery power of the baker. This is the baker's stake relative to the total staked amount.
	BakerLotteryPower float64 `json:"bakerLotteryPower"`
}

type BakerParameters struct {
	MinimumThresholdForBaking *Amount `json:"minimumThresholdForBaking"`
}
