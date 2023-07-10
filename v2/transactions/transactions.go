package transactions

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// TransactionHeaderSize describes size of a transaction Header. This is currently always 60 Bytes.
// Future chain updates might revise this, but this is a big change so this
// is expected to change seldomly.
const TransactionHeaderSize uint64 = 60 // 32 + 8 + 8 + 4 + 8.

// TransactionHeader is a Header of an account transaction that contains basic data to check whether
// the sender and the transaction is valid.
type TransactionHeader struct {
	// Sender account of the transaction.
	Sender types.AccountAddress
	// Sequence number of the transaction.
	Nonce types.Nonce
	// Maximum amount of energy the transaction can take to execute.
	EnergyAmount types.Energy
	// Size of the transaction Payload. This is used to deserialize the Payload.
	PayloadSize types.PayloadSize
	// Latest time the transaction can be included in a block.
	Expiry types.TransactionTime
}

// Serialize returns serialized TransactionHeader.
func (transactionHeader *TransactionHeader) Serialize() []byte {
	buf := make([]byte, 0, TransactionHeaderSize)
	buf = append(buf, transactionHeader.Sender...)
	buf = binary.BigEndian.AppendUint64(buf, uint64(transactionHeader.Nonce))
	buf = binary.BigEndian.AppendUint64(buf, uint64(transactionHeader.EnergyAmount))
	buf = binary.BigEndian.AppendUint32(buf, uint32(transactionHeader.PayloadSize))
	buf = binary.BigEndian.AppendUint64(buf, uint64(transactionHeader.Expiry))

	return buf
}

// PreAccountTransaction describes account transaction before signing.
type PreAccountTransaction struct {
	Header     *TransactionHeader
	Payload    payloads.Payload
	Encoded    payloads.EncodedPayload
	HashToSign types.TransactionSignHash
}

// Sign signs PreAccountTransaction with TransactionSigner and returns AccountTransaction.
func (preAccountTransaction *PreAccountTransaction) Sign(signer TransactionSigner) (*AccountTransaction, error) {
	return signTransaction(signer, preAccountTransaction.Header, preAccountTransaction.Payload)
}

// Serialize returns serialized PreAccountTransaction.
func (preAccountTransaction *PreAccountTransaction) Serialize() []byte {
	buf := make([]byte, 0, int(TransactionHeaderSize)+int(preAccountTransaction.Encoded.Size()))
	buf = append(buf, preAccountTransaction.Header.Serialize()...)
	buf = append(buf, preAccountTransaction.Encoded.Serialize()...)

	return buf
}

// Deserialize parses from bytes PreAccountTransaction.
func (preAccountTransaction *PreAccountTransaction) Deserialize(source []byte) (err error) {
	if len(source) < int(TransactionHeaderSize) {
		return errors.New("could not deserialize PreAccountTransaction: invalid length")
	}

	preAccountTransaction.Header = &TransactionHeader{
		Sender:       source[:32],
		Nonce:        types.Nonce(binary.BigEndian.Uint64(source[32:40])),
		EnergyAmount: types.Energy(binary.BigEndian.Uint64(source[40:48])),
		PayloadSize:  types.PayloadSize(binary.BigEndian.Uint32(source[48:52])),
		Expiry:       types.TransactionTime(binary.BigEndian.Uint64(source[52:TransactionHeaderSize])),
	}

	if len(source) < int(TransactionHeaderSize)+int(preAccountTransaction.Header.PayloadSize) {
		return errors.New("could not deserialize PreAccountTransaction: invalid length")
	}

	preAccountTransaction.Encoded = source[TransactionHeaderSize:]
	preAccountTransaction.Payload, err = preAccountTransaction.Encoded.Decode()
	if err != nil {
		return errors.New(fmt.Sprintf("could not decode Encoded Payload: %v", err))
	}

	preAccountTransaction.HashToSign = ComputeTransactionSignHash(preAccountTransaction.Header, preAccountTransaction.Encoded)

	return nil
}

// TransactionSigner is an interface for signing transactions.
type TransactionSigner interface {
	// SignTransactionHash signs transaction hash and returns signatures in TransactionSignature type.
	SignTransactionHash(hashToSign types.TransactionSignHash) (types.TransactionSignature, error)
}

// ExactSizeTransactionSigner describes TransactionSigner with ability to return number of signers.
type ExactSizeTransactionSigner interface {
	TransactionSigner
	// NumberOfKeys returns number of signers.
	NumberOfKeys() uint32
}

// AccountTransaction describes signed account transaction.
type AccountTransaction struct {
	Signature types.TransactionSignature
	Header    TransactionHeader
	Payload   payloads.PayloadLike
}

// signTransaction signs the Header and Payload, construct the transaction, and return it.
func signTransaction(signer TransactionSigner, header *TransactionHeader, payload payloads.PayloadLike) (*AccountTransaction, error) {
	hashToSign := ComputeTransactionSignHash(header, payload)
	signature, err := signer.SignTransactionHash(hashToSign)
	if err != nil {
		return &AccountTransaction{}, err
	}

	return &AccountTransaction{
		Signature: signature,
		Header:    *header,
		Payload:   payload,
	}, nil
}

// ComputeTransactionSignHash computes the transaction sign hash from an EncodedPayload and Header.
func ComputeTransactionSignHash(header *TransactionHeader, payload payloads.PayloadLike) types.TransactionSignHash {
	encodedPayload := payload.Encode()
	buf := make([]byte, 0, int(TransactionHeaderSize)+len(encodedPayload))
	buf = append(buf, header.Serialize()...)
	buf = append(buf, encodedPayload...)

	return sha256.Sum256(buf)
}
