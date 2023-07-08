package v2

import (
	"encoding/hex"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
	"github.com/btcsuite/btcutil/base58"
)

const (
	AccountAddressLength  = 32
	BlockHashLength       = 32
	TransactionHashLength = 32
	ModuleRefLength       = 32
)

// WalletAccount an account imported from one of the supported export formats.
// This structure implements TransactionSigner and ExactSizeTransactionSigner, so it may be used for sending transactions.
// This structure does not have the encryption key for sending encrypted transfers, it only contains keys for signing transactions.
type WalletAccount struct {
	Address AccountAddress
	Keys    AccountKeys
}

// AccountKeys all account keys indexed by credentials.
type AccountKeys struct {
}

// SignatureThreshold threshold for the number of signatures required.
// The values of this type must maintain the property that they are not 0.
type SignatureThreshold []uint8

// AccountAddress an address of an account.
type AccountAddress [AccountAddressLength]byte

// ToBase58 encodes account address to string.
func (a AccountAddress) ToBase58() string {
	return base58.CheckEncode(a[:], 1)
}

// AccountAddressFromString decodes string to account.
func AccountAddressFromString(s string) AccountAddress {
	return AccountAddressFromBytes(base58.Decode(s))
}

// AccountAddressFromBytes creates account address from given bytes.
func AccountAddressFromBytes(b []byte) AccountAddress {
	var accountAddress AccountAddress
	if len(b) > AccountAddressLength {
		b = b[:AccountAddressLength]
	}
	copy(accountAddress[AccountAddressLength-len(b):], b)
	return accountAddress
}

// BlockHash hash of a block. This is always 32 bytes long.
type BlockHash [BlockHashLength]byte

// Hex encodes block hash to base16 string.
func (b BlockHash) Hex() string {
	return hex.EncodeToString(b[:])
}

type isBlockHashInput interface {
	isBlockHashInput()
}

func convertBlockHashInput(req isBlockHashInput) (_ *pb.BlockHashInput) {
	var res *pb.BlockHashInput
	switch v := req.(type) {
	case BlockHashInputBest:
		res = &pb.BlockHashInput{
			BlockHashInput: &pb.BlockHashInput_Best{},
		}
	case BlockHashInputLastFinal:
		res = &pb.BlockHashInput{
			BlockHashInput: &pb.BlockHashInput_LastFinal{},
		}
	case BlockHashInputGiven:
		res = &pb.BlockHashInput{
			BlockHashInput: &pb.BlockHashInput_Given{
				Given: &pb.BlockHash{
					Value: v.Given[:],
				}}}
	case BlockHashInputAbsoluteHeight:
		res = &pb.BlockHashInput{
			BlockHashInput: &pb.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &pb.AbsoluteBlockHeight{
					Value: uint64(v),
				}}}
	case BlockHashInputRelativeHeight:
		res = &pb.BlockHashInput{
			BlockHashInput: &pb.BlockHashInput_RelativeHeight_{
				RelativeHeight: &pb.BlockHashInput_RelativeHeight{
					GenesisIndex: &pb.GenesisIndex{Value: v.GenesisIndex},
					Height:       &pb.BlockHeight{Value: v.Height},
					Restrict:     v.Restrict,
				}}}
	}

	return res
}

// BlockHashInputBest query for the best block.
type BlockHashInputBest struct{}

func (BlockHashInputBest) isBlockHashInput() {}

// BlockHashInputLastFinal query for the last finalized block.
type BlockHashInputLastFinal struct{}

func (BlockHashInputLastFinal) isBlockHashInput() {}

// BlockHashInputGiven query for the block specified by the hash. This hash should always be 32 bytes.
type BlockHashInputGiven struct {
	Given BlockHash
}

func (BlockHashInputGiven) isBlockHashInput() {}

// BlockHashInputAbsoluteHeight query for a block at absolute height, if a unique block can be identified at that height.
type BlockHashInputAbsoluteHeight uint64

func (BlockHashInputAbsoluteHeight) isBlockHashInput() {}

// BlockHashInputRelativeHeight query for a block at height relative to a genesis index.
type BlockHashInputRelativeHeight struct {
	GenesisIndex uint32
	Height       uint64
	Restrict     bool
}

