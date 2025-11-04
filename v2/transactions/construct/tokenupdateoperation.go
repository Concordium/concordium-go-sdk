package construct

import (
	"github.com/fxamacker/cbor/v2"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
)

func TokenUpdateOperation(numSigs uint32, sender v2.AccountAddress, nonce v2.SequenceNumber, expiry v2.TransactionTime, tokenId v2.TokenId, operations v2.TokenOperations,
) (*v2.PreAccountTransaction, error) {
	txEnergy := operations.TxnEnergy()
	energy := &v2.GivenEnergy{Energy: &v2.AddEnergy{
		NumSigs: numSigs,
		Energy:  *txEnergy,
	}}
	rawOps, err := MarshalTokenOperationsToCBORArray(operations)
	if err != nil {
		return &v2.PreAccountTransaction{}, err
	}

	payload := &v2.TokenUpdate{Payload: &v2.TokenOperationsPayload{
		TokenId:    tokenId,
		Operations: rawOps}}

	return makeTransaction(sender, nonce, expiry, energy, &v2.AccountTransactionPayload{Payload: payload}), nil
}

func MarshalTokenOperationsToCBORArray(ops v2.TokenOperations) (v2.RawCBOR, error) {
	var encodedItems [][]byte
	for _, op := range ops {
		b, err := op.EncodeCBOR()
		if err != nil {
			return v2.RawCBOR{}, err
		}
		encodedItems = append(encodedItems, b)
	}

	finalCBOR, err := cbor.Marshal(encodedItems)
	if err != nil {
		return v2.RawCBOR{}, err
	}

	return v2.RawCBOR{Bytes: finalCBOR}, nil
}
