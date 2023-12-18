package v2

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
	"github.com/btcsuite/btcutil/base58"
)

const (
	AccountAddressLength  = 32
	BlockHashLength       = 32
	TransactionHashLength = 32
	ModuleRefLength       = 32
	hundredThousand       = 100000
)

// WalletAccount an account imported from one of the supported export formats.
// This structure implements TransactionSigner and ExactSizeTransactionSigner, so it may be used for sending transactions.
// This structure does not have the encryption key for sending encrypted transfers, it only contains keys for signing transactions.
type WalletAccount struct {
	Address *AccountAddress
	Keys    *AccountKeys
}

// NewWalletAccount created new WalletAccount from AccountAddress and one KeyPair.
func NewWalletAccount(accountAddress AccountAddress, keyPair KeyPair) *WalletAccount {
	keyPairs := make(map[KeyIndex]*KeyPair, 1)
	keyPairs[0] = &keyPair

	accountKeys := make(map[CredentialIndex]*CredentialData, 1)
	accountKeys[0] = &CredentialData{
		Keys:      keyPairs,
		Threshold: SignatureThreshold{Value: 1},
	}

	return &WalletAccount{
		Address: &accountAddress,
		Keys: &AccountKeys{
			Keys:      accountKeys,
			Threshold: AccountThreshold{Value: 1},
		},
	}
}

// NewWalletAccountFromFile created new WalletAccount from `<account_address>.export` file.
func NewWalletAccountFromFile(pathToFile string) (*WalletAccount, error) {
	data := struct {
		Value struct {
			AccountKeys struct {
				Keys map[string]struct {
					Keys map[string]struct {
						SingKey   string `json:"signKey"`
						VerifyKey string `json:"verifyKey"`
					} `json:"keys"`
					Threshold int `json:"threshold"`
				} `json:"keys"`
				Threshold int `json:"threshold"`
			} `json:"accountKeys"`
			Address string `json:"address"`
		} `json:"value"`
	}{}

	file, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	accountAddress, err := AccountAddressFromString(data.Value.Address)
	if err != nil {
		return nil, err
	}

	credentialData := make(map[CredentialIndex]*CredentialData, data.Value.AccountKeys.Threshold)
	for credIndex := 0; credIndex < data.Value.AccountKeys.Threshold; credIndex++ {
		keyPairs := make(map[KeyIndex]*KeyPair, data.Value.AccountKeys.Keys[strconv.Itoa(credIndex)].Threshold)

		credKeys := data.Value.AccountKeys.Keys[strconv.Itoa(credIndex)]
		for keyIndex := 0; keyIndex < credKeys.Threshold; keyIndex++ {
			keyPair := credKeys.Keys[strconv.Itoa(keyIndex)]

			singKey, err := hex.DecodeString(keyPair.SingKey)
			if err != nil {
				return nil, err
			}

			verifyKey, err := hex.DecodeString(keyPair.VerifyKey)
			if err != nil {
				return nil, err
			}

			keyPairs[KeyIndex(keyIndex)], err = NewKeyPairFromSignKeyAndVerifyKey(singKey, verifyKey)
			if err != nil {
				return nil, err
			}
		}

		credentialData[CredentialIndex(credIndex)] = &CredentialData{
			Keys:      keyPairs,
			Threshold: SignatureThreshold{Value: uint8(len(keyPairs))},
		}
	}

	return &WalletAccount{
		Address: &accountAddress,
		Keys: &AccountKeys{
			Keys:      credentialData,
			Threshold: AccountThreshold{Value: uint8(len(credentialData))},
		},
	}, nil
}

// AddKeyPair adds KeyPair to WalletAccount and updates threshold fields.
func (walletAccount *WalletAccount) AddKeyPair(pair KeyPair) {
	keyPairs := walletAccount.Keys.Keys[CredentialIndex(walletAccount.Keys.Threshold.Value-1)]
	keyPairs.Keys[KeyIndex(keyPairs.Threshold.Value)] = &pair
	keyPairs.Threshold.Value++
	walletAccount.Keys.Keys[CredentialIndex(walletAccount.Keys.Threshold.Value-1)] = keyPairs
}

