package concordium

import (
	"encoding/hex"
	"fmt"
	"time"
)

const blockHashSize = 32

// BlockHash base-16 encoded hash of a block (64 characters)
type BlockHash [blockHashSize]byte

func NewBlockHashFromString(s string) (BlockHash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return BlockHash{}, fmt.Errorf("hex decode: %w", err)
	}
	if len(b) != blockHashSize {
		return BlockHash{}, fmt.Errorf("expect %d bytes but %d given", blockHashSize, len(b))
	}
	var h BlockHash
	copy(h[:], b)
	return h, nil
}

func MustNewBlockHashFromString(s string) BlockHash {
	h, err := NewBlockHashFromString(s)
	if err != nil {
		panic("MustNewBlockHashFromString: " + err.Error())
	}
	return h
}

func (h *BlockHash) String() string {
	return hex.EncodeToString((*h)[:])
}

func (h BlockHash) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(h[:])
	if err != nil {
		return nil, fmt.Errorf("%T: %w", h, err)
	}
	return b, nil
}

func (h *BlockHash) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *h, err)
	}
	var x BlockHash
	copy(x[:], v)
	*h = x
	return nil
}

type BlockHeight uint64

// BlockInfo contains metadata about a given block.
type BlockInfo struct {
	// Time when the block was added to the node's tree. This is a subjective (i.e., node specific) value.
	BlockArriveTime time.Time `json:"blockArriveTime"`
	// Identity of the baker of the block. For non-genesis blocks the value is going to always be `Some`.
	BlockBaker *uint64 `json:"blockBaker"`
	// Hash of the block.
	BlockHash BlockHash `json:"blockHash"`
	// Height of the block from genesis.
	BlockHeight BlockHeight `json:"blockHeight"`
	// Pointer to the last finalized block. Each block has a pointer to a specific finalized block
	// that existed at the time the block was produced.
	BlockLastFinalized BlockHash `json:"blockLastFinalized"`
	// Parent block pointer
	BlockParent BlockHash `json:"blockParent"`
	// Time when the block was first received by the node. This can be in principle quite different from
	// the arrive time if, e.g., block execution takes a long time, or the block must wait for the
	// arrival of its parent.
	BlockReceiveTime time.Time `json:"blockReceiveTime"`
	// Slot number of the slot the block is in.
	BlockSlot uint64 `json:"blockSlot"`
	// Slot time of the slot the block is in. In contrast to [BlockInfo::block_arrive_time] this is an
	// objective value, all nodes agree on it.
	BlockSlotTime time.Time `json:"blockSlotTime"`
	// Hash of the block state at the end of the given block.
	BlockStateHash BlockHash `json:"blockStateHash"`
	// The height of this block relative to the (re)genesis block of its era.
	EraBlockHeight BlockHeight `json:"eraBlockHeight"`
	// Whether the block is finalized or not.
	Finalized bool `json:"finalized"`
	// The genesis index for this block. This counts the number of protocol updates that have
	// preceded this block, and defines the era of the block.
	GenesisIndex uint64 `json:"genesisIndex"`
	// The number of transactions in the block.
	TransactionCount uint64 `json:"transactionCount"`
	// The total energy consumption of transactions in the block.
	TransactionEnergyCost uint64 `json:"transactionEnergyCost"`
	// Size of all the transactions in the block in bytes.
	TransactionsSize uint64 `json:"transactionsSize"`
}

type BlockSummary struct {
	FinalizationData     *FinalizationData     `json:"finalizationData"`
	ProtocolVersion      int                   `json:"protocolVersion"`
	SpecialEvents        []*SpecialEvent       `json:"specialEvents"`
	TransactionSummaries []*TransactionOutcome `json:"transactionSummaries"`
	Updates              *Updates              `json:"updates"`
}

type FinalizationData struct {
	FinalizationBlockPointer BlockHash    `json:"finalizationBlockPointer"`
	FinalizationDelay        int          `json:"finalizationDelay"`
	FinalizationIndex        int          `json:"finalizationIndex"`
	Finalizers               []*Finalizer `json:"finalizers"`
}

type Finalizer struct {
	BakerId BakerId `json:"bakerId"`
	Signed  bool    `json:"signed"`
	Weight  int64   `json:"weight"`
}

type SpecialEvent struct {
	BakerId          BakerId `json:"bakerId"`
	BakerReward      string  `json:"bakerReward"`
	FoundationCharge string  `json:"foundationCharge"`
	NewGASAccount    string  `json:"newGASAccount"`
	OldGASAccount    string  `json:"oldGASAccount"`
	PassiveReward    string  `json:"passiveReward"`
	Tag              string  `json:"tag"`
	TransactionFees  string  `json:"transactionFees"`
}

type Updates struct {
	ChainParameters *ChainParameters `json:"chainParameters"`
	Keys            *UpdateKeys      `json:"keys"`
	UpdateQueues    *UpdateQueues    `json:"updateQueues"`
}

