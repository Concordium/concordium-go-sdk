package concordium

import (
	"testing"
)

const testdataBakerInfo = "testdata/baker_info.json"

var testBakerInfo = &BakerInfo{
	BakerId:           0,
	BakerLotteryPower: 9.96180263437706e-2,
	BakerAccount:      "3Ug5rCqAN2z17MqAyh5KUGDpv6k9eSHu8AN8jCgbmAxmjcu5TM",
}

func TestBakerInfo_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, testBakerInfo, testdataBakerInfo)
}

func TestBakerInfo_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, &BakerInfo{}, testBakerInfo, testdataBakerInfo)
}
