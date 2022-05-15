package concordium

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	BlockItemKindAccountTransaction   BlockItemKind = 0
	BlockItemKindCredentialDeployment BlockItemKind = 1
	BlockItemKindUpdateInstruction    BlockItemKind = 2
)

type BlockItemKind uint8

// TransactionHash base-16 encoded hash of a transaction (64 characters)
type TransactionHash string

func newTransactionHash(kind BlockItemKind, b []byte) TransactionHash {
	h := sha256.New()
	h.Write([]byte{uint8(kind)})
	h.Write(b)
	return TransactionHash(hex.EncodeToString(h.Sum(nil)))
}
