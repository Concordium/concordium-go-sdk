package concordium

type BakerId uint64

type BakerInfo struct {
	BakerId           BakerId        `json:"bakerId"`
	BakerAccount      AccountAddress `json:"bakerAccount"`
	BakerLotteryPower float64        `json:"bakerLotteryPower"`
}
