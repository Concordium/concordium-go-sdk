package account

import "github.com/Concordium/concordium-go-sdk"

type body interface {
	concordium.Serializer
	BaseEnergy() int
}

type baseBody struct{}

func (d *baseBody) serialize(typ uint8, elements ...[]byte) []byte {
	s := 1
	for _, e := range elements {
		s += len(e)
	}
	b := make([]byte, s)
	b[0] = typ
	i := 1
	for _, x := range elements {
		copy(b[i:], x)
		i += len(x)
	}
	return b
}
