package concordium

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	InvokeContractResultSuccess InvokeContractResultTag = "success"
	InvokeContractResultFailure InvokeContractResultTag = "failure"

	moduleRefSize       = 32
	initNamePrefix      = "init_"
	receiveNameSplitter = "."

	RejectReasonTagModuleNotWF                            RejectReasonTag = "ModuleNotWF"
	RejectReasonTagModuleHashAlreadyExists                RejectReasonTag = "ModuleHashAlreadyExists"
	RejectReasonTagInvalidAccountReference                RejectReasonTag = "InvalidAccountReference"
	RejectReasonTagInvalidInitMethod                      RejectReasonTag = "InvalidInitMethod"
	RejectReasonTagInvalidReceiveMethod                   RejectReasonTag = "InvalidReceiveMethod"
	RejectReasonTagInvalidModuleReference                 RejectReasonTag = "InvalidModuleReference"
	RejectReasonTagInvalidContractAddress                 RejectReasonTag = "InvalidContractAddress"
	RejectReasonTagRuntimeFailure                         RejectReasonTag = "RuntimeFailure"
	RejectReasonTagAmountTooLarge                         RejectReasonTag = "AmountTooLarge"
	RejectReasonTagSerializationFailure                   RejectReasonTag = "SerializationFailure"
	RejectReasonTagOutOfEnergy                            RejectReasonTag = "OutOfEnergy"
	RejectReasonTagRejectedInit                           RejectReasonTag = "RejectedInit"
	RejectReasonTagRejectedReceive                        RejectReasonTag = "RejectedReceive"
	RejectReasonTagNonExistentRewardAccount               RejectReasonTag = "NonExistentRewardAccount"
	RejectReasonTagInvalidProof                           RejectReasonTag = "InvalidProof"
	RejectReasonTagAlreadyABaker                          RejectReasonTag = "AlreadyABaker"
	RejectReasonTagNotABaker                              RejectReasonTag = "NotABaker"
	RejectReasonTagInsufficientBalanceForBakerStake       RejectReasonTag = "InsufficientBalanceForBakerStake"
	RejectReasonTagStakeUnderMinimumThresholdForBaking    RejectReasonTag = "StakeUnderMinimumThresholdForBaking"
	RejectReasonTagBakerInCooldown                        RejectReasonTag = "BakerInCooldown"
	RejectReasonTagDuplicateAggregationKey                RejectReasonTag = "DuplicateAggregationKey"
	RejectReasonTagNonExistentCredentialID                RejectReasonTag = "NonExistentCredentialID"
	RejectReasonTagKeyIndexAlreadyInUse                   RejectReasonTag = "KeyIndexAlreadyInUse"
	RejectReasonTagInvalidAccountThreshold                RejectReasonTag = "InvalidAccountThreshold"
	RejectReasonTagInvalidCredentialKeySignThreshold      RejectReasonTag = "InvalidCredentialKeySignThreshold"
	RejectReasonTagInvalidEncryptedAmountTransferProof    RejectReasonTag = "InvalidEncryptedAmountTransferProof"
	RejectReasonTagInvalidTransferToPublicProof           RejectReasonTag = "InvalidTransferToPublicProof"
	RejectReasonTagEncryptedAmountSelfTransfer            RejectReasonTag = "EncryptedAmountSelfTransfer"
	RejectReasonTagInvalidIndexOnEncryptedTransfer        RejectReasonTag = "InvalidIndexOnEncryptedTransfer"
	RejectReasonTagZeroScheduledAmount                    RejectReasonTag = "ZeroScheduledAmount"
	RejectReasonTagNonIncreasingSchedule                  RejectReasonTag = "NonIncreasingSchedule"
	RejectReasonTagFirstScheduledReleaseExpired           RejectReasonTag = "FirstScheduledReleaseExpired"
	RejectReasonTagScheduledSelfTransfer                  RejectReasonTag = "ScheduledSelfTransfer"
	RejectReasonTagInvalidCredentials                     RejectReasonTag = "InvalidCredentials"
	RejectReasonTagDuplicateCredIDs                       RejectReasonTag = "DuplicateCredIDs"
	RejectReasonTagNonExistentCredIDs                     RejectReasonTag = "NonExistentCredIDs"
	RejectReasonTagRemoveFirstCredential                  RejectReasonTag = "RemoveFirstCredential"
	RejectReasonTagCredentialHolderDidNotSign             RejectReasonTag = "CredentialHolderDidNotSign"
	RejectReasonTagNotAllowedMultipleCredentials          RejectReasonTag = "NotAllowedMultipleCredentials"
	RejectReasonTagNotAllowedToReceiveEncrypted           RejectReasonTag = "NotAllowedToReceiveEncrypted"
	RejectReasonTagNotAllowedToHandleEncrypted            RejectReasonTag = "NotAllowedToHandleEncrypted"
	RejectReasonTagMissingBakerAddParameters              RejectReasonTag = "MissingBakerAddParameters"
	RejectReasonTagFinalizationRewardCommissionNotInRange RejectReasonTag = "FinalizationRewardCommissionNotInRange"
	RejectReasonTagBakingRewardCommissionNotInRange       RejectReasonTag = "BakingRewardCommissionNotInRange"
	RejectReasonTagTransactionFeeCommissionNotInRange     RejectReasonTag = "TransactionFeeCommissionNotInRange"
	RejectReasonTagAlreadyADelegator                      RejectReasonTag = "AlreadyADelegator"
	RejectReasonTagInsufficientBalanceForDelegationStake  RejectReasonTag = "InsufficientBalanceForDelegationStake"
	RejectReasonTagMissingDelegationAddParameters         RejectReasonTag = "MissingDelegationAddParameters"
	RejectReasonTagDelegatorInCooldown                    RejectReasonTag = "DelegatorInCooldown"
	RejectReasonTagNotADelegator                          RejectReasonTag = "NotADelegator"
	RejectReasonTagDelegationTargetNotABaker              RejectReasonTag = "DelegationTargetNotABaker"
	RejectReasonTagStakeOverMaximumThresholdForPool       RejectReasonTag = "StakeOverMaximumThresholdForPool"
	RejectReasonTagPoolWouldBecomeOverDelegated           RejectReasonTag = "PoolWouldBecomeOverDelegated"
	RejectReasonTagPoolClosed                             RejectReasonTag = "PoolClosed"
)

