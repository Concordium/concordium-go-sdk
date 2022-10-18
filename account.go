package concordium

import (
	"encoding/json"
	"time"
)

const (
	OpenStatusOpenForAll   OpenStatus = "openForAll"
	OpenStatusClosedForNew OpenStatus = "closedForNew"
	OpenStatusClosedForAll OpenStatus = "closedForAll"
)

// OpenStatus is the status of whether a baking pool allows delegators to join.
type OpenStatus string

type AccountNonce uint64

// NextAccountNonce contains the best guess about the current account nonce,
// with information about reliability.
type NextAccountNonce struct {
	// A flag indicating whether all known transactions are finalized.
	// This can be used as an indicator of how reliable the `nonce` value is.
	AllFinal bool `json:"allFinal"`
	// The nonce that should be used.
	Nonce AccountNonce `json:"nonce"`
}

// EncryptedAmount base-16 encoded string (384 characters)
type EncryptedAmount Hex

// NewEncryptedAmount creates a new EncryptedAmount from string.
func NewEncryptedAmount(s string) (EncryptedAmount, error) {
	v, err := NewHex(s)
	return EncryptedAmount(v), err
}

// MustNewEncryptedAmount calls the NewEncryptedAmount. It panics in case of error.
func MustNewEncryptedAmount(s string) EncryptedAmount {
	a, err := NewEncryptedAmount(s)
	if err != nil {
		panic("MustNewEncryptedAmount: " + err.Error())
	}
	return a
}

func (a *EncryptedAmount) UnmarshalJSON(b []byte) error {
	var h Hex
	err := json.Unmarshal(b, &h)
	if err != nil {
		return err
	}
	*a = EncryptedAmount(h)
	return nil
}

// AccountInfo contains account information exposed via the node's API. This is always
// the state of an account in a specific block.
type AccountInfo struct {
	// Canonical address of the account.
	AccountAddress AccountAddress `json:"accountAddress"`
	// Current (unencrypted) balance of the account.
	AccountAmount *Amount `json:"AccountAmount"`
	// `Some` (non-null) if and only if the account is a baker or delegator. In that case it
	// is the information about the baker or delegator.
	AccountBaker *AccountStakingInfo `json:"accountBaker"`
	// Map of all currently active credentials on the account. This includes public keys
	// that can sign for the given credentials, as well as any revealed attributes. This
	// map always contains a credential with index 0.
	AccountCredentials AccountCredentials `json:"accountCredentials"`
	// The encrypted balance of the account.
	AccountEncryptedAmount *AccountEncryptedAmount `json:"accountEncryptedAmount"`
	// The public key for sending encrypted balances to the account.
	AccountEncryptionKey PublicKey `json:"accountEncryptionKey"`
	// Internal index of the account. Accounts on the chain get sequential indices. These
	// should generally not be used outside of the chain, the account address is meant to
	// be used to refer to accounts, however the account index serves the role of the baker
	// id, if the account is a baker. Hence it is exposed here as well.
	AccountIndex uint64 `json:"accountIndex"`
	// Next nonce to be used for transactions signed from this account.
	AccountNonce AccountNonce `json:"accountNonce"`
	// Release schedule for any locked up amount. This could be an empty release schedule.
	AccountReleaseSchedule *AccountReleaseSchedule `json:"accountReleaseSchedule"`
	// Lower bound on how many credentials must sign any given transaction from this account.
	AccountThreshold uint8 `json:"accountThreshold"`
}

// AccountStakingInfo is different depends on the account type.
// If the account is a baker then only the next is provided:
// 	* AccountStakingInfo.BakerAggregationVerifyKey
// 	* AccountStakingInfo.BakerElectionVerifyKey
// 	* AccountStakingInfo.BakerId
// 	* AccountStakingInfo.PoolInfo
// 	* AccountStakingInfo.PendingChange
// 	* AccountStakingInfo.RestakeEarnings
// 	* AccountStakingInfo.StakedAmount
// If the account is delegating stake to a baker then only the next is provided:
// 	* AccountStakingInfo.DelegationTarget
// 	* AccountStakingInfo.PendingChange
// 	* AccountStakingInfo.RestakeEarnings
// 	* AccountStakingInfo.StakedAmount
type AccountStakingInfo struct {
	// Baker's public key used to check signatures on finalization records. This is
	// only used if the baker has sufficient stake to participate in finalization.
	BakerAggregationVerifyKey string `json:"bakerAggregationVerifyKey"`
	// Baker's public key used to check whether they won the lottery or not.
	BakerElectionVerifyKey string `json:"bakerElectionVerifyKey"`
	// Identity of the baker. This is actually the account index of the account
	// controlling the baker.
	BakerId uint64 `json:"bakerId"`
	// Baker's public key used to check that they are indeed the ones who produced the block.
	BakerSignatureVerifyKey string              `json:"bakerSignatureVerifyKey"`
	PoolInfo                *BakerPoolInfo      `json:"poolInfo"`
	PendingChange           *StakePendingChange `json:"pendingChange"`
	RestakeEarnings         bool                `json:"restakeEarnings"`
	StakedAmount            *Amount             `json:"stakedAmount"`
	DelegationTarget        DelegationTarget    `json:"delegationTarget"`
}

