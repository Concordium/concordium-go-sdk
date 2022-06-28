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
	TransactionStatusAbsent    TransactionStatus = "absent"

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

type TransactionTypeType string

type TransactionTypeContents string

type TransactionResultEventTag string

type TransactionRejectReasonTag string

type ITransactionSummary[E ITransactionResultEvent, R ITransactionRejectReason] interface {
	GetStatus() TransactionStatus
}

type TransactionSummary[E ITransactionResultEvent, R ITransactionRejectReason] struct {
	Status   TransactionStatus         `json:"status"`
	Outcomes TransactionOutcomes[E, R] `json:"outcomes"`
}

func (s *TransactionSummary[E, R]) GetStatus() TransactionStatus {
	return s.Status
}

type TransactionOutcomes[E ITransactionResultEvent, R ITransactionRejectReason] map[BlockHash]*TransactionOutcome[E, R]

func (m *TransactionOutcomes[E, R]) First() (BlockHash, *TransactionOutcome[E, R], bool) {
	for k, v := range *m {
		return k, v, true
	}
	return "", nil, false
}

type TransactionOutcome[E ITransactionResultEvent, R ITransactionRejectReason] struct {
	Sender     AccountAddress           `json:"sender"`
	Hash       TransactionHash          `json:"hash"`
	Cost       *Amount                  `json:"cost"`
	EnergyCost int                      `json:"energyCost"`
	Type       *TransactionType         `json:"type"`
	Result     *TransactionResult[E, R] `json:"result"`
	Index      int                      `json:"index"`
}

func (o *TransactionOutcome[R, E]) Error() error {
	if o.Result.Outcome == TransactionResultSuccess {
		return nil
	}
	return fmt.Errorf("transaction %q was rejected: %w", o.Hash, o.Result.RejectReason.Error())
}

type TransactionType struct {
	Type     TransactionTypeType     `json:"type"`
	Contents TransactionTypeContents `json:"contents"`
}

type TransactionResult[E ITransactionResultEvent, R ITransactionRejectReason] struct {
	Outcome      TransactionResultOutcome   `json:"outcome"`
	Events       TransactionResultEvents[E] `json:"events,omitempty"`
	RejectReason R                          `json:"rejectReason,omitempty"`
}

type TransactionResultEvents[E ITransactionResultEvent] []E

type ITransactionResultEvent interface {
}

type ITransactionRejectReason interface {
	Error() error
}

type TransactionResultEvent struct {
	Tag TransactionResultEventTag `json:"tag"`
}

type TransactionRejectReason struct {
	Tag TransactionRejectReasonTag `json:"tag"`
}

func (r *TransactionRejectReason) Error() error {
	return fmt.Errorf("%s", r.Tag)
}
