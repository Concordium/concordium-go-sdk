package construct

import (
	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/costs"
)

// transactionBuilder is a helper structure to store the intermediate state of a transaction.
// The problem this helps solve is that to compute the exact energy requirements
// for the transaction we need to know its exact size when serialized. For some
// we could compute this manually, but in general it is less error-prone to serialize
// and get the length. `EnergyAmount` value is being computed before transaction signed with signer.
type transactionBuilder struct {
	header  *v2.AccountTransactionHeader
	payload *v2.AccountTransactionPayload
	encoded *v2.RawPayload
}

// newTransactionBuilder is a constructor for transactionBuilder.
func newTransactionBuilder(sender v2.AccountAddress, nonce v2.SequenceNumber, expiry v2.TransactionTime, payload *v2.AccountTransactionPayload) *transactionBuilder {
	encoded := payload.Payload.Encode()
	payloadSize := encoded.Size()
	header := &v2.AccountTransactionHeader{
		Sender:         &sender,
		SequenceNumber: &nonce,
		EnergyAmount:   &v2.Energy{Value: 0},
		PayloadSize:    &payloadSize,
		Expiry:         &expiry,
	}

	return &transactionBuilder{
		header:  header,
		payload: payload,
		encoded: encoded,
	}
}

// size returns size of built transaction.
func (transactionBuilder *transactionBuilder) size() uint64 {
	return v2.TransactionHeaderSize + uint64(transactionBuilder.header.PayloadSize.Value)
}

// construct builds PreAccountTransaction with updated energy amount by transmitted counting function.
func (transactionBuilder *transactionBuilder) construct(countEnergyAmountFunc func(uint64) *v2.Energy) *v2.PreAccountTransaction {
	size := transactionBuilder.size()
	transactionBuilder.header.EnergyAmount = countEnergyAmountFunc(size)
	hashToSign := v2.ComputeTransactionSignHash(transactionBuilder.header,
		&v2.AccountTransactionPayload{Payload: transactionBuilder.encoded})

	return &v2.PreAccountTransaction{
		Header:     transactionBuilder.header,
		Payload:    transactionBuilder.payload,
		Encoded:    transactionBuilder.encoded,
		HashToSign: hashToSign,
	}
}

// makeTransaction returns PreAccountTransaction with computed energy amount and specific Payload.
func makeTransaction(sender v2.AccountAddress, nonce v2.SequenceNumber, expiry v2.TransactionTime,
	energy *v2.GivenEnergy, payload *v2.AccountTransactionPayload) *v2.PreAccountTransaction {
	builder := newTransactionBuilder(sender, nonce, expiry, payload)
	cost := func(size uint64) *v2.Energy {
		switch e := energy.Energy.(type) {
		case *v2.AbsoluteEnergy:
			return &v2.Energy{Value: e.Value}
		case *v2.AddEnergy:
			return &v2.Energy{Value: costs.BaseCost(size, e.NumSigs).Value + e.Energy.Value}
		}
		return &v2.Energy{Value: 0}
	}
	return builder.construct(cost)
}