type RejectReasonTag string

type InvokeContractResultTag string

// ModuleRef base-16 encoded module reference (64 characters)
type ModuleRef [moduleRefSize]byte

func NewModuleRefFromString(s string) (ModuleRef, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return ModuleRef{}, fmt.Errorf("hex decode: %w", err)
	}
	if len(b) != moduleRefSize {
		return ModuleRef{}, fmt.Errorf("expect %d bytes but %d given", moduleRefSize, len(b))
	}
	var h ModuleRef
	copy(h[:], b)
	return h, nil
}

func MustNewModuleRefFromString(s string) ModuleRef {
	h, err := NewModuleRefFromString(s)
	if err != nil {
		panic("MustNewModuleRefFromString: " + err.Error())
	}
	return h
}

func (r *ModuleRef) String() string {
	return hex.EncodeToString((*r)[:])
}

func (r ModuleRef) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(r[:])
	if err != nil {
		return nil, fmt.Errorf("%T: %w", r, err)
	}
	return b, nil
}

func (r *ModuleRef) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *r, err)
	}
	var x ModuleRef
	copy(x[:], v)
	*r = x
	return nil
}

func (r *ModuleRef) Serialize() ([]byte, error) {
	return (*r)[:], nil
}

func (r *ModuleRef) Deserialize(b []byte) error {
	var v ModuleRef
	copy(v[:], b)
	*r = v
	return nil
}

type ContractName string

func (n *ContractName) SerializeModel() ([]byte, error) {
	return nil, fmt.Errorf("use %T instead of %T", InitName(""), *n)
}

