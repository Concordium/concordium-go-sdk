package account

import "encoding/binary"

type deployModuleBody struct {
	baseBody
	version uint32
	size    uint32
	wasm    []byte
}

// newDeployModuleBody returns account transaction body to deploy the module.
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#deploymodule
func newDeployModuleBody(wasm []byte) *deployModuleBody {
	return &deployModuleBody{
		version: 0,
		size:    uint32(len(wasm)),
		wasm:    wasm,
	}
}

func (d *deployModuleBody) BaseEnergy() int {
	return int(d.size) / 10
}

func (d *deployModuleBody) Serialize() ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b, d.version)
	binary.BigEndian.PutUint32(b[4:], d.size)
	return d.serialize(typeDeployModule, b, d.wasm), nil
}
