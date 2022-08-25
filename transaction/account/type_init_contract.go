package account

import (
	"encoding/json"
	"fmt"

	"github.com/Concordium/concordium-go-sdk"
)

type InitContractResultEvent struct {
	Tag             concordium.TransactionResultEventTag `json:"tag"`
	ContractVersion int                                  `json:"contractVersion"`
	Ref             concordium.ModuleRef                 `json:"ref"`
	Address         *concordium.ContractAddress          `json:"address"`
	Amount          *concordium.Amount                   `json:"amount"`
	InitName        concordium.InitName                  `json:"initName"`
	Events          []concordium.Model                   `json:"events"`
}

func NewInitContractResultEvent(origin *concordium.TransactionResultEvent) (*InitContractResultEvent, error) {
	e := &InitContractResultEvent{}
	err := json.Unmarshal(origin.Raw, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

type InitContractRejectReason struct {
	Tag          concordium.TransactionRejectReasonTag `json:"tag"`
	RejectReason int                                   `json:"rejectReason"`
}

func NewInitContractRejectReason(origin *concordium.TransactionRejectReason) (*InitContractRejectReason, error) {
	r := &InitContractRejectReason{}
	err := json.Unmarshal(origin.Raw, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *InitContractRejectReason) Error() error {
	return fmt.Errorf("init rejected with reason %d", r.RejectReason)
}

type initContractBody struct {
	baseBody
	amount    *concordium.Amount
	moduleRef concordium.ModuleRef
	initName  concordium.InitName
	params    transactionParams
}

// newInitContractBody returns account transaction body to init the smart contract.
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#initcontract
func newInitContractBody(
	amount *concordium.Amount,
	moduleRef concordium.ModuleRef,
	contractName concordium.ContractName,
	params ...any,
) *initContractBody {
	return &initContractBody{
		amount:    amount,
		moduleRef: moduleRef,
		initName:  concordium.NewInitName(contractName),
		params:    params,
	}
}

func (d *initContractBody) BaseEnergy() int {
	// TODO
	return 10000
}

// Serialize serializes data.
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#serialization-6
func (d *initContractBody) Serialize() ([]byte, error) {
	a, err := d.amount.Serialize()
	if err != nil {
		return nil, err
	}
	m, err := d.moduleRef.Serialize()
	if err != nil {
		return nil, err
	}
	n, err := d.initName.Serialize()
	if err != nil {
		return nil, err
	}
	p, err := d.params.Serialize()
	if err != nil {
		return nil, err
	}
	return d.serialize(typeInitContract, a, m, n, p), nil
}