// SignTransactionHash returns signed TransactionHash.
func (walletAccount *WalletAccount) SignTransactionHash(hashToSign *TransactionHash) (*AccountTransactionSignature, error) {
	if walletAccount.Address == nil || walletAccount.Keys == nil {
		return nil, errors.New("'AccountAddress' or 'Keys' field is not initialized or empty")
	}
	if walletAccount.Keys.Keys == nil || len(walletAccount.Keys.Keys) == 0 {
		return nil, errors.New("'WalletAccount.Keys' is not initialized or empty")
	}

	signaturesMap := make(map[uint8]*AccountSignatureMap, int(walletAccount.Keys.Threshold.Value))
	for credIdx, credData := range walletAccount.Keys.Keys {
		if credData.Keys == nil || len(credData.Keys) == 0 {
			return nil, errors.New("'WalletAccount.Keys.Keys[ " + strconv.Itoa(int(credIdx)) + "].Keys' is not initialized or empty")
		}

		signatures := make(map[uint8]*Signature, int(credData.Threshold.Value))
		for keyIdx, keyPair := range credData.Keys {
			signature := keyPair.Sign(hashToSign.Value[:])
			signatures[uint8(keyIdx)] = &signature
		}

		signaturesMap[uint8(credIdx)] = &AccountSignatureMap{Signatures: signatures}
	}

	return &AccountTransactionSignature{Signatures: signaturesMap}, nil
}

// NumberOfKeys returns number of signing keys.
func (walletAccount *WalletAccount) NumberOfKeys() uint32 {
	var sum uint32 = 0
	for _, credData := range walletAccount.Keys.Keys {
		sum += uint32(len(credData.Keys))
	}

	return sum
}

// AccountKeys all account keys indexed by credentials.
type AccountKeys struct {
	Keys      map[CredentialIndex]*CredentialData
	Threshold AccountThreshold
}

// CredentialIndex describes index of the credential that is to be used.
type CredentialIndex uint8

// CredentialData describes credential data needed by the account holder to generate proofs to deploy
// the credential object. This contains all the keys on the credential at the moment of its deployment.
// If this creates the account then the account starts with exactly these keys.
type CredentialData struct {
	Keys      map[KeyIndex]*KeyPair
	Threshold SignatureThreshold
}

// KeyIndex describes index of an account key that is to be used.
type KeyIndex uint8

// KeyPair describes ed25519 key pair.
type KeyPair struct {
	// secret describes `signKey`.
	secret ed25519.PrivateKey
	// public describes `verifyKey`.
	public ed25519.PublicKey
}

// NewKeyPairFromSignKey creates new KeyPair from `singKey`.
// You can find this key in `Private Key` field in wallet settings -> `export private key` (key is in hex encoding),
// or copy from export file from filed `signKey` (key is in hex encoding).
func NewKeyPairFromSignKey(signKey []byte) (*KeyPair, error) {
	if signKey == nil || len(signKey) != ed25519.SeedSize {
		return nil, errors.New("sign key should be " + strconv.Itoa(ed25519.SeedSize) + " bytes long")
	}

	privateKey := ed25519.NewKeyFromSeed(signKey)
	return &KeyPair{secret: privateKey.Seed(), public: privateKey.Public().(ed25519.PublicKey)}, nil
}

// NewKeyPairFromSignKeyAndVerifyKey creates new KeyPair from `singKey` and `verifyKey.
// You can find these keys in export file in fields `signKey` and `verifyKey` (keys are in hex encoding).
func NewKeyPairFromSignKeyAndVerifyKey(signKey, verifyKey []byte) (*KeyPair, error) {
	if signKey == nil || len(signKey) != ed25519.SeedSize {
		return nil, errors.New("sign key should be " + strconv.Itoa(ed25519.SeedSize) + " bytes long")
	}

	if verifyKey == nil || len(verifyKey) != ed25519.PublicKeySize {
		return nil, errors.New("verify key should be " + strconv.Itoa(ed25519.PublicKeySize) + " bytes long")
	}

	return &KeyPair{secret: signKey, public: verifyKey}, nil
}

// Secret returns `signKey`.
func (keyPair *KeyPair) Secret() ed25519.PrivateKey {
	return keyPair.secret
}

// Public return `verifyKey`.
func (keyPair *KeyPair) Public() ed25519.PublicKey {
	return keyPair.public
}

// Sign signs the message with private key and returns a signature.
func (keyPair *KeyPair) Sign(msg []byte) Signature {
	privateKey := append(keyPair.Secret(), keyPair.Public()...)
	return Signature{Value: ed25519.Sign(privateKey, msg)}
}

// AccountThreshold describes the minimum number of credentials that need to sign any transaction coming
// from an associated account. The values of this type must maintain the property that they are not 0.
type AccountThreshold struct {
	Value uint8
}

// SignatureThreshold threshold for the number of signatures required.
// The values of this type must maintain the property that they are not 0.
type SignatureThreshold struct {
	Value uint8
}

