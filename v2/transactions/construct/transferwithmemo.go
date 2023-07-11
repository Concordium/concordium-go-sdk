package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
)

// TransferWithMemo constructs a transfer transaction with a memo.
func TransferWithMemo(numSigs uint32, sender v2.AccountAddress, nonce v2.SequenceNumber, expiry v2.TransactionTime,
	receiver v2.AccountAddress, amount v2.Amount, memo v2.Memo) *v2.PreAccountTransaction {
	payload := &v2.AccountTransactionPayload{Payload: &v2.TransferWithMemo{Payload: &v2.TransferWithMemoPayload{
		Receiver: &receiver,
		Memo:     &memo,
		Amount:   &amount,
	}}}
	energy := &v2.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.SimpleTransfer,
	}
	return makeTransaction(sender, nonce, expiry, energy, payload)
}
