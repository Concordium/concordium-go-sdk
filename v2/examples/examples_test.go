package examples_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"concordium-go-sdk/v2"
)

func TestExamples(t *testing.T) {
	client, err := v2.NewClient(v2.Config{
		NodeAddress: "node.testnet.concordium.com:20000",
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	t.Run("GetBlocks", func(t *testing.T) {
		stream, err := client.GetBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, stream)

		blockInfo, err := stream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)
	})

	t.Run("GetFinalizedBlocks", func(t *testing.T) {
		stream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, stream)

		blockInfo, err := stream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)
	})

	t.Run("GetAccountList", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		consensusInfo, err := client.GetConsensusInfo(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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

	t.Run("GetBakerList", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		bakerStream, err := client.GetBakerList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, bakerStream)

		baker, err := bakerStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, baker)
	})

	t.Run("GetPoolInfo", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		bakerStream, err := client.GetBakerList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, bakerStream)

		baker, err := bakerStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, baker)

		poolInfo, err := client.GetPoolInfo(context.Background(), &v2.PoolInfoRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: block.Height.Value,
					}}},
			Baker: baker,
		})
		require.NoError(t, err)
		require.NotNil(t, poolInfo)
	})

	t.Run("GetPassiveDelegationInfo", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		passiveDelegationInfo, err := client.GetPassiveDelegationInfo(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, passiveDelegationInfo)
	})

	t.Run("GetBlocksAtHeight", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blocksAtHeight, err := client.GetBlocksAtHeight(context.Background(), &v2.BlocksAtHeightRequest{
			BlocksAtHeight: &v2.BlocksAtHeightRequest_Absolute_{
				Absolute: &v2.BlocksAtHeightRequest_Absolute{
					Height: &v2.AbsoluteBlockHeight{
						Value: block.Height.Value,
					},
				}}})
		require.NoError(t, err)
		require.NotNil(t, blocksAtHeight)
	})

	t.Run("GetTokenomicsInfo", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		tokenomicsInfo, err := client.GetTokenomicsInfo(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, tokenomicsInfo)
	})

	t.Run("InvokeInstance", func(t *testing.T) {
		t.Skip()

		// TODO: fill with actual test data.
		instance, err := client.InvokeInstance(context.Background(), &v2.InvokeInstanceRequest{
			BlockHash:  nil,
			Invoker:    nil,
			Instance:   nil,
			Amount:     nil,
			Entrypoint: nil,
			Parameter:  nil,
			Energy:     nil,
		})
		require.NoError(t, err)
		require.NotNil(t, instance)
	})

	t.Run("GetPoolDelegators", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		bakerStream, err := client.GetBakerList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, bakerStream)

		baker, err := bakerStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, baker)

		poolDelegators, err := client.GetPoolDelegators(context.Background(), &v2.GetPoolDelegatorsRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: block.Height.Value,
					}}},
			Baker: baker,
		})
		require.NoError(t, err)
		require.NotNil(t, poolDelegators)
	})

	t.Run("GetPoolDelegatorsRewardPeriod", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		bakerStream, err := client.GetBakerList(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, bakerStream)

		baker, err := bakerStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, baker)

		poolDelegatorsRewardPeriod, err := client.GetPoolDelegatorsRewardPeriod(context.Background(), &v2.GetPoolDelegatorsRequest{
			BlockHash: &v2.BlockHashInput{
				BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
					AbsoluteHeight: &v2.AbsoluteBlockHeight{
						Value: block.Height.Value,
					}}},
			Baker: baker,
		})
		require.NoError(t, err)
		require.NotNil(t, poolDelegatorsRewardPeriod)
	})

	t.Run("GetPassiveDelegators", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		passiveDelegators, err := client.GetPassiveDelegators(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, passiveDelegators)
	})

	t.Run("GetPassiveDelegatorsRewardPeriod", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		passiveDelegatorsRewardPeriod, err := client.GetPassiveDelegatorsRewardPeriod(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, passiveDelegatorsRewardPeriod)
	})

	t.Run("GetBranches", func(t *testing.T) {
		branch, err := client.GetBranches(context.Background())
		require.NoError(t, err)
		require.NotNil(t, branch)
	})

	t.Run("GetElectionInfo", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		electionInfo, err := client.GetElectionInfo(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, electionInfo)
	})

	t.Run("GetIdentityProviders", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		identityProviders, err := client.GetIdentityProviders(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, identityProviders)
	})

	t.Run("GetAnonymityRevokers", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		anonymityRevokers, err := client.GetAnonymityRevokers(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, anonymityRevokers)
	})

	t.Run("GetAccountNonFinalizedTransactions", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
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

		accountNonFinalizedTransactions, err := client.GetAccountNonFinalizedTransactions(context.Background(), accCreds)
		require.NoError(t, err)
		require.NotNil(t, accountNonFinalizedTransactions)
	})

	t.Run("GetBlockTransactionEvents", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blockTransactionEvents, err := client.GetBlockTransactionEvents(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, blockTransactionEvents)
	})

	t.Run("GetBlockSpecialEvents", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blockSpecialEvents, err := client.GetBlockSpecialEvents(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, blockSpecialEvents)
	})

	t.Run("GetBlockPendingUpdates", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blockPendingUpdates, err := client.GetBlockPendingUpdates(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, blockPendingUpdates)
	})

	t.Run("GetNextUpdateSequenceNumbers", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		nextUpdateSequenceNumbers, err := client.GetNextUpdateSequenceNumbers(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, nextUpdateSequenceNumbers)
	})

	t.Run("Shutdown", func(t *testing.T) {
		// required error since method is not enabled
		require.Error(t, client.Shutdown(context.Background()))
	})

	t.Run("GetBannedPeers", func(t *testing.T) {
		// required error since method is not enabled
		_, err := client.GetBannedPeers(context.Background())
		require.Error(t, err)
	})

	t.Run("DumpStart", func(t *testing.T) {
		// required error since method is not enabled
		require.Error(t, client.DumpStart(context.Background(), &v2.DumpRequest{
			File: "random path",
			Raw:  false,
		}))
	})

	t.Run("DumpStop", func(t *testing.T) {
		// required error since method is not enabled
		require.Error(t, client.DumpStop(context.Background()))
	})

	t.Run("GetPeersInfo", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)
	})

	t.Run("PeerDisconnect", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)

		// required error since method is not enabled
		require.Error(t, client.PeerDisconnect(context.Background(), &v2.IpSocketAddress{
			Ip:   peersInfo.Peers[0].SocketAddress.Ip,
			Port: peersInfo.Peers[0].SocketAddress.Port,
		}))
	})

	t.Run("PeerConnect", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)

		// required error since method is not enabled
		require.Error(t, client.PeerDisconnect(context.Background(), &v2.IpSocketAddress{
			Ip:   peersInfo.Peers[0].SocketAddress.Ip,
			Port: peersInfo.Peers[0].SocketAddress.Port,
		}))

		// required error since method is not enabled
		require.Error(t, client.PeerConnect(context.Background(), &v2.IpSocketAddress{
			Ip:   peersInfo.Peers[0].SocketAddress.Ip,
			Port: peersInfo.Peers[0].SocketAddress.Port,
		}))
	})

	t.Run("BanPeer", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)

		// required error since method is not enabled
		require.Error(t, client.BanPeer(context.Background(), &v2.PeerToBan{
			IpAddress: &v2.IpAddress{
				Value: peersInfo.Peers[0].SocketAddress.Ip.Value,
			}}))
	})

	t.Run("UnbanPeer", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)

		// required error since method is not enabled
		require.Error(t, client.BanPeer(context.Background(), &v2.PeerToBan{
			IpAddress: &v2.IpAddress{
				Value: peersInfo.Peers[0].SocketAddress.Ip.Value,
			}}))

		// required error since method is not enabled
		require.Error(t, client.UnbanPeer(context.Background(), &v2.BannedPeer{IpAddress: &v2.IpAddress{
			Value: peersInfo.Peers[0].SocketAddress.Ip.Value,
		}}))
	})

	t.Run("GetNodeInfo", func(t *testing.T) {
		nodeInfo, err := client.GetNodeInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, nodeInfo)
	})

	t.Run("GetBlockChainParameters", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blockChainParameters, err := client.GetBlockChainParameters(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, blockChainParameters)
	})

	t.Run("GetBlockFinalizationSummary", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blockFinalizationSummary, err := client.GetBlockFinalizationSummary(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, blockFinalizationSummary)
	})

	t.Run("GetBlockItems", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		block, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, block)

		blockItems, err := client.GetBlockItems(context.Background(), &v2.BlockHashInput{
			BlockHashInput: &v2.BlockHashInput_AbsoluteHeight{
				AbsoluteHeight: &v2.AbsoluteBlockHeight{
					Value: block.Height.Value,
				}}})
		require.NoError(t, err)
		require.NotNil(t, blockItems)
	})
}
