package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// InitContract initializes a smart contract, giving it the given amount of energy for
// execution. The unique parameters are - `energy` -- the amount of energy that can be
// used for contract execution. The base energy amount for transaction verification will
// be added to this cost.
func InitContract(numSigs uint32, sender types.AccountAddress, nonce types.Nonce, expiry types.TransactionTime,
	payload payloads.InitContractPayload, energy types.Energy) *transactions.PreAccountTransaction {
	resultEnergy := types.AddEnergy{
		NumSigs: numSigs,
		Energy:  energy,
	}
	return makeTransaction(sender, nonce, expiry, &resultEnergy, &payload)
}
