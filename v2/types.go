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
type SignatureThreshold struct {
	Value []uint8
}

// AccountAddress an address of an account.
type AccountAddress struct {
	Value [AccountAddressLength]byte
}

// ToBase58 encodes account address to string.
func (a *AccountAddress) ToBase58() string {
	return base58.CheckEncode(a.Value[:], 1)
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
	copy(accountAddress.Value[AccountAddressLength-len(b):], b)
	return accountAddress
}

// BlockHash hash of a block. This is always 32 bytes long.
type BlockHash struct {
	Value [BlockHashLength]byte
}

// Hex encodes block hash to base16 string.
func (b *BlockHash) Hex() string {
	return hex.EncodeToString(b.Value[:])
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
					Value: v.Given.Value[:],
				}}}
	case BlockHashInputAbsoluteHeight:
		res = &pb.BlockHashInput{
			BlockHashInput: &pb.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &pb.AbsoluteBlockHeight{
					Value: v.Value,
				}}}
	case BlockHashInputRelativeHeight:
		res = &pb.BlockHashInput{
			BlockHashInput: &pb.BlockHashInput_RelativeHeight_{
				RelativeHeight: &pb.BlockHashInput_RelativeHeight{
					GenesisIndex: &pb.GenesisIndex{
						Value: v.GenesisIndex,
					},
					Height: &pb.BlockHeight{
						Value: v.Height,
					},
					Restrict: v.Restrict,
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
type BlockHashInputAbsoluteHeight struct {
	Value uint64
}

func (BlockHashInputAbsoluteHeight) isBlockHashInput() {}

// BlockHashInputRelativeHeight query for a block at height relative to a genesis index.
type BlockHashInputRelativeHeight struct {
	GenesisIndex uint32
	Height       uint64
	Restrict     bool
}

func (BlockHashInputRelativeHeight) isBlockHashInput() {}

// TransactionHash hash of a transaction. This is always 32 bytes long.
type TransactionHash struct {
	Value [TransactionHashLength]byte
}

// Hex encodes transaction hash to base16 string.
func (t *TransactionHash) Hex() string {
	return hex.EncodeToString(t.Value[:])
}

// ModuleRef a smart contract module reference. This is always 32 bytes long.
type ModuleRef struct {
	Value [ModuleRefLength]byte
}

// Hex encodes module ref to base16 string.
func (m *ModuleRef) Hex() string {
	return hex.EncodeToString(m.Value[:])
}

// BlockInfo the response for GetBlockInfo.
type BlockInfo struct {
	Hash                   *BlockHash
	Height                 *AbsoluteBlockHeight
	ParentBlock            *BlockHash
	LastFinalizedBlock     *BlockHash
	GenesisIndex           *GenesisIndex
	EraBlockHeight         *BlockHeight
	ReceiveTime            *Timestamp
	ArriveTime             *Timestamp
	SlotNumber             *Slot
	SlotTime               *Timestamp
	Baker                  *BakerId
	Finalized              bool
	TransactionCount       uint32
	TransactionsEnergyCost *Energy
	TransactionsSize       uint32
	StateHash              *StateHash
	ProtocolVersion        ProtocolVersion
}

// ProtocolVersion he different versions of the protocol.
type ProtocolVersion struct {
	Value int32
}

// StateHash hash of the state after some block. This is always 32 bytes long.
type StateHash struct {
	Value []byte
}

// BakerId the ID of a baker, which is the index of its account.
type BakerId struct {
	Value uint64
}

// Slot a number representing a slot for baking a block.
type Slot struct {
	Value uint64
}

// Timestamp unix timestamp in milliseconds.
type Timestamp struct {
	Value uint64
}

// GenesisIndex the number of chain restarts via a protocol update. An effected
// protocol update instruction might not change the protocol version
// specified in the previous field, but it always increments the genesis
// index.
type GenesisIndex struct {
	Value uint32
}

// AbsoluteBlockHeight this is the number of ancestors of a block
// since the genesis block. In particular, the chain genesis block has absolute
// height 0.
type AbsoluteBlockHeight struct {
	Value uint64
}

// BlockHeight the height of a block relative to the last genesis. This differs from the
// absolute block height in that it counts height from the last protocol update.
type BlockHeight struct {
	Value uint64
}

func convertBlockInfo(b *pb.BlockInfo) *BlockInfo {
	var hash, parentBlock, lastFinalizedBlock BlockHash

	copy(hash.Value[:], b.Hash.Value)
	copy(parentBlock.Value[:], b.ParentBlock.Value)
	copy(lastFinalizedBlock.Value[:], b.LastFinalizedBlock.Value)

	return &BlockInfo{
		Hash: &hash,
		Height: &AbsoluteBlockHeight{
			Value: b.Height.Value,
		},
		ParentBlock:        &parentBlock,
		LastFinalizedBlock: &lastFinalizedBlock,
		GenesisIndex: &GenesisIndex{
			Value: b.GenesisIndex.Value,
		},
		EraBlockHeight: &BlockHeight{
			Value: b.EraBlockHeight.Value,
		},
		ReceiveTime: &Timestamp{
			Value: b.ReceiveTime.Value,
		},
		ArriveTime: &Timestamp{
			Value: b.ArriveTime.Value,
		},
		SlotNumber: &Slot{
			Value: b.SlotNumber.Value,
		},
		SlotTime: &Timestamp{
			Value: b.SlotTime.Value,
		},
		Baker: &BakerId{
			Value: b.Baker.Value,
		},
		Finalized:        b.Finalized,
		TransactionCount: b.TransactionCount,
		TransactionsEnergyCost: &Energy{
			Value: b.TransactionsEnergyCost.Value,
		},
		TransactionsSize: b.TransactionsSize,
		StateHash: &StateHash{
			Value: b.StateHash.Value,
		},
		ProtocolVersion: ProtocolVersion{
			Value: int32(b.ProtocolVersion),
		},
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

// SequenceNumber a sequence number that determines the ordering of transactions from the
// account. The minimum sequence number is 1.
type SequenceNumber struct {
	Value uint64
}

// Energy is used to count exact execution cost.
// This cost is then converted to CCD amounts.
type Energy struct {
	Value uint64
}

// GivenEnergy helps handle the fixed costs and allows the user to focus only on the transaction
// specific ones. The most important case for this are smart contract initialisations and updates.
//
// An upper bound on the amount of energy to spend on a transaction. Transaction costs have two
// components, one is based on the size of the transaction and the number of signatures, and then
// there is a transaction specific one.
type GivenEnergy interface {
	isGivenEnergy()
}

// AbsoluteEnergy is an amount of Energy that will be exact used.
type AbsoluteEnergy Energy

// AddEnergy is an amount of Energy that will be added to the base amount.
// The base amount covers transaction size and signature checking.
type AddEnergy struct {
	Energy  Energy
	NumSigs uint32
}

func (*AbsoluteEnergy) isGivenEnergy() {}
func (*AddEnergy) isGivenEnergy()      {}

// CredentialType describes type of credential.
type CredentialType interface {
	isCredentialType()
}

// CredentialTypeInitial is a credential type that is submitted by the identity
// provider on behalf of the user. There is only one initial credential per identity.
type CredentialTypeInitial struct{}

// CredentialTypeNormal is one where the identity behind it is only known to
// the owner of the account, unless the anonymity revocation process was followed.
type CredentialTypeNormal struct{}

func (CredentialTypeInitial) isCredentialType() {}
func (CredentialTypeNormal) isCredentialType()  {}

// PayloadSize describes size of the transaction Payload. This is used to deserialize the Payload.
type PayloadSize struct {
	Value uint32
}

// TransactionTime specified as seconds since unix epoch.
type TransactionTime struct {
	Value uint64
}

// AccountTransactionPayload the payload for an account transaction.
type AccountTransactionPayload struct {
	Payload isAccountTransactionPayload
}

type isAccountTransactionPayload interface {
	Encode() *RawPayload
	isAccountTransactionPayload()
}

// DeployModule deploys a Wasm module with the given source.
type DeployModule struct {
	Payload *DeployModulePayload
}

func (DeployModule) isAccountTransactionPayload() {}
func (deployModule DeployModule) Encode() *RawPayload {
	return deployModule.Payload.Encode()
}

// VersionedModuleSource source bytes of a versioned smart contract module.
type VersionedModuleSource struct {
	Module isVersionedModuleSource
}

// Size returns size of
func (versionedModulePayload VersionedModuleSource) Size() int {
	switch m := versionedModulePayload.Module.(type) {
	case ModuleSourceV0:
		return len(m.Value)
	case ModuleSourceV1:
		return len(m.Value)
	}
	return 0
}

type isVersionedModuleSource interface {
	isVersionedModuleSource()
}

// ModuleSourceV0 v0.
type ModuleSourceV0 struct {
	Value []byte
}

// ModuleSourceV1 v1.
type ModuleSourceV1 struct {
	Value []byte
}

func (ModuleSourceV0) isVersionedModuleSource() {}

func (ModuleSourceV1) isVersionedModuleSource() {}

type moduleVersion uint8

const ModuleVersion0 moduleVersion = 0
const ModuleVersion1 moduleVersion = 1

type InitContract struct {
	Payload *InitContractPayload
}

func (InitContract) isAccountTransactionPayload() {}
func (initContract InitContract) Encode() *RawPayload {
	return initContract.Payload.Encode()
}

// InitName the init name of a smart contract function.
type InitName struct {
	Value string
}

// Amount an amount of microCCD.
type Amount struct {
	Value uint64
}

// Parameter to a smart contract initialization or invocation.
type Parameter struct {
	Value []byte
}

type UpdateContract struct {
	Payload *UpdateContractPayload
}

func (UpdateContract) isAccountTransactionPayload() {}
func (updateContract UpdateContract) Encode() *RawPayload {
	return updateContract.Payload.Encode()
}

// ReceiveName the reception name of a smart contract function. Expected format:
// `<contract_name>.<func_name>`. It must only consist of at most 100 ASCII
// alphanumeric or punctuation characters, and must contain a '.'.
type ReceiveName struct {
	Value string
}

type Transfer struct {
	Payload *TransferPayload
}

func (Transfer) isAccountTransactionPayload() {}
func (transfer Transfer) Encode() *RawPayload {
	return transfer.Payload.Encode()
}

type TransferWithMemo struct {
	Payload *TransferWithMemoPayload
}

func (TransferWithMemo) isAccountTransactionPayload() {}
func (transferWithMemo TransferWithMemo) Encode() *RawPayload {
	return transferWithMemo.Payload.Encode()
}

// Memo a memo which can be included as part of a transfer. Max size is 256 bytes.
type Memo struct {
	Value []byte
}

type RegisterData struct {
	Payload *RegisterDataPayload
}

func (RegisterData) isAccountTransactionPayload() {}
func (registerData RegisterData) Encode() *RawPayload {
	return registerData.Payload.Encode()
}

// RegisteredData data registered on the chain with a register data transaction.
type RegisteredData struct {
	Value []byte
}

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
type UpdateSequenceNumber struct {
	Value uint64
}

// UpdateInstructionPayload payload.
type UpdateInstructionPayload struct {
	Payload isUpdateInstructionPayload
}

type isUpdateInstructionPayload interface {
	isUpdateInstructionPayload()
}

func (RawPayload) isUpdateInstructionPayload() {}

func convertBlockItems(input []*pb.BlockItem) []*BlockItem {
	var result []*BlockItem

	for _, v := range input {
		var blockItem BlockItem
		var hash TransactionHash
		copy(hash.Value[:], v.Hash.Value)

		blockItem.Hash = &hash

		switch k := v.BlockItem.(type) {
		case *pb.BlockItem_AccountTransaction:
			signaturesMap := make(map[uint32]*AccountSignatureMap)

			for i, v := range k.AccountTransaction.Signature.Signatures {
				signatures := make(map[uint32]*Signature)

				for j, k := range v.Signatures {
					signatures[j] = &Signature{
						Value: k.Value,
					}
				}

				signaturesMap[i] = &AccountSignatureMap{
					Signatures: signatures,
				}
			}

			var accountTransactionPayload AccountTransactionPayload

			switch payload := k.AccountTransaction.Payload.Payload.(type) {
			case *pb.AccountTransactionPayload_RawPayload:
				accountTransactionPayload.Payload = &RawPayload{
					Value: payload.RawPayload,
				}
			case *pb.AccountTransactionPayload_DeployModule:
				switch dm := payload.DeployModule.Module.(type) {
				case *pb.VersionedModuleSource_V0:
					accountTransactionPayload.Payload = &DeployModule{
						Payload: &DeployModulePayload{
							DeployModule: &VersionedModuleSource{
								&ModuleSourceV0{
									Value: dm.V0.Value,
								}}}}
				case *pb.VersionedModuleSource_V1:
					accountTransactionPayload.Payload = &DeployModule{
						Payload: &DeployModulePayload{
							DeployModule: &VersionedModuleSource{
								&ModuleSourceV1{
									Value: dm.V1.Value,
								}}}}
				}
			case *pb.AccountTransactionPayload_InitContract:
				var moduleRef ModuleRef
				copy(moduleRef.Value[:], payload.InitContract.ModuleRef.Value)

				accountTransactionPayload.Payload = &InitContract{Payload: &InitContractPayload{
					Amount: &Amount{
						Value: payload.InitContract.Amount.Value,
					},
					ModuleRef: &moduleRef,
					InitName: &InitName{
						Value: payload.InitContract.InitName.Value,
					},
					Parameter: &Parameter{
						Value: payload.InitContract.Parameter.Value,
					},
				}}
			case *pb.AccountTransactionPayload_UpdateContract:
				accountTransactionPayload.Payload = &UpdateContract{
					Payload: &UpdateContractPayload{
						Amount: &Amount{
							Value: payload.UpdateContract.Amount.Value,
						},
						Address: &ContractAddress{
							Index:    payload.UpdateContract.Address.Index,
							Subindex: payload.UpdateContract.Address.Subindex,
						},
						ReceiveName: &ReceiveName{
							Value: payload.UpdateContract.ReceiveName.Value,
						},
						Parameter: &Parameter{
							Value: payload.UpdateContract.Parameter.Value,
						},
					}}
			case *pb.AccountTransactionPayload_Transfer:
				var accAddress AccountAddress
				copy(accAddress.Value[:], payload.Transfer.Receiver.Value)

				accountTransactionPayload.Payload = &Transfer{
					Payload: &TransferPayload{
						Amount: &Amount{
							Value: payload.Transfer.Amount.Value,
						},
						Receiver: &accAddress,
					}}
			case *pb.AccountTransactionPayload_TransferWithMemo:
				var receiver AccountAddress
				copy(receiver.Value[:], payload.TransferWithMemo.Receiver.Value)

				accountTransactionPayload.Payload = &TransferWithMemo{
					Payload: &TransferWithMemoPayload{
						Amount: &Amount{
							Value: payload.TransferWithMemo.Amount.Value,
						},
						Receiver: &receiver,
						Memo: &Memo{
							Value: payload.TransferWithMemo.Memo.Value,
						},
					}}
			case *pb.AccountTransactionPayload_RegisterData:
				accountTransactionPayload.Payload = &RegisterData{
					Payload: &RegisterDataPayload{
						Data: &RegisteredData{
							Value: payload.RegisterData.Value,
						},
					}}
			}

			var accountAddress AccountAddress
			copy(accountAddress.Value[:], k.AccountTransaction.Header.Sender.Value)

			blockItem.BlockItem = &AccountTransaction{
				Signature: &AccountTransactionSignature{
					Signatures: signaturesMap,
				},
				Header: &AccountTransactionHeader{
					Sender: &accountAddress,
					SequenceNumber: &SequenceNumber{
						Value: k.AccountTransaction.Header.SequenceNumber.Value,
					},
					EnergyAmount: &Energy{
						Value: k.AccountTransaction.Header.EnergyAmount.Value,
					},
					Expiry: &TransactionTime{
						Value: k.AccountTransaction.Header.Expiry.Value,
					},
				},
				Payload: &accountTransactionPayload,
			}
		case *pb.BlockItem_CredentialDeployment:
			var payload RawPayload
			switch v := k.CredentialDeployment.Payload.(type) {
			case *pb.CredentialDeployment_RawPayload:
				payload.Value = v.RawPayload
			}

			blockItem.BlockItem = &CredentialDeployment{
				MessageExpiry: &TransactionTime{
					Value: k.CredentialDeployment.MessageExpiry.Value,
				},
				Payload: &payload,
			}
		case *pb.BlockItem_UpdateInstruction:
			var signatureMap SignatureMap
			for i, v := range k.UpdateInstruction.Signatures.Signatures {
				signatureMap.Signatures[i] = &Signature{
					Value: v.Value,
				}
			}

			var updInstructionPayload = UpdateInstructionPayload{}

			switch t := k.UpdateInstruction.Payload.Payload.(type) {
			case *pb.UpdateInstructionPayload_RawPayload:
				updInstructionPayload.Payload = &RawPayload{
					Value: t.RawPayload,
				}
			}

			blockItem.BlockItem = &UpdateInstruction{
				Signatures: &signatureMap,
				Header: &UpdateInstructionHeader{
					SequenceNumber: &UpdateSequenceNumber{
						Value: k.UpdateInstruction.Header.SequenceNumber.Value,
					},
					EffectiveTime: &TransactionTime{
						Value: k.UpdateInstruction.Header.EffectiveTime.Value,
					},
					Timeout: &TransactionTime{
						Value: k.UpdateInstruction.Header.Timeout.Value,
					},
				},
				Payload: &updInstructionPayload,
			}
		}

		result = append(result, &blockItem)
	}

	return result
}
