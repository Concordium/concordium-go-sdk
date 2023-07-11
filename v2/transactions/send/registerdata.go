package send

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
)

// RegisterData construct a transaction to register the given piece of data.
func RegisterData(signer v2.ExactSizeTransactionSigner, sender v2.AccountAddress, nonce v2.SequenceNumber,
	expiry v2.TransactionTime, data v2.RegisteredData) (*v2.AccountTransaction, error) {
	return construct.RegisterData(signer.NumberOfKeys(), sender, nonce, expiry, data).Sign(signer)
}
