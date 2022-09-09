package concordium

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	moduleRefSize       = 32
	initNamePrefix      = "init_"
	receiveNameSplitter = "."
)

// ModuleRef base-16 encoded module reference (64 characters)
type ModuleRef [moduleRefSize]byte

// NewModuleRef creates a new ModuleRef from string.
func NewModuleRef(s string) (ModuleRef, error) {
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

// MustNewModuleRef calls the NewModuleRef. It panics in case of error.
func MustNewModuleRef(s string) ModuleRef {
	h, err := NewModuleRef(s)
	if err != nil {
		panic("MustNewModuleRef: " + err.Error())
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

// InvokeContractResultTag describes InvokeContractResult type. See related constants.
type InvokeContractResultTag string

const (
	InvokeContractResultSuccess InvokeContractResultTag = "success"
	InvokeContractResultFailure InvokeContractResultTag = "failure"
)

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

// RejectReasonTag describes RejectReason type. See related constants.
type RejectReasonTag string

const (
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

// RejectReason is the reason for why a transaction was rejected. Rejected means included in a block,
// but the desired action was not achieved. The only effect of a rejected transaction is payment.
// NOTE: Some of the variant definitions can look peculiar, but they are made to be compatible with
// the serialization of the Haskell datatype.
type RejectReason struct {
	Tag RejectReasonTag `json:"tag"`
	// Error raised when validating the Wasm module.
	ModuleNotWF struct{} `json:"-"`
	// As the name says.
	ModuleHashAlreadyExists struct {
		Contents ModuleRef `json:"contents"`
	} `json:"-"`
	// Account does not exist.
	InvalidAccountReference struct {
		Contents AccountAddress `json:"contents"`
	} `json:"-"`
	// Reference to a non-existing contract init method.
	InvalidInitMethod struct {
		Contents [2]string `json:"contents"`
	} `json:"-"`
	// Reference to a non-existing contract receive method.
	InvalidReceiveMethod struct {
		Contents [2]string `json:"contents"`
	} `json:"-"`
	// Reference to a non-existing module.
	InvalidModuleReference struct {
		Contents ModuleRef `json:"contents"`
	} `json:"-"`
	// Contract instance does not exist.
	InvalidContractAddress struct {
		Contents *ContractAddress `json:"contents"`
	} `json:"-"`
	// Runtime exception occurred when running either the init or receive method.
	RuntimeFailure struct{} `json:"-"`
	// When one wishes to transfer an amount from A to B but there are not enough
	// funds on account/contract A to make this possible. The data are the from
	// address and the amount to transfer.
	AmountTooLarge struct {
		Contents [2]any `json:"contents"` // Address and Amount
	} `json:"-"`
	// Serialization of the body failed.
	SerializationFailure struct{} `json:"-"`
	// We ran of out energy to process this transaction.
	OutOfEnergy struct{} `json:"-"`
	// Rejected due to contract logic in init function of a contract.
	RejectedInit struct {
		RejectReason int32 `json:"rejectReason"`
	} `json:"-"`
	RejectedReceive struct {
		ContractAddress *ContractAddress `json:"contractAddress"`
		Parameter       Model            `json:"parameter"`
		ReceiveName     string           `json:"receiveName"`
		RejectReason    int32            `json:"rejectReason"`
	} `json:"-"`
	// Reward account desired by the baker does not exist.
	NonExistentRewardAccount struct {
		Contents AccountAddress `json:"contents"`
	} `json:"-"`
	// Proof that the baker owns relevant private keys is not valid.
	InvalidProof struct{} `json:"-"`
	// Tried to add baker for an account that already has a baker
	AlreadyABaker struct {
		Contents uint64 `json:"contents"`
	} `json:"-"`
	// Tried to remove a baker for an account that has no baker
	NotABaker struct {
		Contents AccountAddress `json:"contents"`
	} `json:"-"`
	// The amount on the account was insufficient to cover the proposed stake
	InsufficientBalanceForBakerStake struct{} `json:"-"`
	// The amount provided is under the threshold required for becoming a baker
	StakeUnderMinimumThresholdForBaking struct{} `json:"-"`
	// The change could not be made because the baker is in cooldown for another change
	BakerInCooldown struct{} `json:"-"`
	// A baker with the given aggregation key already exists
	DuplicateAggregationKey struct {
		Contents string `json:"contents"`
	} `json:"-"`
	// Encountered credential ID that does not exist
	NonExistentCredentialID struct{} `json:"-"`
	// Attempted to add an account key to a key index already in use
	KeyIndexAlreadyInUse struct{} `json:"-"`
	// When the account threshold is updated, it must not exceed the amount of existing keys
	InvalidAccountThreshold struct{} `json:"-"`
	// When the credential key threshold is updated, it must not exceed the amount of existing keys
	InvalidCredentialKeySignThreshold struct{} `json:"-"`
	// Proof for an encrypted amount transfer did not validate.
	InvalidEncryptedAmountTransferProof struct{} `json:"-"`
	// Proof for a secret to public transfer did not validate.
	InvalidTransferToPublicProof struct{} `json:"-"`
	// Account tried to transfer an encrypted amount to itself, that's not allowed.
	EncryptedAmountSelfTransfer struct {
		Contents AccountAddress `json:"contents"`
	} `json:"-"`
	// The provided index is below the start index or above `startIndex + length incomingAmounts`
	InvalidIndexOnEncryptedTransfer struct{} `json:"-"`
	// The transfer with schedule is going to send 0 tokens
	ZeroScheduledAmount struct{} `json:"-"`
	// The transfer with schedule has a non strictly increasing schedule
	NonIncreasingSchedule struct{} `json:"-"`
	// The first scheduled release in a transfer with schedule has already expired
	FirstScheduledReleaseExpired struct{} `json:"-"`
	// Account tried to transfer with schedule to itself, that's not allowed.
	ScheduledSelfTransfer struct {
		Contents AccountAddress `json:"contents"`
	} `json:"-"`
	// At least one of the credentials was either malformed or its proof was incorrect.
	InvalidCredentials struct{} `json:"-"`
	// Some of the credential IDs already exist or are duplicated in the transaction.
	DuplicateCredIDs struct {
		Contents []string `json:"contents"`
	} `json:"-"`
	// A credential id that was to be removed is not part of the account.
	NonExistentCredIDs struct {
		Contents []string `json:"contents"`
	} `json:"-"`
	// Attemp to remove the first credential
	RemoveFirstCredential struct{} `json:"-"`
	// The credential holder of the keys to be updated did not sign the transaction
	CredentialHolderDidNotSign struct{} `json:"-"`
	// Account is not allowed to have multiple credentials because it contains a non-zero encrypted transfer.
	NotAllowedMultipleCredentials struct{} `json:"-"`
	// The account is not allowed to receive encrypted transfers because it has multiple credentials.
	NotAllowedToReceiveEncrypted struct{} `json:"-"`
	// The account is not allowed to send encrypted transfers (or transfer from/to public to/from encrypted)
	NotAllowedToHandleEncrypted struct{} `json:"-"`
	// A configure baker transaction is missing one or more arguments in order to add a baker.
	MissingBakerAddParameters struct{} `json:"-"`
	// Finalization reward commission is not in the valid range for a baker
	FinalizationRewardCommissionNotInRange struct{} `json:"-"`
	// Baking reward commission is not in the valid range for a baker
	BakingRewardCommissionNotInRange struct{} `json:"-"`
	// Transaction fee commission is not in the valid range for a baker
	TransactionFeeCommissionNotInRange struct{} `json:"-"`
	// Tried to add baker for an account that already has a delegator.
	AlreadyADelegator struct{} `json:"-"`
	// The amount on the account was insufficient to cover the proposed stake.
	InsufficientBalanceForDelegationStake struct{} `json:"-"`
	// A configure delegation transaction is missing one or more arguments in order to add a delegator.
	MissingDelegationAddParameters struct{} `json:"-"`
	// Account is not a delegation account.
	DelegatorInCooldown struct{} `json:"-"`
	// Account is not a delegation account.
	NotADelegator struct {
		Contents AccountAddress `json:"contents"`
	} `json:"-"`
	// Delegation target is not a baker
	DelegationTargetNotABaker struct {
		Contents uint64 `json:"contents"`
	} `json:"-"`
	// The amount would result in pool capital higher than the maximum threshold.
	StakeOverMaximumThresholdForPool struct{} `json:"-"`
	// The amount would result in pool with a too high fraction of delegated capital.
	PoolWouldBecomeOverDelegated struct{} `json:"-"`
	// The pool is not open to delegators.
	PoolClosed struct{} `json:"-"`
}
