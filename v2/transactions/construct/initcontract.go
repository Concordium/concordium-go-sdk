package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

// InitContract initializes a smart contract, giving it the given amount of energy for
// execution. The unique parameters are - `energy` -- the amount of energy that can be
// used for contract execution. The base energy amount for transaction verification will
// be added to this cost.
func InitContract(numSigs uint32, sender v2.AccountAddress, nonce v2.SequenceNumber, expiry v2.TransactionTime,
	payload v2.InitContractPayload, energy v2.Energy) *v2.PreAccountTransaction {
	accountPayload := &v2.AccountTransactionPayload{Payload: &v2.InitContract{Payload: &payload}}
	resultEnergy := &v2.GivenEnergy{Energy: &v2.AddEnergy{
		NumSigs: numSigs,
		Energy:  energy,
	}}

	return makeTransaction(sender, nonce, expiry, resultEnergy, accountPayload)
}
