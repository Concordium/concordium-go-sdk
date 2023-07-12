package send

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
)

// UpdateContract updates a smart contract instance, giving it the given amount of energy
// for execution. The unique parameters are - `energy` -- the amount of energy that can be
// used for contract execution. The base energy amount for transaction verification will be
// added to this cost.
func UpdateContract(signer v2.ExactSizeTransactionSigner, sender v2.AccountAddress, nonce v2.SequenceNumber,
	expiry v2.TransactionTime, payload v2.UpdateContractPayload, energy v2.Energy) (*v2.AccountTransaction, error) {
	return construct.UpdateContract(signer.NumberOfKeys(), sender, nonce, expiry, payload, energy).Sign(signer)
}
