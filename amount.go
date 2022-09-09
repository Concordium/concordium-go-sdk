package concordium

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
)

const amountSize = 8

// Amount is a CCD amount.
type Amount struct {
	microCCD uint64
}

// NewAmountZero creates an empty amount.
func NewAmountZero() *Amount {
	return &Amount{}
}

// NewAmountFromMicroCCD created an amount with given micro CCD.
func NewAmountFromMicroCCD(v uint64) *Amount {
	return &Amount{microCCD: v}
}

// MicroCCD returns micro CCD value.
func (a *Amount) MicroCCD() uint64 {
	return a.microCCD
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(strconv.Itoa(int(a.microCCD)))
	if err != nil {
		return nil, fmt.Errorf("%T: %w", *a, err)
	}
	return b, nil
}

func (a *Amount) UnmarshalJSON(b []byte) error {
	var str json.Number
	err := json.Unmarshal(b, &str)
	if err != nil {
		return fmt.Errorf("%T: %w", *a, err)
	}
	num, err := str.Int64()
	if err != nil {
		return fmt.Errorf("%T: %w", *a, err)
	}
	a.microCCD = uint64(num)
	return nil
}

func (a *Amount) Serialize() ([]byte, error) {
	b := make([]byte, amountSize)
	binary.BigEndian.PutUint64(b, a.microCCD)
	return b, nil
}

func (a *Amount) Deserialize(b []byte) error {
	if len(b) < amountSize {
		return fmt.Errorf("%T requires %d bytes", *a, amountSize)
	}
	a.microCCD = binary.BigEndian.Uint64(b)
	return nil
}

func (a *Amount) SerializeModel() ([]byte, error) {
	b := make([]byte, amountSize)
	binary.LittleEndian.PutUint64(b, a.microCCD)
	return b, nil
}

func (a *Amount) DeserializeModel(b []byte) (int, error) {
	if len(b) < amountSize {
		return 0, fmt.Errorf("%T requires %d bytes", *a, amountSize)
	}
	a.microCCD = binary.LittleEndian.Uint64(b)
	return amountSize, nil
}
