package concordium

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Amount struct {
	microCCD uint64
}

func NewAmountZero() *Amount {
	return &Amount{}
}

func NewAmountFromMicroCCD(v int) *Amount {
	return &Amount{microCCD: uint64(v)}
}

func NewAmountFromCCD(v float64) *Amount {
	return &Amount{microCCD: uint64(v * 1e6)}
}

func (a *Amount) MicroCCD() int {
	return int(a.microCCD)
}

func (a *Amount) CCD() float64 {
	return float64(a.microCCD) / 1e6
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
