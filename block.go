package concordium

import (
	"encoding/hex"
	"fmt"
	"time"
)

const (
	SpecialTransactionOutcomeTagBakingRewards          SpecialTransactionOutcomeTag = "BakingRewards"
	SpecialTransactionOutcomeTagMint                   SpecialTransactionOutcomeTag = "Mint"
	SpecialTransactionOutcomeTagFinalizationRewards    SpecialTransactionOutcomeTag = "FinalizationRewards"
	SpecialTransactionOutcomeTagBlockReward            SpecialTransactionOutcomeTag = "BlockReward"
	SpecialTransactionOutcomeTagPaydayFoundationReward SpecialTransactionOutcomeTag = "PaydayFoundationReward"
	SpecialTransactionOutcomeTagPaydayAccountReward    SpecialTransactionOutcomeTag = "PaydayAccountReward"
	SpecialTransactionOutcomeTagBlockAccrueReward      SpecialTransactionOutcomeTag = "BlockAccrueReward"
	SpecialTransactionOutcomeTagPaydayPoolReward       SpecialTransactionOutcomeTag = "PaydayPoolReward"

	// BlockItemResultOutcomeSuccess means that the intended action was completed. The
	// sender was charged, if applicable. Some events were generated describing the
	// changes that happened on the chain.
	BlockItemResultOutcomeSuccess BlockItemResultOutcome = "success"
	// BlockItemResultOutcomeReject means that the intended action was not completed
	// due to an error. The sender was charged, but no other effect is seen on the chain.
	BlockItemResultOutcomeReject BlockItemResultOutcome = "reject"

	blockHashSize = 32
)

type SpecialTransactionOutcomeTag string

type BlockItemResultOutcome string

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

// BlockSummary is summary of transactions, protocol generated transfers, and chain parameters in
// a given block.
type BlockSummary struct {
	// If the block contains a finalization record this contains its summary. Otherwise [None].
	FinalizationData *FinalizationSummary `json:"finalizationData"`
	// Protocol version at which this block was baked. This is no more than [ProtocolVersion::P3]
	ProtocolVersion uint64 `json:"protocolVersion"`
	// Any special events generated as part of this block. Special events are protocol defined
	// transfers, e.g., rewards, minting.
	SpecialEvents []*SpecialTransactionOutcome `json:"specialEvents"`
	// Outcomes of transactions in this block, ordered as in the block.
	TransactionSummaries []*BlockItemSummary `json:"transactionSummaries"`
	// Chain parameters, and any scheduled updates to chain parameters or the protocol.
	Updates *UpdateState `json:"updates"`
}

// FinalizationSummary is summary of the finalization record in a block, if any.
type FinalizationSummary struct {
	FinalizationBlockPointer BlockHash                   `json:"finalizationBlockPointer"`
	FinalizationDelay        uint64                      `json:"finalizationDelay"`
	FinalizationIndex        uint64                      `json:"finalizationIndex"`
	Finalizers               []*FinalizationSummaryParty `json:"finalizers"`
}

// FinalizationSummaryParty contains details of a party in a finalization.
type FinalizationSummaryParty struct {
	// The identity of the baker.
	BakerId BakerId `json:"bakerId"`
	// Whether the party's signature is present
	Signed bool `json:"signed"`
	// The party's relative weight in the committee
	Weight uint64 `json:"weight"`
}