type isAddress interface {
	isAddress()
}

// AccountAddress an address of an account.
type AccountAddress struct {
	Value [AccountAddressLength]byte
}

func (a *AccountAddress) isAddress() {}

// ToBase58 encodes account address to string.
func (a *AccountAddress) ToBase58() string {
	return base58.CheckEncode(a.Value[:], 1)
}

// AccountAddressFromString decodes string to account.
func AccountAddressFromString(s string) (AccountAddress, error) {
	b, _, err := base58.CheckDecode(s)
	if err != nil {
		return AccountAddress{}, err
	}

	return AccountAddressFromBytes(b)
}

// AccountAddressFromBytes creates account address from given bytes.
func AccountAddressFromBytes(b []byte) (AccountAddress, error) {
	if len(b) != 32 {
		return AccountAddress{}, errors.New("account address must be exactly 32 bytes")
	}

	var accountAddress AccountAddress
	copy(accountAddress.Value[AccountAddressLength-len(b):], b)
	return accountAddress, nil
}

// BlockHash hash of a block. This is always 32 bytes long.
type BlockHash struct {
	Value [BlockHashLength]byte
}

// BlockHashFromBytes creates a BlockHash from given []byte. Length of given []byte must be excactly 32 bytes.
func BlockHashFromBytes(b []byte) (BlockHash, error) {
	if len(b) != BlockHashLength {
		return BlockHash{}, errors.New("BlockHash must be excactly 32 bytes")
	}
	var blockHash BlockHash
	copy(blockHash.Value[:], b)
	return blockHash, nil
}

