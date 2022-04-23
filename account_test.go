package concordium

import (
	"testing"
)

const testdataNextAccountNonce = "testdata/next_account_nonce.json"
const testdataAccountInfo = "testdata/account_info.json"

var testNextAccountNonce = &NextAccountNonce{
	Nonce:    2,
	AllFinal: true,
}

var testAccountInfo = &AccountInfo{
	AccountNonce:  2,
	AccountAmount: NewAmountFromMicroGTU(2166),
	AccountReleaseSchedule: &AccountReleaseSchedule{
		Total:    NewAmountFromMicroGTU(0),
		Schedule: []*Release{},
	},
	AccountThreshold: 1,
	AccountEncryptedAmount: &AccountEncryptedAmount{
		IncomingAmounts: []EncryptedAmount{},
		SelfAmount:      "c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		StartIndex:      0,
	},
	AccountEncryptionKey: "b14cbfe44a02c6b1f78711176d5f437295367aa4f2a8c2551ee10d25a03adc69d61a332a058971919dad7312e1fc94c58b62436ed271a973b0cb294b7515e8be5d085bf89c54f6ad2c7c54c3b5b1a6c51872d80ff5953a2e8f284148351fef13",
	AccountIndex:         1731,
	AccountAddress:       "4tUQrKhVKPN5pEev5joobH6n8RR5sXX6u6REdWZxi7NVvoZhVc",
}

func TestNextAccountNonce_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, testNextAccountNonce, testdataNextAccountNonce)
}

func TestNextAccountNonce_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, &NextAccountNonce{}, testNextAccountNonce, testdataNextAccountNonce)
}

func TestAccountInfo_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, testAccountInfo, testdataAccountInfo)
}

func TestAccountInfo_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, &AccountInfo{}, testAccountInfo, testdataAccountInfo)
}
