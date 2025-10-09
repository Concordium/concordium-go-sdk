package construct

import (
	v2 "github.com/Concordium/concordium-go-sdk/v2"
)

// Update prepares PreUpdateInstruction to construct UpdateInstruction.
func Update(seqNumber v2.UpdateSequenceNumber, effectiveTime v2.TransactionTime, timeout v2.TransactionTime, payload *v2.UpdateInstructionPayload) (*v2.PreUpdateInstruction, error) {
	encoded := payload.Payload.Encode()
	payloadSize := encoded.Size()
	header := &v2.UpdateInstructionHeader{
		SequenceNumber: &seqNumber,
		EffectiveTime:  &effectiveTime,
		Timeout:        &timeout,
		PayloadSize:    &payloadSize,
	}

	hashToSign, err := v2.ComputeUpdateSignHash(header, payload)
	if err != nil {
		return nil, err
	}

	return &v2.PreUpdateInstruction{
		Header:     header,
		Payload:    payload,
		Encoded:    encoded,
		HashToSign: hashToSign,
	}, nil
}
