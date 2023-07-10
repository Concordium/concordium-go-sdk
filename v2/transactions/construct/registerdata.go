package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// RegisterData construct a transaction to register the given piece of data.
func RegisterData(numSigs uint32, sender types.AccountAddress, nonce types.Nonce,
	expiry types.TransactionTime, data types.RegisterData) *transactions.PreAccountTransaction {
	payload := payloads.RegisterDataPayload{
		Data: data,
	}
	energy := types.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.RegisterData,
	}
	return makeTransaction(sender, nonce, expiry, &energy, &payload)
}