// SpecialTransactionOutcome in addition to the user initiated transactions the protocol
// generates some events which are deemed "Special outcomes". These are rewards for
// running the consensus and finalization protocols.
//
// Cases:
//
// Reward issued to all the bakers at the end of an epoch for baking blocks in the epoch:
// 	* SpecialTransactionOutcome.Tag - has SpecialTransactionOutcomeTagBakingRewards value
// 	* SpecialTransactionOutcome.BakerRewards
// 	* SpecialTransactionOutcome.Remainder - remaining balance of the baking account. This will be
//		transferred to the next epoch's reward account. It exists since it is not possible to
//		perfectly distribute the accumulated rewards. The reason this is not possible is that
//		amounts are integers.
//
// Distribution of newly minted CCD:
// 	* SpecialTransactionOutcome.Tag - has SpecialTransactionOutcomeTagMint value
// 	* SpecialTransactionOutcome.FoundationAccount - the address of the foundation account that
//		the newly minted CCD goes to.
// 	* SpecialTransactionOutcome.MintBakingReward
// 	* SpecialTransactionOutcome.MintFinalizationReward
// 	* SpecialTransactionOutcome.MintPlatformDevelopmentCharge
//
// Distribution of finalization rewards:
// 	* SpecialTransactionOutcome.Tag has SpecialTransactionOutcomeTagFinalizationRewards value
// 	* SpecialTransactionOutcome.FinalizationRewards - the finalization reward at payday to the account.
// 	* SpecialTransactionOutcome.Remainder - remaining balance of the finalization reward account.
//		It exists since it is not possible to perfectly distribute the accumulated rewards. The
//		reason this is not possible is that amounts are integers.
//
// Reward for including transactions in a block:
// 	* SpecialTransactionOutcome.Tag has SpecialTransactionOutcomeTagBlockReward value
// 	* SpecialTransactionOutcome.Baker
// 	* SpecialTransactionOutcome.BakerReward - the amount of CCD that goes to the baker.
// 	* SpecialTransactionOutcome.FoundationAccount - the account address where the foundation receives the tax.
// 	* SpecialTransactionOutcome.FoundationCharge - the amount of CCD that goes to the foundation.
// 	* SpecialTransactionOutcome.NewGASAccount - new balance of the GAS account.
// 	* SpecialTransactionOutcome.OldGASAccount - previous balance of the GAS account.
// 	* SpecialTransactionOutcome.TransactionFees - total amount of transaction fees in the block.
//
// Payment for the foundation:
// 	* SpecialTransactionOutcome.Tag has SpecialTransactionOutcomeTagPaydayFoundationReward value
// 	* SpecialTransactionOutcome.DevelopmentCharge
// 	* SpecialTransactionOutcome.FoundationAccount - address of the foundation account.
//
// Payment for a particular account. When listed in a block summary, the delegated
// pool of the account is given by the last PaydayPoolReward outcome included before
// this outcome:
// 	* SpecialTransactionOutcome.Tag has SpecialTransactionOutcomeTagPaydayAccountReward value
// 	* SpecialTransactionOutcome.Account
// 	* SpecialTransactionOutcome.BakerReward - the baking reward at payday to the account.
// 	* SpecialTransactionOutcome.FinalizationReward
// 	* SpecialTransactionOutcome.TransactionFees - the transaction fee reward at payday to the account.
//
// Amounts accrued to accounts for each baked block:
// 	* SpecialTransactionOutcome.Tag has SpecialTransactionOutcomeTagBlockAccrueReward value
// 	* SpecialTransactionOutcome.BakerId
// 	* SpecialTransactionOutcome.BakerReward - the amount awarded to the baker.
// 	* SpecialTransactionOutcome.FoundationCharge - the amount awarded to the foundation.
// 	* SpecialTransactionOutcome.NewGASAccount - the new balance of the GAS account.
// 	* SpecialTransactionOutcome.OldGASAccount - the old balance of the GAS account.
// 	* SpecialTransactionOutcome.PassiveReward
// 	* SpecialTransactionOutcome.TransactionFees - the total fees paid for transactions in the block.
//
// Payment distributed to a pool or passive delegators:
// 	* SpecialTransactionOutcome.Tag has SpecialTransactionOutcomeTagPaydayPoolReward value
// 	* SpecialTransactionOutcome.BakerReward - accrued baking rewards for pool.
// 	* SpecialTransactionOutcome.FinalizationReward - accrued finalization rewards for pool.
// 	* SpecialTransactionOutcome.PoolOwner
// 	* SpecialTransactionOutcome.TransactionFees - accrued transaction fees for pool.
type SpecialTransactionOutcome struct {
	BakerRewards      *AccountAmount `json:"bakerRewards"`
	Remainder         *Amount        `json:"remainder"`
	Tag               string         `json:"tag"`
	FoundationAccount AccountAddress `json:"foundationAccount"`
	// The portion of the newly minted CCD that goes to the baking reward account.
	MintBakingReward *Amount `json:"mintBakingReward"`
	// The portion that goes to the finalization reward account.
	MintFinalizationReward *Amount `json:"mintFinalizationReward"`
	// The portion that goes to the foundation, as foundation tax.
	MintPlatformDevelopmentCharge *Amount        `json:"mintPlatformDevelopmentCharge"`
	FinalizationRewards           *AccountAmount `json:"finalizationRewards"`
	// The account address where the baker receives the reward.
	Baker            AccountAddress `json:"baker"`
	BakerReward      *Amount        `json:"bakerReward"`
	FoundationCharge *Amount        `json:"foundationCharge"`
	NewGASAccount    *Amount        `json:"newGASAccount"`
	OldGASAccount    *Amount        `json:"oldGASAccount"`
	TransactionFees  *Amount        `json:"transactionFees"`
	// Amount rewarded.
	DevelopmentCharge *Amount `json:"developmentCharge"`
	// The account that got rewarded.
	Account            AccountAddress `json:"account"`
	FinalizationReward *Amount        `json:"finalizationReward"`
	// The baker of the block, who will receive the award.
	BakerId uint64 `json:"bakerId"`
	// The amount awarded to the passive delegators.
	PassiveReward *Amount `json:"passiveReward"`
	// The pool owner (passive delegators when 'None').
	PoolOwner *uint64 `json:"poolOwner"`
}

