package send

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// UpdateContract updates a smart contract instance, giving it the given amount of energy
// for execution. The unique parameters are - `energy` -- the amount of energy that can be
// used for contract execution. The base energy amount for transaction verification will be
// added to this cost.
func UpdateContract(signer transactions.ExactSizeTransactionSigner, sender types.AccountAddress, nonce types.Nonce,
	expiry types.TransactionTime, payload payloads.UpdateContractPayload, energy types.Energy) (*transactions.AccountTransaction, error) {
	return construct.UpdateContract(signer.NumberOfKeys(), sender, nonce, expiry, payload, energy).Sign(signer)
}
