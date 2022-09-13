package concordium

import "testing"

const testdataRewardStatus = "testdata/reward_status.json"

var testRewardStatus = &RewardStatus{
	TotalAmount:               NewAmountFromMicroCCD(10872562806746186),
	TotalEncryptedAmount:      NewAmountFromMicroCCD(82443855514),
	BakingRewardAccount:       NewAmountFromMicroCCD(68951065940),
	FinalizationRewardAccount: NewAmountFromMicroCCD(5),
	GasAccount:                NewAmountFromMicroCCD(3),
}

func TestRewardStatus_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &RewardStatus{}, testRewardStatus, testdataRewardStatus)
}
