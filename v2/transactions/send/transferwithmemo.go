package send

import (
	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
)

// TransferWithMemo constructs a transfer transaction with a memo.
func TransferWithMemo(signer v2.ExactSizeTransactionSigner, sender v2.AccountAddress, nonce v2.SequenceNumber,
	expiry v2.TransactionTime, receiver v2.AccountAddress, amount v2.Amount, memo v2.Memo) (*v2.AccountTransaction, error) {
	return construct.TransferWithMemo(signer.NumberOfKeys(), sender, nonce, expiry, receiver, amount, memo).Sign(signer)
}