func (BlockHashInputRelativeHeight) isBlockHashInput() {}

// TransactionHash hash of a transaction. This is always 32 bytes long.
type TransactionHash [TransactionHashLength]byte

// Hex encodes transaction hash to base16 string.
func (t TransactionHash) Hex() string {
	return hex.EncodeToString(t[:])
}

// ModuleRef a smart contract module reference. This is always 32 bytes long.
type ModuleRef [ModuleRefLength]byte

// Hex encodes module ref to base16 string.
func (m ModuleRef) Hex() string {
	return hex.EncodeToString(m[:])
}

// BlockInfo the response for GetBlockInfo.
type BlockInfo struct {
	Hash                   BlockHash
	Height                 AbsoluteBlockHeight
	ParentBlock            BlockHash
	LastFinalizedBlock     BlockHash
	GenesisIndex           GenesisIndex
	EraBlockHeight         BlockHeight
	ReceiveTime            Timestamp
	ArriveTime             Timestamp
	SlotNumber             Slot
	SlotTime               Timestamp
	Baker                  BakerId
	Finalized              bool
	TransactionCount       uint32
	TransactionsEnergyCost Energy
	TransactionsSize       uint32
	StateHash              StateHash
	ProtocolVersion        ProtocolVersion
}

// ProtocolVersion he different versions of the protocol.
type ProtocolVersion int32

// StateHash hash of the state after some block. This is always 32 bytes long.
type StateHash []byte

// BakerId the ID of a baker, which is the index of its account.
type BakerId uint64

// Slot a number representing a slot for baking a block.
type Slot uint64

// Timestamp unix timestamp in milliseconds.
type Timestamp uint64

// GenesisIndex the number of chain restarts via a protocol update. An effected
// protocol update instruction might not change the protocol version
// specified in the previous field, but it always increments the genesis
// index.
type GenesisIndex uint32

// AbsoluteBlockHeight this is the number of ancestors of a block
// since the genesis block. In particular, the chain genesis block has absolute
// height 0.
type AbsoluteBlockHeight uint64

// BlockHeight the height of a block relative to the last genesis. This differs from the
// absolute block height in that it counts height from the last protocol update.
type BlockHeight uint64

func convertBlockInfo(b *pb.BlockInfo) BlockInfo {
	var hash, parentBlock, lastFinalizedBlock BlockHash

	copy(hash[:], b.Hash.Value)
	copy(parentBlock[:], b.ParentBlock.Value)
	copy(lastFinalizedBlock[:], b.LastFinalizedBlock.Value)

	return BlockInfo{
		Hash:                   hash,
		Height:                 AbsoluteBlockHeight(b.Height.Value),
		ParentBlock:            parentBlock,
		LastFinalizedBlock:     lastFinalizedBlock,
		GenesisIndex:           GenesisIndex(b.GenesisIndex.Value),
		EraBlockHeight:         BlockHeight(b.EraBlockHeight.Value),
		ReceiveTime:            Timestamp(b.ReceiveTime.Value),
		ArriveTime:             Timestamp(b.ArriveTime.Value),
		SlotNumber:             Slot(b.SlotNumber.Value),
		SlotTime:               Timestamp(b.SlotTime.Value),
		Baker:                  BakerId(b.Baker.Value),
		Finalized:              b.Finalized,
		TransactionCount:       b.TransactionCount,
		TransactionsEnergyCost: Energy(b.TransactionsEnergyCost.Value),
		TransactionsSize:       b.TransactionsSize,
		StateHash:              b.StateHash.Value,
		ProtocolVersion:        ProtocolVersion(b.ProtocolVersion),
	}
}

// ContractAddress address of a smart contract instance.
type ContractAddress struct {
	Index    uint64
	Subindex uint64
}

// BlockItem is account transaction or credential deployment or update instruction item.
type BlockItem struct {
	Hash      *TransactionHash
	BlockItem isBlockItem
}

type isBlockItem interface {
	isBlockItem()
}

// AccountTransaction messages which are signed and paid for by the sender account.
type AccountTransaction struct {
	Signature *AccountTransactionSignature
	Header    *AccountTransactionHeader
	Payload   *AccountTransactionPayload
}

func (AccountTransaction) isBlockItem() {}

