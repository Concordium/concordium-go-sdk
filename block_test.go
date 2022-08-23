package concordium

import (
	"reflect"
	"testing"
	"time"
)

const testdataBlockInfo = "testdata/block_info.json"

var (
	testBlockHashString = "7c9260d78b8cb702d25c9c9bb166276204f0ef617e867e597c6c7d53b8049267"
	testBlockHashJSON   = []byte(`"` + testBlockHashString + `"`)
	testBlockHash       = MustNewBlockHashFromString(testBlockHashString)

	testBlockInfo = &BlockInfo{
		BlockHash:             testBlockHash,
		BlockParent:           MustNewBlockHashFromString("1dd06682da815e86bc1589ec743103fda2039d0d9fa2ed919f400dc5a7c5b70b"),
		BlockLastFinalized:    MustNewBlockHashFromString("1dd06682da815e86bc1589ec743103fda2039d0d9fa2ed919f400dc5a7c5b70b"),
		BlockHeight:           2705072,
		GenesisIndex:          2,
		EraBlockHeight:        916798,
		BlockReceiveTime:      testTimeMustParse(time.RFC3339, "2022-03-23T20:58:16Z"),
		BlockArriveTime:       testTimeMustParse(time.RFC3339, "2022-03-23T20:58:16Z"),
		BlockSlot:             37122774,
		BlockSlotTime:         testTimeMustParse(time.RFC3339, "2022-03-23T20:58:16.5Z"),
		BlockBaker:            5,
		Finalized:             true,
		TransactionCount:      0,
		TransactionEnergyCost: 0,
		TransactionsSize:      0,
		BlockStateHash:        MustNewBlockHashFromString("e9730d53eed8dafda69c0b0440faabfd4a89e536553a151c657b62c3af827545"),
	}
)

func TestBlockInfo_MarshalJSON(t *testing.T) {
	testFileMarshalJSON(t, testBlockInfo, testdataBlockInfo)
}

func TestBlockInfo_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &BlockInfo{}, testBlockInfo, testdataBlockInfo)
}

func TestBlockHash_MarshalJSON(t *testing.T) {
	got, err := testBlockHash.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(got, testBlockHashJSON) {
		t.Errorf("MarshalJSON() got = %s, w %s", got, testBlockHashJSON)
	}
}

func TestBlockHash_UnmarshalJSON(t *testing.T) {
	var h BlockHash
	err := h.UnmarshalJSON(testBlockHashJSON)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if !reflect.DeepEqual(h, testBlockHash) {
		t.Errorf("UnmarshalJSON() got = %v, w %v", h, testBlockHash)
	}
}
