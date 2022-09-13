package account

import (
	"fmt"
	"time"

	"github.com/Concordium/concordium-go-sdk"
)

type request struct {
	signature *signature
	header    *header
	body      body
}

func newRequest(
	cred concordium.Credentials,
	addr concordium.AccountAddress,
	nonce concordium.AccountNonce,
	expiry time.Time,
	body body,
) *request {
	return &request{
		signature: &signature{
			cred: cred,
		},
		header: &header{
			accountAddress: addr,
			nonce:          nonce,
			expiry:         expiry,
		},
		body: body,
	}
}

func (r *request) Version() uint8 {
	return 0
}

func (r *request) Kind() uint8 {
	return 0
}

func (r *request) ExpiredAt() time.Time {
	return r.header.expiry
}

func (r *request) Serialize() ([]byte, error) {
	bodyBytes, err := r.body.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to serialize body: %w", err)
	}
	bodySize := len(bodyBytes)
	var signatureCount int
	for _, c := range r.signature.cred {
		signatureCount += len(c)
	}
	r.header.energy = calculateTransactionEnergy(signatureCount, headerSize+bodySize, r.body.BaseEnergy())
	r.header.bodySize = uint32(bodySize)
	headerBytes, err := r.header.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to serialize header: %w", err)
	}
	r.signature.bodyBytes = bodyBytes
	r.signature.headerBytes = headerBytes
	signatureBytes, err := r.signature.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to serialize signature: %w", err)
	}
	headerSize, signatureSize := len(headerBytes), len(signatureBytes)
	b := make([]byte, signatureSize+headerSize+bodySize)
	copy(b, signatureBytes)
	copy(b[signatureSize:], headerBytes)
	copy(b[signatureSize+headerSize:], bodyBytes)
	return b, nil
}
