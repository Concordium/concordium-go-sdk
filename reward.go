package concordium

import "time"

// RewardStatus is common to both V0 and V1 rewards.
type RewardStatus struct {
	// The amount in the baking reward account.
	BakingRewardAccount *Amount `json:"bakingRewardAccount"`
	// The amount in the finalization reward account.
	FinalizationRewardAccount *Amount `json:"finalizationRewardAccount"`
	// The transaction reward fraction accruing to the foundation (to be paid at next payday).
	FoundationTransactionRewards *Amount `json:"foundationTransactionRewards"`
	// The amount in the GAS account.
	GasAccount *Amount `json:"gasAccount"`
	// The rate at which CCD will be minted (as a proportion of the total supply) at the next payday
	NextPaydayMintRate *float64 `json:"nextPaydayMintRate"`
	// The time of the next payday.
	NextPaydayTime *time.Time `json:"nextPaydayTime"`
	// Protocol version that applies to these rewards. V0 variant only exists for protocol versions 1, 2, and 3.
	ProtocolVersion uint64 `json:"protocolVersion"`
	// The total CCD in existence.
	TotalAmount *Amount `json:"totalAmount"`
	// The total CCD in encrypted balances.
	TotalEncryptedAmount *Amount `json:"totalEncryptedAmount"`
	// The total capital put up as stake by bakers and delegators
	TotalStakedCapital *Amount `json:"totalStakedCapital"`
}
