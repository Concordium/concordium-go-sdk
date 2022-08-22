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

type IdentityProvider struct {
	IpIdentity     int                          `json:"ipIdentity"`
	IpDescription  *IdentityProviderDescription `json:"ipDescription"`
	IpVerifyKey    string                       `json:"ipVerifyKey"`
	IpCdiVerifyKey string                       `json:"ipCdiVerifyKey"`
}

type IdentityProviderDescription struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

type AnonymityRevoker struct {
	ArIdentity    int                          `json:"arIdentity"`
	ArDescription *AnonymityRevokerDescription `json:"arDescription"`
	ArPublicKey   string                       `json:"arPublicKey"`
}

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
	if len(b) < 2 {
		return nil, fmt.Errorf("expect at least 2 bytes but %d given", len(b))
	}
	v, err := hex.DecodeString(string(b[1 : len(b)-1]))
	if err != nil {
		return nil, err
	}
	return v, nil
}
