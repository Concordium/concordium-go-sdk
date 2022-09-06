package concordium

// CryptographicParameters is versioned cryptographic parameters.
type CryptographicParameters struct {
	V     uint32                        `json:"v"`
	Value *CryptographicParametersValue `json:"value"`
}

// CryptographicParametersValue is a set of cryptographic parameters that are particular
// to the chain and shared by everybody that interacts with the chain.
type CryptographicParametersValue struct {
	// Generators for the bulletproofs. It is unclear what length we will require here,
	// or whether we'll allow dynamic generation.
	BulletproofGenerators PublicKey `json:"bulletproofGenerators"`
	// A free-form string used to distinguish between different chains even if they share other parameters.
	GenesisString string `json:"genesisString"`
	// A shared commitment key known to the chain and the account holder (and therefore it is public).
	// The account holder uses this commitment key to generate commitments to values in the attribute list.
	OnChainCommitmentKey PublicKey `json:"onChainCommitmentKey"`
}