// BakerPoolInfo contains additional information about a baking pool. This information
// is added with the introduction of delegation.
type BakerPoolInfo struct {
	// The commission rates charged by the pool owner.
	CommissionRates *CommissionRates `json:"commissionRates"`
	// The URL that links to the metadata about the pool.
	MetadataUrl string `json:"metadataUrl"`
	// Whether the pool allows delegators.
	OpenStatus OpenStatus `json:"openStatus"`
}

// CommissionRates contains info about commission rates.
type CommissionRates struct {
	// Fraction of baking rewards charged by the pool owner.
	BakingCommission float64 `json:"bakingCommission"`
	// Fraction of finalization rewards charged by the pool owner.
	FinalizationCommission float64 `json:"finalizationCommission"`
	// Fraction of transaction rewards charged by the pool owner.
	TransactionCommission float64 `json:"transactionCommission"`
}

// StakePendingChangeType is type of StakePendingChange
type StakePendingChangeType string

const (
	StakePendingChangeTypeReduce StakePendingChangeType = "ReduceStake"
	StakePendingChangeTypeRemove StakePendingChangeType = "RemoveStake"
)

// StakePendingChange is the pending change in the baker's stake.
type StakePendingChange struct {
	// Is StakePendingChangeTypeReduce if the stake is being reduced and StakePendingChangeTypeRemove if baker
	// will be removed at the end of the given epoch
	Change        StakePendingChangeType `json:"change"`
	EffectiveTime time.Time              `json:"effectiveTime"`
	// Provided only if StakePendingChange.Change is StakePendingChangeTypeReduce
	NewStake *Amount `json:"newStake"`
}

// DelegationTargetType is type of the DelegationTarget
type DelegationTargetType string

const (
	// DelegationTargetTypePassive when delegate passively, i.e., to no specific baker.
	DelegationTargetTypePassive DelegationTargetType = "Passive"
	// DelegationTargetTypeBaker when delegate to a specific baker.
	DelegationTargetTypeBaker DelegationTargetType = "Baker"
)

type DelegationTarget struct {
	DelegateType DelegationTargetType `json:"delegateType"`
	// Only if DelegationTarget.DelegateType is DelegationTargetTypeBaker
	BakerId BakerId `json:"bakerId"`
}

// AccountReleaseSchedule is the state of the account's release schedule. This is the balance
// of the account that is owned by the account, but cannot be used until the release point.
type AccountReleaseSchedule struct {
	// List of timestamped releases. In increasing order of timestamps.
	Schedule []*Release `json:"schedule"`
	// Total amount that is locked up in releases.
	Total *Amount `json:"total"`
}

// Release is an individual release of a locked balance.
type Release struct {
	// Effective time of release.
	Timestamp int64 `json:"timestamp"`
	// Amount to be released.
	Amount *Amount `json:"amount"`
	// List of transaction hashes that contribute a balance to this release.
	Transactions []TransactionHash `json:"transactions"`
}

// AccountEncryptedAmount is the state of the encrypted balance of an account.
type AccountEncryptedAmount struct {
	// If ['Some'], the amount that has resulted from aggregating other amounts and the
	// number of aggregated amounts (must be at least 2 if present).
	AggregatedAmount *PairTuple[EncryptedAmount, uint32] `json:"aggregatedAmount"`
	// Amounts starting at `start_index` (or at `start_index + 1` if there is an aggregated
	// amount present). They are assumed to be numbered sequentially. The length of this
	// list is bounded by the maximum number of incoming amounts on the accounts, which is
	// currently 32. After that aggregation kicks in.
	IncomingAmounts []EncryptedAmount `json:"incomingAmounts"`
	// Encrypted amount that is a result of this accounts' actions. In particular this list includes
	// the aggregate of
	//	* remaining amounts that result when transferring to public balance
	//	* remaining amounts when transferring to another account
	// 	* encrypted amounts that are transferred from public balance
	//
	// When a transfer is made all of these must always be used.
	SelfAmount EncryptedAmount `json:"selfAmount"`
	// Starting index for incoming encrypted amounts. If an aggregated amount is present then
	// this index is associated with such an amount and the list of incoming encrypted amounts
	// starts at the index `start_index + 1`.
	StartIndex uint64 `json:"startIndex"`
}

