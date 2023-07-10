package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// DeployModule deploys the given Wasm module. The module is given
// as a binary source, and no processing is done to the module.
func DeployModule(numSigs uint32, sender types.AccountAddress, nonce types.Nonce, expiry types.TransactionTime,
	module types.WasmModule) *transactions.PreAccountTransaction {
	moduleSize := module.Source.Size()
	payload := payloads.DeployModulePayload{Module: module}
	energy := types.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.DeployModuleCost(moduleSize),
	}
	return makeTransaction(sender, nonce, expiry, &energy, &payload)
}
