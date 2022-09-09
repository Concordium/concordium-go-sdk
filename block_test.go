package concordium

import (
	"reflect"
	"testing"
	"time"
)

const (
	testdataBlockInfo    = "testdata/block_info.json"
	testdataBlockSummary = "testdata/block_summary.json"
)

var (
	testBlockHashString = "7c9260d78b8cb702d25c9c9bb166276204f0ef617e867e597c6c7d53b8049267"
	testBlockHashJSON   = []byte(`"` + testBlockHashString + `"`)
	testBlockHash       = MustNewBlockHash(testBlockHashString)

	testBlockInfo = &BlockInfo{
		BlockHash:             testBlockHash,
		BlockParent:           MustNewBlockHash("1dd06682da815e86bc1589ec743103fda2039d0d9fa2ed919f400dc5a7c5b70b"),
		BlockLastFinalized:    MustNewBlockHash("1dd06682da815e86bc1589ec743103fda2039d0d9fa2ed919f400dc5a7c5b70b"),
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
		BlockStateHash:        MustNewBlockHash("e9730d53eed8dafda69c0b0440faabfd4a89e536553a151c657b62c3af827545"),
	}

	testBlockSummary = &BlockSummary{
		FinalizationData: &FinalizationData{
			FinalizationBlockPointer: MustNewBlockHash("c7077bf80cefef660eb5ae99ad3fa9137c1f9cc7ed3ea019a8058d3e5c81aa07"),
			FinalizationDelay:        0,
			FinalizationIndex:        530252,
			Finalizers: []*Finalizer{
				{
					BakerId: 0,
					Signed:  true,
					Weight:  7119286292406144,
				},
			},
		},
		ProtocolVersion: 4,
		SpecialEvents: []*SpecialEvent{
			{
				BakerId:          5,
				BakerReward:      "758296",
				FoundationCharge: "167207",
				NewGASAccount:    "769692",
				OldGASAccount:    "23138",
				PassiveReward:    "10",
				Tag:              "BlockAccrueReward",
				TransactionFees:  "1672067",
			},
		},
		TransactionSummaries: []*TransactionOutcome{
			{
				Cost:       NewAmountFromMicroCCD(1672067),
				EnergyCost: 999,
				Hash:       "71eda267f1e717fad5867c64eb2e362dc6467abbb01807528e91edcb3dc65e41",
				Index:      0,
				Result: &TransactionResult{
					Events: TransactionResultEvents{
						{
							Tag: "ContractInitialized",
							Raw: []byte(`{"tag": "ContractInitialized"}`),
						},
					},
					Outcome: "success",
				},
				Sender: MustNewAccountAddress("4hvvPeHb9HY4Lur7eUZv4KfL3tYBug8DRc4X9cVU8mpJLa1f2X"),
				Type: &TransactionType{
					Contents: "initContract",
					Type:     "accountTransaction",
				},
			},
		},
		Updates: &Updates{
			ChainParameters: &ChainParameters{
				AccountCreationLimit: 10,
				BakingCommissionRange: &CommissionRange{
					Max: 0.1,
					Min: 0.1,
				},
				CapitalBound:       0.1,
				DelegatorCooldown:  1209600,
				ElectionDifficulty: 0.025,
				EuroPerEnergy: &Fraction{
					Denominator: 50000,
					Numerator:   1,
				},
				FinalizationCommissionRange: &CommissionRange{
					Max: 1.0,
					Min: 1.0,
				},
				FoundationAccountIndex: 11,
				LeverageBound: &Fraction{
					Denominator: 1,
					Numerator:   3,
				},
				MicroGTUPerEuro: &Fraction{
					Denominator: 103736215559,
					Numerator:   8681372131440148480,
				},
				MinimumEquityCapital:          "14000000000",
				MintPerPayday:                 0.000261157877,
				PassiveBakingCommission:       0.12,
				PassiveFinalizationCommission: 1.0,
				PassiveTransactionCommission:  0.12,
				PoolOwnerCooldown:             1814400,
				RewardParameters: &RewardParameters{
					GASRewards: &GASRewards{
						AccountCreation:   0.02,
						Baker:             0.25,
						ChainUpdate:       0.005,
						FinalizationProof: 0.005,
					},
					MintDistribution: &MintDistribution{
						BakingReward:       0.6,
						FinalizationReward: 0.3,
					},
					TransactionFeeDistribution: &TransactionFeeDistribution{
						Baker:      0.45,
						GasAccount: 0.45,
					},
				},
				RewardPeriodLength: 24,
				TransactionCommissionRange: &CommissionRange{
					Max: 0.1,
					Min: 0.1,
				},
			},
			Keys: &UpdateKeys{
				Level1Keys: &Level1Keys{
					Keys: []*Level1Key{
						{
							SchemeId:  "Ed25519",
							VerifyKey: "55721372d942742db382c3680737a424fe3234cf82a80d42da008d6b47179500",
						},
					},
					Threshold: 7,
				},
				Level2Keys: &Level2Keys{
					AddAnonymityRevoker: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					AddIdentityProvider: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					CooldownParameters: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					ElectionDifficulty: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					Emergency: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					EuroPerEnergy: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					FoundationAccount: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					Keys: []*Level1Key{
						{
							SchemeId:  "Ed25519",
							VerifyKey: "b8ddf4505a37eee2c046671f634b74cf3630f3958ad70a04f39dc843041965be",
						},
					},
					MicroGTUPerEuro: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					MintDistribution: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					ParamGASRewards: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					PoolParameters: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					Protocol: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					TimeParameters: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
					TransactionFeeDistribution: &Level2Key{
						AuthorizedKeys: []int{
							0,
						},
						Threshold: 7,
					},
				},
				RootKeys: &Level1Keys{
					Keys: []*Level1Key{
						{
							SchemeId:  "Ed25519",
							VerifyKey: "f4c9fb9da8d2b00cb6e1fd241b6271a0c4afcb3784e7e6c323bda7ed0b80a478",
						},
					},
					Threshold: 7,
				},
			},
			UpdateQueues: &UpdateQueues{
				AddAnonymityRevoker: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				AddIdentityProvider: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				CooldownParameters: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				ElectionDifficulty: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				EuroPerEnergy: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				FoundationAccount: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				GasRewards: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				Level1Keys: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				Level2Keys: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				MicroGTUPerEuro: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				MintDistribution: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				PoolParameters: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				Protocol: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				RootKeys: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				TimeParameters: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
				TransactionFeeDistribution: &UpdateQueue{
					NextSequenceNumber: 1,
					Queue:              []any{},
				},
			},
		},
	}
)

func TestBlockInfo_MarshalJSON(t *testing.T) {
	testFileMarshalJSON(t, testBlockInfo, testdataBlockInfo)
}

func TestBlockInfo_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &BlockInfo{}, testBlockInfo, testdataBlockInfo)
}

func TestBlockSummary_MarshalJSON(t *testing.T) {
	testFileMarshalJSON(t, testBlockSummary, testdataBlockSummary)
}

func TestBlockSummary_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &BlockSummary{}, testBlockSummary, testdataBlockSummary)
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
