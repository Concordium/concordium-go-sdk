package concordium

import (
	"encoding/hex"
	"fmt"
)

type CryptographicParameters struct {
	V     int                           `json:"v"`
	Value *CryptographicParametersValue `json:"value"`
}

type CryptographicParametersValue struct {
	BulletproofGenerators BulletproofGenerators `json:"bulletproofGenerators"`
	GenesisString         string                `json:"genesisString"`
	OnChainCommitmentKey  OnChainCommitmentKey  `json:"onChainCommitmentKey"`
}

type BulletproofGenerators []byte

func NewBulletproofGeneratorsFromString(s string) (BulletproofGenerators, error) {
	g, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("hex decode: %w", err)
	}
	return g, nil
}

func MustNewBulletproofGeneratorsFromString(s string) BulletproofGenerators {
	g, err := NewBulletproofGeneratorsFromString(s)
	if err != nil {
		panic("MustNewBulletproofGeneratorsFromString: " + err.Error())
	}
	return g
}

func (g BulletproofGenerators) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(g)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", g, err)
	}
	return b, nil
}

func (g *BulletproofGenerators) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *g, err)
	}
	*g = v
	return nil
}

type OnChainCommitmentKey []byte

func NewOnChainCommitmentKeyFromString(s string) (OnChainCommitmentKey, error) {
	g, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("hex decode: %w", err)
	}
	return g, nil
}

func MustNewOnChainCommitmentKeyFromString(s string) OnChainCommitmentKey {
	k, err := NewOnChainCommitmentKeyFromString(s)
	if err != nil {
		panic("MustNewOnChainCommitmentKeyFromString: " + err.Error())
	}
	return k
}

func (k OnChainCommitmentKey) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(k)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", k, err)
	}
	return b, nil
}

func (k *OnChainCommitmentKey) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *k, err)
	}
	*k = v
	return nil
}
