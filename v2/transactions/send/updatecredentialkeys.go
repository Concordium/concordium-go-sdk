package send

import (
	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
)

// UpdateCredentialKeys constructs and signs a transaction that replaces the keys
// and signature threshold of an existing credential.
func UpdateCredentialKeys(
	signer v2.ExactSizeTransactionSigner,
	sender v2.AccountAddress,
	nonce v2.SequenceNumber,
	expiry v2.TransactionTime,
	credID v2.CredentialRegistrationID,
	keys v2.CredentialPublicKeys,
	numExistingCredentials uint16,
) (*v2.AccountTransaction, error) {
	pre, err := construct.UpdateCredentialKeys(
		signer.NumberOfKeys(),
		sender,
		nonce,
		expiry,
		credID,
		keys,
		numExistingCredentials,
	)
	if err != nil {
		return nil, err
	}

	return pre.Sign(signer)
}
