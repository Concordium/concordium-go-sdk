package concordium

type BirkParameters struct {
	ElectionDifficulty float64      `json:"electionDifficulty"`
	ElectionNonce      string       `json:"electionNonce"`
	Bakers             []*BakerInfo `json:"bakers"`
}
