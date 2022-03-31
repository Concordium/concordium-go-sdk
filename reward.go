package concordium

type RewardStatus struct {
	TotalAmount               *Amount `json:"totalAmount"`
	TotalEncryptedAmount      *Amount `json:"totalEncryptedAmount"`
	GasAccount                *Amount `json:"gasAccount"`
	BakingRewardAccount       *Amount `json:"bakingRewardAccount"`
	FinalizationRewardAccount *Amount `json:"finalizationRewardAccount"`
}
