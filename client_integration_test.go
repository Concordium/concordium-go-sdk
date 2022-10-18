package concordium

import (
	"context"
	"os"
	"testing"
)

var (
	// TODO move to env
	testIntegrationGrpcTarget      = "34.71.98.161:10001"
	testIntegrationBlockHash       = MustNewBlockHash("7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	testIntegrationBlockHeight     = BlockHeight(3335575)
	testIntegrationAccountAddress  = MustNewAccountAddress("4hvvPeHb9HY4Lur7eUZv4KfL3tYBug8DRc4X9cVU8mpJLa1f2X")
	testIntegrationContractAddress = &ContractAddress{Index: 888, SubIndex: 0}
	testIntegrationModuleRef       = MustNewModuleRef("935d17711a4dea10ba5a851df4f19cfdd7cdbd79c8d6ec9abfe5cacff873f6d0")

	testIntegrationBaseClient Client
)

func TestMain(m *testing.M) {
	var err error
	testIntegrationBaseClient, err = NewClient(context.Background(), testIntegrationGrpcTarget)
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func Test_BaseClient_PeerConnect(t *testing.T) {
	// TODO
}

func Test_BaseClient_PeerDisconnect(t *testing.T) {
	// TODO
}

func Test_BaseClient_PeerUptime(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.PeerUptime(ctx)
	if err != nil {
		t.Fatalf("PeerUptime() error = %v", err)
	}
}

func Test_BaseClient_PeerTotalSent(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.PeerTotalSent(ctx)
	if err != nil {
		t.Fatalf("PeerTotalSent() error = %v", err)
	}
}

func Test_BaseClient_PeerTotalReceived(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.PeerTotalReceived(ctx)
	if err != nil {
		t.Fatalf("PeerTotalReceived() error = %v", err)
	}
}

func Test_BaseClient_PeerVersion(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.PeerVersion(ctx)
	if err != nil {
		t.Fatalf("PeerVersion() error = %v", err)
	}
}

func Test_BaseClient_PeerStats(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.PeerStats(ctx, true)
	if err != nil {
		t.Fatalf("PeerStats() error = %v", err)
	}
}

func Test_BaseClient_PeerList(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.PeerList(ctx, true)
	if err != nil {
		t.Fatalf("PeerList() error = %v", err)
	}
}

func Test_BaseClient_BanNode(t *testing.T) {
	// TODO
}

func Test_BaseClient_UnbanNode(t *testing.T) {
	// TODO
}

func Test_BaseClient_JoinNetwork(t *testing.T) {
	// TODO
}

func Test_BaseClient_LeaveNetwork(t *testing.T) {
	// TODO
}

func Test_BaseClient_NodeInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.NodeInfo(ctx)
	if err != nil {
		t.Fatalf("NodeInfo() error = %v", err)
	}
}

func Test_BaseClient_GetConsensusStatus(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetConsensusStatus(ctx)
	if err != nil {
		t.Fatalf("GetConsensusStatus() error = %v", err)
	}
}

func Test_BaseClient_GetBlockInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetBlockInfo(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetBlockInfo() error = %v", err)
	}
}

func Test_BaseClient_GetAncestors(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetAncestors(ctx, testIntegrationBlockHash, 10)
	if err != nil {
		t.Fatalf("GetAncestors() error = %v", err)
	}
}

func Test_BaseClient_GetBranches(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetBranches(ctx)
	if err != nil {
		t.Fatalf("GetBranches() error = %v", err)
	}
}

func Test_BaseClient_GetBlocksAtHeight(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetBlocksAtHeight(ctx, testIntegrationBlockHeight)
	if err != nil {
		t.Fatalf("() error = %v", err)
	}
}

func Test_BaseClient_SendTransaction(t *testing.T) {
	// TODO
}

func Test_BaseClient_StartBaker(t *testing.T) {
	// TODO
}

func Test_BaseClient_StopBaker(t *testing.T) {
	// TODO
}

func Test_BaseClient_GetAccountList(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetAccountList(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetAccountList() error = %v", err)
	}
}

func Test_BaseClient_GetInstances(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetInstances(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetInstances() error = %v", err)
	}
}

func Test_BaseClient_GetAccountInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetAccountInfo(ctx, testIntegrationBlockHash, testIntegrationAccountAddress)
	if err != nil {
		t.Fatalf("GetAccountInfo() error = %v", err)
	}
}

func Test_BaseClient_GetInstanceInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetInstanceInfo(ctx, testIntegrationBlockHash, testIntegrationContractAddress)
	if err != nil {
		t.Fatalf("GetInstanceInfo() error = %v", err)
	}
}

func Test_BaseClient_InvokeContract(t *testing.T) {
	// TODO
}

func Test_BaseClient_GetRewardStatus(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetRewardStatus(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetRewardStatus() error = %v", err)
	}
}

func Test_BaseClient_GetBirkParameters(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetBirkParameters(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetBirkParameters() error = %v", err)
	}
}

func Test_BaseClient_GetModuleList(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetModuleList(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetModuleList() error = %v", err)
	}
}

func Test_BaseClient_GetModuleSource(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetModuleSource(ctx, testIntegrationBlockHash, testIntegrationModuleRef)
	if err != nil {
		t.Fatalf("GetModuleSource() error = %v", err)
	}
}

func Test_BaseClient_GetIdentityProviders(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetIdentityProviders(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetIdentityProviders() error = %v", err)
	}
}

func Test_BaseClient_GetAnonymityRevokers(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetAnonymityRevokers(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetAnonymityRevokers() error = %v", err)
	}
}

func Test_BaseClient_GetCryptographicParameters(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetCryptographicParameters(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetCryptographicParameters() error = %v", err)
	}
}

func Test_BaseClient_GetBannedPeers(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetBannedPeers(ctx)
	if err != nil {
		t.Fatalf("GetBannedPeers() error = %v", err)
	}
}

func Test_BaseClient_Shutdown(t *testing.T) {
	// TODO
}

func Test_BaseClient_DumpStart(t *testing.T) {
	// TODO
}

func Test_BaseClient_DumpStop(t *testing.T) {
	// TODO
}

func Test_BaseClient_GetTransactionStatus(t *testing.T) {
	// TODO
}

func Test_BaseClient_GetTransactionStatusInBlock(t *testing.T) {
	// TODO
}

func Test_BaseClient_GetAccountNonFinalizedTransactions(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetAccountNonFinalizedTransactions(ctx, testIntegrationAccountAddress)
	if err != nil {
		t.Fatalf("GetAccountNonFinalizedTransactions() error = %v", err)
	}
}

func Test_BaseClient_GetBlockSummary(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetBlockSummary(ctx, testIntegrationBlockHash)
	if err != nil {
		t.Fatalf("GetBlockSummary() error = %v", err)
	}
}

func Test_BaseClient_GetNextAccountNonce(t *testing.T) {
	ctx := context.Background()
	_, err := testIntegrationBaseClient.GetNextAccountNonce(ctx, testIntegrationAccountAddress)
	if err != nil {
		t.Fatalf("GetNextAccountNonce() error = %v", err)
	}
}
