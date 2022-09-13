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
		BlockBaker:            pointer(uint64(5)),
		Finalized:             true,
		TransactionCount:      0,
		TransactionEnergyCost: 0,
		TransactionsSize:      0,
		BlockStateHash:        MustNewBlockHash("e9730d53eed8dafda69c0b0440faabfd4a89e536553a151c657b62c3af827545"),
	}

	testBlockSummary = &BlockSummary{
		FinalizationData: &FinalizationSummary{
			FinalizationBlockPointer: MustNewBlockHash("c7077bf80cefef660eb5ae99ad3fa9137c1f9cc7ed3ea019a8058d3e5c81aa07"),
			FinalizationDelay:        0,
			FinalizationIndex:        530252,
			Finalizers: []*FinalizationSummaryParty{
				{
					BakerId: 0,
					Signed:  true,
					Weight:  7119286292406144,
				},
			},
		},
		ProtocolVersion: 4,
		SpecialEvents: []*SpecialTransactionOutcome{
			{
				BakerId:          5,
				BakerReward:      NewAmountFromMicroCCD(758296),
				FoundationCharge: NewAmountFromMicroCCD(167207),
				NewGASAccount:    NewAmountFromMicroCCD(769692),
				OldGASAccount:    NewAmountFromMicroCCD(23138),
				PassiveReward:    NewAmountFromMicroCCD(10),
				Tag:              "BlockAccrueReward",
				TransactionFees:  NewAmountFromMicroCCD(1672067),
			},
		},
		TransactionSummaries: []*BlockItemSummary{
			{
				Cost:       NewAmountFromMicroCCD(1672067),
				EnergyCost: 999,
				Hash:       "71eda267f1e717fad5867c64eb2e362dc6467abbb01807528e91edcb3dc65e41",
				Index:      0,
				Result: &BlockItemResult{
					Events: Events{
						{
							Tag: "ContractInitialized",
							ContractInitialized: &EventContractInitialized{
								Address: &ContractAddress{
									Index:    888,
									SubIndex: 0,
								},
								Amount:          NewAmountFromMicroCCD(0),
								ContractVersion: 1,
								Events:          []Model{},
								InitName:        "init_a",
								Ref:             MustNewModuleRef("935d17711a4dea10ba5a851df4f19cfdd7cdbd79c8d6ec9abfe5cacff873f6d0"),
							},
						},
					},
					Outcome: "success",
				},
				Sender: pointer(MustNewAccountAddress("4hvvPeHb9HY4Lur7eUZv4KfL3tYBug8DRc4X9cVU8mpJLa1f2X")),
				Type: &BlockItemType{
					Type:               "accountTransaction",
					AccountTransaction: "initContract",
				},
			},
		},
		Updates: &UpdateState{
			ChainParameters: &ChainParameters{
				AccountCreationLimit: 10,
				BakingCommissionRange: &InclusiveRange{
					Max: 0.1,
					Min: 0.1,
				},
				CapitalBound:       0.1,
				DelegatorCooldown:  1209600,
				ElectionDifficulty: 2.5e-2,
				EuroPerEnergy: &ExchangeRate{
					Denominator: 50000,
					Numerator:   1,
				},
				FinalizationCommissionRange: &InclusiveRange{
					Max: 1.0,
					Min: 1.0,
				},
				FoundationAccountIndex: 11,
				LeverageBound: &LeverageFactor{
					Denominator: 1,
					Numerator:   3,
				},
				MicroGTUPerEuro: &ExchangeRate{
					Denominator: 103736215559,
					Numerator:   8681372131440148480,
				},
				MinimumEquityCapital:          NewAmountFromMicroCCD(14000000000),
				MintPerPayday:                 2.61157877e-4,
				PassiveBakingCommission:       0.12,
				PassiveFinalizationCommission: 1.0,
				PassiveTransactionCommission:  0.12,
				PoolOwnerCooldown:             1814400,
				RewardParameters: &RewardParameters{
					GASRewards: &GASRewards{
						AccountCreation:   2.0e-2,
						Baker:             0.25,
						ChainUpdate:       5.0e-3,
						FinalizationProof: 5.0e-3,
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
				TransactionCommissionRange: &InclusiveRange{
					Max: 0.1,
					Min: 0.1,
				},
			},
			Keys: &UpdateKeys{
				Level1Keys: &HigherLevelKeys{
					Keys: []*VerifyKey{
						{
							SchemeId:  "Ed25519",
							VerifyKey: "55721372d942742db382c3680737a424fe3234cf82a80d42da008d6b47179500",
						},
					},
					Threshold: 7,
				},
				Level2Keys: &Authorizations{
					Keys: []*VerifyKey{
						{
							SchemeId:  "Ed25519",
							VerifyKey: "b8ddf4505a37eee2c046671f634b74cf3630f3958ad70a04f39dc843041965be",
						},
					},
					AddAnonymityRevoker: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					AddIdentityProvider: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					CooldownParameters: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					ElectionDifficulty: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					Emergency: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					EuroPerEnergy: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					FoundationAccount: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					MicroGTUPerEuro: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					MintDistribution: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					ParamGASRewards: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					PoolParameters: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					Protocol: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
					TransactionFeeDistribution: &AccessStructure{
						AuthorizedKeys: []uint16{
							0,
						},
						Threshold: 7,
					},
				},
				RootKeys: &HigherLevelKeys{
					Keys: []*VerifyKey{
						{
							SchemeId:  "Ed25519",
							VerifyKey: "f4c9fb9da8d2b00cb6e1fd241b6271a0c4afcb3784e7e6c323bda7ed0b80a478",
						},
					},
					Threshold: 5,
				},
			},
			UpdateQueues: &PendingUpdates{
				AddAnonymityRevoker: &UpdateQueue[*AnonymityRevoker]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*AnonymityRevoker]{},
				},
				AddIdentityProvider: &UpdateQueue[*IdentityProvider]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*IdentityProvider]{},
				},
				CooldownParameters: &UpdateQueue[*CooldownParameters]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*CooldownParameters]{},
				},
				ElectionDifficulty: &UpdateQueue[float64]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[float64]{},
				},
				EuroPerEnergy: &UpdateQueue[*ExchangeRate]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*ExchangeRate]{},
				},
				FoundationAccount: &UpdateQueue[uint64]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[uint64]{},
				},
				GasRewards: &UpdateQueue[*GASRewards]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*GASRewards]{},
				},
				Level1Keys: &UpdateQueue[*HigherLevelKeys]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*HigherLevelKeys]{},
				},
				Level2Keys: &UpdateQueue[*Authorizations]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*Authorizations]{},
				},
				MicroGTUPerEuro: &UpdateQueue[*ExchangeRate]{
					NextSequenceNumber: 3354,
					Queue:              []*ScheduledUpdate[*ExchangeRate]{},
				},
				MintDistribution: &UpdateQueue[*MintDistribution]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*MintDistribution]{},
				},
				Protocol: &UpdateQueue[*ProtocolUpdate]{
					NextSequenceNumber: 2,
					Queue:              []*ScheduledUpdate[*ProtocolUpdate]{},
				},
				RootKeys: &UpdateQueue[*HigherLevelKeys]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*HigherLevelKeys]{},
				},
				TransactionFeeDistribution: &UpdateQueue[*TransactionFeeDistribution]{
					NextSequenceNumber: 1,
					Queue:              []*ScheduledUpdate[*TransactionFeeDistribution]{},
				},
			},
		},
	}
)

func TestBlockInfo_UnmarshalJSON(t *testing.T) {
	testFileUnmarshalJSON(t, &BlockInfo{}, testBlockInfo, testdataBlockInfo)
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
