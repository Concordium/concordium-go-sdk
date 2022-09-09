package concordium

import (
	"testing"
)

const testdataBranch = "testdata/branch.json"

var testBranch = &Branch{
	BlockHash: MustNewBlockHash("030cf811f9e706188e4b2fbcb36fee7146d104a5b1e98162883bdbe63be4071e"),
	Children:  []*Branch{},
}

func TestBranch_MarshalJSON(t *testing.T) {
	testFileMarshalJSON(t, testBranch, testdataBranch)
}

func TestBranch_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &Branch{}, testBranch, testdataBranch)
}
