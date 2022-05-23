package account

import (
	"encoding/binary"
	"fmt"
	"github.com/Concordium/concordium-go-sdk"
	"time"
)

// header
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#transactionheader
type header struct {
	accountAddress concordium.AccountAddress
	nonce          concordium.AccountNonce
	expiry         time.Time

	energy   uint64
	bodySize uint32
}

func (h *header) Serialize() ([]byte, error) {
	b := make([]byte, headerSize)

	a, err := h.accountAddress.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to serialize account address: %w", err)
	}
	copy(b, a)

	i := len(a)
	binary.BigEndian.PutUint64(b[i:], uint64(h.nonce))
	i += 8
	binary.BigEndian.PutUint64(b[i:], h.energy)
	i += 8
	binary.BigEndian.PutUint32(b[i:], h.bodySize)
	i += 4
	binary.BigEndian.PutUint64(b[i:], uint64(h.expiry.Unix()))

	return b, nil
}
