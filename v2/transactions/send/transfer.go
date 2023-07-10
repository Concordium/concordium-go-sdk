package send

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Transfer constructs a transfer transaction.
func Transfer(signer transactions.ExactSizeTransactionSigner, sender types.AccountAddress, nonce types.Nonce,
	expiry types.TransactionTime, receiver types.AccountAddress, amount types.Amount) (*transactions.AccountTransaction, error) {
	return construct.Transfer(signer.NumberOfKeys(), sender, nonce, expiry, receiver, amount).Sign(signer)
}
