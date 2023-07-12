package v2

import (
	"crypto/ed25519"
	"errors"
)

// SimpleSigner implements ExactSizeTransactionSigner and TransactionSigner for one private key.
type SimpleSigner struct {
	privateKey []byte
}

// NewSimpleSigner is a constructor for SimpleSigner.
func NewSimpleSigner(privateKey []byte) *SimpleSigner {
	return &SimpleSigner{privateKey: privateKey}
}

// SignTransactionHash returns signed TransactionHash.
func (simpleSigner *SimpleSigner) SignTransactionHash(hashToSign *TransactionHash) (*AccountTransactionSignature, error) {
	if simpleSigner.privateKey == nil || len(simpleSigner.privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key size")
	}

	signature := make([]byte, 0, ed25519.SignatureSize)
	signature = append(signature, ed25519.Sign(simpleSigner.privateKey, hashToSign.Value[:])...)

	externalMap := make(map[uint32]*AccountSignatureMap, 1)
	internalMap := make(map[uint32]*Signature, 1)
	internalMap[0] = &Signature{Value: signature}
	externalMap[0] = &AccountSignatureMap{Signatures: internalMap}

	return &AccountTransactionSignature{
		Signatures: externalMap,
	}, nil
}

// NumberOfKeys returns number of signing keys.
func (simpleSigner *SimpleSigner) NumberOfKeys() uint32 {
	return 1
}