// AccountCredentials is versioned account credentials.
type AccountCredentials struct {
	V     uint32                         `json:"v"`
	Value AccountCredentialWithoutProofs `json:"value"`
}

// AccountCredentialWithoutProofsType is type of AccountCredentialWithoutProofs
type AccountCredentialWithoutProofsType string

const (
	AccountCredentialWithoutProofsTypeInitial AccountCredentialWithoutProofsType = "initial"
	AccountCredentialWithoutProofsTypeNormal  AccountCredentialWithoutProofsType = "normal"
)

// AccountCredentialWithoutProofs is account credential with values and commitments, but without proofs.
// Serialization must match the serialization of AccountCredential in Haskell.
type AccountCredentialWithoutProofs struct {
	Type     AccountCredentialWithoutProofsType     `json:"type"`
	Contents *AccountCredentialWithoutProofsContent `json:"contents"`
}

// AccountCredentialWithoutProofsContent is different depends on the credentials type.
// Values in initial credential deployment:
// 	* AccountCredentialWithoutProofsContent.CredentialPublicKeys
// 	* AccountCredentialWithoutProofsContent.IpIdentity
// 	* AccountCredentialWithoutProofsContent.Policy
// 	* AccountCredentialWithoutProofsContent.RegId
// Values (as opposed to proofs) in credential deployment:
// 	* AccountCredentialWithoutProofsContent.ArData
// 	* AccountCredentialWithoutProofsContent.Commitments
// 	* AccountCredentialWithoutProofsContent.CredId
// 	* AccountCredentialWithoutProofsContent.CredentialPublicKeys
// 	* AccountCredentialWithoutProofsContent.IpIdentity
// 	* AccountCredentialWithoutProofsContent.Policy
// 	* AccountCredentialWithoutProofsContent.RevocationThreshold
type AccountCredentialWithoutProofsContent struct {
	// Account this credential belongs to.
	CredentialPublicKeys *CredentialPublicKeys `json:"credentialPublicKeys"`
	// Identity of the identity provider who signed the identity object from
	// which this credential is derived.
	IpIdentity uint32 `json:"ipIdentity"`
	// Policy of this credential object.
	Policy *CredentialPolicy `json:"policy"`
	// Credential registration id of the credential.
	RegId Hex `json:"regId"`
	// Anonymity revocation data. List of anonymity revokers which can revoke identity.
	// NB: The order is important since it is the same order as that signed by the identity provider,
	// and permuting the list will invalidate the signature from the identity provider.
	ArData      *ChainArData                     `json:"arData"`
	Commitments *CredentialDeploymentCommitments `json:"commitments"`
	// Credential registration id of the credential.
	CredId Hex `json:"credId"`
	// Credential keys (i.e. account holder keys).
	// Anonymity revocation threshold. Must be <= length of ar_data.
	RevocationThreshold uint8 `json:"revocationThreshold"`
}

// CredentialPublicKeys is public credential keys currently on the account, together
// with the threshold needed for a valid signature on a transaction.
type CredentialPublicKeys struct {
	Keys      *VerifyKey `json:"keys"`
	Threshold uint8      `json:"threshold"`
}

type VerifyKey struct {
	SchemeId  string `json:"schemeId"`
	VerifyKey string `json:"verifyKey"`
}

// CredentialPolicy is a policy is (currently) revealed values of attributes
// that are part of the identity object. Policies are part of credentials.
type CredentialPolicy struct {
	CreatedAt string `json:"createdAt"`
	ValidTo   string `json:"validTo"`
	// Revealed attributes for now. In the future we might have additional items with (Tag, Property, Proof).
	RevealedAttributes string `json:"revealedAttributes"`
}

// ChainArData is data relating to a single anonymity revoker sent by the account holder to the chain.
// Typically a vector of these will be sent to the chain.
type ChainArData struct {
	// encrypted share of id cred pub
	EncIdCredPubShare string `json:"encIdCredPubShare"`
}

// CredentialDeploymentCommitments is the commitments sent by the account holder
// to the chain in order to deploy credentials
type CredentialDeploymentCommitments struct {
	// List of commitments to the attributes that are not revealed. For the purposes of
	// checking signatures, the commitments to those that are revealed as part of the
	// policy are going to be computed by the verifier.
	CmmAttributes Hex `json:"cmmAttributes"`
	// commitment to credential counter
	CmmCredCounter Hex `json:"cmmCredCounter"`
	// commitments to the coefficients of the polynomial used to share
	// id_cred_sec S + b1 X + b2 X^2... where S is id_cred_sec"
	CmmIdCredSecSharingCoeff []Hex `json:"cmmIdCredSecSharingCoeff"`
	// commitment to the max account number.
	CmmMaxAccounts Hex `json:"cmmMaxAccounts"`
	// commitment to the prf key
	CmmPrf Hex `json:"cmmPrf"`
}
