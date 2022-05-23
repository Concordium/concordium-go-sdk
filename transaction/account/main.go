package account

import (
	"encoding/binary"
	"fmt"
	"github.com/Concordium/concordium-go-sdk"
)

const (
	typeDeployModule   uint8 = 0
	typeInitContract   uint8 = 1
	typeUpdateContract uint8 = 2
	typeSimpleTransfer uint8 = 3
	//typeAddBaker                     uint8 = 4
	//typeRemoveBaker                  uint8 = 5
	//typeUpdateBakerStake             uint8 = 6
	//typeUpdateBakerReStakeEarnings   uint8 = 7
	//typeUpdateBakerKeys              uint8 = 8
	//typeUpdateCredentialKeys         uint8 = 13
	//typeEncryptedTransfer            uint8 = 16
	//typeTransferToEncrypted          uint8 = 17
	//typeTransferToPublic             uint8 = 18
	//typeTransferWithSchedule         uint8 = 19
	//typeUpdateCredentials            uint8 = 20
	//typeRegisterData                 uint8 = 21
	//typeSimpleTransferWithMemo       uint8 = 22
	//typeEncryptedTransferWithMemo    uint8 = 23
	//typeTransferWithScheduleWithMemo uint8 = 24

	headerSize = 32 + 8 + 8 + 4 + 8 // address + nonce + energy + body size + expired at

	// These constants must be consistent with constA and constB in:
	// https://github.com/Concordium/concordium-base/blob/main/haskell-src/Concordium/Cost.hs
	energyConstA = 100
	energyConstB = 1
)

// calculateTransactionEnergy returns energy value according to the formula:
// A * signatureCount + B * size + C_t, where C_t is a transaction specific cost.
// The transaction specific cost can be found at
// https://github.com/Concordium/concordium-base/blob/main/haskell-src/Concordium/Cost.hs.
func calculateTransactionEnergy(signatureCount, txSize, txBaseEnergy int) uint64 {
	return uint64(energyConstA*signatureCount + energyConstB*txSize + txBaseEnergy)
}

type transactionParams []any

func (s transactionParams) Serialize() ([]byte, error) {
	size := 0
	params := make([][]byte, len(s))
	for i, p := range s {
		b, err := concordium.SerializeModel(p)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize param: %w", err)
		}
		size += len(b)
		params[i] = b
	}
	i := 2
	b := make([]byte, size+i)
	binary.BigEndian.PutUint16(b, uint16(size))
	for _, p := range params {
		copy(b[i:], p)
		i += len(p)
	}
	return b, nil
}
