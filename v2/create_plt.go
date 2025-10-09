package v2

import (
	"bytes"
	"github.com/fxamacker/cbor/v2"
)

const CreatePLTByte = 24

// CreatePLTPayload Update payload for creating a new protocol-level token
type CreatePLTPayload struct {
	TokenId                  TokenId
	TokenModule              TokenModuleRef
	Decimals                 uint32
	InitializationParameters RawCBOR
}

func (c *CreatePLTPayload) EncodeCBOR() ([]byte, error) {
	return cbor.Marshal(c)
}

type CreatePLT struct {
	Payload *CreatePLTPayload
}

func (c *CreatePLT) Encode() (*UpdateInstructionPayload, error) {
	rawBytes, err := c.Payload.EncodeCBOR()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.WriteByte(CreatePLTByte)
	buf.Write(rawBytes)

	return &UpdateInstructionPayload{
		Payload: &RawPayload{Value: buf.Bytes()},
	}, nil
}
