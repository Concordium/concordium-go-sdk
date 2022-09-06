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

type PublicKey []byte

func NewPublicKeyFromString(s string) (PublicKey, error) {
	g, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("hex decode: %w", err)
	}
	return g, nil
}

func MustNewPublicKeyFromString(s string) PublicKey {
	g, err := NewPublicKeyFromString(s)
	if err != nil {
		panic("MustNewPublicKeyFromString: " + err.Error())
	}
	return g
}

func (k PublicKey) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(k)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", k, err)
	}
	return b, nil
}

func (k *PublicKey) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *k, err)
	}
	*k = v
	return nil
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
