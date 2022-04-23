package concordium

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Amount struct {
	microGTU uint64
}

func NewAmountZero() *Amount {
	return &Amount{}
}

func NewAmountFromMicroGTU(v int) *Amount {
	return &Amount{microGTU: uint64(v)}
}

func NewAmountFromGTU(v float64) *Amount {
	return &Amount{microGTU: uint64(v * 1e6)}
}

func (a *Amount) MicroGTU() int {
	return int(a.microGTU)
}

func (a *Amount) GTU() float64 {
	return float64(a.microGTU) / 1e6
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(strconv.Itoa(int(a.microGTU)))
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
	a.microGTU = uint64(num)
	return nil
}
