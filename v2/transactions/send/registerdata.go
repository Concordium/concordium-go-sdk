package send

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// RegisterData construct a transaction to register the given piece of data.
func RegisterData(signer transactions.ExactSizeTransactionSigner, sender types.AccountAddress, nonce types.Nonce,
	expiry types.TransactionTime, data types.RegisterData) (*transactions.AccountTransaction, error) {
	return construct.RegisterData(signer.NumberOfKeys(), sender, nonce, expiry, data).Sign(signer)
}
