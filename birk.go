package concordium

import (
	"encoding/hex"
	"fmt"
)

const electionNonceSize = 32

type BirkParameters struct {
	ElectionDifficulty float64       `json:"electionDifficulty"`
	ElectionNonce      ElectionNonce `json:"electionNonce"`
	Bakers             []*BakerInfo  `json:"bakers"`
}

type ElectionNonce [electionNonceSize]byte

func NewElectionNonceFromString(s string) (ElectionNonce, error) {
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

func MustNewElectionNonceFromString(s string) ElectionNonce {
	a, err := NewElectionNonceFromString(s)
	if err != nil {
		panic("MustNewElectionNonceFromString: " + err.Error())
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
