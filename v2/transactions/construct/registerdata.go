package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
)

// RegisterData construct a transaction to register the given piece of data.
func RegisterData(numSigs uint32, sender v2.AccountAddress, nonce v2.SequenceNumber,
	expiry v2.TransactionTime, data v2.RegisteredData) *v2.PreAccountTransaction {
	payload := &v2.AccountTransactionPayload{Payload: &v2.RegisterData{Payload: &v2.RegisterDataPayload{Data: &data}}}
	energy := &v2.GivenEnergy{Energy: &v2.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.RegisterData,
	}}

	return makeTransaction(sender, nonce, expiry, energy, payload)
}
