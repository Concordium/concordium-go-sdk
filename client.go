package concordium

import (
	"context"
	"io"
)

type BaseClient interface {
	// PeerConnect Suggest to a peer to connect to the submitted peer details.
	// This, if successful, adds the peer to the list of given addresses.
	PeerConnect(ctx context.Context, ip string, port int) (bool, error)

	// PeerDisconnect Disconnect from the peer and remove them from the given addresses list
	// if they are on it. Return if the request was processed successfully.
	PeerDisconnect(ctx context.Context, ip string, port int) (bool, error)

	// PeerUptime Peer uptime
	PeerUptime(ctx context.Context) (int64, error)

	// PeerTotalSent Peer total number of sent packets
	PeerTotalSent(ctx context.Context) (int, error)

	// PeerTotalReceived Peer total number of received packets
	PeerTotalReceived(ctx context.Context) (int, error)

	// PeerVersion Peer client software version
	PeerVersion(ctx context.Context) (string, error)

	// PeerStats Stats for connected peers
	PeerStats(ctx context.Context, includeBootstrappers bool) (*PeerStats, error)

	// PeerList List of connected peers
	PeerList(ctx context.Context, includeBootstrappers bool) (*PeerList, error)

	// BanNode ...
	BanNode(ctx context.Context, element PeerElement) (bool, error)

	// UnbanNode ...
	UnbanNode(ctx context.Context, element PeerElement) (bool, error)

	// JoinNetwork ...
	JoinNetwork(ctx context.Context, id NetworkId) (bool, error)

	// LeaveNetwork ...
	LeaveNetwork(ctx context.Context, id NetworkId) (bool, error)

	// NodeInfo Get information about the running Node
	NodeInfo(ctx context.Context) (*NodeInfo, error)

	// GetConsensusStatus see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getconsensusstatus--consensusstatus
	GetConsensusStatus(ctx context.Context) (*ConsensusStatus, error)

	// GetBlockInfo see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getblockinfo--blockhash---blockinfo
	GetBlockInfo(ctx context.Context, hash BlockHash) (*BlockInfo, error)

	// GetAncestors see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getancestors--blockhash---blockhash
	GetAncestors(ctx context.Context, hash BlockHash, amount int) ([]BlockHash, error)

	// GetBranches see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getbranches--branch
	GetBranches(ctx context.Context) (*Branch, error)

	// GetBlocksAtHeight Get the blocks at the given height
	GetBlocksAtHeight(ctx context.Context, height BlockHeight) ([]BlockHash, error)

	// SendTransaction Submit a local transaction
	SendTransaction(ctx context.Context, id NetworkId, payload io.Reader) (bool, error)

	// StartBaker Start the baker in the consensus module
	StartBaker(ctx context.Context) (bool, error)

	// StopBaker Stop the baker in the consensus module
	StopBaker(ctx context.Context) (bool, error)

	// GetAccountList see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getaccountlist--blockhash---accountaddress
	GetAccountList(ctx context.Context, hash BlockHash) ([]AccountAddress, error)

	// GetInstances see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getinstances--blockhash---contractaddress
	GetInstances(ctx context.Context, hash BlockHash) ([]*ContractAddress, error)

	// GetAccountInfo see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getaccountinfo--blockhash---accountaddress---accountinfo
	GetAccountInfo(ctx context.Context, hash BlockHash, address AccountAddress) (*AccountInfo, error)

	// GetInstanceInfo see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getinstanceinfo--blockhash---contractaddress---contractinfo
	GetInstanceInfo(ctx context.Context, hash BlockHash, address *ContractAddress) (*InstanceInfo, error)

	// InvokeContract see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#invokecontract--blockhash---contractcontext---invokecontractresult
	InvokeContract(interface{}) ([]byte, error) // TODO

	// GetRewardStatus see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getrewardstatus-blockhash---rewardstatus
	GetRewardStatus(ctx context.Context, hash BlockHash) (*RewardStatus, error)

	// GetBirkParameters see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getbirkparameters--blockhash---birkparameters
	GetBirkParameters(ctx context.Context, hash BlockHash) (*BirkParameters, error)

	// GetModuleList see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getmodulelist--blockhash---moduleref
	GetModuleList(ctx context.Context, hash BlockHash) ([]ModuleRef, error)

	// GetModuleSource see https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getmodulesource--blockhash---moduleref---modulesource
	GetModuleSource(ctx context.Context, hash BlockHash, ref ModuleRef) (io.Reader, error)

	// GetIdentityProviders ...
	GetIdentityProviders(ctx context.Context, hash BlockHash) ([]*IdentityProvider, error)

	// GetAnonymityRevokers ...
	GetAnonymityRevokers(ctx context.Context, hash BlockHash) ([]*AnonymityRevoker, error)

	// GetCryptographicParameters ...
	GetCryptographicParameters(ctx context.Context, hash BlockHash) ([]byte, error) // TODO

	// GetBannedPeers ...
	GetBannedPeers(ctx context.Context) (*PeerList, error)

	// Shutdown ...
	Shutdown(ctx context.Context) (bool, error)

	// DumpStart ...
	DumpStart(ctx context.Context, file string, raw bool) (bool, error)

	// DumpStop ...
	DumpStop(ctx context.Context) (bool, error)

	// GetTransactionStatus Query for the status of a transaction by its hash
	GetTransactionStatus(ctx context.Context, hash TransactionHash) ([]byte, error) // TODO

	// GetTransactionStatusInBlock Query for transactions in a block by its hash
	GetTransactionStatusInBlock(ctx context.Context, hash TransactionHash, blockHash BlockHash) ([]byte, error) // TODO

	// GetAccountNonFinalizedTransactions Query for non-finalized
	// transactions present on an account by the account address
	GetAccountNonFinalizedTransactions(ctx context.Context, address AccountAddress) ([]byte, error) // TODO

	// GetBlockSummary Request a summary for a block by its hash
	GetBlockSummary(ctx context.Context, hash BlockHash) ([]byte, error) // TODO

	// GetNextAccountNonce Request next nonce information for an account
	GetNextAccountNonce(ctx context.Context, address AccountAddress) (*NextAccountNonce, error)
}
