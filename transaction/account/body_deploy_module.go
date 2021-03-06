package account

type deployModuleBody struct {
	baseBody
	wasm []byte
}

// newDeployModuleBody returns account transaction body to deploy the module.
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#deploymodule
func newDeployModuleBody(wasm []byte) *deployModuleBody {
	return &deployModuleBody{
		wasm: wasm,
	}
}

func (d *deployModuleBody) BaseEnergy() int {
	return len(d.wasm) / 10
}

func (d *deployModuleBody) Serialize() ([]byte, error) {
	return d.serialize(typeDeployModule, d.wasm), nil
}
