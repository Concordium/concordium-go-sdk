package send

import (
	v2 "github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
)

// TokenUpdateOperation
// / Construct and sign a protocol level token update transaction consisting
// / of the token update operations encoded in the given CBOR.
func TokenUpdateOperation(sender v2.AccountAddress, signer v2.ExactSizeTransactionSigner, nonce v2.SequenceNumber,
	expiry v2.TransactionTime, tokenId v2.TokenId, operations v2.TokenOperations) (*v2.AccountTransaction, error) {
	preUpdInstruction, err := construct.TokenUpdateOperation(signer.NumberOfKeys(), sender, nonce, expiry, tokenId, operations)
	if err != nil {
		return nil, err
	}

	return preUpdInstruction.Sign(signer)
}
