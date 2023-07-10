package send

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// TransferWithMemo constructs a transfer transaction with a memo.
func TransferWithMemo(signer transactions.ExactSizeTransactionSigner, sender types.AccountAddress, nonce types.Nonce,
	expiry types.TransactionTime, receiver types.AccountAddress, amount types.Amount, memo types.Memo) (*transactions.AccountTransaction, error) {
	return construct.TransferWithMemo(signer.NumberOfKeys(), sender, nonce, expiry, receiver, amount, memo).Sign(signer)
}