type ChainParameters struct {
	AccountCreationLimit          int               `json:"accountCreationLimit"`
	BakingCommissionRange         *CommissionRange  `json:"bakingCommissionRange"`
	CapitalBound                  float64           `json:"capitalBound"`
	DelegatorCooldown             int               `json:"delegatorCooldown"`
	ElectionDifficulty            float64           `json:"electionDifficulty"`
	EuroPerEnergy                 *Fraction         `json:"euroPerEnergy"`
	FinalizationCommissionRange   *CommissionRange  `json:"finalizationCommissionRange"`
	FoundationAccountIndex        int               `json:"foundationAccountIndex"`
	LeverageBound                 *Fraction         `json:"leverageBound"`
	MicroGTUPerEuro               *Fraction         `json:"microGTUPerEuro"`
	MinimumEquityCapital          string            `json:"minimumEquityCapital"`
	MintPerPayday                 float64           `json:"mintPerPayday"`
	PassiveBakingCommission       float64           `json:"passiveBakingCommission"`
	PassiveFinalizationCommission float64           `json:"passiveFinalizationCommission"`
	PassiveTransactionCommission  float64           `json:"passiveTransactionCommission"`
	PoolOwnerCooldown             int               `json:"poolOwnerCooldown"`
	RewardParameters              *RewardParameters `json:"rewardParameters"`
	RewardPeriodLength            int               `json:"rewardPeriodLength"`
	TransactionCommissionRange    *CommissionRange  `json:"transactionCommissionRange"`
}

type CommissionRange struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}

type Fraction struct {
	Denominator int `json:"denominator"`
	Numerator   int `json:"numerator"`
}

type RewardParameters struct {
	GASRewards                 *GASRewards                 `json:"gASRewards"`
	MintDistribution           *MintDistribution           `json:"mintDistribution"`
	TransactionFeeDistribution *TransactionFeeDistribution `json:"transactionFeeDistribution"`
}

type GASRewards struct {
	AccountCreation   float64 `json:"accountCreation"`
	Baker             float64 `json:"baker"`
	ChainUpdate       float64 `json:"chainUpdate"`
	FinalizationProof float64 `json:"finalizationProof"`
}

type MintDistribution struct {
	BakingReward       float64 `json:"bakingReward"`
	FinalizationReward float64 `json:"finalizationReward"`
}

type TransactionFeeDistribution struct {
	Baker      float64 `json:"baker"`
	GasAccount float64 `json:"gasAccount"`
}

type UpdateKeys struct {
	RootKeys   *Level1Keys `json:"rootKeys"`
	Level1Keys *Level1Keys `json:"level1Keys"`
	Level2Keys *Level2Keys `json:"level2Keys"`
}

type Level1Keys struct {
	Keys      []*Level1Key `json:"keys"`
	Threshold int          `json:"threshold"`
}

type Level1Key struct {
	SchemeId  string `json:"schemeId"`
	VerifyKey string `json:"verifyKey"`
}

type Level2Keys struct {
	Keys                       []*Level1Key `json:"keys"`
	AddAnonymityRevoker        *Level2Key   `json:"addAnonymityRevoker"`
	AddIdentityProvider        *Level2Key   `json:"addIdentityProvider"`
	CooldownParameters         *Level2Key   `json:"cooldownParameters"`
	ElectionDifficulty         *Level2Key   `json:"electionDifficulty"`
	Emergency                  *Level2Key   `json:"emergency"`
	EuroPerEnergy              *Level2Key   `json:"euroPerEnergy"`
	FoundationAccount          *Level2Key   `json:"foundationAccount"`
	MicroGTUPerEuro            *Level2Key   `json:"microGTUPerEuro"`
	MintDistribution           *Level2Key   `json:"mintDistribution"`
	ParamGASRewards            *Level2Key   `json:"paramGASRewards"`
	PoolParameters             *Level2Key   `json:"poolParameters"`
	Protocol                   *Level2Key   `json:"protocol"`
	TimeParameters             *Level2Key   `json:"timeParameters"`
	TransactionFeeDistribution *Level2Key   `json:"transactionFeeDistribution"`
}

type Level2Key struct {
	AuthorizedKeys []int `json:"authorizedKeys"`
	Threshold      int   `json:"threshold"`
}

type UpdateQueues struct {
	AddAnonymityRevoker        *UpdateQueue `json:"addAnonymityRevoker"`
	AddIdentityProvider        *UpdateQueue `json:"addIdentityProvider"`
	CooldownParameters         *UpdateQueue `json:"cooldownParameters"`
	ElectionDifficulty         *UpdateQueue `json:"electionDifficulty"`
	EuroPerEnergy              *UpdateQueue `json:"euroPerEnergy"`
	FoundationAccount          *UpdateQueue `json:"foundationAccount"`
	GasRewards                 *UpdateQueue `json:"gasRewards"`
	Level1Keys                 *UpdateQueue `json:"level1Keys"`
	Level2Keys                 *UpdateQueue `json:"level2Keys"`
	MicroGTUPerEuro            *UpdateQueue `json:"microGTUPerEuro"`
	MintDistribution           *UpdateQueue `json:"mintDistribution"`
	PoolParameters             *UpdateQueue `json:"poolParameters"`
	Protocol                   *UpdateQueue `json:"protocol"`
	RootKeys                   *UpdateQueue `json:"rootKeys"`
	TimeParameters             *UpdateQueue `json:"timeParameters"`
	TransactionFeeDistribution *UpdateQueue `json:"transactionFeeDistribution"`
}

type UpdateQueue struct {
	NextSequenceNumber int   `json:"nextSequenceNumber"`
	Queue              []any `json:"queue"`
}