// AccountTransactionSignature transaction signature.
type AccountTransactionSignature struct {
	Signatures map[uint32]*AccountSignatureMap
}

// AccountSignatureMap wrapper for a map from indexes to signatures.
// Needed because protobuf doesn't allow nested maps directly.
// The keys in the SignatureMap must not exceed 2^8.
type AccountSignatureMap struct {
	Signatures map[uint32]*Signature
}

// Signature a single signature. Used when sending block items to a node with `SendBlockItem`.
type Signature []byte

// AccountTransactionHeader header of an account transaction that contains basic data to check whether
// the sender and the transaction are valid. The header is shared by all transaction types.
type AccountTransactionHeader struct {
	Sender         *AccountAddress
	SequenceNumber *SequenceNumber
	EnergyAmount   *Energy
	Expiry         *TransactionTime
}

// SequenceNumber a sequence number that determines the ordering of transactions from the
// account. The minimum sequence number is 1.
type SequenceNumber uint64

// Energy is used to count exact execution cost.
// This cost is then converted to CCD amounts.
type Energy uint64

// TransactionTime specified as seconds since unix epoch.
type TransactionTime uint64

// AccountTransactionPayload the payload for an account transaction.
type AccountTransactionPayload struct {
	Payload isAccountTransactionPayload
}

type isAccountTransactionPayload interface {
	isAccountTransactionPayload()
}

// RawPayload a pre-serialized payload in the binary serialization format defined by the protocol.
type RawPayload []byte

func (RawPayload) isAccountTransactionPayload() {}

// DeployModule a transfer between two accounts. With an optional memo.
type DeployModule struct {
	DeployModule *VersionedModuleSource
}

func (DeployModule) isAccountTransactionPayload() {}

// VersionedModuleSource source bytes of a versioned smart contract module.
type VersionedModuleSource struct {
	Module isVersionedModuleSource
}

type isVersionedModuleSource interface {
	isVersionedModuleSource()
}

// ModuleSourceV0 v0.
type ModuleSourceV0 []byte

// ModuleSourceV1 v1
type ModuleSourceV1 []byte

func (ModuleSourceV0) isVersionedModuleSource() {}

func (ModuleSourceV1) isVersionedModuleSource() {}

type InitContract struct {
	Payload *InitContractPayload
}

func (InitContract) isAccountTransactionPayload() {}

// InitContractPayload data required to initialize a new contract instance.
type InitContractPayload struct {
	Amount    *Amount
	ModuleRef *ModuleRef
	InitName  *InitName
	Parameter *Parameter
}

// InitName the init name of a smart contract function
type InitName string

// Amount an amount of microCCD.
type Amount uint64

// Parameter to a smart contract initialization or invocation.
type Parameter []byte

type UpdateContract struct {
	Payload *UpdateContractPayload
}

func (UpdateContract) isAccountTransactionPayload() {}

// UpdateContractPayload data required to update a contract instance.
type UpdateContractPayload struct {
	Amount      *Amount
	Address     *ContractAddress
	ReceiveName *ReceiveName
	Parameter   *Parameter
}

// ReceiveName the reception name of a smart contract function. Expected format:
// `<contract_name>.<func_name>`. It must only consist of atmost 100 ASCII
// alphanumeric or punctuation characters, and must contain a '.'.
type ReceiveName string

type Transfer struct {
	Payload *TransferPayload
}

func (Transfer) isAccountTransactionPayload() {}

// TransferPayload payload of a transfer between two accounts.
type TransferPayload struct {
	Amount   *Amount
	Receiver *AccountAddress
}

type TransferWithMemo struct {
	Payload *TransferWithMemoPayload
}

func (TransferWithMemo) isAccountTransactionPayload() {}

// TransferWithMemoPayload payload of a transfer between two accounts with a memo.
type TransferWithMemoPayload struct {
	Amount   *Amount
	Receiver *AccountAddress
	Memo     *Memo
}

// Memo a memo which can be included as part of a transfer. Max size is 256 bytes.
type Memo []byte

type RegisterData struct {
	Payload *RegisteredData
}

func (RegisterData) isAccountTransactionPayload() {}

// RegisteredData data registered on the chain with a register data transaction.
type RegisteredData []byte

