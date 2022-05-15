package concordium

type Serialize interface {
	Serialize() ([]byte, error)
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
