package account

import (
	"github.com/Concordium/concordium-go-sdk"
)

type SimpleTransferParams struct {
	// Recipient
	To concordium.AccountAddress
	// Amount to transfer
	Amount *concordium.Amount
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
