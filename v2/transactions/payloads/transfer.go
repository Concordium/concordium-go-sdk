package payloads

import (
	"encoding/binary"

	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Ensure that TransferPayload implements Payload.
var _ Payload = (*TransferPayload)(nil)

// TransferPayload transfers CCD to an account.
type TransferPayload struct {
	// Address to send to.
	ToAddress types.AccountAddress
	// Amount to send.
	Amount types.Amount
}

// Encode encodes Payload into EncodedPayload.
func (payload *TransferPayload) Encode() EncodedPayload {
	// Payload type byte + payload size.
	buf := make([]byte, 0, payload.Size()+1)
	buf = append(buf, byte(TransferPayloadType))
	buf = append(buf, payload.ToAddress...)
	buf = binary.BigEndian.AppendUint64(buf, uint64(payload.Amount))
	return buf
}

// Decode decodes bytes into TransferPayload.
func (payload *TransferPayload) Decode(source []byte) error {
	if len(source) != payload.Size() {
		return InvalidEncodedPayloadSize
	}

	copy(payload.ToAddress[:], source[:32])
	payload.Amount = types.Amount(binary.BigEndian.Uint64(source[32:]))

	return nil
}

// Size returns the size of the payload in number of bytes.
func (payload *TransferPayload) Size() int {
	// 32 bytes (account address) + 8 bytes (amount).
	return 40
}

func (*TransferPayload) isPayload()     {}
func (*TransferPayload) isPayloadLike() {}
