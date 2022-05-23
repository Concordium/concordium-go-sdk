package account

import "github.com/Concordium/concordium-go-sdk"

var (
	testRandomBytes = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	testCredentials = concordium.Credentials{
		0: {
			0: concordium.KeyPair{
				SignKey:   concordium.DecryptedSignKey("b53af4521a678b015bbae217277933e87b978a48a9a07d55cc369cdf5e1ac215"),
				VerifyKey: "056a0ab31cd169ae545becad4bf7cc0609945c76877f72c0e8574a85d4818260",
			},
		},
	}
)
