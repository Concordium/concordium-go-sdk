package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// Transfer constructs a transfer transaction.
func Transfer(numSigs uint32, sender types.AccountAddress, nonce types.Nonce, expiry types.TransactionTime,
	receiver types.AccountAddress, amount types.Amount) *transactions.PreAccountTransaction {
	payload := payloads.TransferPayload{
		ToAddress: receiver,
		Amount:    amount,
	}
	energy := types.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.SimpleTransfer,
	}
	return makeTransaction(sender, nonce, expiry, &energy, &payload)
}
