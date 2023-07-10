package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// transactionBuilder is a helper structure to store the intermediate state of a transaction.
// The problem this helps solve is that to compute the exact energy requirements
// for the transaction we need to know its exact size when serialized. For some
// we could compute this manually, but in general it is less error-prone to serialize
// and get the length. `EnergyAmount` value is being computed before transaction signed with signer.
type transactionBuilder struct {
	header  *transactions.TransactionHeader
	payload payloads.Payload
	encoded payloads.EncodedPayload
}

// newTransactionBuilder is a constructor for transactionBuilder.
func newTransactionBuilder(sender types.AccountAddress, nonce types.Nonce, expiry types.TransactionTime, payload payloads.Payload) *transactionBuilder {
	encoded := payload.Encode()
	header := &transactions.TransactionHeader{
		Sender:       sender,
		Nonce:        nonce,
		EnergyAmount: 0,
		PayloadSize:  encoded.Size(),
		Expiry:       expiry,
	}

	return &transactionBuilder{
		header:  header,
		payload: payload,
		encoded: encoded,
	}
}

// size returns size of built transaction.
func (transactionBuilder *transactionBuilder) size() uint64 {
	return transactions.TransactionHeaderSize + uint64(transactionBuilder.header.PayloadSize)
}

// construct builds PreAccountTransaction with updated energy amount by transmitted counting function.
func (transactionBuilder *transactionBuilder) construct(countEnergyAmountFunc func(uint64) types.Energy) *transactions.PreAccountTransaction {
	size := transactionBuilder.size()
	transactionBuilder.header.EnergyAmount = countEnergyAmountFunc(size)
	hashToSign := transactions.ComputeTransactionSignHash(transactionBuilder.header, transactionBuilder.encoded)

	return &transactions.PreAccountTransaction{
		Header:     transactionBuilder.header,
		Payload:    transactionBuilder.payload,
		Encoded:    transactionBuilder.encoded,
		HashToSign: hashToSign,
	}
}

// makeTransaction returns PreAccountTransaction with computed energy amount and specific Payload.
func makeTransaction(sender types.AccountAddress, nonce types.Nonce, expiry types.TransactionTime,
	energy types.GivenEnergy, payload payloads.Payload) *transactions.PreAccountTransaction {
	builder := newTransactionBuilder(sender, nonce, expiry, payload)
	cost := func(size uint64) types.Energy {
		switch energy.(type) {
		case *types.AbsoluteEnergy:
			return types.Energy(*energy.(*types.AbsoluteEnergy))
		case *types.AddEnergy:
			addEnergy := *energy.(*types.AddEnergy)
			return costs.BaseCost(size, addEnergy.NumSigs) + addEnergy.Energy
		}
		return 0
	}
	return builder.construct(cost)
}
