package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
)

// Transfer constructs a transfer transaction.
func Transfer(numSigs uint32, sender v2.AccountAddress, nonce v2.SequenceNumber, expiry v2.TransactionTime,
	receiver v2.AccountAddress, amount v2.Amount) *v2.PreAccountTransaction {
	payload := &v2.AccountTransactionPayload{Payload: &v2.Transfer{Payload: &v2.TransferPayload{
		Receiver: &receiver,
		Amount:   &amount,
	}}}
	energy := &v2.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.SimpleTransfer,
	}
	return makeTransaction(sender, nonce, expiry, energy, payload)
}
