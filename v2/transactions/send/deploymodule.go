package send

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// DeployModule deploys the given Wasm module. The module is given
// as a binary source, and no processing is done to the module.
func DeployModule(signer transactions.ExactSizeTransactionSigner, sender types.AccountAddress, nonce types.Nonce,
	expiry types.TransactionTime, module types.WasmModule) (*transactions.AccountTransaction, error) {
	return construct.DeployModule(signer.NumberOfKeys(), sender, nonce, expiry, module).Sign(signer)
}
