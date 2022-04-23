package concordium

import "testing"

const testdataRewardStatus = "testdata/reward_status.json"

var testRewardStatus = &RewardStatus{
	TotalAmount:               NewAmountFromMicroGTU(10872562806746186),
	TotalEncryptedAmount:      NewAmountFromMicroGTU(82443855514),
	BakingRewardAccount:       NewAmountFromMicroGTU(68951065940),
	FinalizationRewardAccount: NewAmountFromMicroGTU(5),
	GasAccount:                NewAmountFromMicroGTU(3),
}

func TestRewardStatus_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, testRewardStatus, testdataRewardStatus)
}

func TestRewardStatus_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, &RewardStatus{}, testRewardStatus, testdataRewardStatus)
}
