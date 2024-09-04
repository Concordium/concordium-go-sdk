package v2

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// TransactionHeaderSize describes size of a transaction Header. This is currently always 60 Bytes.
// Future chain updates might revise this, but this is a big change so this is expected to change seldom.
const TransactionHeaderSize uint64 = 60 // 32 + 8 + 8 + 4 + 8.

// AccountTransactionHeader header of an account transaction that contains basic data to check whether
// the sender and the transaction are valid. The header is shared by all transaction types.
type AccountTransactionHeader struct {
	// Sender account of the transaction.
	Sender *AccountAddress
	// Sequence number of the transaction.
	SequenceNumber *SequenceNumber
	// Maximum amount of energy the transaction can take to execute.
	EnergyAmount *Energy
	// Size of the transaction Payload. This is used to deserialize the Payload.
	PayloadSize *PayloadSize
	// Latest time the transaction can be included in a block.
	Expiry *TransactionTime
}

// Serialize returns serialized AccountTransactionHeader.
func (transactionHeader *AccountTransactionHeader) Serialize() []byte {
	buf := make([]byte, 0, TransactionHeaderSize)
	buf = append(buf, transactionHeader.Sender.Value[:]...)
	buf = binary.BigEndian.AppendUint64(buf, transactionHeader.SequenceNumber.Value)
	buf = binary.BigEndian.AppendUint64(buf, transactionHeader.EnergyAmount.Value)
	buf = binary.BigEndian.AppendUint32(buf, transactionHeader.PayloadSize.Value)
	buf = binary.BigEndian.AppendUint64(buf, transactionHeader.Expiry.Value)

	return buf
}

// PreAccountTransaction describes account transaction before signing.
type PreAccountTransaction struct {
	Header     *AccountTransactionHeader
	Payload    *AccountTransactionPayload
	Encoded    *RawPayload
	HashToSign *TransactionHash
}

// Sign signs PreAccountTransaction with TransactionSigner and returns AccountTransaction.
func (preAccountTransaction *PreAccountTransaction) Sign(signer TransactionSigner) (*AccountTransaction, error) {
	return signTransaction(signer, preAccountTransaction.Header, preAccountTransaction.Payload)
}

// Serialize returns serialized PreAccountTransaction.
func (preAccountTransaction *PreAccountTransaction) Serialize() []byte {
	buf := make([]byte, 0, int(TransactionHeaderSize)+int(preAccountTransaction.Encoded.Size().Value))
	buf = append(buf, preAccountTransaction.Header.Serialize()...)
	buf = append(buf, preAccountTransaction.Encoded.Serialize()...)

	return buf
}

// Deserialize parses from bytes PreAccountTransaction.
func (preAccountTransaction *PreAccountTransaction) Deserialize(source []byte) (err error) {
	if len(source) < int(TransactionHeaderSize) {
		return errors.New("could not deserialize PreAccountTransaction: invalid length")
	}

	sender, err := AccountAddressFromBytes(source[:32])
	if err != nil {
		return fmt.Errorf("could not receive address from bytes: %v", err)
	}
	preAccountTransaction.Header = &AccountTransactionHeader{
		Sender:         &sender,
		SequenceNumber: &SequenceNumber{Value: binary.BigEndian.Uint64(source[32:40])},
		EnergyAmount:   &Energy{Value: binary.BigEndian.Uint64(source[40:48])},
		PayloadSize:    &PayloadSize{Value: binary.BigEndian.Uint32(source[48:52])},
		Expiry:         &TransactionTime{Value: binary.BigEndian.Uint64(source[52:TransactionHeaderSize])},
	}

	if len(source) < int(TransactionHeaderSize)+int(preAccountTransaction.Header.PayloadSize.Value) {
		return errors.New("could not deserialize PreAccountTransaction: invalid length")
	}

	preAccountTransaction.Encoded.Value = source[TransactionHeaderSize:]
	preAccountTransaction.Payload, err = preAccountTransaction.Encoded.Decode()
	if err != nil {
		return fmt.Errorf("could not decode Encoded Payload: %v", err)
	}

	preAccountTransaction.HashToSign = ComputeTransactionSignHash(preAccountTransaction.Header,
		&AccountTransactionPayload{Payload: preAccountTransaction.Encoded})

	return nil
}

