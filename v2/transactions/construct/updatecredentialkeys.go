package construct

import (
	"fmt"

	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/costs"
)

// UpdateCredentialKeys constructs a transaction that replaces the keys and
// signature threshold of an existing credential (identified by credID), e.g.
// turning a single-key account into an M-of-N multisig.
//
// numExistingCredentials is the number of credentials currently on the account
// (used only for energy accounting); it is 1 for a typical single-credential
// account. The signature threshold must be in [1, len(keys.Keys)].
func UpdateCredentialKeys(
	numSigs uint32,
	sender v2.AccountAddress,
	nonce v2.SequenceNumber,
	expiry v2.TransactionTime,
	credID v2.CredentialRegistrationID,
	keys v2.CredentialPublicKeys,
	numExistingCredentials uint16,
) (*v2.PreAccountTransaction, error) {
	if len(credID.Value) != 48 {
		return nil, fmt.Errorf("credential registration id must be 48 bytes, got %d", len(credID.Value))
	}
	numKeys := len(keys.Keys)
	if numKeys == 0 {
		return nil, fmt.Errorf("credential must have at least one key")
	}
	if numKeys > 255 {
		return nil, fmt.Errorf("credential has too many keys: %d (max 255)", numKeys)
	}
	if keys.Threshold.Value == 0 || int(keys.Threshold.Value) > numKeys {
		return nil, fmt.Errorf("signature threshold %d must be in [1, %d]", keys.Threshold.Value, numKeys)
	}
	for idx, vk := range keys.Keys {
		if vk.Scheme != v2.SchemeEd25519 {
			return nil, fmt.Errorf("key %d uses unsupported signature scheme %d", idx, vk.Scheme)
		}
		if len(vk.Key) != 32 {
			return nil, fmt.Errorf("key %d must be a 32-byte ed25519 public key, got %d bytes", idx, len(vk.Key))
		}
	}

	payload := &v2.AccountTransactionPayload{Payload: &v2.UpdateCredentialKeys{
		Payload: &v2.UpdateCredentialKeysPayload{CredID: credID, Keys: keys},
	}}
	energy := &v2.GivenEnergy{Energy: &v2.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.UpdateCredentialKeys(numExistingCredentials, uint16(numKeys)),
	}}

	return makeTransaction(sender, nonce, expiry, energy, payload), nil
}
