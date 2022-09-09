package concordium

import (
	"encoding/hex"
	"fmt"
)

const electionNonceSize = 32

// BirkParameters is the state of consensus parameters, and allowed participants (i.e., bakers).
type BirkParameters struct {
	// The list of active bakers.
	Bakers []*BakerInfo `json:"bakers"`
	// Current election difficulty.
	ElectionDifficulty float64 `json:"electionDifficulty"`
	// "Leadership election nonce for the current epoch.
	ElectionNonce ElectionNonce `json:"electionNonce"`
}

type ElectionNonce [electionNonceSize]byte

// NewElectionNonce creates a new ElectionNonce from string.
func NewElectionNonce(s string) (ElectionNonce, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return ElectionNonce{}, fmt.Errorf("hex decode: %w", err)
	}
	if len(b) != electionNonceSize {
		return ElectionNonce{}, fmt.Errorf("expect %d bytes but %d given", electionNonceSize, len(b))
	}
	var n ElectionNonce
	copy(n[:], b)
	return n, nil
}

// MustNewElectionNonce calls the NewElectionNonce. It panics in case of error.
func MustNewElectionNonce(s string) ElectionNonce {
	a, err := NewElectionNonce(s)
	if err != nil {
		panic("MustNewElectionNonce: " + err.Error())
	}
	return a
}

func (n ElectionNonce) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(n[:])
	if err != nil {
		return nil, fmt.Errorf("%T: %w", n, err)
	}
	return b, nil
}

func (n *ElectionNonce) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *n, err)
	}
	var x ElectionNonce
	copy(x[:], v)
	*n = x
	return nil
}
