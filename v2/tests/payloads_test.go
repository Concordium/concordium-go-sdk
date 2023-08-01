package tests_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

func TestPayloads(t *testing.T) {
	receiver, err := v2.AccountAddressFromBytes(bytes.Repeat([]byte{1}, 32))
	require.NoError(t, err)
	amount := &v2.Amount{
		Value: 100000,
	}
	memo := &v2.Memo{
		Value: []byte{1, 2, 3, 4, 5, 6},
	}
	initName := &v2.InitName{
		Value: "some_name",
	}
	contractAddress := &v2.ContractAddress{
		Index:    1234567890,
		Subindex: 0,
	}
	receiveName := &v2.ReceiveName{
		Value: "some_name",
	}
	parameter := &v2.Parameter{
		Value: bytes.Repeat([]byte{3}, 48),
	}
	registerData := &v2.RegisteredData{
		Value: []byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1},
	}
	deployModule0 := &v2.VersionedModuleSource{
		Module: &v2.ModuleSourceV0{
			Value: bytes.Repeat([]byte{0}, 64),
		}}
	deployModule1 := &v2.VersionedModuleSource{
		Module: &v2.ModuleSourceV1{
			Value: bytes.Repeat([]byte{1}, 64),
		}}
	moduleRef := new(v2.ModuleRef)
	copy(moduleRef.Value[:], bytes.Repeat([]byte{1, 2}, 16))

	t.Run("initContract encode/decode", func(t *testing.T) {
		initContractPayload := &v2.InitContractPayload{
			Amount:    amount,
			ModuleRef: moduleRef,
			InitName:  initName,
			Parameter: parameter,
		}
		encodedInitContractPayload := initContractPayload.Encode()
		require.NotNil(t, encodedInitContractPayload)
		require.NotEmpty(t, encodedInitContractPayload.Value)

		newInitContractPayload, err := encodedInitContractPayload.Decode()
		require.NoError(t, err)
		require.NotNil(t, newInitContractPayload)

		newEncodedInitContractPayload := newInitContractPayload.Payload.Encode()
		require.Equal(t, encodedInitContractPayload.Size(), newEncodedInitContractPayload.Size())
		require.Equal(t, encodedInitContractPayload.Value, newEncodedInitContractPayload.Value)
	})

	t.Run("updateContract encode/decode", func(t *testing.T) {
		updateContractPayload := &v2.UpdateContractPayload{
			Amount:      amount,
			Address:     contractAddress,
			ReceiveName: receiveName,
			Parameter:   parameter,
		}
		encodedUpdateContractPayload := updateContractPayload.Encode()
		require.NotNil(t, encodedUpdateContractPayload)
		require.NotEmpty(t, encodedUpdateContractPayload.Value)

		newUpdateContractPayload, err := encodedUpdateContractPayload.Decode()
		require.NoError(t, err)
		require.NotNil(t, newUpdateContractPayload)

		newEncodedUpdateContractPayload := newUpdateContractPayload.Payload.Encode()
		require.Equal(t, encodedUpdateContractPayload.Size(), newEncodedUpdateContractPayload.Size())
		require.Equal(t, encodedUpdateContractPayload.Value, newEncodedUpdateContractPayload.Value)
	})

	t.Run("transfer encode/decode", func(t *testing.T) {
		transferPayload := &v2.TransferPayload{
			Receiver: &receiver,
			Amount:   amount,
		}
		encodedTransferPayload := transferPayload.Encode()
		require.NotNil(t, encodedTransferPayload)
		require.NotEmpty(t, encodedTransferPayload.Value)

		newTransferPayload, err := encodedTransferPayload.Decode()
		require.NoError(t, err)
		require.NotNil(t, newTransferPayload)

		newEncodedTransferPayload := newTransferPayload.Payload.Encode()
		require.Equal(t, encodedTransferPayload.Size(), newEncodedTransferPayload.Size())
		require.Equal(t, encodedTransferPayload.Value, newEncodedTransferPayload.Value)
	})

	t.Run("transferWithMemo encode/decode", func(t *testing.T) {
		transferWithMemoPayload := &v2.TransferWithMemoPayload{
			Receiver: &receiver,
			Memo:     memo,
			Amount:   amount,
		}
		encodedTransferWithMemoPayload := transferWithMemoPayload.Encode()
		require.NotNil(t, encodedTransferWithMemoPayload)
		require.NotEmpty(t, encodedTransferWithMemoPayload.Value)

		newTransferWithMemoPayload, err := encodedTransferWithMemoPayload.Decode()
		require.NoError(t, err)
		require.NotNil(t, newTransferWithMemoPayload)

		newEncodedTransferWithMemoPayload := newTransferWithMemoPayload.Payload.Encode()
		require.Equal(t, encodedTransferWithMemoPayload.Size(), newEncodedTransferWithMemoPayload.Size())
		require.Equal(t, encodedTransferWithMemoPayload.Value, newEncodedTransferWithMemoPayload.Value)
	})

	t.Run("registerData encode/decode", func(t *testing.T) {
		registerDataPayload := &v2.RegisterDataPayload{
			Data: registerData,
		}
		encodedRegisterDataPayload := registerDataPayload.Encode()
		require.NotNil(t, encodedRegisterDataPayload)
		require.NotEmpty(t, encodedRegisterDataPayload.Value)

		newRegisterDataPayload, err := encodedRegisterDataPayload.Decode()
		require.NoError(t, err)
		require.NotNil(t, newRegisterDataPayload)

		newEncodedRegisterDataPayload := newRegisterDataPayload.Payload.Encode()
		require.Equal(t, encodedRegisterDataPayload.Size(), newEncodedRegisterDataPayload.Size())
		require.Equal(t, encodedRegisterDataPayload.Value, newEncodedRegisterDataPayload.Value)
	})

	t.Run("deployModule V0 encode/decode", func(t *testing.T) {
		deployModulePayload := &v2.DeployModulePayload{
			DeployModule: deployModule0,
		}
		encodedDeployModulePayload := deployModulePayload.Encode()
		require.NotNil(t, encodedDeployModulePayload)
		require.NotEmpty(t, encodedDeployModulePayload.Value)

		newDeployModulePayload, err := encodedDeployModulePayload.Decode()
		require.NoError(t, err)
		require.NotNil(t, newDeployModulePayload)

		newEncodedDeployModulePayload := newDeployModulePayload.Payload.Encode()
		require.Equal(t, encodedDeployModulePayload.Size(), newEncodedDeployModulePayload.Size())
		require.Equal(t, encodedDeployModulePayload.Value, newEncodedDeployModulePayload.Value)
	})

	t.Run("deployModule V1 encode/decode", func(t *testing.T) {
		deployModulePayload := &v2.DeployModulePayload{
			DeployModule: deployModule1,
		}
		encodedDeployModulePayload := deployModulePayload.Encode()
		require.NotNil(t, encodedDeployModulePayload)
		require.NotEmpty(t, encodedDeployModulePayload.Value)

		newDeployModulePayload, err := encodedDeployModulePayload.Decode()
		require.NoError(t, err)
		require.NotNil(t, newDeployModulePayload)

		newEncodedDeployModulePayload := newDeployModulePayload.Payload.Encode()
		require.Equal(t, encodedDeployModulePayload.Size(), newEncodedDeployModulePayload.Size())
		require.Equal(t, encodedDeployModulePayload.Value, newEncodedDeployModulePayload.Value)
	})
}
