package concordium

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Serializer interface {
	Serialize() ([]byte, error)
}

type Deserializer interface {
	Deserialize(b []byte) error
}

// Hex represents base-16 encoded data.
type Hex []byte

// NewHex creates a new Hex from string.
func NewHex(s string) (Hex, error) {
	g, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("hex decode: %w", err)
	}
	return g, nil
}

// MustNewHex calls the NewHex. It panics in case of error.
func MustNewHex(s string) Hex {
	h, err := NewHex(s)
	if err != nil {
		panic("MustNewHex: " + err.Error())
	}
	return h
}

func (h Hex) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(hex.EncodeToString(h))
	if err != nil {
		return nil, fmt.Errorf("%T: %w", h, err)
	}
	return b, nil
}

func (h *Hex) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("%T: %w", *h, err)
	}
	v, err := hex.DecodeString(s)
	if err != nil {
		return fmt.Errorf("%T: %w", *h, err)
	}
	*h = v
	return nil
}

type PublicKey Hex

// NewPublicKey creates a new PublicKey from string.
func NewPublicKey(s string) (PublicKey, error) {
	v, err := NewHex(s)
	return PublicKey(v), err
}

// MustNewPublicKey calls the NewPublicKey. It panics in case of error.
func MustNewPublicKey(s string) PublicKey {
	v, err := NewPublicKey(s)
	if err != nil {
		panic("MustNewPublicKey: " + err.Error())
	}
	return v
}

// IdentityProvider is public information about an identity provider.
type IdentityProvider struct {
	// Ed public key of the IP
	IpCdiVerifyKey PublicKey `json:"ipCdiVerifyKey"`
	// Free form description, e.g., how to contact them off-chain
	IpDescription *IdentityProviderDescription `json:"ipDescription"`
	// Unique identifier of the identity provider.
	IpIdentity uint32 `json:"ipIdentity"`
	// PS public key of the IP
	IpVerifyKey PublicKey `json:"ipVerifyKey"`
}

// IdentityProviderDescription is description either of an anonymity revoker or identity provider.
// Metadata that should be visible on the chain.
type IdentityProviderDescription struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

// AnonymityRevoker is information on a single anonymity revoker held by the IP.
// Typically an IP will hold a more than one.
type AnonymityRevoker struct {
	// description of the anonymity revoker (e.g. name, contact number)
	ArDescription *AnonymityRevokerDescription `json:"arDescription"`
	// unique identifier of the anonymity revoker
	ArIdentity uint32 `json:"arIdentity"`
	// elgamal encryption key of the anonymity revoker
	ArPublicKey PublicKey `json:"arPublicKey"`
}

// AnonymityRevokerDescription is description either of an anonymity revoker or identity provider.
// Metadata that should be visible on the chain.
type AnonymityRevokerDescription struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

func hexMarshalJSON(v []byte) ([]byte, error) {
	b, err := json.Marshal(hex.EncodeToString(v))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func hexUnmarshalJSON(b []byte) ([]byte, error) {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	v, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return v, nil
}
