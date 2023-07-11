package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
)

// DeployModule deploys the given Wasm module. The module is given
// as a binary source, and no processing is done to the module.
func DeployModule(numSigs uint32, sender v2.AccountAddress, nonce v2.SequenceNumber, expiry v2.TransactionTime,
	module v2.VersionedModuleSource) *v2.PreAccountTransaction {
	moduleSize := module.Size()
	payload := &v2.AccountTransactionPayload{Payload: &v2.DeployModule{Payload: &v2.DeployModulePayload{DeployModule: &module}}}
	energy := &v2.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.DeployModuleCost(uint64(moduleSize)),
	}

	return makeTransaction(sender, nonce, expiry, energy, payload)
}