// TransactionSigner is an interface for signing transactions.
type TransactionSigner interface {
	// SignTransactionHash signs transaction hash and returns signatures in AccountTransactionSignature type.
	SignTransactionHash(hashToSign *TransactionHash) (*AccountTransactionSignature, error)
}

// ExactSizeTransactionSigner describes TransactionSigner with ability to return number of signers.
type ExactSizeTransactionSigner interface {
	TransactionSigner
	// NumberOfKeys returns number of signers.
	NumberOfKeys() uint32
}

// AccountTransaction messages which are signed and paid for by the sender account.
type AccountTransaction struct {
	Signature *AccountTransactionSignature
	Header    *AccountTransactionHeader
	Payload   *AccountTransactionPayload
}

func (*AccountTransaction) isBlockItem() {}

// Send sends BlockItem with AccountTransaction using provided Client and returns TransactionHash.
func (accountTransaction *AccountTransaction) Send(ctx context.Context, client *Client) (*TransactionHash, error) {
	signaturesMap := make(map[uint32]*pb.AccountSignatureMap, len(accountTransaction.Signature.Signatures))
	for extKey, signatureMap := range accountTransaction.Signature.Signatures {
		signatures := make(map[uint32]*pb.Signature, len(signatureMap.Signatures))
		for innKey, signature := range signatureMap.Signatures {
			signatures[uint32(innKey)] = &pb.Signature{Value: signature.Value}
		}
		signaturesMap[uint32(extKey)] = &pb.AccountSignatureMap{Signatures: signatures}
	}

	return client.SendBlockItem(ctx, &pb.SendBlockItemRequest{
		BlockItem: &pb.SendBlockItemRequest_AccountTransaction{AccountTransaction: &pb.AccountTransaction{
			Signature: &pb.AccountTransactionSignature{Signatures: signaturesMap},
			Header: &pb.AccountTransactionHeader{
				Sender:         &pb.AccountAddress{Value: accountTransaction.Header.Sender.Value[:]},
				SequenceNumber: &pb.SequenceNumber{Value: accountTransaction.Header.SequenceNumber.Value},
				EnergyAmount:   &pb.Energy{Value: accountTransaction.Header.EnergyAmount.Value},
				Expiry:         &pb.TransactionTime{Value: accountTransaction.Header.Expiry.Value},
			},
			Payload: &pb.AccountTransactionPayload{Payload: &pb.AccountTransactionPayload_RawPayload{
				RawPayload: accountTransaction.Payload.Payload.Encode().Value,
			}},
		}},
	})
}

// AccountTransactionSignature transaction signature.
type AccountTransactionSignature struct {
	Signatures map[uint8]*AccountSignatureMap
}

// AccountSignatureMap wrapper for a map from indexes to signatures.
// Needed because protobuf doesn't allow nested maps directly.
// The keys in the SignatureMap must not exceed 2^8.
type AccountSignatureMap struct {
	Signatures map[uint8]*Signature
}

// Signature a single signature. Used when sending block items to a node with `SendBlockItem`.
type Signature struct {
	Value []byte
}

// signTransaction signs the AccountTransactionHeader and AccountTransactionPayload, construct the transaction, and return it.
func signTransaction(signer TransactionSigner, header *AccountTransactionHeader, payload *AccountTransactionPayload) (*AccountTransaction, error) {
	hashToSign := ComputeTransactionSignHash(header, payload)
	signature, err := signer.SignTransactionHash(hashToSign)
	if err != nil {
		return &AccountTransaction{}, err
	}

	return &AccountTransaction{
		Signature: signature,
		Header:    header,
		Payload:   payload,
	}, nil
}

// ComputeTransactionSignHash computes the transaction sign hash from an AccountTransactionHeader and AccountTransactionPayload.
func ComputeTransactionSignHash(header *AccountTransactionHeader, payload *AccountTransactionPayload) *TransactionHash {
	encodedPayload := payload.Payload.Encode()
	buf := make([]byte, 0, int(TransactionHeaderSize)+len(encodedPayload.Value))
	buf = append(buf, header.Serialize()...)
	buf = append(buf, encodedPayload.Value...)

	hash := new(TransactionHash)
	hash.Value = sha256.Sum256(buf)

	return hash
}
