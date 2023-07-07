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
	Height                 uint64
	ParentBlock            BlockHash
	LastFinalizedBlock     BlockHash
	GenesisIndex           uint32
	EraBlockHeight         uint64
	ReceiveTime            uint64
	ArriveTime             uint64
	SlotNumber             uint64
	SlotTime               uint64
	Baker                  uint64
	Finalized              bool
	TransactionCount       uint32
	TransactionsEnergyCost uint64
	TransactionsSize       uint32
	StateHash              []byte
	ProtocolVersion        int32
}

func convertBlockInfo(b *pb.BlockInfo) BlockInfo {
	var hash, parentBlock, lastFinalizedBlock BlockHash

	copy(hash[:], b.Hash.Value)
	copy(parentBlock[:], b.ParentBlock.Value)
	copy(lastFinalizedBlock[:], b.LastFinalizedBlock.Value)

	return BlockInfo{
		Hash:                   hash,
		Height:                 b.Height.Value,
		ParentBlock:            parentBlock,
		LastFinalizedBlock:     lastFinalizedBlock,
		GenesisIndex:           b.GenesisIndex.Value,
		EraBlockHeight:         b.EraBlockHeight.Value,
		ReceiveTime:            b.ReceiveTime.Value,
		ArriveTime:             b.ArriveTime.Value,
		SlotNumber:             b.SlotNumber.Value,
		SlotTime:               b.SlotTime.Value,
		Baker:                  b.Baker.Value,
		Finalized:              b.Finalized,
		TransactionCount:       b.TransactionCount,
		TransactionsEnergyCost: b.TransactionsEnergyCost.Value,
		TransactionsSize:       b.TransactionsSize,
		StateHash:              b.StateHash.Value,
		ProtocolVersion:        int32(b.ProtocolVersion),
	}
}

// ContractAddress address of a smart contract instance.
type ContractAddress struct {
	Index    uint64
	Subindex uint64
}