type AccountAmount struct {
	Address AccountAddress `json:"address"`
	Amount  *Amount        `json:"amount"`
}

// BlockItemSummary is summary of the outcome of a block item.
type BlockItemSummary struct {
	// The amount of CCD the transaction was charged to the sender.
	Cost *Amount `json:"cost"`
	// The amount of NRG the transaction cost.
	EnergyCost uint64 `json:"energyCost"`
	// Hash of the transaction.
	Hash BlockHash `json:"hash"`
	// Index of the transaction in the block where it is included.
	Index uint64 `json:"index"`
	// What is the outcome of this particular block item.
	Result *BlockItemResult `json:"result"`
	// Sender, if available. The sender is always available for account transactions.
	Sender *AccountAddress `json:"sender"`
	// Which type of block item this is.
	Type *BlockItemType `json:"type"`
}

// BlockItemResult contains outcome of a block item execution.
type BlockItemResult struct {
	Outcome BlockItemResultOutcome `json:"outcome"`
	// Exists if BlockItemResult.Outcome is BlockItemResultOutcomeSuccess
	Events []Event `json:"events"`
	// Exists if BlockItemResult.Outcome is BlockItemResultOutcomeReject
	RejectReason *RejectReason `json:"rejectReason"`
}

// BlockItemType is the type of the block item.
type BlockItemType struct {
	// Account transactions are transactions that are signed by an account.
	// Most transactions are account transactions.

	Contents *TransactionType `json:"contents"`
	Type     string           `json:"type"` // accountTransaction

	// Credential deployments that create accounts are special kinds of transactions.
	// They are not signed by the account in the usual way, and they are not paid for
	// directly by the sender.
	Contents CredentialType `json:"contents"`
	Type     string         `json:"type"` // credentialDeploymentTransaction

	// Chain updates are signed by the governance keys. They affect the core parameters of the chain.
	Contents UpdateType `json:"contents"`
	Type     string     `json:"type"` // updateTransaction
}

