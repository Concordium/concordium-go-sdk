package concordium

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	DefaultExpiry = 10 * time.Minute

	TransactionStatusFinalized TransactionStatus = "finalized"
	TransactionStatusCommitted TransactionStatus = "committed"
	TransactionStatusReceived  TransactionStatus = "received"

	TransactionResultSuccess TransactionResultOutcome = "success"
	TransactionResultReject  TransactionResultOutcome = "reject"
)

const (
	BlockItemVersion0 BlockItemVersion = 0

	BlockItemKindAccountTransaction   BlockItemKind = 0
	BlockItemKindCredentialDeployment BlockItemKind = 1
	BlockItemKindUpdateInstruction    BlockItemKind = 2
)

type BlockItemVersion uint8

type BlockItemKind uint8

type TransactionRequest interface {
	Serializer
	Version() BlockItemVersion
	Kind() BlockItemKind
	ExpiredAt() time.Time
}

// TransactionHash base-16 encoded hash of a transaction (64 characters)
type TransactionHash string

func newTransactionHash(kind BlockItemKind, b []byte) TransactionHash {
	h := sha256.New()
	h.Write([]byte{uint8(kind)})
	h.Write(b)
	return TransactionHash(hex.EncodeToString(h.Sum(nil)))
}

type TransactionStatus string

type TransactionResultOutcome string

type TransactionSummary struct {
	Status   TransactionStatus             `json:"status"`
	Outcomes map[string]TransactionOutcome `json:"outcomes"`
}

type TransactionOutcome struct {
	Hash   TransactionHash   `json:"hash"`
	Result TransactionResult `json:"result"`
}

func (o *TransactionOutcome) Error() error {
	if o.Result.Outcome == TransactionResultSuccess {
		return nil
	}
	return fmt.Errorf("transaction %q was rejected", o.Hash)
}

type TransactionResult struct {
	Outcome TransactionResultOutcome `json:"outcome"`
	Events  TransactionResultEvents  `json:"events,omitempty"`
}

type TransactionResultEvents []*TransactionEvent

type TransactionEvent struct {
	Address  ContractAddress `json:"address"`
	Contents string          `json:"contents"`
}