// Parses *pb.BlockHash to BlockHash
func parseBlockHash(h *pb.BlockHash) (BlockHash, error) {
	return BlockHashFromBytes(h.Value)
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

// BlockInfo information about given block, contains height, timings, transaction count, state, etc.
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

// ProtocolVersion the different versions of the protocol.
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

// Parses *pb.BakerId to BakerId
func parseBakerId(b *pb.BakerId) BakerId {
	return BakerId{Value: b.Value}
}

// Slot a number representing a slot for baking a block.
type Slot struct {
	Value uint64
}

// Timestamp unix timestamp in milliseconds.
type Timestamp struct {
	Value uint64
}

// Parses *pb.Timestamp to Timestamp
func parseTimestamp(t *pb.Timestamp) Timestamp {
	return Timestamp{Value: t.Value}
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

	var slotTime, slotNumber uint64
	if b.SlotTime != nil {
		slotTime = b.SlotTime.Value
	}
	if b.SlotNumber != nil {
		slotNumber = b.SlotNumber.Value
	}

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
			Value: slotNumber,
		},
		SlotTime: &Timestamp{
			Value: slotTime,
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

func (c *ContractAddress) isAddress() {}

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
type GivenEnergy struct {
	Energy isGivenEnergy
}

type isGivenEnergy interface {
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
type CredentialType struct {
	Type isCredentialType
}

type isCredentialType interface {
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

// Size returns size of module bytes.
func (versionedModulePayload VersionedModuleSource) Size() int {
	switch m := versionedModulePayload.Module.(type) {
	case ModuleSourceV0:
		return len(m.Value)
	case *ModuleSourceV0:
		return len(m.Value)
	case ModuleSourceV1:
		return len(m.Value)
	case *ModuleSourceV1:
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

// InitContract contains data needed to initialize a smart contract.
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

func parseAmount(a *pb.Amount) Amount {
	return Amount{Value: a.Value}
}

// Parameter to a smart contract initialization or invocation.
type Parameter struct {
	Value []byte
}

// UpdateContract updates a smart contract instance by invoking a specific function.
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

// Transfer transfers CCD to an account.
type Transfer struct {
	Payload *TransferPayload
}

func (Transfer) isAccountTransactionPayload() {}
func (transfer Transfer) Encode() *RawPayload {
	return transfer.Payload.Encode()
}

// TransferWithMemo payload of a transfer between two accounts with a memo.
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

// RegisterData registers the given data on the chain.
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

func ConvertBlockItems(input []*pb.BlockItem) []*BlockItem {
	var result []*BlockItem

	for _, v := range input {
		var blockItem BlockItem
		var hash TransactionHash
		copy(hash.Value[:], v.Hash.Value)

		blockItem.Hash = &hash

		switch k := v.BlockItem.(type) {
		case *pb.BlockItem_AccountTransaction:
			signaturesMap := make(map[uint8]*AccountSignatureMap)

			for i, v := range k.AccountTransaction.Signature.Signatures {
				signatures := make(map[uint8]*Signature)

				for j, k := range v.Signatures {
					signatures[uint8(j)] = &Signature{
						Value: k.Value,
					}
				}

				signaturesMap[uint8(i)] = &AccountSignatureMap{
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
			signatureMap := make(map[uint32]*Signature)
			for i, v := range k.UpdateInstruction.Signatures.Signatures {
				signatureMap[i] = &Signature{
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
				Signatures: &SignatureMap{Signatures: signatureMap},
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

// A chain epoch.
type Epoch struct {
	Value uint64
}

// Input to queries which take an epoch as a parameter.
type isEpochRequest interface {
	isEpochRequest()
}

// Query for the epoch of a specified block.
type EpochRequestBlockHash struct {
	// The block to query at.
	BlockHash isBlockHashInput
}

func (EpochRequestBlockHash) isEpochRequest() {}

// Request an epoch by number at a given genesis index.
type EpochRequestRelativeEpoch struct {
	// The genesis index to query at. The query is restricted to this genesis idex, and
	// will not return results for other indices even if the epoch number is out of bounds.
	GenesisIndex GenesisIndex
	// The epoch number to query at.
	Epoch Epoch
}

func (EpochRequestRelativeEpoch) isEpochRequest() {}

func convertEpochRequest(req isEpochRequest) (_ *pb.EpochRequest) {
	var res *pb.EpochRequest
	switch v := req.(type) {
	case EpochRequestBlockHash:
		res = &pb.EpochRequest{
			EpochRequestInput: &pb.EpochRequest_BlockHash{
				BlockHash: convertBlockHashInput(v.BlockHash),
			},
		}
	case EpochRequestRelativeEpoch:
		res = &pb.EpochRequest{
			EpochRequestInput: &pb.EpochRequest_RelativeEpoch_{
				RelativeEpoch: &pb.EpochRequest_RelativeEpoch{
					GenesisIndex: &pb.GenesisIndex{
						Value: v.GenesisIndex.Value,
					},
					Epoch: &pb.Epoch{
						Value: v.Epoch.Value,
					},
				},
			},
		}
	}

	return res
}

// Return type of GetBakersRewardPeriod. Parses the returned *pb.BakerRewardPeriodInfo to BakerRewardPeriodInfo when Recv() is called.
type BakerRewardPeriodInfoStream struct {
	stream pb.Queries_GetBakersRewardPeriodClient
}

// Recv retrieves the next BakerRewardPeriodInfo.
func (s *BakerRewardPeriodInfoStream) Recv() (BakerRewardPeriodInfo, error) {
	info, err := s.stream.Recv()
	if err != nil {
		return BakerRewardPeriodInfo{}, err
	}
	return parseBakerRewardPeriodInfo(info)
}

// Information about a particular baker with respect to the current reward period.
type BakerRewardPeriodInfo struct {

	// The baker id and public keys for the baker.
	Baker BakerInfo

	// The effective stake of the baker for the consensus protocol.
	// The returned amount accounts for delegation, capital bounds and leverage bounds.
	EffectiveStake Amount

	// The effective commission rate for the baker that applies for the reward period.
	CommissionRates CommissionRates

	// The amount staked by the baker itself.
	EquityCapital Amount

	// The total amount of capital delegated to this baker pool.
	DelegatedCapital Amount

	// Whether the baker is a finalizer or not.
	IsFinalizer bool
}

// Parses *pb.BakerRewardPeriodInfo to BakerRewardPeriodInfo
func parseBakerRewardPeriodInfo(b *pb.BakerRewardPeriodInfo) (BakerRewardPeriodInfo, error) {
	baker := parseBakerInfo(b.GetBaker())
	effectiveStake := parseAmount(b.GetEffectiveStake())
	CommissionRates, err := parseCommissionRates(b.CommissionRates)
	if err != nil {
		return BakerRewardPeriodInfo{}, err
	}
	equityCapital := parseAmount(b.EquityCapital)
	delegatedCapital := parseAmount(b.DelegatedCapital)
	isFinalizer := b.IsFinalizer
	return BakerRewardPeriodInfo{
		Baker:            baker,
		EffectiveStake:   effectiveStake,
		CommissionRates:  CommissionRates,
		EquityCapital:    equityCapital,
		DelegatedCapital: delegatedCapital,
		IsFinalizer:      isFinalizer,
	}, nil
}

// Information about a baker.
type BakerInfo struct {

	// Identity of the baker. This is actually the account index of the account controlling the baker.
	BakerId BakerId

	// Baker's public key used to check whether they won the lottery or not.
	ElectionKey BakerElectionVerifyKey

	// Baker's public key used to check that they are indeed the ones who produced the block.
	SignatureKey BakerSignatureVerifyKey

	// Baker's public key used to check signatures on finalization records.
	// This is only used if the baker has sufficient stake to participate in finalization.
	AggregationKey BakerAggregationVerifyKey
}

// Parses *pb.BakerInfo to BakerInfo
func parseBakerInfo(b *pb.BakerInfo) BakerInfo {
	bakerId := parseBakerId(b.BakerId)
	electionKey := parseBakerElectionVerifyKey(b.ElectionKey)
	signatureKey := parseBakerSignatureVerifyKey(b.SignatureKey)
	aggregationKey := parseBakerAggregationVerifyKey(b.AggregationKey)
	return BakerInfo{BakerId: bakerId, ElectionKey: electionKey, SignatureKey: signatureKey, AggregationKey: aggregationKey}

}

// Baker's public key used to check whether they won the lottery or not.
type BakerElectionVerifyKey struct {
	Value []byte
}

// Parses *pb.BakerElectionVerifyKey to BakerElectionVerifyKey
func parseBakerElectionVerifyKey(k *pb.BakerElectionVerifyKey) BakerElectionVerifyKey {
	return BakerElectionVerifyKey{Value: k.Value}
}

// Baker's public key used to check that they are indeed the ones who produced the block.
type BakerSignatureVerifyKey struct {
	Value []byte
}

// Parses *pb.BakerSignatureVerifyKey to BakerSignatureVerifyKey
func parseBakerSignatureVerifyKey(k *pb.BakerSignatureVerifyKey) BakerSignatureVerifyKey {
	return BakerSignatureVerifyKey{Value: k.Value}
}

// Baker's public key used to check signatures on finalization records.
// This is only used if the baker has sufficient stake to participate in finalization.
type BakerAggregationVerifyKey struct {
	Value []byte
}

// Parses *pb.BakerAggregationVerifyKey to BakerAggregationVerifyKey
func parseBakerAggregationVerifyKey(k *pb.BakerAggregationVerifyKey) BakerAggregationVerifyKey {
	return BakerAggregationVerifyKey{Value: k.Value}
}

// Distribution of the rewards for the particular pool.
type CommissionRates struct {

	// Fraction of finalization rewards charged by the pool owner.
	Finalization AmountFraction

	// Fraction of baking rewards charged by the pool owner.
	Baking AmountFraction

	// Fraction of transaction rewards charged by the pool owner.
	Transaction AmountFraction
}

// Parses *pb.CommissionRates to *CommissionRates.
func parseCommissionRates(cr *pb.CommissionRates) (CommissionRates, error) {
	finalization, err := parseAmountFraction(cr.Finalization)
	if err != nil {
		return CommissionRates{}, errors.New("Error parsing CommissionRates: " + err.Error())
	}

	baking, err := parseAmountFraction(cr.Baking)
	if err != nil {
		return CommissionRates{}, errors.New("Error parsing CommissionRates: " + err.Error())
	}

	transaction, err := parseAmountFraction(cr.Transaction)
	if err != nil {
		return CommissionRates{}, errors.New("Error parsing CommissionRates: " + err.Error())
	}
	return CommissionRates{Finalization: finalization, Baking: baking, Transaction: transaction}, nil
}

// A fraction of an amount with a precision of 1/100_000
type AmountFraction struct {
	// Must not exceed 100_000
	partsPerHundredThousand uint32
}

// GetValue returns the value of the AmountFraction, e.g. 'partsPerHundredThousand/100_000'.
func (a *AmountFraction) GetValue() uint32 {
	return a.partsPerHundredThousand / hundredThousand
}

// AmountFractionFromUInt32 constructs an AmountFraction from a uint32 value. The value must not exceed 100_000.
func AmountFractionFromUInt32(value uint32) (AmountFraction, error) {
	if value > hundredThousand {
		return AmountFraction{}, errors.New("PartsPerHundredThousand must not exceed 100_000")
	}
	return AmountFraction{partsPerHundredThousand: value}, nil
}

// Parses *pb.AmountFraction to AmountFraction.
func parseAmountFraction(a *pb.AmountFraction) (AmountFraction, error) {
	res, err := AmountFractionFromUInt32(a.PartsPerHundredThousand)
	if err != nil {
		return AmountFraction{}, errors.New("Error parsing AmountFraction: " + err.Error())
	}
	return res, nil
}
