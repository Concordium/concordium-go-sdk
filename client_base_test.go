package concordium

import (
	"context"
	"testing"
)

func Test_BaseClient_PeerConnect(t *testing.T) {
	// TODO
}

func Test_BaseClient_PeerDisconnect(t *testing.T) {
	// TODO
}

func Test_BaseClient_PeerUptime(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.PeerUptime(ctx)
	if err != nil {
		t.Fatalf("PeerUptime() error = %v", err)
	}
}

func Test_BaseClient_PeerTotalSent(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.PeerTotalSent(ctx)
	if err != nil {
		t.Fatalf("PeerTotalSent() error = %v", err)
	}
}

func Test_BaseClient_PeerTotalReceived(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.PeerTotalReceived(ctx)
	if err != nil {
		t.Fatalf("PeerTotalReceived() error = %v", err)
	}
}

func Test_BaseClient_PeerVersion(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.PeerVersion(ctx)
	if err != nil {
		t.Fatalf("PeerVersion() error = %v", err)
	}
}

func Test_BaseClient_PeerStats(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.PeerStats(ctx, true)
	if err != nil {
		t.Fatalf("PeerStats() error = %v", err)
	}
}

func Test_BaseClient_PeerList(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.PeerList(ctx, true)
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
	_, err := testBaseClient.NodeInfo(ctx)
	if err != nil {
		t.Fatalf("NodeInfo() error = %v", err)
	}
}

func Test_BaseClient_GetConsensusStatus(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetConsensusStatus(ctx)
	if err != nil {
		t.Fatalf("GetConsensusStatus() error = %v", err)
	}
}

func Test_BaseClient_GetBlockInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetBlockInfo(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetBlockInfo() error = %v", err)
	}
}

func Test_BaseClient_GetAncestors(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetAncestors(ctx, testBlockHash, 10)
	if err != nil {
		t.Fatalf("GetAncestors() error = %v", err)
	}
}

func Test_BaseClient_GetBranches(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetBranches(ctx)
	if err != nil {
		t.Fatalf("GetBranches() error = %v", err)
	}
}

func Test_BaseClient_GetBlocksAtHeight(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetBlocksAtHeight(ctx, testBlockHeight)
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
	_, err := testBaseClient.GetAccountList(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetAccountList() error = %v", err)
	}
}

func Test_BaseClient_GetInstances(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetInstances(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetInstances() error = %v", err)
	}
}

func Test_BaseClient_GetAccountInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetAccountInfo(ctx, testBlockHash, testAccountAddress)
	if err != nil {
		t.Fatalf("GetAccountInfo() error = %v", err)
	}
}

func Test_BaseClient_GetInstanceInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetInstanceInfo(ctx, testBlockHash, testContractAddress)
	if err != nil {
		t.Fatalf("GetInstanceInfo() error = %v", err)
	}
}

func Test_BaseClient_InvokeContract(t *testing.T) {
	// TODO
}

func Test_BaseClient_GetRewardStatus(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetRewardStatus(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetRewardStatus() error = %v", err)
	}
}

func Test_BaseClient_GetBirkParameters(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetBirkParameters(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetBirkParameters() error = %v", err)
	}
}

func Test_BaseClient_GetModuleList(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetModuleList(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetModuleList() error = %v", err)
	}
}

func Test_BaseClient_GetModuleSource(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetModuleSource(ctx, testBlockHash, testModuleRef)
	if err != nil {
		t.Fatalf("GetModuleSource() error = %v", err)
	}
}

func Test_BaseClient_GetIdentityProviders(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetIdentityProviders(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetIdentityProviders() error = %v", err)
	}
}

func Test_BaseClient_GetAnonymityRevokers(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetAnonymityRevokers(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetAnonymityRevokers() error = %v", err)
	}
}

func Test_BaseClient_GetCryptographicParameters(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetCryptographicParameters(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetCryptographicParameters() error = %v", err)
	}
}

func Test_BaseClient_GetBannedPeers(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetBannedPeers(ctx)
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
	_, err := testBaseClient.GetAccountNonFinalizedTransactions(ctx, testAccountAddress)
	if err != nil {
		t.Fatalf("GetAccountNonFinalizedTransactions() error = %v", err)
	}
}

func Test_BaseClient_GetBlockSummary(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetBlockSummary(ctx, testBlockHash)
	if err != nil {
		t.Fatalf("GetBlockSummary() error = %v", err)
	}
}

func Test_BaseClient_GetNextAccountNonce(t *testing.T) {
	ctx := context.Background()
	_, err := testBaseClient.GetNextAccountNonce(ctx, testAccountAddress)
	if err != nil {
		t.Fatalf("GetNextAccountNonce() error = %v", err)
	}
}
