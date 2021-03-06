package concordium

import "testing"

const testdataInstanceInfo = "testdata/instance_info.json"

var testInstanceInfo = &InstanceInfo{
	Model:  "00000000030000000101000000620000000002010000006300000000000100000061000000000000006ab2da7f010000d106000000000000000000000000000002",
	Owner:  "3f5mdmn3zVSoC3vKkDNP9YcDJjcK1A2UUiPSsMpZhKLZyKRHNY",
	Amount: NewAmountZero(),
	Methods: []ReceiveName{
		"govogo.cancelVote",
		"govogo.claimToVote",
		"govogo.giveRightToVote",
		"govogo.vote",
		"govogo.winningProposal",
	},
	Name:         "init_govogo",
	SourceModule: "1d40f9366f6fcb4586ac8e09ed391b5832cfd752fb63ee7bd38da0f3e77c4204",
}

func TestInstanceInfo_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, testInstanceInfo, testdataInstanceInfo)
}

func TestInstanceInfo_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, &InstanceInfo{}, testInstanceInfo, testdataInstanceInfo)
}