// CredentialDeployment create new accounts. They are not paid for
// directly by the sender. Instead, bakers are rewarded by the protocol for
// including them.
type CredentialDeployment struct {
	MessageExpiry *TransactionTime
	Payload       isCredentialDeploymentPayload
}

func (CredentialDeployment) isBlockItem() {}

type isCredentialDeploymentPayload interface {
	isCredentialDeploymentPayload()
}

func (RawPayload) isCredentialDeploymentPayload() {}

// UpdateInstruction messages which can update the chain parameters. Including which keys are allowed
// to make future update instructions.
type UpdateInstruction struct {
	Signatures *SignatureMap
	Header     *UpdateInstructionHeader
	Payload    *UpdateInstructionPayload
}

func (UpdateInstruction) isBlockItem() {}

// SignatureMap wrapper for a map from indexes to signatures.
// Needed because protobuf doesn't allow nested maps directly.
type SignatureMap struct {
	Signatures map[uint32]*Signature
}

type UpdateInstructionHeader struct {
	SequenceNumber *UpdateSequenceNumber
	EffectiveTime  *TransactionTime
	Timeout        *TransactionTime
}

// UpdateSequenceNumber determines the ordering of update transactions.
// Equivalent to `SequenceNumber` for account transactions.
// Update sequence numbers are per update type and the minimum value is 1.
type UpdateSequenceNumber uint64

// UpdateInstructionPayload payload.
type UpdateInstructionPayload struct {
	Payload isUpdateInstructionPayload
}

type isUpdateInstructionPayload interface {
	isUpdateInstructionPayload()
}

func (RawPayload) isUpdateInstructionPayload() {}