func (n *ContractName) DeserializeModel([]byte) (int, error) {
	return 0, fmt.Errorf("use %T instead of %T", InitName(""), *n)
}

type InitName string

func NewInitName(contractName ContractName) InitName {
	return InitName(initNamePrefix + contractName)
}

func (n *InitName) Serialize() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.BigEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *InitName) SerializeModel() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.LittleEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *InitName) DeserializeModel(b []byte) (int, error) {
	i := 2
	if len(b) < i {
		return 0, fmt.Errorf("%T requires %d bytes", *n, i)
	}
	l := int(binary.LittleEndian.Uint16(b))
	*n = InitName(b[i : i+l])
	return i + l, nil
}

type ReceiveName string

func NewReceiveName(contractName ContractName, receiver string) ReceiveName {
	return ReceiveName(string(contractName) + receiveNameSplitter + receiver)
}

func (n *ReceiveName) Serialize() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.BigEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *ReceiveName) SerializeModel() ([]byte, error) {
	nLen := len(*n)
	b := make([]byte, 2+nLen)
	binary.LittleEndian.PutUint16(b, uint16(nLen))
	copy(b[2:], *n)
	return b, nil
}

func (n *ReceiveName) DeserializeModel(b []byte) (int, error) {
	i := 2
	if len(b) < i {
		return 0, fmt.Errorf("%T requires %d bytes", *n, i)
	}
	l := int(binary.LittleEndian.Uint16(b))
	*n = ReceiveName(b[i : i+l])
	return i + l, nil
}

// ContractContext contains data needed to invoke the contract.
type ContractContext struct {
	// Amount to invoke the contract with.
	Amount *Amount `json:"amount"`
	// Contract to invoke.
	Contract *ContractAddress `json:"contract"`
	// The amount of energy to allow for execution. This should be small enough so that
	// it can be converted to interpreter energy. Default 10000000.
	Energy uint64 `json:"energy"`
	// Invoker of the contract. If this is not supplied then the contract will be invoked
	// by an account with address 0, no credentials and sufficient amount of CCD to cover
	// the transfer amount. If given, the relevant address must exist in the blockstate.
	Invoker *Address `json:"invoker"`
	// Which entrypoint to invoke.
	Method ReceiveName `json:"method"`
	// Parameter to invoke with.
	Parameter Model `json:"parameter"`
}

type InvokeContractResult struct {
	Tag         InvokeContractResultTag `json:"tag"`
	ReturnValue Model                   `json:"returnValue"`
	UsedEnergy  uint64                  `json:"usedEnergy"`
	// In case when InvokeContractResult.Tag is InvokeContractResultSuccess
	Events []*ContractTraceElement `json:"events"`
	// In case when InvokeContractResult.Tag is InvokeContractResultFailure
	Reason RejectReason `json:"reason"`
}

// ContractTraceElement a successful contract invocation produces a sequence of effects on
// smart contracts and possibly accounts (if any contract transfers CCD to an account).
type ContractTraceElement struct {
	// A contract instance was updated.
	Updated struct {
		Data *InstanceUpdatedEvent `json:"data"`
	} `json:"Updated"`
	// A contract transferred am amount to the account,
	Transferred struct {
		// Amount transferred.
		Amount *Amount `json:"amount"`
		// Sender contract.
		From *ContractAddress `json:"from"`
		// Receiver account.
		To AccountAddress `json:"to"`
	} `json:"Transferred"`
	Interrupted struct {
		Address *ContractAddress `json:"address"`
		Events  []Model          `json:"events"`
	} `json:"Interrupted"`
	Resumed struct {
		Address *ContractAddress `json:"address"`
		Success bool             `json:"success"`
	} `json:"Resumed"`
}