// UpdateState is a state of updates. This includes current values of parameters as well as any scheduled updates.
type UpdateState struct {
	// Values of chain parameters.
	ChainParameters *ChainParameters `json:"chainParameters"`
	// Keys allowed to perform updates.
	Keys *UpdateKeys `json:"keys"`
	// Possibly pending protocol update.
	ProtocolUpdate *ProtocolUpdate `json:"protocolUpdate"`
	// Any scheduled updates.
	UpdateQueues *PendingUpdates `json:"updateQueues"`
}

// ChainParameters contains values of chain parameters that can be updated via chain updates.
// Only for V0:
// 	* ChainParameters.BakerCooldownEpochs
// 	* ChainParameters.MinimumThresholdForBaking
// Only for V1:
// 	* ChainParameters.BakingCommissionRange
//	* ChainParameters.CapitalBound
//	* ChainParameters.DelegatorCooldown
//	* ChainParameters.FinalizationCommissionRange
//	* ChainParameters.LeverageBound
//	* ChainParameters.MinimumEquityCapital
//	* ChainParameters.MintPerPayday
//	* ChainParameters.PassiveBakingCommission
//	* ChainParameters.PassiveFinalizationCommission
//	* ChainParameters.PassiveTransactionCommission
//	* ChainParameters.PoolOwnerCooldown
//	* ChainParameters.RewardPeriodLength
//	* ChainParameters.TransactionCommissionRange
type ChainParameters struct {
	// The limit for the number of account creations in a block.
	AccountCreationLimit uint16 `json:"accountCreationLimit"`
	// Election difficulty for consensus lottery.
	ElectionDifficulty float64 `json:"electionDifficulty"`
	// Euro per energy exchange rate.
	EuroPerEnergy *ExchangeRate `json:"euroPerEnergy"`
	// Index of the foundation account.
	FoundationAccountIndex uint64 `json:"foundationAccountIndex"`
	// Micro ccd per euro exchange rate.
	MicroGTUPerEuro *ExchangeRate `json:"microGTUPerEuro"`
	// Current reward parameters.
	RewardParameters *RewardParameters `json:"rewardParameters"`

	// Extra number of epochs before reduction in stake, or baker deregistration is completed.
	BakerCooldownEpochs uint64 `json:"bakerCooldownEpochs"`
	// Minimum threshold for becoming a baker.
	MinimumThresholdForBaking *Amount `json:"minimumThresholdForBaking"`

	// The range of allowed baker commissions.
	BakingCommissionRange *InclusiveRange `json:"bakingCommissionRange"`
	// Maximum fraction of the total staked capital of that a new baker can have.
	CapitalBound float64 `json:"capitalBound"`
	// Number of seconds that a delegator must cooldown when reducing their delegated stake.
	DelegatorCooldown uint64 `json:"delegatorCooldown"`
	// The range of allowed finalization commissions.
	FinalizationCommissionRange *InclusiveRange `json:"finalizationCommissionRange"`
	// Index of the foundation account.
	LeverageBound *LeverageFactor `json:"leverageBound"`
	// Minimum equity capital required for a new baker.
	MinimumEquityCapital *Amount `json:"minimumEquityCapital"`
	MintPerPayday        float64 `json:"mintPerPayday"`
	// Fraction of baking rewards charged by the passive delegation.
	PassiveBakingCommission float64 `json:"passiveBakingCommission"`
	// Fraction of finalization rewards charged by the passive delegation.
	PassiveFinalizationCommission float64 `json:"passiveFinalizationCommission"`
	// Fraction of transaction rewards charged by the L-pool.
	PassiveTransactionCommission float64 `json:"passiveTransactionCommission"`
	// Number of seconds that pool owners must cooldown when reducing their equity capital or closing the pool.
	PoolOwnerCooldown  uint64 `json:"poolOwnerCooldown"`
	RewardPeriodLength uint64 `json:"rewardPeriodLength"`
	// The range of allowed transaction commissions.
	TransactionCommissionRange *InclusiveRange `json:"transactionCommissionRange"`
}

