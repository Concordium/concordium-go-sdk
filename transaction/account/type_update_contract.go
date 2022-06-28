package account

import (
	"fmt"

	"github.com/Concordium/concordium-go-sdk"
)

type UpdateContractResultEvent struct {
	*concordium.TransactionResultEvent `json:""`

	ContractVersion int                         `json:"contractVersion"` // Updated
	Address         *concordium.ContractAddress `json:"address"`         // Updated, Interrupted, Resumed
	Instigator      *concordium.Address         `json:"instigator"`      // Updated
	Amount          *concordium.Amount          `json:"amount"`          // Updated, Transferred
	Message         concordium.Model            `json:"message"`         // Updated
	ReceiveName     concordium.ReceiveName      `json:"receiveName"`     // Updated
	Events          []concordium.Model          `json:"events"`          // Updated, Interrupted
	From            *concordium.Address         `json:"from"`            // Transferred
	To              *concordium.Address         `json:"to"`              // Transferred
	Success         bool                        `json:"success"`         // Resumed
}

type UpdateContractRejectReason struct {
	*concordium.TransactionRejectReason `json:""`
	RejectReason                        int                         `json:"rejectReason"`
	ContractAddress                     *concordium.ContractAddress `json:"contractAddress"`
	ReceiveName                         concordium.ReceiveName      `json:"receiveName"`
	Parameter                           concordium.Model            `json:"parameter"`
}

func (r *UpdateContractRejectReason) Error() error {
	return fmt.Errorf("%q rejected with reason %d", r.ReceiveName, r.RejectReason)
}

type updateContractBody struct {
	baseBody
	amount          *concordium.Amount
	contractAddress *concordium.ContractAddress
	receiveName     concordium.ReceiveName
	params          transactionParams
}

// newUpdateContractBody returns account transaction body to update the smart contract.
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#update-contract
func newUpdateContractBody(
	amount *concordium.Amount,
	contractAddress *concordium.ContractAddress,
	receiveName concordium.ReceiveName,
	params ...any,
) *updateContractBody {
	return &updateContractBody{
		amount:          amount,
		contractAddress: contractAddress,
		receiveName:     receiveName,
		params:          params,
	}
}

func (d *updateContractBody) BaseEnergy() int {
	// TODO
	return 10000
}

// Serialize serializes data.
// See https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#serialization-7
func (d *updateContractBody) Serialize() ([]byte, error) {
	g, err := d.amount.Serialize()
	if err != nil {
		return nil, err
	}
	a, err := d.contractAddress.Serialize()
	if err != nil {
		return nil, err
	}
	n, err := d.receiveName.Serialize()
	if err != nil {
		return nil, err
	}
	p, err := d.params.Serialize()
	if err != nil {
		return nil, err
	}
	return d.serialize(typeUpdateContract, g, a, n, p), nil
}
