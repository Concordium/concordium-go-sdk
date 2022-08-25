package account

import (
	"encoding/json"
	"fmt"

	"github.com/Concordium/concordium-go-sdk"
)

type DeployModuleResultEvent struct {
	Tag      concordium.TransactionResultEventTag `json:"tag"`
	Contents concordium.ModuleRef                 `json:"contents"`
}

func NewDeployModuleResultEvent(origin *concordium.TransactionResultEvent) (*DeployModuleResultEvent, error) {
	e := &DeployModuleResultEvent{}
	err := json.Unmarshal(origin.Raw, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

type DeployModuleRejectReason struct {
	Tag      concordium.TransactionRejectReasonTag `json:"tag"`
	Contents concordium.ModuleRef                  `json:"contents"`
}

func NewDeployModuleRejectReason(origin *concordium.TransactionRejectReason) (*DeployModuleRejectReason, error) {
	r := &DeployModuleRejectReason{}
	err := json.Unmarshal(origin.Raw, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *DeployModuleRejectReason) Error() error {
	return fmt.Errorf("module %q already exists", r.Contents)
}

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
