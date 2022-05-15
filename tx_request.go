package concordium

import "time"

type TransactionRequest interface {
	Serialize
	Kind() BlockItemKind
	ExpiredAt() time.Time
}
