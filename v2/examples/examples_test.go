package examples_test

import (
	"concordium-go-sdk/v2"
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

	"concordium-go-sdk/v2/pb"
)

func TestExamples(t *testing.T) {
	conn, err := grpc.Dial(
		"node.testnet.concordium.com:20000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	require.NotNil(t, conn)

	defer conn.Close()

	client := pb.NewQueriesClient(conn)

	t.Run("GetBlocks", func(t *testing.T) {
		stream, err := client.GetBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, stream)

		blockInfo, err := stream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)
	})

	t.Run("GetFinalizedBlocks", func(t *testing.T) {
		stream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, stream)

		blockInfo, err := stream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)
	})

	t.Run("GetAccountList", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		accStream, err := client.GetAccountList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, accStream)
	})

	t.Run("GetAccountInfo", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		accStream, err := client.GetAccountList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, accStream)

		accCreds, err := accStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, accCreds)

		accInfo, err := client.GetAccountInfo(context.Background(), &v2.AccountInfoRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: blockInfo.Height.Value,
					}}},
			AccountIdentifier: &v2.AccountIdentifierInput{
				AccountIdentifierInput: &v2.AccountIdentifierInput_Address{
					Address: &v2.AccountAddress{
						Value: accCreds.Value,
					}}},
		})
		require.NoError(t, err)
		require.NotNil(t, accInfo)
	})

	t.Run("GetModuleList", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		modulesStream, err := client.GetModuleList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, modulesStream)

		module, err := modulesStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, module)
	})

	t.Run("GetAncestors", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		var amount uint64
		// TODO: swap to rand value (add rand func to internal)
		amount = 5

		ancestorsStream, err := client.GetAncestors(context.Background(), &v2.AncestorsRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: blockInfo.Height.Value,
					}}},
			Amount: amount,
		})
		require.NoError(t, err)
		require.NotNil(t, ancestorsStream)

		for i := 0; i < int(amount); i++ {
			blockHash, err := ancestorsStream.Recv()
			require.NoError(t, err)
			require.NotNil(t, blockHash)
		}
	})

	t.Run("invalid options", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		modulesStream, err := client.GetModuleList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, modulesStream)

		module, err := modulesStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, module)

		source, err := client.GetModuleSource(context.Background(), &v2.ModuleSourceRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: blockInfo.Height.Value,
					}}},
			ModuleRef: &v2.ModuleRef{Value: module.Value},
		})
		require.NoError(t, err)
		require.NotNil(t, source)
	})

	t.Run("GetInstanceList", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		instanceStream, err := client.GetInstanceList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, instanceStream)

		instance, err := instanceStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, instance)
	})

	t.Run("GetInstanceInfo", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		instanceStream, err := client.GetInstanceList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, instanceStream)

		instance, err := instanceStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, instance)

		info, err := client.GetInstanceInfo(context.Background(), &v2.InstanceInfoRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: blockInfo.Height.Value,
					}}},
			Address: &v2.ContractAddress{
				Index:    instance.Index,
				Subindex: instance.Subindex,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, info)
	})

	t.Run("GetInstanceState", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		instanceStream, err := client.GetInstanceList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, instanceStream)

		instance, err := instanceStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, instance)

		stateStream, err := client.GetInstanceState(context.Background(), &v2.InstanceInfoRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: blockInfo.Height.Value,
					}}},
			Address: &v2.ContractAddress{
				Index:    instance.Index,
				Subindex: instance.Subindex,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, stateStream)

		state, err := stateStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, state)
	})

	t.Run("InstanceStateLookup", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		instanceStream, err := client.GetInstanceList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, instanceStream)

		instance, err := instanceStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, instance)

		stateStream, err := client.GetInstanceState(context.Background(), &v2.InstanceInfoRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: blockInfo.Height.Value,
					}}},
			Address: &v2.ContractAddress{
				Index:    instance.Index,
				Subindex: instance.Subindex,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, stateStream)

		state, err := stateStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, state)

		valueAtKey, err := client.InstanceStateLookup(context.Background(), &v2.InstanceStateLookupRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: blockInfo.Height.Value,
					}}},
			Address: &v2.ContractAddress{
				Index:    instance.Index,
				Subindex: instance.Subindex,
			},
			Key: state.Key,
		})
		require.NoError(t, err)
		require.NotNil(t, valueAtKey)
	})

	t.Run("GetNextAccountSequenceNumber", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		accStream, err := client.GetAccountList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, accStream)

		accCreds, err := accStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, accCreds)

		accountNum, err := client.GetNextAccountSequenceNumber(context.Background(), &v2.AccountAddress{
			Value: accCreds.Value,
		})
		require.NoError(t, err)
		require.NotNil(t, accountNum)
	})

	t.Run("GetConsensusInfo", func(t *testing.T) {
		consensusInfo, err := client.GetConsensusInfo(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, consensusInfo)
	})

	t.Run("GetBlockItemStatus", func(t *testing.T) {
		t.Skip()

		// TODO: swap for real input when method will be ready
		value := []byte("input")

		status, err := client.GetBlockItemStatus(context.Background(), &v2.TransactionHash{
			Value: value,
		})
		require.NoError(t, err)
		require.NotNil(t, status)
	})

	t.Run("GetCryptographicParameters", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		cryptographicParams, err := client.GetCryptographicParameters(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: blockInfo.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, cryptographicParams)
	})

	t.Run("GetBlockInfo", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background(), new(v2.Empty))
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blockInfo, err := client.GetBlockInfo(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, blockInfo)
	})

	t.Run("GetBlockItemStatus", func(t *testing.T) {
		t.Skip()

		// TODO: swap for real input when method will be ready
		value := []byte("input")

		itemStatus, err := client.GetBlockItemStatus(context.Background(), &v2.TransactionHash{
			Value: value,
		})
		require.NoError(t, err)
		require.NotNil(t, itemStatus)
	})
}
