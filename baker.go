package concordium

type BakerId uint64

type BakerInfo struct {
	BakerId           BakerId        `json:"bakerId"`
	BakerAccount      AccountAddress `json:"bakerAccount"`
	BakerLotteryPower float64        `json:"bakerLotteryPower"`
}

type BirkParameters struct {
	ElectionDifficulty float64      `json:"electionDifficulty"`
	ElectionNonce      string       `json:"electionNonce"`
	Bakers             []*BakerInfo `json:"bakers"`
}
