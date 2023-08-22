package send

import (
	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
)

// Transfer constructs a transfer transaction.
func Transfer(signer v2.ExactSizeTransactionSigner, sender v2.AccountAddress, nonce v2.SequenceNumber,
	expiry v2.TransactionTime, receiver v2.AccountAddress, amount v2.Amount) (*v2.AccountTransaction, error) {
	return construct.Transfer(signer.NumberOfKeys(), sender, nonce, expiry, receiver, amount).Sign(signer)
}
