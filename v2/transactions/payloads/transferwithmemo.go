package payloads

import (
	"encoding/binary"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Ensure that TransferWithMemoPayload implements Payload.
var _ Payload = (*TransferWithMemoPayload)(nil)

// TransferWithMemoPayload transfers CCD to an account.
type TransferWithMemoPayload struct {
	// Address to send to.
	ToAddress types.AccountAddress
	// Memo to include in the transfer.
	Memo types.Memo
	// Amount to send.
	Amount types.Amount
}

// Encode encodes Payload into EncodedPayload.
func (payload *TransferWithMemoPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(TransferPayloadType))
	buf = append(buf, payload.ToAddress...)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(payload.Memo)))
	buf = append(buf, payload.Memo...)
	buf = binary.BigEndian.AppendUint64(buf, uint64(payload.Amount))
	return buf
}

// Decode decodes bytes into TransferWithMemoPayload.
func (payload *TransferWithMemoPayload) Decode(source []byte) error {
	if len(source) <= 34 {
		return InvalidEncodedPayloadSize
	}

	copy(payload.ToAddress[:], source[:32])
	memoSize := binary.BigEndian.Uint16(source[32:34])
	if len(source) != int(42+memoSize) {
		return InvalidEncodedPayloadSize
	}

	copy(payload.Memo[:], source[34:34+memoSize])
	payload.Amount = types.Amount(binary.BigEndian.Uint64(source[34+memoSize:]))

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *TransferWithMemoPayload) Size() int {
	// 32 bytes (account address) + 2 bytes (memo size) + memo bytes + 8 bytes (amount).
	return 42 + len(payload.Memo)
}

func (*TransferWithMemoPayload) isPayload()     {}
func (*TransferWithMemoPayload) isPayloadLike() {}
