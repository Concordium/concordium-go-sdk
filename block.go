package concordium

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	blockHashSize = 32
)

// BlockHash base-16 encoded hash of a block (64 characters)
type BlockHash [blockHashSize]byte

// NewBlockHash creates a new BlockHash from string.
func NewBlockHash(s string) (BlockHash, error) {
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

// MustNewBlockHash calls the NewBlockHash. It panics in case of error.
func MustNewBlockHash(s string) BlockHash {
	h, err := NewBlockHash(s)
	if err != nil {
		panic("MustNewBlockHash: " + err.Error())
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

// SpecialTransactionOutcomeTag describes SpecialTransactionOutcome type. See related constants.
type SpecialTransactionOutcomeTag string

const (
	SpecialTransactionOutcomeTagBakingRewards          SpecialTransactionOutcomeTag = "BakingRewards"
	SpecialTransactionOutcomeTagMint                   SpecialTransactionOutcomeTag = "Mint"
	SpecialTransactionOutcomeTagFinalizationRewards    SpecialTransactionOutcomeTag = "FinalizationRewards"
	SpecialTransactionOutcomeTagBlockReward            SpecialTransactionOutcomeTag = "BlockReward"
	SpecialTransactionOutcomeTagPaydayFoundationReward SpecialTransactionOutcomeTag = "PaydayFoundationReward"
	SpecialTransactionOutcomeTagPaydayAccountReward    SpecialTransactionOutcomeTag = "PaydayAccountReward"
	SpecialTransactionOutcomeTagBlockAccrueReward      SpecialTransactionOutcomeTag = "BlockAccrueReward"
	SpecialTransactionOutcomeTagPaydayPoolReward       SpecialTransactionOutcomeTag = "PaydayPoolReward"
)

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
	Tag               SpecialTransactionOutcomeTag `json:"tag"`
	BakerRewards      *AccountAmount               `json:"bakerRewards"`
	Remainder         *Amount                      `json:"remainder"`
	FoundationAccount AccountAddress               `json:"foundationAccount"`
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
	Hash TransactionHash `json:"hash"`
	// Index of the transaction in the block where it is included.
	Index uint64 `json:"index"`
	// What is the outcome of this particular block item.
	Result *BlockItemResult `json:"result"`
	// Sender, if available. The sender is always available for account transactions.
	Sender *AccountAddress `json:"sender"`
	// Which type of block item this is.
	Type *BlockItemType `json:"type"`
}

// BlockItemResultOutcome describes BlockItemResult type. See related constants.
type BlockItemResultOutcome string

const (
	// BlockItemResultOutcomeSuccess means that the intended action was completed. The
	// sender was charged, if applicable. Some events were generated describing the
	// changes that happened on the chain.
	BlockItemResultOutcomeSuccess BlockItemResultOutcome = "success"
	// BlockItemResultOutcomeReject means that the intended action was not completed
	// due to an error. The sender was charged, but no other effect is seen on the chain.
	BlockItemResultOutcomeReject BlockItemResultOutcome = "reject"
)

// BlockItemResult contains outcome of a block item execution.
type BlockItemResult struct {
	Outcome BlockItemResultOutcome `json:"outcome"`
	// Exists if BlockItemResult.Outcome is BlockItemResultOutcomeSuccess
	Events Events `json:"events"`
	// Exists if BlockItemResult.Outcome is BlockItemResultOutcomeReject
	RejectReason *RejectReason `json:"rejectReason"`
}

// EventTag describes Event type. See related constants.
type EventTag string

const (
	EventTagModuleDeployed                       EventTag = "ModuleDeployed"
	EventTagContractInitialized                  EventTag = "ContractInitialized"
	EventTagUpdated                              EventTag = "Updated"
	EventTagTransferred                          EventTag = "Transferred"
	EventTagAccountCreated                       EventTag = "AccountCreated"
	EventTagCredentialDeployed                   EventTag = "CredentialDeployed"
	EventTagBakerAdded                           EventTag = "BakerAdded"
	EventTagBakerRemoved                         EventTag = "BakerRemoved"
	EventTagBakerStakeIncreased                  EventTag = "BakerStakeIncreased"
	EventTagBakerStakeDecreased                  EventTag = "BakerStakeDecreased"
	EventTagBakerSetRestakeEarnings              EventTag = "BakerSetRestakeEarnings"
	EventTagBakerKeysUpdated                     EventTag = "BakerKeysUpdated"
	EventTagCredentialKeysUpdated                EventTag = "CredentialKeysUpdated"
	EventTagNewEncryptedAmount                   EventTag = "NewEncryptedAmount"
	EventTagEncryptedAmountsRemoved              EventTag = "EncryptedAmountsRemoved"
	EventTagAmountAddedByDecryption              EventTag = "AmountAddedByDecryption"
	EventTagEncryptedSelfAmountAdded             EventTag = "EncryptedSelfAmountAdded"
	EventTagUpdateEnqueued                       EventTag = "UpdateEnqueued"
	EventTagTransferredWithSchedule              EventTag = "TransferredWithSchedule"
	EventTagCredentialsUpdated                   EventTag = "CredentialsUpdated"
	EventTagDataRegistered                       EventTag = "DataRegistered"
	EventTagTransferMemo                         EventTag = "TransferMemo"
	EventTagInterrupted                          EventTag = "Interrupted"
	EventTagResumed                              EventTag = "Resumed"
	EventTagBakerSetOpenStatus                   EventTag = "BakerSetOpenStatus"
	EventTagBakerSetMetadataURL                  EventTag = "BakerSetMetadataURL"
	EventTagBakerSetTransactionFeeCommission     EventTag = "BakerSetTransactionFeeCommission"
	EventTagBakerSetBakingRewardCommission       EventTag = "BakerSetBakingRewardCommission"
	EventTagBakerSetFinalizationRewardCommission EventTag = "BakerSetFinalizationRewardCommission"
	EventTagDelegationStakeIncreased             EventTag = "DelegationStakeIncreased"
	EventTagDelegationStakeDecreased             EventTag = "DelegationStakeDecreased"
	EventTagDelegationSetRestakeEarnings         EventTag = "DelegationSetRestakeEarnings"
	EventTagDelegationSetDelegationTarget        EventTag = "DelegationSetDelegationTarget"
	EventTagDelegationAdded                      EventTag = "DelegationAdded"
	EventTagDelegationRemoved                    EventTag = "DelegationRemoved"
)

type Events []*Event

// Event describing the changes that occurred to the state of the chain.
type Event struct {
	Tag EventTag `json:"tag"`
	// A smart contract module was successfully deployed.
	ModuleDeployed *EventModuleDeployed `json:"-"`
	// A new smart contract instance was created.
	ContractInitialized *EventContractInitialized `json:"-"`
	// A smart contract instance was updated.
	Updated *EventUpdated `json:"-"`
	// An amount of CCD was transferred.
	Transferred *EventTransferred `json:"-"`
	// An account with the given address was created.
	AccountCreated *EventAccountCreated `json:"-"`
	// A new credential with the given ID was deployed onto an account.
	// This is used only when a new account is created. See [Event::CredentialsUpdated]
	// for when an existing account's credentials are updated.
	CredentialDeployed *EventCredentialDeployed `json:"-"`
	// A new baker was registered, with the given ID and keys.
	BakerAdded *EventBakerAdded `json:"-"`
	// A baker was scheduled to be removed.
	BakerRemoved *EventBakerRemoved `json:"-"`
	// A baker's stake was increased. This has effect immediately.
	BakerStakeIncreased *EventBakerStakeIncreased `json:"-"`
	// A baker's stake was scheduled to be decreased. This will have an effect on the stake
	// after a number of epochs, controlled by the baker cooldown period.
	BakerStakeDecreased *EventBakerStakeDecreased `json:"-"`
	// The setting for whether rewards are added to stake immediately or not was changed to the given value.
	BakerSetRestakeEarnings *EventBakerSetRestakeEarnings `json:"-"`
	// The baker keys were updated. The new keys are listed.
	BakerKeysUpdated *EventBakerKeysUpdated `json:"-"`
	// Keys of the given credential were updated.
	CredentialKeysUpdated *EventCredentialKeysUpdated `json:"-"`
	// A new encrypted amount was added to the account.
	NewEncryptedAmount *EventNewEncryptedAmount `json:"-"`
	// One or more encrypted amounts were removed from an account as part of a transfer or decryption.
	EncryptedAmountsRemoved *EventEncryptedAmountsRemoved `json:"-"`
	// The public balance of the account was increased via a transfer from encrypted to public balance.
	AmountAddedByDecryption *EventAmountAddedByDecryption `json:"-"`
	// The encrypted balance of the account was updated due to transfer from public to encrypted balance of the account.
	EncryptedSelfAmountAdded *EventEncryptedSelfAmountAdded `json:"-"`
	// An update was enqueued for the given time.
	UpdateEnqueued *EventUpdateEnqueued `json:"-"`
	// A transfer with schedule was enqueued.
	TransferredWithSchedule *EventTransferredWithSchedule `json:"-"`
	// The credentials of the account were updated. Either added, removed, or both.
	CredentialsUpdated *EventCredentialsUpdated `json:"-"`
	// Data was registered.
	DataRegistered *EventDataRegistered `json:"-"`
	// Memo
	TransferMemo *EventTransferMemo `json:"-"`
	// A V1 contract was interrupted.
	Interrupted *EventInterrupted `json:"-"`
	// A V1 contract resumed execution.
	Resumed *EventResumed `json:"-"`
	// Updated open status for a baker pool
	BakerSetOpenStatus *EventBakerSetOpenStatus `json:"-"`
	// Updated metadata url for baker pool
	BakerSetMetadataURL *EventBakerSetMetadataURL `json:"-"`
	// Updated transaction fee commission for baker pool
	BakerSetTransactionFeeCommission *EventBakerSetTransactionFeeCommission `json:"-"`
	// Updated baking reward commission for baker pool
	BakerSetBakingRewardCommission       *EventBakerSetBakingRewardCommission       `json:"-"`
	BakerSetFinalizationRewardCommission *EventBakerSetFinalizationRewardCommission `json:"-"`
	DelegationStakeIncreased             *EventDelegationStakeIncreased             `json:"-"`
	DelegationStakeDecreased             *EventDelegationStakeDecreased             `json:"-"`
	DelegationSetRestakeEarnings         *EventDelegationSetRestakeEarnings         `json:"-"`
	DelegationSetDelegationTarget        *EventDelegationSetDelegationTarget        `json:"-"`
	DelegationAdded                      *EventDelegationAdded                      `json:"-"`
	DelegationRemoved                    *EventDelegationRemoved                    `json:"-"`
}

func (e *Event) UnmarshalJSON(b []byte) error {
	v := &struct {
		Tag EventTag `json:"tag"`
	}{}
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	e.Tag = v.Tag
	var w any
	switch e.Tag {
	case EventTagModuleDeployed:
		e.ModuleDeployed = &EventModuleDeployed{}
		w = e.ModuleDeployed
	case EventTagContractInitialized:
		e.ContractInitialized = &EventContractInitialized{}
		w = e.ContractInitialized
	case EventTagUpdated:
		e.Updated = &EventUpdated{}
		w = e.Updated
	case EventTagTransferred:
		e.Transferred = &EventTransferred{}
		w = e.Transferred
	case EventTagAccountCreated:
		e.AccountCreated = &EventAccountCreated{}
		w = e.AccountCreated
	case EventTagCredentialDeployed:
		e.CredentialDeployed = &EventCredentialDeployed{}
		w = e.CredentialDeployed
	case EventTagBakerAdded:
		e.BakerAdded = &EventBakerAdded{}
		w = e.BakerAdded
	case EventTagBakerRemoved:
		e.BakerRemoved = &EventBakerRemoved{}
		w = e.BakerRemoved
	case EventTagBakerStakeIncreased:
		e.BakerStakeIncreased = &EventBakerStakeIncreased{}
		w = e.BakerStakeIncreased
	case EventTagBakerStakeDecreased:
		e.BakerStakeDecreased = &EventBakerStakeDecreased{}
		w = e.BakerStakeDecreased
	case EventTagBakerSetRestakeEarnings:
		e.BakerSetRestakeEarnings = &EventBakerSetRestakeEarnings{}
		w = e.BakerSetRestakeEarnings
	case EventTagBakerKeysUpdated:
		e.BakerKeysUpdated = &EventBakerKeysUpdated{}
		w = e.BakerKeysUpdated
	case EventTagCredentialKeysUpdated:
		e.CredentialKeysUpdated = &EventCredentialKeysUpdated{}
		w = e.CredentialKeysUpdated
	case EventTagNewEncryptedAmount:
		e.NewEncryptedAmount = &EventNewEncryptedAmount{}
		w = e.NewEncryptedAmount
	case EventTagEncryptedAmountsRemoved:
		e.EncryptedAmountsRemoved = &EventEncryptedAmountsRemoved{}
		w = e.EncryptedAmountsRemoved
	case EventTagAmountAddedByDecryption:
		e.AmountAddedByDecryption = &EventAmountAddedByDecryption{}
		w = e.AmountAddedByDecryption
	case EventTagEncryptedSelfAmountAdded:
		e.EncryptedSelfAmountAdded = &EventEncryptedSelfAmountAdded{}
		w = e.EncryptedSelfAmountAdded
	case EventTagUpdateEnqueued:
		e.UpdateEnqueued = &EventUpdateEnqueued{}
		w = e.UpdateEnqueued
	case EventTagTransferredWithSchedule:
		e.TransferredWithSchedule = &EventTransferredWithSchedule{}
		w = e.TransferredWithSchedule
	case EventTagCredentialsUpdated:
		e.CredentialsUpdated = &EventCredentialsUpdated{}
		w = e.CredentialsUpdated
	case EventTagDataRegistered:
		e.DataRegistered = &EventDataRegistered{}
		w = e.DataRegistered
	case EventTagTransferMemo:
		e.TransferMemo = &EventTransferMemo{}
		w = e.TransferMemo
	case EventTagInterrupted:
		e.Interrupted = &EventInterrupted{}
		w = e.Interrupted
	case EventTagResumed:
		e.Resumed = &EventResumed{}
		w = e.Resumed
	case EventTagBakerSetOpenStatus:
		e.BakerSetOpenStatus = &EventBakerSetOpenStatus{}
		w = e.BakerSetOpenStatus
	case EventTagBakerSetMetadataURL:
		e.BakerSetMetadataURL = &EventBakerSetMetadataURL{}
		w = e.BakerSetMetadataURL
	case EventTagBakerSetTransactionFeeCommission:
		e.BakerSetTransactionFeeCommission = &EventBakerSetTransactionFeeCommission{}
		w = e.BakerSetTransactionFeeCommission
	case EventTagBakerSetBakingRewardCommission:
		e.BakerSetBakingRewardCommission = &EventBakerSetBakingRewardCommission{}
		w = e.BakerSetBakingRewardCommission
	case EventTagBakerSetFinalizationRewardCommission:
		e.BakerSetFinalizationRewardCommission = &EventBakerSetFinalizationRewardCommission{}
		w = e.BakerSetFinalizationRewardCommission
	case EventTagDelegationStakeIncreased:
		e.DelegationStakeIncreased = &EventDelegationStakeIncreased{}
		w = e.DelegationStakeIncreased
	case EventTagDelegationStakeDecreased:
		e.DelegationStakeDecreased = &EventDelegationStakeDecreased{}
		w = e.DelegationStakeDecreased
	case EventTagDelegationSetRestakeEarnings:
		e.DelegationSetRestakeEarnings = &EventDelegationSetRestakeEarnings{}
		w = e.DelegationSetRestakeEarnings
	case EventTagDelegationSetDelegationTarget:
		e.DelegationSetDelegationTarget = &EventDelegationSetDelegationTarget{}
		w = e.DelegationSetDelegationTarget
	case EventTagDelegationAdded:
		e.DelegationAdded = &EventDelegationAdded{}
		w = e.DelegationAdded
	case EventTagDelegationRemoved:
		e.DelegationRemoved = &EventDelegationRemoved{}
		w = e.DelegationRemoved
	default:
		return fmt.Errorf("%T: unexpected tag: %s", *e, e.Tag)
	}
	return json.Unmarshal(b, w)
}

type EventModuleDeployed struct {
	Contents ModuleRef `json:"contents"`
}

type EventContractInitialized struct {
	Address *ContractAddress `json:"address"`
	// The amount the instance was initialized with.
	Amount          *Amount `json:"amount"`
	ContractVersion uint8   `json:"contractVersion"`
	// Any contract events that might have been generated by the contract initialization.
	Events []Model `json:"events"`
	// The name of the contract.
	InitName string `json:"initName"`
	// Module with the source code of the contract.
	Ref ModuleRef `json:"ref"`
}

type EventUpdated struct {
	// Address of the affected instance.
	Address *ContractAddress `json:"address"`
	// The amount the method was invoked with.
	Amount          *Amount `json:"amount"`
	ContractVersion uint8   `json:"contractVersion"`
	// Any contract events that might have been generated by the contract execution.
	Events []Model `json:"events"`
	// The origin of the message to the smart contract. This can be either an account or a smart contract.
	Instigator *Address `json:"instigator"`
	// The message passed to method.
	Message Model `json:"message"`
	// The name of the method that was executed.
	ReceiveName string `json:"receiveName"`
}

type EventTransferred struct {
	// Amount that was transferred.
	Amount *Amount `json:"amount"`
	// Sender, either smart contract instance or account.
	From *Address `json:"from"`
	// Receiver. This will currently always be an account. Transferring to a smart contract is always an update.
	To *Address `json:"to"`
}

type EventAccountCreated struct {
	Contents AccountAddress `json:"contents"`
}

type EventCredentialDeployed struct {
	Account AccountAddress `json:"account"`
	RegId   string         `json:"regId"`
}

type EventBakerAdded struct {
	// Account address of the baker.
	Account AccountAddress `json:"account"`
	// The new public key for verifying finalization records.
	AggregationKey string `json:"aggregationKey"`
	// ID of the baker whose keys were changed.
	BakerId BakerId `json:"bakerId"`
	// The new public key for verifying whether the baker won the block lottery.
	ElectionKey string `json:"electionKey"`
	// Whether the baker will automatically add earnings to their stake or not.
	RestakeEarnings bool `json:"restakeEarnings"`
	// The new public key for verifying block signatures.
	SignKey string `json:"signKey"`
	// The amount the account staked to become a baker. This amount is locked.
	Stake *Amount `json:"stake"`
}

type EventBakerRemoved struct {
	Account AccountAddress `json:"account"`
	BakerId BakerId        `json:"bakerId"`
}

type EventBakerStakeIncreased struct {
	Account  AccountAddress `json:"account"`
	BakerId  BakerId        `json:"bakerId"`
	NewStake *Amount        `json:"newStake"`
}

type EventBakerStakeDecreased struct {
	Account  AccountAddress `json:"account"`
	BakerId  BakerId        `json:"bakerId"`
	NewStake *Amount        `json:"newStake"`
}

type EventBakerSetRestakeEarnings struct {
	Account AccountAddress `json:"account"`
	BakerId BakerId        `json:"bakerId"`
	// description": "The new value of the flag.
	RestakeEarnings bool `json:"restakeEarnings"`
}

type EventBakerKeysUpdated struct {
	// Account address of the baker.
	Account AccountAddress `json:"account"`
	// The new public key for verifying finalization records.
	AggregationKey string `json:"aggregationKey"`
	// ID of the baker whose keys were changed.
	BakerId BakerId `json:"bakerId"`
	// The new public key for verifying whether the baker won the block lottery.
	ElectionKey string `json:"electionKey"`
	// The new public key for verifying block signatures.
	SignKey string `json:"signKey"`
}

type EventCredentialKeysUpdated struct {
	CredId string `json:"credId"`
}

type EventNewEncryptedAmount struct {
	// The account onto which the amount was added.
	Account AccountAddress `json:"account"`
	// The encrypted amount that was added.
	EncryptedAmount EncryptedAmount `json:"encryptedAmount"`
	// The index the amount was assigned.
	NewIndex uint64 `json:"newIndex"`
}

type EventEncryptedAmountsRemoved struct {
	// The affected account.
	Account AccountAddress `json:"account"`
	// The input encrypted amount that was removed.
	InputAmount EncryptedAmount `json:"inputAmount"`
	// The new self encrypted amount on the affected account.
	NewAmount EncryptedAmount `json:"newAmount"`
	// The index indicating which amounts were used.
	UpToIndex uint64 `json:"upToIndex"`
}

type EventAmountAddedByDecryption struct {
	Account AccountAddress `json:"account"`
	Amount  *Amount        `json:"amount"`
}

type EventEncryptedSelfAmountAdded struct {
	// The affected account.
	Account AccountAddress `json:"account"`
	// The amount that was transferred from public to encrypted balance.
	Amount *Amount `json:"amount"`
	// The new self encrypted amount of the account.
	NewAmount EncryptedAmount `json:"newAmount"`
}

type EventUpdateEnqueued struct {
	EffectiveTime uint64         `json:"effectiveTime"`
	Payload       *UpdatePayload `json:"payload"`
}

type EventTransferredWithSchedule struct {
	// The list of releases. Ordered by increasing timestamp.
	Amount *PairTuple[uint64, *Amount] `json:"amount"`
	// Sender account.
	From AccountAddress `json:"from"`
	// Receiver account.
	To AccountAddress `json:"to"`
}

type EventCredentialsUpdated struct {
	// The affected account.
	Account AccountAddress `json:"account"`
	// The credential ids that were added.
	NewCredIds string `json:"newCredIds"`
	// The (possibly) updated account threshold.
	NewThreshold uint8 `json:"newThreshold"`
	// The credentials that were removed.
	RemovedCredIds string `json:"removedCredIds"`
}

type EventDataRegistered struct {
	Data string `json:"data"`
}

type EventTransferMemo struct {
	Memo string `json:"memo"`
}

type EventInterrupted struct {
	// Address of the contract that was interrupted.
	Address ContractAddress `json:"address"`
	// Events generated up to the interrupt.
	Events []Model `json:"events"`
}

type EventResumed struct {
	// Address of the contract that is resuming.
	Address ContractAddress `json:"address"`
	// Whether the interrupt succeeded or not.
	Success bool `json:"success"`
}

type EventBakerSetOpenStatus struct {
	// Baker account
	Account AccountAddress `json:"account"`
	// Baker's id
	BakerId BakerId `json:"bakerId"`
	// The open status.
	OpenStatus OpenStatus `json:"openStatus"`
}

type EventBakerSetMetadataURL struct {
	// Baker account
	Account AccountAddress `json:"account"`
	// Baker's id
	BakerId BakerId `json:"bakerId"`
	// The URL.
	MetadataURL string `json:"metadataURL"`
}

type EventBakerSetTransactionFeeCommission struct {
	// Baker account
	Account AccountAddress `json:"account"`
	// Baker's id
	BakerId BakerId `json:"bakerId"`
	// The transaction fee commission.
	TransactionFeeCommission int `json:"transactionFeeCommission"`
}

type EventBakerSetBakingRewardCommission struct {
	// Baker account
	Account AccountAddress `json:"account"`
	// Baker's id
	BakerId BakerId `json:"bakerId"`
	// The baking reward commission
	BakingRewardCommission int `json:"bakingRewardCommission"`
}
type EventBakerSetFinalizationRewardCommission struct {
	// Baker account
	Account AccountAddress `json:"account"`
	// Baker's id
	BakerId BakerId `json:"bakerId"`
	// The finalization reward commission
	FinalizationRewardCommission int `json:"finalizationRewardCommission"`
}
type EventDelegationStakeIncreased struct {
	// Delegator account
	Account AccountAddress `json:"account"`
	// Delegator's id
	DelegatorId uint64 `json:"delegatorId"`
	// New stake
	NewStake *Amount `json:"newStake"`
}
type EventDelegationStakeDecreased struct {
	// Delegator account
	Account AccountAddress `json:"account"`
	// Delegator's id
	DelegatorId uint64 `json:"delegatorId"`
	// New stake
	NewStake *Amount `json:"newStake"`
}
type EventDelegationSetRestakeEarnings struct {
	// Delegator account
	Account AccountAddress `json:"account"`
	// Delegator's id
	DelegatorId uint64 `json:"delegatorId"`
	// Whether earnings will be restaked
	RestakeEarnings bool `json:"restakeEarnings"`
}
type EventDelegationSetDelegationTarget struct {
	// Delegator account
	Account AccountAddress `json:"account"`
	// Delegator's id
	DelegatorId uint64 `json:"delegatorId"`
	// New delegation target
	DelegationTarget DelegationTarget `json:"delegationTarget"`
}
type EventDelegationAdded struct {
	// Delegator account
	Account AccountAddress `json:"account"`
	// Delegator's id
	DelegatorId uint64 `json:"delegatorId"`
}
type EventDelegationRemoved struct {
	// Delegator account
	Account AccountAddress `json:"account"`
	// Delegator's id
	DelegatorId uint64 `json:"delegatorId"`
}

// UpdatePayloadType describes UpdatePayload type. See related constants.
type UpdatePayloadType string

const (
	UpdatePayloadTypeProtocol                   UpdatePayloadType = "protocol"
	UpdatePayloadTypeElectionDifficulty         UpdatePayloadType = "electionDifficulty"
	UpdatePayloadTypeEuroPerEnergy              UpdatePayloadType = "euroPerEnergy"
	UpdatePayloadTypeMicroGTUPerEuro            UpdatePayloadType = "microGTUPerEuro"
	UpdatePayloadTypeFoundationAccount          UpdatePayloadType = "foundationAccount"
	UpdatePayloadTypeMintDistribution           UpdatePayloadType = "mintDistribution"
	UpdatePayloadTypeTransactionFeeDistribution UpdatePayloadType = "transactionFeeDistribution"
	UpdatePayloadTypeGASRewards                 UpdatePayloadType = "gASRewards"
	UpdatePayloadTypeBakerStakeThreshold        UpdatePayloadType = "bakerStakeThreshold"
	UpdatePayloadTypeRoot                       UpdatePayloadType = "root"
	UpdatePayloadTypeLevel1                     UpdatePayloadType = "level1"
	UpdatePayloadTypeAddAnonymityRevoker        UpdatePayloadType = "addAnonymityRevoker"
	UpdatePayloadTypeAddIdentityProvider        UpdatePayloadType = "addIdentityProvider"
	UpdatePayloadTypeCooldownParametersCPV1     UpdatePayloadType = "cooldownParametersCPV1"
	UpdatePayloadTypePoolParametersCPV1         UpdatePayloadType = "poolParametersCPV1"
	UpdatePayloadTypeTimeParametersCPV1         UpdatePayloadType = "timeParametersCPV1"
	UpdatePayloadTypeMintDistributionCPV1       UpdatePayloadType = "mintDistributionCPV1"
)

// UpdatePayload is the type of an update payload.
type UpdatePayload struct {
	UpdateType UpdatePayloadType `json:"updateType"`

	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeProtocol
	Protocol *ProtocolUpdate `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeElectionDifficulty
	ElectionDifficulty float64 `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeEuroPerEnergy
	EuroPerEnergy *ExchangeRate `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeMicroGTUPerEuro
	MicroGTUPerEuro *ExchangeRate `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeFoundationAccount
	FoundationAccount AccountAddress `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeMintDistribution
	MintDistribution *MintDistribution `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeTransactionFeeDistribution
	TransactionFeeDistribution *TransactionFeeDistribution `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeGASRewards
	GASRewards *GASRewards `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeBakerStakeThreshold
	BakerStakeThreshold *BakerParameters `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeRoot
	Root *RootUpdate `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeLevel1
	Level1 *Level1Update `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeAddAnonymityRevoker
	AddAnonymityRevoker *AnonymityRevoker `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeAddIdentityProvider
	AddIdentityProvider *IdentityProvider `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeCooldownParametersCPV1
	CooldownParametersCPV1 *CooldownParameters `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypePoolParametersCPV1
	PoolParametersCPV1 *PoolParameters `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeTimeParametersCPV1
	TimeParametersCPV1 *TimeParameters `json:"-"`
	// Makes sense only if UpdatePayload.UpdateType is UpdatePayloadTypeMintDistributionCPV1
	MintDistributionCPV1 *MintDistribution `json:"-"`
}

func (p *UpdatePayload) UnmarshalJSON(b []byte) error {
	v := &struct {
		UpdateType UpdatePayloadType `json:"updateType"`
	}{}
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	p.UpdateType = v.UpdateType
	var w any
	switch p.UpdateType {
	case UpdatePayloadTypeProtocol:
		p.Protocol = &ProtocolUpdate{}
		w = p.Protocol
	case UpdatePayloadTypeElectionDifficulty:
		p.ElectionDifficulty = 0
		w = &p.ElectionDifficulty
	case UpdatePayloadTypeEuroPerEnergy:
		p.EuroPerEnergy = &ExchangeRate{}
		w = p.EuroPerEnergy
	case UpdatePayloadTypeMicroGTUPerEuro:
		p.MicroGTUPerEuro = &ExchangeRate{}
		w = p.MicroGTUPerEuro
	case UpdatePayloadTypeFoundationAccount:
		p.FoundationAccount = AccountAddress{}
		w = &p.FoundationAccount
	case UpdatePayloadTypeMintDistribution:
		p.MintDistribution = &MintDistribution{}
		w = p.MintDistribution
	case UpdatePayloadTypeTransactionFeeDistribution:
		p.TransactionFeeDistribution = &TransactionFeeDistribution{}
		w = p.TransactionFeeDistribution
	case UpdatePayloadTypeGASRewards:
		p.GASRewards = &GASRewards{}
		w = p.GASRewards
	case UpdatePayloadTypeBakerStakeThreshold:
		p.BakerStakeThreshold = &BakerParameters{}
		w = p.BakerStakeThreshold
	case UpdatePayloadTypeRoot:
		p.Root = &RootUpdate{}
		w = p.Root
	case UpdatePayloadTypeLevel1:
		p.Level1 = &Level1Update{}
		w = p.Level1
	case UpdatePayloadTypeAddAnonymityRevoker:
		p.AddAnonymityRevoker = &AnonymityRevoker{}
		w = p.AddAnonymityRevoker
	case UpdatePayloadTypeAddIdentityProvider:
		p.AddIdentityProvider = &IdentityProvider{}
		w = p.AddIdentityProvider
	case UpdatePayloadTypeCooldownParametersCPV1:
		p.CooldownParametersCPV1 = &CooldownParameters{}
		w = p.CooldownParametersCPV1
	case UpdatePayloadTypePoolParametersCPV1:
		p.PoolParametersCPV1 = &PoolParameters{}
		w = p.PoolParametersCPV1
	case UpdatePayloadTypeTimeParametersCPV1:
		p.TimeParametersCPV1 = &TimeParameters{}
		w = p.TimeParametersCPV1
	case UpdatePayloadTypeMintDistributionCPV1:
		p.MintDistributionCPV1 = &MintDistribution{}
		w = p.MintDistributionCPV1
	default:
		return fmt.Errorf("%T: unexpected type: %s", *p, p.UpdateType)
	}
	return json.Unmarshal(b, w)
}

// RootUpdateType describes RootUpdate type. See related constants.
type RootUpdateType string

const (
	RootUpdateTypeRootKeysUpdate     RootUpdateType = "rootKeysUpdate"
	RootUpdateTypeLevel1KeysUpdate   RootUpdateType = "level1KeysUpdate"
	RootUpdateTypeLevel2KeysUpdate   RootUpdateType = "level2KeysUpdate"
	RootUpdateTypeLevel2KeysUpdateV1 RootUpdateType = "level2KeysUpdateV1"
)

// RootUpdate is an update with root keys of some other set of governance keys, or the root
// keys themselves. Each update is a separate transaction.
type RootUpdate struct {
	TypeOfUpdate RootUpdateType `json:"typeOfUpdate"`
	// Makes sense if RootUpdate.TypeOfUpdate is:
	// 	* RootUpdateTypeRootKeysUpdate
	// 	* RootUpdateTypeLevel1KeysUpdate
	HigherLevelKeys *HigherLevelKeys `json:"-"`
	// Makes sense if RootUpdate.TypeOfUpdate is:
	// 	* RootUpdateTypeLevel2KeysUpdate
	// 	* RootUpdateTypeLevel2KeysUpdateV1
	Authorizations *Authorizations `json:"-"`
}

func (u *RootUpdate) UnmarshalJSON(b []byte) error {
	v := &struct {
		TypeOfUpdate RootUpdateType `json:"typeOfUpdate"`
	}{}
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	u.TypeOfUpdate = v.TypeOfUpdate
	var w any
	switch u.TypeOfUpdate {
	case RootUpdateTypeRootKeysUpdate, RootUpdateTypeLevel1KeysUpdate:
		u.HigherLevelKeys = &HigherLevelKeys{}
		w = u.HigherLevelKeys
	case RootUpdateTypeLevel2KeysUpdate, RootUpdateTypeLevel2KeysUpdateV1:
		u.Authorizations = &Authorizations{}
		w = u.Authorizations
	default:
		return fmt.Errorf("%T: unexpected type: %s", *u, u.TypeOfUpdate)
	}
	return json.Unmarshal(b, w)
}

// Level1UpdateType describes Level1Update type. See related constants.
type Level1UpdateType string

const (
	Level1UpdateTypeLevel1KeysUpdate   Level1UpdateType = "level1KeysUpdate"
	Level1UpdateTypeLevel2KeysUpdate   Level1UpdateType = "level2KeysUpdate"
	Level1UpdateTypeLevel2KeysUpdateV1 Level1UpdateType = "level2KeysUpdateV1"
)

// Level1Update is an update with level 1 keys of either level 1 or level 2 keys. Each of the
// updates must be a separate transaction.
type Level1Update struct {
	TypeOfUpdate Level1UpdateType `json:"typeOfUpdate"`
	// Makes sense if Level1Update.TypeOfUpdate is:
	// 	* Level1UpdateTypeLevel1KeysUpdate
	HigherLevelKeys *HigherLevelKeys `json:"-"`
	// Makes sense if Level1Update.TypeOfUpdate is:
	// 	* Level1UpdateTypeLevel2KeysUpdate
	// 	* Level1UpdateTypeLevel2KeysUpdateV1
	Authorizations *Authorizations `json:"-"`
}

func (u *Level1Update) UnmarshalJSON(b []byte) error {
	v := &struct {
		TypeOfUpdate Level1UpdateType `json:"typeOfUpdate"`
	}{}
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	u.TypeOfUpdate = v.TypeOfUpdate
	var w any
	switch u.TypeOfUpdate {
	case Level1UpdateTypeLevel1KeysUpdate:
		u.HigherLevelKeys = &HigherLevelKeys{}
		w = u.HigherLevelKeys
	case Level1UpdateTypeLevel2KeysUpdate, Level1UpdateTypeLevel2KeysUpdateV1:
		u.Authorizations = &Authorizations{}
		w = u.Authorizations
	default:
		return fmt.Errorf("%T: unexpected type: %s", *u, u.TypeOfUpdate)
	}
	return json.Unmarshal(b, w)
}

// PoolParameters are the parameters related to staking pools.
type PoolParameters struct {
	// The range of allowed baker commissions.
	BakingCommissionRange *InclusiveRange `json:"bakingCommissionRange"`
	// Maximum fraction of the total staked capital of that a new baker can have.
	CapitalBound float64 `json:"capitalBound"`
	// The range of allowed finalization commissions.
	FinalizationCommissionRange *InclusiveRange `json:"finalizationCommissionRange"`
	// The maximum leverage that a baker can have as a ratio of total stake to equity capital.
	LeverageBound *LeverageFactor `json:"leverageBound"`
	// Minimum equity capital required for a new baker.
	MinimumEquityCapital *Amount `json:"minimumEquityCapital"`
	// Fraction of baking rewards charged by the passive delegation.
	PassiveBakingCommission float64 `json:"passiveBakingCommission"`
	// Fraction of finalization rewards charged by the passive delegation.
	PassiveFinalizationCommission float64 `json:"passiveFinalizationCommission"`
	// Fraction of transaction rewards charged by the L-pool.
	PassiveTransactionCommission float64 `json:"passiveTransactionCommission"`
	// The range of allowed transaction commissions.
	TransactionCommissionRange *InclusiveRange `json:"transactionCommissionRange"`
}

// TimeParameters is the time parameters are introduced as of protocol version 4, and consist
// of the reward period length and the mint rate per payday. These are coupled as a change to
// either affects the overall rate of minting.
type TimeParameters struct {
	MintPerPayday      float64 `json:"mintPerPayday"`
	RewardPeriodLength uint64  `json:"rewardPeriodLength"`
}

// BlockItemTypeType describes BlockItemType type. See related constants.
type BlockItemTypeType string

const (
	// BlockItemTypeTypeAccountTransaction means an Account transactions are transactions
	// that are signed by an account. Most transactions are account transactions.
	BlockItemTypeTypeAccountTransaction BlockItemTypeType = "accountTransaction"
	// BlockItemTypeTypeCredentialDeploymentTransaction means credential deployments that create
	// accounts are special kinds of transactions. They are not signed by the account in the
	// usual way, and they are not paid for  directly by the sender.
	BlockItemTypeTypeCredentialDeploymentTransaction BlockItemTypeType = "credentialDeploymentTransaction"
	// BlockItemTypeTypeUpdateTransaction means chain updates are signed by the governance keys.
	// They affect the core parameters of the chain.
	BlockItemTypeTypeUpdateTransaction BlockItemTypeType = "updateTransaction"
)

// TransactionType describes types of account transactions. See related constants.
type TransactionType string

const (
	TransactionTypeDeployModule                    TransactionType = "deployModule"
	TransactionTypeInitContract                    TransactionType = "initContract"
	TransactionTypeUpdate                          TransactionType = "update"
	TransactionTypeTransfer                        TransactionType = "transfer"
	TransactionTypeAddBaker                        TransactionType = "addBaker"
	TransactionTypeRemoveBaker                     TransactionType = "removeBaker"
	TransactionTypeUpdateBakerStake                TransactionType = "updateBakerStake"
	TransactionTypeUpdateBakerRestakeEarnings      TransactionType = "updateBakerRestakeEarnings"
	TransactionTypeUpdateBakerKeys                 TransactionType = "updateBakerKeys"
	TransactionTypeUpdateCredentialKeys            TransactionType = "updateCredentialKeys"
	TransactionTypeEncryptedAmountTransfer         TransactionType = "encryptedAmountTransfer"
	TransactionTypeTransferToEncrypted             TransactionType = "transferToEncrypted"
	TransactionTypeTransferToPublic                TransactionType = "transferToPublic"
	TransactionTypeTransferWithSchedule            TransactionType = "transferWithSchedule"
	TransactionTypeUpdateCredentials               TransactionType = "updateCredentials"
	TransactionTypeRegisterData                    TransactionType = "registerData"
	TransactionTypeTransferWithMemo                TransactionType = "transferWithMemo"
	TransactionTypeEncryptedAmountTransferWithMemo TransactionType = "encryptedAmountTransferWithMemo"
	TransactionTypeTransferWithScheduleAndMemo     TransactionType = "transferWithScheduleAndMemo"
	TransactionTypeConfigureBaker                  TransactionType = "configureBaker"
	TransactionTypeConfigureDelegation             TransactionType = "configureDelegation"
)

// CredentialType is enumeration of the types of credentials. See related constants.
type CredentialType string

const (
	CredentialTypeInitial CredentialType = "initial"
	CredentialTypeNormal  CredentialType = "normal"
)

// UpdateType is enumeration of the types of updates that are possible. See related constants.
type UpdateType string

const (
	UpdateTypeUpdateProtocol                   UpdateType = "updateProtocol"
	UpdateTypeUpdateElectionDifficulty         UpdateType = "updateElectionDifficulty"
	UpdateTypeUpdateEuroPerEnergy              UpdateType = "updateEuroPerEnergy"
	UpdateTypeUpdateMicroGTUPerEuro            UpdateType = "updateMicroGTUPerEuro"
	UpdateTypeUpdateFoundationAccount          UpdateType = "updateFoundationAccount"
	UpdateTypeUpdateMintDistribution           UpdateType = "updateMintDistribution"
	UpdateTypeUpdateTransactionFeeDistribution UpdateType = "updateTransactionFeeDistribution"
	UpdateTypeUpdateGASRewards                 UpdateType = "updateGASRewards"
	UpdateTypeUpdateAddAnonymityRevoker        UpdateType = "updateAddAnonymityRevoker"
	UpdateTypeUpdateAddIdentityProvider        UpdateType = "updateAddIdentityProvider"
	UpdateTypeUpdateRootKeys                   UpdateType = "updateRootKeys"
	UpdateTypeUpdateLevel1Keys                 UpdateType = "updateLevel1Keys"
	UpdateTypeUpdateLevel2Keys                 UpdateType = "updateLevel2Keys"
	UpdateTypeUpdatePoolParameters             UpdateType = "updatePoolParameters"
	UpdateTypeUpdateCooldownParameters         UpdateType = "updateCooldownParameters"
	UpdateTypeUpdateTimeParameters             UpdateType = "updateTimeParameters"
)

// BlockItemType is the type of the block item.
type BlockItemType struct {
	Type BlockItemTypeType `json:"type"`
	// Makes sense only if BlockItemType.Type is BlockItemTypeTypeAccountTransaction
	AccountTransaction TransactionType `json:"-"`
	// Makes sense only if BlockItemType.Type is BlockItemTypeTypeCredentialDeploymentTransaction
	CredentialDeploymentTransaction CredentialType `json:"-"`
	// Makes sense only if BlockItemType.Type is BlockItemTypeTypeUpdateTransaction
	UpdateTransaction UpdateType `json:"-"`
}

func (t *BlockItemType) UnmarshalJSON(b []byte) error {
	v := &struct {
		Type     BlockItemTypeType `json:"type"`
		Contents string            `json:"contents"`
	}{}
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	t.Type = v.Type
	switch t.Type {
	case BlockItemTypeTypeAccountTransaction:
		t.AccountTransaction = TransactionType(v.Contents)
	case BlockItemTypeTypeCredentialDeploymentTransaction:
		t.CredentialDeploymentTransaction = CredentialType(v.Contents)
	case BlockItemTypeTypeUpdateTransaction:
		t.UpdateTransaction = UpdateType(v.Contents)
	default:
		return fmt.Errorf("%T: unexpected type: %s", *t, t.Type)
	}
	return nil
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
