package tests_test

import (
	"context"
	"testing"

	v2 "github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/pb"
	"github.com/stretchr/testify/require"
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
		accList, err := client.GetAccountList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, accList)
	})

	t.Run("GetAccountInfo", func(t *testing.T) {
		accList, err := client.GetAccountList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, accList)

		accInfo, err := client.GetAccountInfo(context.Background(),
			&pb.AccountIdentifierInput{
				AccountIdentifierInput: &pb.AccountIdentifierInput_Address{
					Address: &pb.AccountAddress{
						Value: accList[0].Value[:],
					}}},
			v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, accInfo)
	})

	t.Run("GetModuleList", func(t *testing.T) {
		modules, err := client.GetModuleList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, modules)
	})

	t.Run("GetAncestors", func(t *testing.T) {
		ancestorsHash, err := client.GetAncestors(context.Background(), 5, v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, ancestorsHash)
	})

	t.Run("invalid options", func(t *testing.T) {
		modules, err := client.GetModuleList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, modules)

		source, err := client.GetModuleSource(context.Background(), &pb.ModuleSourceRequest{
			BlockHash: &pb.BlockHashInput{
				BlockHashInput: &pb.BlockHashInput_Best{},
			},
			ModuleRef: &pb.ModuleRef{
				Value: modules[0].Value[:],
			},
		})
		require.NoError(t, err)
		require.NotNil(t, source)
	})

	t.Run("GetInstanceList", func(t *testing.T) {
		instanceList, err := client.GetInstanceList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, instanceList)
	})

	t.Run("GetInstanceInfo", func(t *testing.T) {
		instanceList, err := client.GetInstanceList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, instanceList)

		info, err := client.GetInstanceInfo(context.Background(), v2.BlockHashInputBest{}, v2.ContractAddress{
			Index:    instanceList[0].Index,
			Subindex: instanceList[0].Subindex,
		})
		require.NoError(t, err)
		require.NotNil(t, info)
	})

	t.Run("GetInstanceState", func(t *testing.T) {
		instanceList, err := client.GetInstanceList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, instanceList)

		states, err := client.GetInstanceState(context.Background(), v2.BlockHashInputBest{}, v2.ContractAddress{
			Index:    instanceList[0].Index,
			Subindex: instanceList[0].Subindex,
		})
		require.NoError(t, err)
		require.NotNil(t, states)
	})

	t.Run("InstanceStateLookup", func(t *testing.T) {
		instanceList, err := client.GetInstanceList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, instanceList)

		states, err := client.GetInstanceState(context.Background(), v2.BlockHashInputBest{}, v2.ContractAddress{
			Index:    instanceList[0].Index,
			Subindex: instanceList[0].Subindex,
		})
		require.NoError(t, err)
		require.NotNil(t, states)

		state, err := states.Recv()
		require.NoError(t, err)
		require.NotNil(t, state)

		valueAtKey, err := client.InstanceStateLookup(context.Background(), &pb.InstanceStateLookupRequest{
			BlockHash: &pb.BlockHashInput{
				BlockHashInput: &pb.BlockHashInput_Best{}},
			Address: &pb.ContractAddress{
				Index:    instanceList[0].Index,
				Subindex: instanceList[0].Subindex,
			},
			Key: state.Key,
		})
		require.NoError(t, err)
		require.NotNil(t, valueAtKey)
	})

	t.Run("GetNextAccountSequenceNumber", func(t *testing.T) {
		accounts, err := client.GetAccountList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, accounts)

		accountNum, err := client.GetNextAccountSequenceNumber(context.Background(), accounts[0])
		require.NoError(t, err)
		require.NotNil(t, accountNum)
	})

	t.Run("GetConsensusInfo", func(t *testing.T) {
		consensusInfo, err := client.GetConsensusInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, consensusInfo)
	})

	t.Run("GetCryptographicParameters", func(t *testing.T) {
		cryptographicParams, err := client.GetCryptographicParameters(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, cryptographicParams)
	})

	t.Run("GetBlockInfo", func(t *testing.T) {
		blockInfo, err := client.GetBlockInfo(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, blockInfo)
	})

	t.Run("GetBakerList", func(t *testing.T) {
		bakerList, err := client.GetBakerList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, bakerList)
	})

	t.Run("GetPoolInfo", func(t *testing.T) {
		bakerList, err := client.GetBakerList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, bakerList)

		baker, err := bakerList.Recv()
		require.NoError(t, err)
		require.NotNil(t, baker)

		poolInfo, err := client.GetPoolInfo(context.Background(), &pb.PoolInfoRequest{
			BlockHash: &pb.BlockHashInput{
				BlockHashInput: &pb.BlockHashInput_Best{},
			},
			Baker: baker,
		})
		require.NoError(t, err)
		require.NotNil(t, poolInfo)
	})

	t.Run("GetPassiveDelegationInfo", func(t *testing.T) {
		passiveDelegationInfo, err := client.GetPassiveDelegationInfo(context.Background(), v2.BlockHashInputBest{})
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

		blocksAtHeight, err := client.GetBlocksAtHeight(context.Background(), &pb.BlocksAtHeightRequest{
			BlocksAtHeight: &pb.BlocksAtHeightRequest_Absolute_{
				Absolute: &pb.BlocksAtHeightRequest_Absolute{
					Height: &pb.AbsoluteBlockHeight{
						Value: block.Height.Value,
					},
				}}})
		require.NoError(t, err)
		require.NotNil(t, blocksAtHeight)
	})

	t.Run("GetTokenomicsInfo", func(t *testing.T) {
		tokenomicsInfo, err := client.GetTokenomicsInfo(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, tokenomicsInfo)
	})

	t.Run("GetPoolDelegators", func(t *testing.T) {
		bakerList, err := client.GetBakerList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, bakerList)

		baker, err := bakerList.Recv()
		require.NoError(t, err)
		require.NotNil(t, baker)

		poolDelegators, err := client.GetPoolDelegators(context.Background(), &pb.GetPoolDelegatorsRequest{
			BlockHash: &pb.BlockHashInput{
				BlockHashInput: &pb.BlockHashInput_Best{},
			},
			Baker: baker,
		})
		require.NoError(t, err)
		require.NotNil(t, poolDelegators)
	})

	t.Run("GetPoolDelegatorsRewardPeriod", func(t *testing.T) {
		bakerList, err := client.GetBakerList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, bakerList)

		baker, err := bakerList.Recv()
		require.NoError(t, err)
		require.NotNil(t, baker)

		poolDelegatorsRewardPeriod, err := client.GetPoolDelegatorsRewardPeriod(context.Background(), &pb.GetPoolDelegatorsRequest{
			BlockHash: &pb.BlockHashInput{
				BlockHashInput: &pb.BlockHashInput_Best{},
			},
			Baker: baker,
		})
		require.NoError(t, err)
		require.NotNil(t, poolDelegatorsRewardPeriod)
	})

	t.Run("GetPassiveDelegators", func(t *testing.T) {
		passiveDelegators, err := client.GetPassiveDelegators(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, passiveDelegators)
	})

	t.Run("GetPassiveDelegatorsRewardPeriod", func(t *testing.T) {
		passiveDelegatorsRewardPeriod, err := client.GetPassiveDelegatorsRewardPeriod(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, passiveDelegatorsRewardPeriod)
	})

	t.Run("GetBranches", func(t *testing.T) {
		branch, err := client.GetBranches(context.Background())
		require.NoError(t, err)
		require.NotNil(t, branch)
	})

	t.Run("GetElectionInfo", func(t *testing.T) {
		electionInfo, err := client.GetElectionInfo(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, electionInfo)
	})

	t.Run("GetIdentityProviders", func(t *testing.T) {
		identityProviders, err := client.GetIdentityProviders(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, identityProviders)
	})

	t.Run("GetAnonymityRevokers", func(t *testing.T) {
		anonymityRevokers, err := client.GetAnonymityRevokers(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, anonymityRevokers)
	})

	t.Run("GetAccountNonFinalizedTransactions", func(t *testing.T) {
		accounts, err := client.GetAccountList(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, accounts)

		accountNonFinalizedTransactions, err := client.GetAccountNonFinalizedTransactions(context.Background(), accounts[0])
		require.NoError(t, err)
		require.Nil(t, accountNonFinalizedTransactions)
	})

	t.Run("GetBlockTransactionEvents", func(t *testing.T) {
		blocksStream, err := client.GetFinalizedBlocks(context.Background())
		require.NoError(t, err)
		require.NotNil(t, blocksStream)

		blockInfo, err := blocksStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, blockInfo)

		var blockHash v2.BlockHash
		copy(blockHash.Value[:], blockInfo.Hash.Value)

		blockTransactionEvents, err := client.GetBlockTransactionEvents(context.Background(), v2.BlockHashInputGiven{
			Given: blockHash,
		})
		require.NoError(t, err)
		require.NotNil(t, blockTransactionEvents)
	})

	t.Run("GetBlockSpecialEvents", func(t *testing.T) {
		blockSpecialEvents, err := client.GetBlockSpecialEvents(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, blockSpecialEvents)
	})

	t.Run("GetBlockPendingUpdates", func(t *testing.T) {
		blockPendingUpdates, err := client.GetBlockPendingUpdates(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, blockPendingUpdates)
	})

	t.Run("GetNextUpdateSequenceNumbers", func(t *testing.T) {
		nextUpdateSequenceNumbers, err := client.GetNextUpdateSequenceNumbers(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, nextUpdateSequenceNumbers)
	})

	t.Run("Shutdown", func(t *testing.T) {
		// required error since method is not enabled.
		require.Error(t, client.Shutdown(context.Background()))
	})

	t.Run("GetBannedPeers", func(t *testing.T) {
		// required error since method is not enabled.
		_, err := client.GetBannedPeers(context.Background())
		require.Error(t, err)
	})

	t.Run("DumpStart", func(t *testing.T) {
		// required error since method is not enabled.
		require.Error(t, client.DumpStart(context.Background(), &pb.DumpRequest{
			File: "random path",
			Raw:  false,
		}))
	})

	t.Run("DumpStop", func(t *testing.T) {
		// required error since method is not enabled.
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

		// required error since method is not enabled.
		require.Error(t, client.PeerDisconnect(context.Background(), &pb.IpSocketAddress{
			Ip:   peersInfo.Peers[0].SocketAddress.Ip,
			Port: peersInfo.Peers[0].SocketAddress.Port,
		}))
	})

	t.Run("PeerConnect", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)

		// required error since method is not enabled.
		require.Error(t, client.PeerDisconnect(context.Background(), &pb.IpSocketAddress{
			Ip:   peersInfo.Peers[0].SocketAddress.Ip,
			Port: peersInfo.Peers[0].SocketAddress.Port,
		}))

		// required error since method is not enabled.
		require.Error(t, client.PeerConnect(context.Background(), &pb.IpSocketAddress{
			Ip:   peersInfo.Peers[0].SocketAddress.Ip,
			Port: peersInfo.Peers[0].SocketAddress.Port,
		}))
	})

	t.Run("BanPeer", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)

		// required error since method is not enabled.
		require.Error(t, client.BanPeer(context.Background(), &pb.PeerToBan{
			IpAddress: &pb.IpAddress{
				Value: peersInfo.Peers[0].SocketAddress.Ip.Value,
			}}))
	})

	t.Run("UnbanPeer", func(t *testing.T) {
		peersInfo, err := client.GetPeersInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, peersInfo)

		// required error since method is not enabled.
		require.Error(t, client.BanPeer(context.Background(), &pb.PeerToBan{
			IpAddress: &pb.IpAddress{
				Value: peersInfo.Peers[0].SocketAddress.Ip.Value,
			}}))

		// required error since method is not enabled.
		require.Error(t, client.UnbanPeer(context.Background(), &pb.BannedPeer{IpAddress: &pb.IpAddress{
			Value: peersInfo.Peers[0].SocketAddress.Ip.Value,
		}}))
	})

	t.Run("GetNodeInfo", func(t *testing.T) {
		nodeInfo, err := client.GetNodeInfo(context.Background())
		require.NoError(t, err)
		require.NotNil(t, nodeInfo)
	})

	t.Run("GetBlockChainParameters", func(t *testing.T) {
		blockChainParameters, err := client.GetBlockChainParameters(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, blockChainParameters)
	})

	t.Run("GetBlockFinalizationSummary", func(t *testing.T) {
		blockFinalizationSummary, err := client.GetBlockFinalizationSummary(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, blockFinalizationSummary)
	})

	t.Run("GetBlockItems", func(t *testing.T) {
		blockItems, err := client.GetBlockItems(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.Nil(t, blockItems)
	})

	t.Run("GetFirstBlockEpoch", func(t *testing.T) {
		blochHash, err := client.GetFirstBlockEpoch(context.Background(), v2.EpochRequestRelativeEpoch{
			GenesisIndex: v2.GenesisIndex{Value: 3},
			Epoch:        v2.Epoch{Value: 5},
		})
		require.NoError(t, err)
		require.NotNil(t, blochHash)
	})

	t.Run("GetWinningBakersEpoch", func(t *testing.T) {
		winningBakerStream, err := client.GetWinningBakersEpoch(context.Background(), v2.EpochRequestRelativeEpoch{
			GenesisIndex: v2.GenesisIndex{Value: 3},
			Epoch:        v2.Epoch{Value: 5},
		})
		require.NoError(t, err)
		require.NotNil(t, winningBakerStream)

		winningBaker, err := winningBakerStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, winningBaker)
	})

	t.Run("GetBlockCertificates", func(t *testing.T) {
		certificate, err := client.GetBlockCertificates(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, certificate)
	})

	t.Run("GetBakersRewardPeriod", func(t *testing.T) {
		bakerRewardPeriodInfoStream, err := client.GetBakersRewardPeriod(context.Background(), v2.BlockHashInputBest{})
		require.NoError(t, err)
		require.NotNil(t, bakerRewardPeriodInfoStream)

		bakerRewardPeriodInfo, err := bakerRewardPeriodInfoStream.Recv()
		require.NoError(t, err)
		require.NotNil(t, bakerRewardPeriodInfo)
	})

	t.Run("GetBakerEarliestWintime", func(t *testing.T) {
		timestamp, err := client.GetBakerEarliestWinTime(context.Background(), &pb.BakerId{Value: 1})
		require.NoError(t, err)
		require.NotNil(t, timestamp)
	})
}