func convertBlockItems(c []*pb.BlockItem) []*BlockItem {
	var result []*BlockItem

	for _, v := range c {
		var blockItem BlockItem
		var t TransactionHash
		copy(t[:], v.Hash.Value)
		blockItem.Hash = &t

		switch k := v.BlockItem.(type) {
		case *pb.BlockItem_AccountTransaction:
			signatures := make(map[uint32]*AccountSignatureMap)

			for i, v := range k.AccountTransaction.Signature.Signatures {
				sig := make(map[uint32]*Signature)

				for j, k := range v.Signatures {
					var s Signature
					copy(s, k.Value)
					sig[j] = &s
				}

				signatures[i] = &AccountSignatureMap{Signatures: sig}
			}

			var accountAddress AccountAddress
			copy(accountAddress[:], k.AccountTransaction.Header.Sender.Value)

			var sequenceNumber = SequenceNumber(k.AccountTransaction.Header.SequenceNumber.Value)
			var energy = Energy(k.AccountTransaction.Header.EnergyAmount.Value)
			var expiry = TransactionTime(k.AccountTransaction.Header.Expiry.Value)

			var accountTransactionPayload AccountTransactionPayload

			switch payload := k.AccountTransaction.Payload.Payload.(type) {
			case *pb.AccountTransactionPayload_RawPayload:
				var rawPayload RawPayload
				copy(rawPayload[:], payload.RawPayload)
				accountTransactionPayload.Payload = &rawPayload
			case *pb.AccountTransactionPayload_DeployModule:
				var deployModule DeployModule
				switch dm := payload.DeployModule.Module.(type) {
				case *pb.VersionedModuleSource_V0:
					var mod ModuleSourceV0
					copy(mod, dm.V0.Value)
					deployModule.DeployModule.Module = &mod
				case *pb.VersionedModuleSource_V1:
					var mod ModuleSourceV1
					copy(mod, dm.V1.Value)
					deployModule.DeployModule.Module = &mod
				}
				accountTransactionPayload.Payload = &deployModule
			case *pb.AccountTransactionPayload_InitContract:
				var initContract InitContract

				var initName = InitName(payload.InitContract.InitName.Value)
				var amount = Amount(payload.InitContract.Amount.Value)
				var parameter Parameter
				copy(parameter, payload.InitContract.Parameter.Value)
				var moduleRef ModuleRef
				copy(moduleRef[:], payload.InitContract.ModuleRef.Value)

				initContract.Payload.InitName = &initName
				initContract.Payload.Amount = &amount
				initContract.Payload.Parameter = &parameter
				initContract.Payload.ModuleRef = &moduleRef

				accountTransactionPayload.Payload = initContract
			case *pb.AccountTransactionPayload_UpdateContract:
				var updateContract UpdateContract

				var amount = Amount(payload.UpdateContract.Amount.Value)
				var parameter Parameter
				copy(parameter, payload.UpdateContract.Parameter.Value)
				var address ContractAddress
				address.Subindex = payload.UpdateContract.Address.Subindex
				address.Index = payload.UpdateContract.Address.Index
				var receiveName = ReceiveName(payload.UpdateContract.ReceiveName.Value)

				updateContract.Payload.Address = &address
				updateContract.Payload.ReceiveName = &receiveName
				updateContract.Payload.Amount = &amount
				updateContract.Payload.Parameter = &parameter

				accountTransactionPayload.Payload = updateContract
			case *pb.AccountTransactionPayload_Transfer:
				var transfer Transfer

				var amount = Amount(payload.Transfer.Amount.Value)
				var receiver AccountAddress
				copy(receiver[:], payload.Transfer.Receiver.Value)

				transfer.Payload.Receiver = &receiver
				transfer.Payload.Amount = &amount

				accountTransactionPayload.Payload = transfer
			case *pb.AccountTransactionPayload_TransferWithMemo:
				var transferWithMemo TransferWithMemo

				var amount = Amount(payload.TransferWithMemo.Amount.Value)
				var receiver AccountAddress
				copy(receiver[:], payload.TransferWithMemo.Receiver.Value)
				var memo Memo
				copy(memo, payload.TransferWithMemo.Memo.Value)

				transferWithMemo.Payload.Memo = &memo
				transferWithMemo.Payload.Amount = &amount
				transferWithMemo.Payload.Receiver = &receiver

				accountTransactionPayload.Payload = transferWithMemo
			case *pb.AccountTransactionPayload_RegisterData:
				var registerData RegisterData
				var r RegisteredData
				copy(r, payload.RegisterData.Value)

				registerData.Payload = &r
				accountTransactionPayload.Payload = &registerData
			}

			blockItem.BlockItem = &AccountTransaction{
				Signature: &AccountTransactionSignature{
					Signatures: signatures,
				},
				Header: &AccountTransactionHeader{
					Sender:         &accountAddress,
					SequenceNumber: &sequenceNumber,
					EnergyAmount:   &energy,
					Expiry:         &expiry,
				},
				Payload: &accountTransactionPayload,
			}
		case *pb.BlockItem_CredentialDeployment:
			var credentialDeployment CredentialDeployment
			var transactionTime = TransactionTime(k.CredentialDeployment.MessageExpiry.Value)
			credentialDeployment.MessageExpiry = &transactionTime

			switch v := k.CredentialDeployment.Payload.(type) {
			case *pb.CredentialDeployment_RawPayload:
				var payload = RawPayload(v.RawPayload)
				credentialDeployment.Payload = &payload
			}

			blockItem.BlockItem = &credentialDeployment
		case *pb.BlockItem_UpdateInstruction:
			var signatureMap SignatureMap
			for i, v := range k.UpdateInstruction.Signatures.Signatures {
				signature := Signature(v.Value)
				signatureMap.Signatures[i] = &signature
			}

			sequenceNumber := UpdateSequenceNumber(k.UpdateInstruction.Header.SequenceNumber.Value)
			effectiveTime := TransactionTime(k.UpdateInstruction.Header.EffectiveTime.Value)
			timeout := TransactionTime(k.UpdateInstruction.Header.Timeout.Value)

			var updInstructionPayload = UpdateInstructionPayload{}

			switch t := k.UpdateInstruction.Payload.Payload.(type) {
			case *pb.UpdateInstructionPayload_RawPayload:
				updInstructionPayload.Payload = RawPayload(t.RawPayload)
			}

			blockItem.BlockItem = &UpdateInstruction{
				Signatures: &signatureMap,
				Header: &UpdateInstructionHeader{
					SequenceNumber: &sequenceNumber,
					EffectiveTime:  &effectiveTime,
					Timeout:        &timeout,
				},
				Payload: &updInstructionPayload,
			}
		}

		result = append(result, &blockItem)
	}

	return result
}
