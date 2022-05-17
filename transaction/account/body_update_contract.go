package account

import "github.com/Concordium/concordium-go-sdk"

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