// RewardParameters Values of reward parameters.
type RewardParameters struct {
	GASRewards                 *GASRewards                 `json:"gASRewards"`
	MintDistribution           *MintDistribution           `json:"mintDistribution"`
	TransactionFeeDistribution *TransactionFeeDistribution `json:"transactionFeeDistribution"`
}

// GASRewards is the reward fractions related to the gas account and inclusion of special transactions.
type GASRewards struct {
	// `FeeAccountCreation`: fraction paid for including each account creation transaction in a block.
	AccountCreation float64 `json:"accountCreation"`
	// `BakerPrevTransFrac`: fraction of the previous gas account paid to the baker.
	Baker float64 `json:"baker"`
	// `FeeUpdate`: fraction paid for including an update transaction in a block.
	ChainUpdate float64 `json:"chainUpdate"`
	// `FeeAddFinalisationProof`: fraction paid for including a finalization proof in a block.
	FinalizationProof float64 `json:"finalizationProof"`
}

type MintDistribution struct {
	BakingReward       float64 `json:"bakingReward"`
	FinalizationReward float64 `json:"finalizationReward"`
	// Only for V0
	MintPerSlot float64 `json:"mintPerSlot"`
}

// TransactionFeeDistribution update the transaction fee distribution to the specified value.
type TransactionFeeDistribution struct {
	// The fraction that goes to the baker of the block.
	Baker float64 `json:"baker"`
	// The fraction that goes to the gas account. The remaining fraction will go to the foundation.
	GasAccount float64 `json:"gasAccount"`
}

// ProtocolUpdate is a generic protocol update. This is essentially an announcement
// of the update. The details of the update will be communicated in some off-chain
// way, and bakers will need to update their node software to support the update.
type ProtocolUpdate struct {
	Message                    string `json:"message"`
	SpecificationAuxiliaryData string `json:"specificationAuxiliaryData"` // HexString
	SpecificationHash          string `json:"specificationHash"`          // SHA256Hash
	SpecificationURL           string `json:"specificationURL"`
}

type ExchangeRate struct {
	Denominator uint64 `json:"denominator"`
	Numerator   uint64 `json:"numerator"`
}

type LeverageFactor struct {
	Denominator uint64 `json:"denominator"`
	Numerator   uint64 `json:"numerator"`
}

type InclusiveRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// UpdateKeys is the current collection of keys allowed to do updates.
// Parametrized by the chain parameter version.
type UpdateKeys struct {
	RootKeys   *HigherLevelKeys `json:"rootKeys"`
	Level1Keys *HigherLevelKeys `json:"level1Keys"`
	Level2Keys *Authorizations  `json:"level2Keys"`
}

// HigherLevelKeys is either root, level1, or level 2 access structure. They all have the same
// structure, keys and a threshold. The phantom type parameter is used for added type safety to
// distinguish different access structures in different contexts.
type HigherLevelKeys struct {
	Keys      []*VerifyKey `json:"keys"`
	Threshold uint16       `json:"threshold"`
}

