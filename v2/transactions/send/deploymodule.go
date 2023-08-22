package send

import (
	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
)

// DeployModule deploys the given Wasm module. The module is given
// as a binary source, and no processing is done to the module.
func DeployModule(signer v2.ExactSizeTransactionSigner, sender v2.AccountAddress, nonce v2.SequenceNumber,
	expiry v2.TransactionTime, module v2.VersionedModuleSource) (*v2.AccountTransaction, error) {
	return construct.DeployModule(signer.NumberOfKeys(), sender, nonce, expiry, module).Sign(signer)
}