// InstanceUpdatedEvent is the data generated as part of updating a single contract instance.
// In general a single [Update](transactions::Payload::Update) transaction will generate one
// or more of these events, together with possibly some transfers.
type InstanceUpdatedEvent struct {
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

// RejectReason is the reason for why a transaction was rejected. Rejected means included in a block,
// but the desired action was not achieved. The only effect of a rejected transaction is payment.
// NOTE: Some of the variant definitions can look peculiar, but they are made to be compatible with
// the serialization of the Haskell datatype.
type RejectReason struct {
	// Error raised when validating the Wasm module.
	Tag string `json:"tag"` // ModuleNotWF

	// As the name says.
	Tag      string    `json:"tag"` // ModuleHashAlreadyExists
	Contents ModuleRef `json:"contents"`

	// Account does not exist.
	Tag      string         `json:"tag"` // InvalidAccountReference
	Contents AccountAddress `json:"contents"`

	// Reference to a non-existing contract init method.
	Tag      string    `json:"tag"` // InvalidInitMethod
	Contents [2]string `json:"contents"`

	// Reference to a non-existing contract receive method.
	Tag      string    `json:"tag"` // InvalidReceiveMethod
	Contents [2]string `json:"contents"`

	// Reference to a non-existing module.
	Tag      string    `json:"tag"` // InvalidModuleReference
	Contents ModuleRef `json:"contents"`

	// Contract instance does not exist.
	Tag      string           `json:"tag"` // InvalidContractAddress
	Contents *ContractAddress `json:"contents"`

	// Runtime exception occurred when running either the init or receive method.
	Tag string `json:"tag"` // RuntimeFailure

	// When one wishes to transfer an amount from A to B but there are not enough
	// funds on account/contract A to make this possible. The data are the from
	// address and the amount to transfer.
	Tag      string `json:"tag"`      // AmountTooLarge
	Contents [2]any `json:"contents"` // Address and Amount

	// Serialization of the body failed.
	Tag string `json:"tag"` // SerializationFailure

	// We ran of out energy to process this transaction.
	Tag string `json:"tag"` // OutOfEnergy

	// Rejected due to contract logic in init function of a contract.
	Tag          string `json:"tag"` // RejectedInit
	RejectReason int32  `json:"rejectReason"`

	Tag             string           `json:"tag"` // RejectedReceive
	ContractAddress *ContractAddress `json:"contractAddress"`
	Parameter       Model            `json:"parameter"`
	ReceiveName     string           `json:"receiveName"`
	RejectReason    int32            `json:"rejectReason"`

	// Reward account desired by the baker does not exist.
	Tag      string         `json:"tag"` // NonExistentRewardAccount
	Contents AccountAddress `json:"contents"`

	// Proof that the baker owns relevant private keys is not valid.
	Tag string `json:"tag"` // InvalidProof

	// Tried to add baker for an account that already has a baker
	Tag      string `json:"tag"` // AlreadyABaker
	Contents uint64 `json:"contents"`

	// Tried to remove a baker for an account that has no baker
	Tag      string         `json:"tag"` // NotABaker
	Contents AccountAddress `json:"contents"`

	// The amount on the account was insufficient to cover the proposed stake
	Tag string `json:"tag"` // InsufficientBalanceForBakerStake

	// The amount provided is under the threshold required for becoming a baker
	Tag string `json:"tag"` // StakeUnderMinimumThresholdForBaking

	// The change could not be made because the baker is in cooldown for another change
	Tag string `json:"tag"` // BakerInCooldown

	// A baker with the given aggregation key already exists
	Tag      string `json:"tag"` // DuplicateAggregationKey
	Contents string `json:"contents"`

	// Encountered credential ID that does not exist
	Tag string `json:"tag"` // NonExistentCredentialID

	// Attempted to add an account key to a key index already in use
	Tag string `json:"tag"` // KeyIndexAlreadyInUse

	// When the account threshold is updated, it must not exceed the amount of existing keys
	Tag string `json:"tag"` // InvalidAccountThreshold

	// When the credential key threshold is updated, it must not exceed the amount of existing keys
	Tag string `json:"tag"` // InvalidCredentialKeySignThreshold

	// Proof for an encrypted amount transfer did not validate.
	Tag string `json:"tag"` // InvalidEncryptedAmountTransferProof

	// Proof for a secret to public transfer did not validate.
	Tag string `json:"tag"` // InvalidTransferToPublicProof

	// Account tried to transfer an encrypted amount to itself, that's not allowed.
	Tag      string         `json:"tag"` // EncryptedAmountSelfTransfer
	Contents AccountAddress `json:"contents"`

	// The provided index is below the start index or above `startIndex + length incomingAmounts`
	Tag string `json:"tag"` // InvalidIndexOnEncryptedTransfer

	// The transfer with schedule is going to send 0 tokens
	Tag string `json:"tag"` // ZeroScheduledAmount

	// The transfer with schedule has a non strictly increasing schedule
	Tag string `json:"tag"` // NonIncreasingSchedule

	// The first scheduled release in a transfer with schedule has already expired
	Tag string `json:"tag"` // FirstScheduledReleaseExpired

	// Account tried to transfer with schedule to itself, that's not allowed.
	Tag      string         `json:"tag"` // ScheduledSelfTransfer
	Contents AccountAddress `json:"contents"`

	// At least one of the credentials was either malformed or its proof was incorrect.
	Tag string `json:"tag"` // InvalidCredentials

	// Some of the credential IDs already exist or are duplicated in the transaction.
	Tag      string   `json:"tag"` // DuplicateCredIDs
	Contents []string `json:"contents"`

	// A credential id that was to be removed is not part of the account.
	Tag      string   `json:"tag"` // NonExistentCredIDs
	Contents []string `json:"contents"`

	// Attemp to remove the first credential
	Tag string `json:"tag"` // RemoveFirstCredential

	// The credential holder of the keys to be updated did not sign the transaction
	Tag string `json:"tag"` // CredentialHolderDidNotSign

	// Account is not allowed to have multiple credentials because it contains a non-zero encrypted transfer.
	Tag string `json:"tag"` // NotAllowedMultipleCredentials

	// The account is not allowed to receive encrypted transfers because it has multiple credentials.
	Tag string `json:"tag"` // NotAllowedToReceiveEncrypted

	// The account is not allowed to send encrypted transfers (or transfer from/to public to/from encrypted)
	Tag string `json:"tag"` // NotAllowedToHandleEncrypted

	// A configure baker transaction is missing one or more arguments in order to add a baker.
	Tag string `json:"tag"` // MissingBakerAddParameters

	// Finalization reward commission is not in the valid range for a baker
	Tag string `json:"tag"` // FinalizationRewardCommissionNotInRange

	// Baking reward commission is not in the valid range for a baker
	Tag string `json:"tag"` // BakingRewardCommissionNotInRange

	// Transaction fee commission is not in the valid range for a baker
	Tag string `json:"tag"` // TransactionFeeCommissionNotInRange

	// Tried to add baker for an account that already has a delegator.
	Tag string `json:"tag"` // AlreadyADelegator

	// The amount on the account was insufficient to cover the proposed stake.
	Tag string `json:"tag"` // InsufficientBalanceForDelegationStake

	// A configure delegation transaction is missing one or more arguments in order to add a delegator.
	Tag string `json:"tag"` // MissingDelegationAddParameters

	// Account is not a delegation account.
	Tag string `json:"tag"` // DelegatorInCooldown

	// Account is not a delegation account.
	Tag      string         `json:"tag"` // NotADelegator
	Contents AccountAddress `json:"contents"`

	// Delegation target is not a baker
	Tag      string `json:"tag"` // DelegationTargetNotABaker
	Contents uint64 `json:"contents"`

	// The amount would result in pool capital higher than the maximum threshold.
	Tag string `json:"tag"` // StakeOverMaximumThresholdForPool

	// The amount would result in pool with a too high fraction of delegated capital.
	Tag string `json:"tag"` // PoolWouldBecomeOverDelegated

	// The pool is not open to delegators.
	Tag string `json:"tag"` // PoolClosed
}