// Authorizations contains Access structures for each of the different possible chain updates,
// together with the context giving all the possible keys.
type Authorizations struct {
	// The list of all keys that are currently authorized to perform updates.
	Keys []*VerifyKey `json:"keys"`
	// Access structure for adding new anonymity revokers.
	AddAnonymityRevoker *AccessStructure `json:"addAnonymityRevoker"`
	// Access structure for adding new identity providers.
	AddIdentityProvider *AccessStructure `json:"addIdentityProvider"`
	// Access structure for updating the election difficulty.
	ElectionDifficulty *AccessStructure `json:"electionDifficulty"`
	// Access structure for emergency updates.
	Emergency *AccessStructure `json:"emergency"`
	// Access structure for updating the euro to energy exchange rate.
	EuroPerEnergy *AccessStructure `json:"euroPerEnergy"`
	// Access structure for updating the foundation account address.
	FoundationAccount *AccessStructure `json:"foundationAccount"`
	// Access structure for updating the microccd per euro exchange rate.
	MicroGTUPerEuro *AccessStructure `json:"microGTUPerEuro"`
	// Access structure for updating the mint distribution parameters.
	MintDistribution *AccessStructure `json:"mintDistribution"`
	// Access structure for updating the gas reward distribution parameters.
	ParamGASRewards *AccessStructure `json:"paramGASRewards"`
	// Access structure for updating the pool parameters. For V0 this is only the baker stake threshold, for V1 there are more.
	PoolParameters *AccessStructure `json:"poolParameters"`
	// Access structure for protocol updates.
	Protocol *AccessStructure `json:"protocol"`
	// Access structure for updating the transaction fee distribution."
	TransactionFeeDistribution *AccessStructure `json:"transactionFeeDistribution"`

	// Only for V1. Keys for changing cooldown periods related to baking and delegating.
	CooldownParameters *AccessStructure `json:"cooldownParameters"`
}

// AccessStructure is only meaningful in the context of a list of update keys to which the indices refer to.
type AccessStructure struct {
	AuthorizedKeys []uint16 `json:"authorizedKeys"`
	Threshold      uint16   `json:"threshold"`
}

type PendingUpdates struct {
	AddAnonymityRevoker        *UpdateQueue[*AnonymityRevoker]           `json:"addAnonymityRevoker"`
	AddIdentityProvider        *UpdateQueue[*IdentityProvider]           `json:"addIdentityProvider"`
	BakerStakeThreshold        *UpdateQueue[*BakerParameters]            `json:"bakerStakeThreshold"`
	ElectionDifficulty         *UpdateQueue[float64]                     `json:"electionDifficulty"`
	EuroPerEnergy              *UpdateQueue[*ExchangeRate]               `json:"euroPerEnergy"`
	FoundationAccount          *UpdateQueue[uint64]                      `json:"foundationAccount"`
	GasRewards                 *UpdateQueue[*GASRewards]                 `json:"gasRewards"`
	Level1Keys                 *UpdateQueue[*HigherLevelKeys]            `json:"level1Keys"`
	Level2Keys                 *UpdateQueue[*Authorizations]             `json:"level2Keys"`
	MicroGTUPerEuro            *UpdateQueue[*ExchangeRate]               `json:"microGTUPerEuro"`
	MintDistribution           *UpdateQueue[*MintDistribution]           `json:"mintDistribution"`
	Protocol                   *UpdateQueue[*ProtocolUpdate]             `json:"protocol"`
	RootKeys                   *UpdateQueue[*HigherLevelKeys]            `json:"rootKeys"`
	TransactionFeeDistribution *UpdateQueue[*TransactionFeeDistribution] `json:"transactionFeeDistribution"`

	// Only for V1.
	CooldownParameters *UpdateQueue[*CooldownParameters] `json:"cooldownParameters"`
}

// UpdateQueue is a queue of updates of a given type.
type UpdateQueue[T any] struct {
	// Next available sequence number for the update type.
	NextSequenceNumber uint64 `json:"nextSequenceNumber"`
	// Queue of updates, ordered by effective time.
	Queue []*ScheduledUpdate[T] `json:"queue"`
}

// ScheduledUpdate is a scheduled update of a given type.
type ScheduledUpdate[T any] struct {
	EffectiveTime uint64 `json:"effectiveTime"`
	Update        T      `json:"update"`
}

type CooldownParameters struct {
	// Number of seconds that a delegator must cooldown when reducing their delegated stake.
	DelegatorCooldown uint64 `json:"delegatorCooldown"`
	// Number of seconds that pool owners must cooldown when reducing their equity capital or closing the pool.
	PoolOwnerCooldown uint64 `json:"poolOwnerCooldown"`
}
