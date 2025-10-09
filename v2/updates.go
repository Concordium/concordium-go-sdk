package v2

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fxamacker/cbor/v2"

	"github.com/Concordium/concordium-go-sdk/v2/pb"
)

// UpdateInstructionHeaderSize describes size of an update instruction header.
// 8 (seq) + 8 (effective time) + 8 (timeout) + 4 (payload size) = 28 bytes.
const UpdateInstructionHeaderSize uint64 = 28

// UpdateInstructionHeader is the header of an update instruction transaction.
// Contains basic data to check whether the update is valid.
type UpdateInstructionHeader struct {
	// Sequence number of the update instruction.
	SequenceNumber *UpdateSequenceNumber
	// Time at which this update becomes effective.
	EffectiveTime *TransactionTime
	// Latest time the update instruction can be included in a block.
	Timeout *TransactionTime
	// Size of the transaction payload.
	PayloadSize *PayloadSize
}

// Serialize returns serialized UpdateInstructionHeader.
func (header *UpdateInstructionHeader) Serialize() []byte {
	buf := make([]byte, 0, UpdateInstructionHeaderSize)
	buf = binary.BigEndian.AppendUint64(buf, header.SequenceNumber.Value)
	buf = binary.BigEndian.AppendUint64(buf, header.EffectiveTime.Value)
	buf = binary.BigEndian.AppendUint64(buf, header.Timeout.Value)
	buf = binary.BigEndian.AppendUint32(buf, header.PayloadSize.Value)

	return buf
}

func (header *UpdateInstructionHeader) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(header)
}

// PreUpdateInstruction describes update instruction before signing.
type PreUpdateInstruction struct {
	Header     *UpdateInstructionHeader
	Payload    *UpdateInstructionPayload
	Encoded    *RawPayload
	HashToSign *UpdateSignHash
}

// UpdateSigner is an interface for signing.
type UpdateSigner interface {
	// SignUpdateHash signs transaction hash and returns signatures in UpdateInstructionSignature type.
	SignUpdateHash(hashToSign *UpdateSignHash) (*UpdateInstructionSignature, error)
}

// UpdateInstruction messages which can update the chain parameters. Including which keys are allowed
// to make future update instructions.
type UpdateInstruction struct {
	Signatures *SignatureMap
	Header     *UpdateInstructionHeader
	Payload    *UpdateInstructionPayload
}

func (*UpdateInstruction) isBlockItem() {}

func (updateInstruction *UpdateInstruction) Send(ctx context.Context, client *Client) (*TransactionHash, error) {
	// Convert SignatureMap to pb.SignatureMap.
	signatures := make(map[uint32]*pb.Signature, len(updateInstruction.Signatures.Signatures))
	for key, sig := range updateInstruction.Signatures.Signatures {
		signatures[key] = &pb.Signature{Value: sig.Value}
	}

	// Constructing pb.UpdateInstruction.
	pbUpdateInstr := &pb.UpdateInstruction{
		Signatures: &pb.SignatureMap{Signatures: signatures},
		Header: &pb.UpdateInstructionHeader{
			SequenceNumber: &pb.UpdateSequenceNumber{Value: updateInstruction.Header.SequenceNumber.Value},
			EffectiveTime:  &pb.TransactionTime{Value: updateInstruction.Header.EffectiveTime.Value},
			Timeout:        &pb.TransactionTime{Value: updateInstruction.Header.Timeout.Value},
		},
		Payload: &pb.UpdateInstructionPayload{
			Payload: &pb.UpdateInstructionPayload_RawPayload{
				RawPayload: updateInstruction.Payload.Payload.Encode().Value,
			},
		},
	}

	return client.SendBlockItem(ctx, &pb.SendBlockItemRequest{
		BlockItem: &pb.SendBlockItemRequest_UpdateInstruction{
			UpdateInstruction: pbUpdateInstr,
		},
	})
}

// UpdateKeysIndex describes index of an update key.
type UpdateKeysIndex uint8

// UpdateKeyPair is similar to Account KeyPair but for update instructions.
type UpdateKeyPair struct {
	// secret describes `signKey`.
	secret ed25519.PrivateKey
	// public describes `verifyKey`.
	public ed25519.PublicKey
}

// UpdateInstructionSignature transaction signature.
type UpdateInstructionSignature struct {
	Signatures SignatureMap
}

// Sign signs the message with update key.
func (kp *UpdateKeyPair) Sign(msg []byte) Signature {
	privateKey := append(kp.secret, kp.public...)
	return Signature{Value: ed25519.Sign(privateKey, msg)}
}

// WalletUpdateSigner holds update keys and implements UpdateSigner.
type WalletUpdateSigner struct {
	Keys map[UpdateKeysIndex]*UpdateKeyPair
}

// signer TransactionSigner, header *AccountTransactionHeader, payload *AccountTransactionPayload
func signUpdate(signer UpdateSigner, header *UpdateInstructionHeader, payload *UpdateInstructionPayload) (*UpdateInstruction, error) {
	hashToSign, err := ComputeUpdateSignHash(header, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to compute update sign hash: %w", err)
	}

	signature, err := signer.SignUpdateHash(hashToSign)
	if err != nil {
		return nil, fmt.Errorf("failed to sign update hash: %w", err)
	}

	return &UpdateInstruction{
		Signatures: &signature.Signatures,
		Header:     header,
		Payload:    payload,
	}, nil
}

// Sign signs PreAccountTransaction with TransactionSigner and returns AccountTransaction.
func (preUpdateInstruction *PreUpdateInstruction) Sign(signer UpdateSigner) (*UpdateInstruction, error) {
	return signUpdate(signer, preUpdateInstruction.Header, preUpdateInstruction.Payload)
}

// SignUpdateHash implements UpdateSigner.
// It signs hashToSign with all update keys and returns signatures map.
func (signer *WalletUpdateSigner) SignUpdateHash(hashToSign *UpdateSignHash) (*UpdateInstructionSignature, error) {
	if signer.Keys == nil || len(signer.Keys) == 0 {
		return nil, errors.New("WalletUpdateSigner.Keys is not initialized or empty")
	}

	signatures := make(map[uint32]*Signature, len(signer.Keys))
	for idx, kp := range signer.Keys {
		sig := kp.Sign(hashToSign.Value[:])
		signatures[uint32(idx)] = &sig
	}

	return &UpdateInstructionSignature{
		Signatures: SignatureMap{Signatures: signatures},
	}, nil
}

// ComputeUpdateSignHash computes the transaction sign hash from an UpdateInstructionHeader and UpdateInstructionPayload.
func ComputeUpdateSignHash(header *UpdateInstructionHeader, payload *UpdateInstructionPayload) (*UpdateSignHash, error) {
	encodedPayload := payload.Payload.Encode()
	buf := make([]byte, 0, int(UpdateInstructionHeaderSize)+len(encodedPayload.Value))
	buf = append(buf, header.Serialize()...)
	buf = append(buf, encodedPayload.Value...)

	var result UpdateSignHash
	result.Value = sha256.Sum256(buf)

	return &result, nil
}
