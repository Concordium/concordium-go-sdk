package account

import (
	"encoding/json"
	"fmt"

	"github.com/Concordium/concordium-go-sdk"
)

type SimpleTransferResultEvent struct {
	Tag    concordium.TransactionResultEventTag `json:"tag"`
	From   *concordium.Address                  `json:"from"`
	Amount *concordium.Amount                   `json:"amount"`
	To     *concordium.Address                  `json:"to"`
}

func NewSimpleTransferResultEvent(origin *concordium.TransactionResultEvent) (*SimpleTransferResultEvent, error) {
	e := &SimpleTransferResultEvent{}
	err := json.Unmarshal(origin.Raw, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

type SimpleTransferRejectReason struct {
	Tag concordium.TransactionRejectReasonTag `json:"tag"`

	// TODO
	// Array of different types? Really???!!!
	// [
	//   {
	//     "type": "AddressAccount",
	//     "address": "3rsc7HNLVKnFz9vmKkAaEMVpNkFA4hZxJpZinCtUTJbBh58yYi"
	//   },
	//   "999999999999999999"
	// ]
	Contents []any `json:"contents"`
}

func NewSimpleTransferRejectReason(origin *concordium.TransactionRejectReason) (*SimpleTransferRejectReason, error) {
	r := &SimpleTransferRejectReason{}
	err := json.Unmarshal(origin.Raw, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *SimpleTransferRejectReason) Error() error {
	return fmt.Errorf("%s", r.Tag)
}

type simpleTransferBody struct {
	baseBody
	to     concordium.AccountAddress
	amount *concordium.Amount
}

// newSimpleTransferBody returns account transaction body to simple transfer.
func newSimpleTransferBody(to concordium.AccountAddress, amount *concordium.Amount) *simpleTransferBody {
	return &simpleTransferBody{
		to:     to,
		amount: amount,
	}
}

func (d *simpleTransferBody) BaseEnergy() int {
	return 300
}

func (d *simpleTransferBody) Serialize() ([]byte, error) {
	a, err := d.to.Serialize()
	if err != nil {
		return nil, err
	}
	g, err := d.amount.Serialize()
	if err != nil {
		return nil, err
	}
	return d.serialize(typeSimpleTransfer, a, g), nil
}
