package v2

import (
	"github.com/fxamacker/cbor/v2"
)

const (
	// PltTransfer
	// Additional cost of a PLT transfer.
	PltTransfer uint64 = 100

	// PltMint
	// Additional cost of a PLT mint.
	PltMint uint64 = 50

	// PltBurn
	// Additional cost of a PLT burn.
	PltBurn uint64 = 50

	// PltListUpdate
	// Additional cost of a PLT allow or deny list update.
	PltListUpdate uint64 = 50

	// PltPause
	// Additional cost of a PLT pausing.
	PltPause uint64 = 50

	// PltOperationsTransactions
	// Additional cost of a transaction consisting of protocol level token
	// operations.
	PltOperationsTransactions uint64 = 300

	// CCDCoinInfo 919 is the BIP-0044 SLIP coin type code for CCD.
	CCDCoinInfo CoinInfo = 919
	// CoinInfoTag CoinInfo Tag
	CoinInfoTag uint64 = 40305
	// EncodedCBORDataTag Encoded CBOR data tag
	EncodedCBORDataTag = 24
)

// TokenOperations enum.
type TokenOperations []TokenOperation

// TxnEnergy returns the total energy cost of all operations in the list.
func (ops TokenOperations) TxnEnergy() *Energy {
	total := PltOperationsTransactions
	for _, op := range ops {
		if op != nil {
			total += op.TxnEnergy().Value
		}
	}
	return &Energy{Value: total}
}

type TokenOperation interface {
	TxnEnergy() *Energy
	EncodeCBOR() ([]byte, error)
}

type TransferOperation struct {
	Transfer TokenTransfer `cbor:"transfer"`
}

func (*TransferOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltTransfer}
}
func (op *TransferOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type MintOperation struct {
	Details TokenSupplyUpdateDetails `cbor:"mint"`
}

func (*MintOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltMint}
}
func (op *MintOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type BurnOperation struct {
	Details TokenSupplyUpdateDetails `cbor:"burn"`
}

func (*BurnOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltBurn}
}
func (op *BurnOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type AddAllowListOperation struct {
	Details TokenListUpdateDetails `cbor:"add_allow_list"`
}

func (*AddAllowListOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltListUpdate}
}
func (op *AddAllowListOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type RemoveAllowListOperation struct {
	Details TokenListUpdateDetails `cbor:"remove_allow_list"`
}

func (*RemoveAllowListOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltListUpdate}
}
func (op *RemoveAllowListOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type AddDenyListOperation struct {
	Details TokenListUpdateDetails `cbor:"add_deny_list"`
}

func (*AddDenyListOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltListUpdate}
}
func (op *AddDenyListOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type RemoveDenyListOperation struct {
	Details TokenListUpdateDetails `cbor:"remove_deny_list"`
}

func (*RemoveDenyListOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltListUpdate}
}
func (op *RemoveDenyListOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type PauseOperation struct {
	Details TokenPauseDetails `cbor:"pause"`
}

func (*PauseOperation) TxnEnergy() *Energy {
	return &Energy{Value: PltPause}
}
func (op *PauseOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type UnpauseOperation struct {
	Details TokenPauseDetails `cbor:"unpause"`
}

func (*UnpauseOperation) TxnEnergy() *Energy {
	// TODO: Should have some fixed energy cost.
	return &Energy{}
}
func (op *UnpauseOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type UnknownOperation struct{}

func (*UnknownOperation) TxnEnergy() *Energy {
	return &Energy{}
}
func (op *UnknownOperation) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(op)
}

type TokenSupplyUpdateDetails struct {
	Amount TokenAmount `cbor:"amount"`
}

type TokenPauseDetails struct{}
type TokenListUpdateDetails struct {
	Target CborHolderAccount
}
type TokenTransfer struct {
	Amount    TokenAmount       `cbor:"amount"`
	Recipient CborHolderAccount `cbor:"recipient"`
	Memo      CborMemo          `cbor:"memo,omitempty"`
}

type CborHolderAccount struct {
	CoinInfo *CoinInfo      `cbor:"1,keyasint,omitempty"`
	Address  AccountAddress `cbor:"3,keyasint"`
}

func (c CoinInfo) MarshalCBOR() ([]byte, error) {
	enc, _ := cbor.EncOptions{}.EncMode()
	return enc.Marshal(cbor.Tag{
		Number: CoinInfoTag,
		Content: map[uint64]uint16{
			1: uint16(c),
		},
	})
}
func (h *CborHolderAccount) MarshalCBOR() ([]byte, error) {
	m := make(map[any]any)
	if h.CoinInfo != nil {
		m[uint64(1)] = *h.CoinInfo
	}
	m[uint64(3)] = h.Address.Value[:]

	enc, _ := cbor.EncOptions{}.EncMode()
	return enc.Marshal(cbor.Tag{
		Number:  AccountHolderTag,
		Content: m,
	})
}

type CoinInfo uint64 // constant ccd 919

type CborMemo struct {
	Raw  *Memo
	Cbor *Memo
}

func (m *CborMemo) MarshalCBOR() ([]byte, error) {
	if m.Raw == nil || m.Raw.Value == nil {
		enc, _ := cbor.EncOptions{}.EncMode()
		return enc.Marshal(cbor.Tag{
			Number:  EncodedCBORDataTag,
			Content: []byte{},
		})
	}
	enc, _ := cbor.EncOptions{}.EncMode()
	return enc.Marshal(cbor.Tag{
		Number:  EncodedCBORDataTag,
		Content: m.Raw.Value,
	})
}
