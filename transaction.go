package concordium

import (
	"time"
)

const (
	DefaultExpiry = 10 * time.Minute
)

type TransactionRequest interface {
	Serializer
	Version() uint8
	Kind() uint8
	ExpiredAt() time.Time
}

// TransactionHash base-16 encoded hash of a transaction (64 characters)
type TransactionHash string

type TransactionStatusStatus string

const (
	// TransactionStatusStatusFinalized means that transaction is finalized in the given block,
	// with the given summary. If the finalization committee is not corrupt then this will
	// always be a singleton map.
	TransactionStatusStatusFinalized TransactionStatusStatus = "finalized"
	// TransactionStatusStatusCommitted means that Transaction is committed to one or more blocks.
	// The outcomes are listed for each block. Note that in the vast majority of cases the outcome
	// of a transaction should not be dependent on the block it is in, but this can in principle happen.
	TransactionStatusStatusCommitted TransactionStatusStatus = "committed"
	// TransactionStatusStatusReceived means that transaction is received, but not yet in any blocks.
	TransactionStatusStatusReceived TransactionStatusStatus = "received"
)

// TransactionStatus is the status of a transaction known to the node.
type TransactionStatus struct {
	Status TransactionStatusStatus `json:"status"`
	// Absents is TransactionStatus.Status is TransactionStatusStatusReceived
	Outcomes map[TransactionHash]*BlockItemSummary `json:"outcomes"`
}
