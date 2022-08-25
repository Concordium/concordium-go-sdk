package concordium

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

type TransactionSummary struct {
	Status   TransactionStatus   `json:"status"`
	Outcomes TransactionOutcomes `json:"outcomes"`
}

type TransactionOutcomes map[BlockHash]*TransactionOutcome

func (m *TransactionOutcomes) First() (BlockHash, *TransactionOutcome, bool) {
	for k, v := range *m {
		return k, v, true
	}
	return BlockHash{}, nil, false
}

type TransactionOutcome struct {
	Sender     AccountAddress     `json:"sender"`
	Hash       TransactionHash    `json:"hash"`
	Cost       *Amount            `json:"cost"`
	EnergyCost int                `json:"energyCost"`
	Type       *TransactionType   `json:"type"`
	Result     *TransactionResult `json:"result"`
	Index      int                `json:"index"`
}

type TransactionType struct {
	Type     TransactionTypeType     `json:"type"`
	Contents TransactionTypeContents `json:"contents"`
}

type TransactionResult struct {
	Outcome      TransactionResultOutcome `json:"outcome"`
	Events       TransactionResultEvents  `json:"events,omitempty"`
	RejectReason *TransactionRejectReason `json:"rejectReason,omitempty"`
}

type TransactionResultEvents []*TransactionResultEvent

type TransactionResultEvent struct {
	Raw []byte                    `json:"-"`
	Tag TransactionResultEventTag `json:"tag"`
}

func (e *TransactionResultEvent) MarshalJSON() ([]byte, error) {
	return e.Raw, nil
}

func (e *TransactionResultEvent) UnmarshalJSON(b []byte) error {
	e.Raw = b
	s := &struct {
		Tag TransactionResultEventTag `json:"tag"`
	}{}
	err := json.Unmarshal(b, s)
	if err != nil {
		return fmt.Errorf("%T: %w", e, err)
	}
	e.Tag = s.Tag
	return nil
}

type TransactionRejectReason struct {
	Raw []byte                     `json:"-"`
	Tag TransactionRejectReasonTag `json:"tag"`
}

func (r *TransactionRejectReason) MarshalJSON() ([]byte, error) {
	return r.Raw, nil
}

func (r *TransactionRejectReason) UnmarshalJSON(b []byte) error {
	r.Raw = b
	s := &struct {
		Tag TransactionRejectReasonTag `json:"tag"`
	}{}
	err := json.Unmarshal(b, s)
	if err != nil {
		return fmt.Errorf("%T: %w", r, err)
	}
	r.Tag = s.Tag
	return nil
}

func (r *TransactionRejectReason) Error() error {
	return fmt.Errorf("%s", r.Tag)
}
