package concordium

import (
	"testing"
)

const testdataNextAccountNonce = "testdata/next_account_nonce.json"
const testdataAccountInfo = "testdata/account_info.json"

var (
	testNextAccountNonce = &NextAccountNonce{
		Nonce:    2,
		AllFinal: true,
	}

	testEncryptionKeyString = "b14cbfe44a02c6b1f78711176d5f437295367aa4f2a8c2551ee10d25a03adc69d61a332a058971919dad7312e1fc94c58b62436ed271a973b0cb294b7515e8be5d085bf89c54f6ad2c7c54c3b5b1a6c51872d80ff5953a2e8f284148351fef13"
	testEncryptionKeyJSON   = []byte(`"` + testEncryptionKeyString + `"`)
	testEncryptionKey       = MustNewAccountEncryptionKeyFromString(testEncryptionKeyString)

	testEncryptedAmountString = "c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	testEncryptedAmountJSON   = []byte(`"` + testEncryptedAmountString + `"`)
	testEncryptedAmount       = MustNewEncryptedAmountFromString(testEncryptedAmountString)
)

var testAccountInfo = &AccountInfo{
	AccountNonce:  2,
	AccountAmount: NewAmountFromMicroCCD(2166),
	AccountReleaseSchedule: &AccountReleaseSchedule{
		Total:    NewAmountFromMicroCCD(0),
		Schedule: []*Release{},
	},
	AccountThreshold: 1,
	AccountEncryptedAmount: &AccountEncryptedAmount{
		IncomingAmounts: []EncryptedAmount{},
		SelfAmount:      testEncryptedAmount,
		StartIndex:      0,
	},
	AccountEncryptionKey: testEncryptionKey,
	AccountIndex:         1731,
	AccountAddress:       MustNewAccountAddressFromString("4tUQrKhVKPN5pEev5joobH6n8RR5sXX6u6REdWZxi7NVvoZhVc"),
}

func TestNextAccountNonce_MarshalJSON(t *testing.T) {
	testFileMarshalJSON(t, testNextAccountNonce, testdataNextAccountNonce)
}

func TestNextAccountNonce_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &NextAccountNonce{}, testNextAccountNonce, testdataNextAccountNonce)
}

func TestAccountInfo_MarshalJSON(t *testing.T) {
	testFileMarshalJSON(t, testAccountInfo, testdataAccountInfo)
}

func TestAccountInfo_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &AccountInfo{}, testAccountInfo, testdataAccountInfo)
}

func TestAccountEncryptionKey_MarshalJSON(t *testing.T) {
	testHexMarshalJSON(t, testEncryptionKey, testEncryptionKeyJSON)
}

func TestAccountEncryptionKey_UnmarshalJSON(t *testing.T) {
	var a AccountEncryptionKey
	testHexUnmarshalJSON(t, &a, &testEncryptionKey, testEncryptionKeyJSON)
}

func TestEncryptedAmount_MarshalJSON(t *testing.T) {
	testHexMarshalJSON(t, testEncryptedAmount, testEncryptedAmountJSON)
}

func TestEncryptedAmount_UnmarshalJSON(t *testing.T) {
	var a EncryptedAmount
	testHexUnmarshalJSON(t, &a, &testEncryptedAmount, testEncryptedAmountJSON)
}
