package account

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/Concordium/concordium-go-sdk"
)

// signature
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#transactionsignature
type signature struct {
	cred concordium.Credentials

	headerBytes []byte
	bodyBytes   []byte
}

func (s *signature) Serialize() ([]byte, error) {
	size := 1 // Length of outer map (uint8)
	for _, c := range s.cred {
		cs := len(c)
		size += 1                          // CredentialIndex (uint8)
		size += 1                          // Length of inner map (uint8)
		size += cs                         // KeyIndex (uint8) of each keyPair
		size += cs * 2                     // Length of Signature (uint16) of each keyPair
		size += cs * ed25519.SignatureSize // Signature of each keyPair
	}
	h := sha256.New()
	h.Write(s.headerBytes)
	h.Write(s.bodyBytes)
	d := h.Sum(nil)

	b := make([]byte, size)
	i := 0
	b[0] = uint8(len(s.cred))
	i++
	for ci, c := range s.cred {
		b[i] = ci
		i++
		b[i] = uint8(len(c))
		i++
		for ki, kp := range c {
			b[i] = ki
			i++
			binary.BigEndian.PutUint16(b[i:], ed25519.SignatureSize)
			i += 2
			sk, err := kp.SignKey.Decode()
			if err != nil {
				return nil, fmt.Errorf("unable to decode sign key: %w", err)
			}
			vk, err := hex.DecodeString(kp.VerifyKey)
			if err != nil {
				return nil, fmt.Errorf("unable to decode verify key: %w", err)
			}
			k := make([]byte, ed25519.PrivateKeySize)
			copy(k, sk)
			copy(k[len(sk):], vk)
			copy(b[i:], ed25519.Sign(k, d))
			i += ed25519.SignatureSize
		}
	}
	return b, nil
}
