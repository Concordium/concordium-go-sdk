package costs

import "github.com/Concordium/concordium-go-sdk/v2"

const (
	// A is the constant for NRG assignment. This scales the effect of the number of signatures on the energy.
	A uint64 = 100

	// B is the constant for NRG assignment. This scales the effect of transaction size on the energy.
	B uint64 = 1
)

// BaseCost returns base cost of a transaction, which is the minimum cost
// that accounts pays for transaction size and signature checking. In addition
// to base cost each transaction has a transaction-type specific cost.
func BaseCost(transactionSize uint64, numSignatures uint32) v2.Energy {
	return v2.Energy{Value: B*transactionSize + A*uint64(numSignatures)}
}

var (
	// SimpleTransfer is an additional cost of a normal, account to account, transfer.
	SimpleTransfer = v2.Energy{Value: 300}

	// EncryptedTransfer is an additional cost of an encrypted transfer.
	EncryptedTransfer = v2.Energy{Value: 27000}

	// TransferToEncrypted is an additional cost of a transfer from public to encrypted balance.
	TransferToEncrypted = v2.Energy{Value: 600}

	// TransferToPublic is an additional cost of a transfer from encrypted to public balance.
	TransferToPublic = v2.Energy{Value: 14850}

	// AddBaker is an additional cost of registerding the account as a baker.
	AddBaker = v2.Energy{Value: 4050}

	// UpdateBakerKeys is an additional cost of updating baker's keys.
	UpdateBakerKeys = v2.Energy{Value: 4050}

	// UpdateBakerStake is an additional cost of updating the baker's stake, either increasing or lowering it.
	UpdateBakerStake = v2.Energy{Value: 300}

	// UpdateBakerRestake is an additional cost of updating the baker's restake flag.
	UpdateBakerRestake = v2.Energy{Value: 300}

	// RemoveBaker is an additional cost of removing a baker.
	RemoveBaker = v2.Energy{Value: 300}

	// RegisterData is an additional cost of registering a piece of data.
	RegisterData = v2.Energy{Value: 300}

	// ConfigureBakerWithKeys is an additional cost of configuring a baker if new keys are registered.
	ConfigureBakerWithKeys = v2.Energy{Value: 4050}

	// ConfigureBakerWithoutKeys is an additional cost of configuring a baker if new keys are **not** registered.
	ConfigureBakerWithoutKeys = v2.Energy{Value: 300}

	// ConfigureDelegation is an additional cost of configuring delegation.
	ConfigureDelegation = v2.Energy{Value: 300}

	// UpdateCredentialsBase cost is going to be negligible compared to
	// verifying the credential, if he credential updates are genuine.
	//
	// There is a non-trivial amount of lookup
	// that needs to be done before we can start any checking. This ensures
	// that those lookups are not a problem.
	UpdateCredentialsBase = v2.Energy{Value: 500}
)

// ScheduledTransfer returns cost of a scheduled transfer, parametrized by the number of releases.
func ScheduledTransfer(numReleases uint16) v2.Energy {
	return v2.Energy{Value: uint64(numReleases) * (300 + 64)}
}

// UpdateCredentialKeys returns an additional cost of updating existing credential
// keys. Parametrised by amount of existing credentials and new keys. Due to the way
// the accounts are stored a new copy of all credentials will be created, so we need to account for that storage increase.
func UpdateCredentialKeys(numCredentialsBefore uint16, numKeys uint16) v2.Energy {
	return v2.Energy{Value: 500*uint64(numCredentialsBefore) + 100*uint64(numKeys)}
}

// DeployModuleCost returns additional cost of deploying a smart contract module,
// parametrized by the size of the module, which is defined to be the size of
// the binary `.wasm` file that is sent as part of the transaction.
func DeployModuleCost(moduleSize uint64) v2.Energy {
	return v2.Energy{Value: moduleSize / 10}
}

// DeployCredential returns additional cost of deploying a credential
// of the given type and with the given number of keys.
func DeployCredential(credentialType v2.CredentialType, numKeys uint16) v2.Energy {
	switch credentialType.Type.(type) {
	case v2.CredentialTypeInitial:
		return v2.Energy{Value: 1000 + 100*uint64(numKeys)}
	case v2.CredentialTypeNormal:
		return v2.Energy{Value: 54000 + 100*uint64(numKeys)}
	}
	return v2.Energy{Value: 0}
}

// UpdateCredentials returns additional cost of updating account's credentials, parametrized by
// - the number of credentials on the account before the update
// - list of keys of credentials to be added.
func UpdateCredentials(numCredentialsBefore uint16, numKeys []uint16) v2.Energy {
	return v2.Energy{
		Value: UpdateCredentialsBase.Value + UpdateCredentialsVariable(numCredentialsBefore, numKeys).Value,
	}
}

// UpdateCredentialsVariable determines the cost of updating credentials on an account.
func UpdateCredentialsVariable(numCredentialsBefore uint16, numKeys []uint16) v2.Energy {
	// the 500 * num_credentials_before is to account for transactions which do
	// nothing, e.g., don't add don't remove, and don't update the
	// threshold. These still have a cost since the way the accounts are
	// stored it will update the stored account data, which does take up
	// quite a bit of space per credential.
	energy := 500 * uint64(numCredentialsBefore)
	for _, key := range numKeys {
		energy += DeployCredential(v2.CredentialType{Type: v2.CredentialTypeNormal{}}, key).Value
	}
	return v2.Energy{Value: energy}
}
