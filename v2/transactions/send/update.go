package send

import (
	v2 "github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
)

// Update constructs UpdateInstruction.
func Update(signer v2.UpdateSigner, seqNumber v2.UpdateSequenceNumber, effectiveTime v2.TransactionTime,
	timeout v2.TransactionTime, payload v2.UpdateInstructionPayload) (*v2.UpdateInstruction, error) {
	preUpdInstruction, err := construct.Update(seqNumber, effectiveTime, timeout, &payload)
	if err != nil {
		return nil, err
	}

	signatures, err := signer.SignUpdateHash(preUpdInstruction.HashToSign)
	if err != nil {
		return nil, err
	}

	var updInstruction = v2.UpdateInstruction{
		Signatures: &signatures.Signatures,
		Header:     preUpdInstruction.Header,
		Payload:    preUpdInstruction.Payload}

	return &updInstruction, err
}
