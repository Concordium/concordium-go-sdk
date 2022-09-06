package concordium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	grpc_api "github.com/Concordium/concordium-go-sdk/grpc-api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Client interface {
	// PeerConnect suggest the node to connect to the submitted peer. If successful,
	// this adds the peer to the list of peers.
	PeerConnect(ctx context.Context, ip string, port int32) (bool, error)

	// PeerDisconnect disconnect from the peer and remove them from the given
	// addresses list if they are on it.
	PeerDisconnect(ctx context.Context, ip string, port int32) (bool, error)

	// PeerUptime get the uptime of the node in milliseconds.
	PeerUptime(ctx context.Context) (uint64, error)

	// PeerTotalSent get the total number of packets sent by the node.
	PeerTotalSent(ctx context.Context) (uint64, error)

	// PeerTotalReceived get the total number of packets received by the node.
	PeerTotalReceived(ctx context.Context) (uint64, error)

	// PeerVersion get the version of the node software.
	PeerVersion(ctx context.Context) (string, error)

	// PeerStats get information on the peers that the node is connected to.
	PeerStats(ctx context.Context, includeBootstrappers bool) (*PeerStats, error)

	// PeerList get a list of the peers that the node is connected to.
	PeerList(ctx context.Context, includeBootstrappers bool) (*PeerList, error)

	// BanNode ban a node from being a peer. Note that you should provide a nodeId
	// OR an ip, but not both. Use an empty value for the option not chosen.
	BanNode(ctx context.Context, nodeId NodeId, ip string) (bool, error)

	// UnbanNode unban a previously banned node. Note that you should provide a nodeId
	// OR an ip, but not both. Use an empty value for the option not chosen.
	UnbanNode(ctx context.Context, nodeId NodeId, ip string) (bool, error)

	// JoinNetwork attempt to join the specified network.
	JoinNetwork(ctx context.Context, networkId NetworkId) (bool, error)

	// LeaveNetwork attempt to leave the specified network.
	LeaveNetwork(ctx context.Context, networkId NetworkId) (bool, error)

	// NodeInfo get information about the running node.
	NodeInfo(ctx context.Context) (*NodeInfo, error)

	// GetConsensusStatus get the information about the consensus.
	GetConsensusStatus(ctx context.Context) (*ConsensusStatus, error)

	// GetBlockInfo get information, such as height, timings, and transaction
	// counts for the given block.
	GetBlockInfo(ctx context.Context, blockHash BlockHash) (*BlockInfo, error)

	// GetAncestors get a list of the blocks preceding the given block. The list will contain
	// at most amount blocks.
	GetAncestors(ctx context.Context, blockHash BlockHash, amount int) ([]BlockHash, error)

	// GetBranches get the branches of the tree. This is the part of the tree above the last
	// finalized block.
	GetBranches(ctx context.Context) (*Branch, error)

	// GetBlocksAtHeight get a list of the blocks at the given height.
	GetBlocksAtHeight(ctx context.Context, blockHeight BlockHeight) ([]BlockHash, error)

	// SendTransactionAsync sends a transaction to the given network. The node will do basic
	// transaction validation, such as signature checks and account nonce checks, and if these
	// fail, the call will return an error.
	SendTransactionAsync(ctx context.Context, networkId NetworkId, request TransactionRequest) (TransactionHash, error)

	// SendTransactionAwait sends a transaction to the given network. and await its finalization.
	SendTransactionAwait(ctx context.Context, networkId NetworkId, request TransactionRequest) (*TransactionSummary, TransactionHash, error)

	// StartBaker start the baker.
	StartBaker(ctx context.Context) (bool, error)

	// StopBaker stop the baker.
	StopBaker(ctx context.Context) (bool, error)

	// GetAccountList get a list of all accounts that exist in the state at the end
	// of the given block.
	GetAccountList(ctx context.Context, blockHash BlockHash) ([]AccountAddress, error)

	// GetInstances get a list of all smart contract instances that exist in the
	// state at the end of the given block.
	GetInstances(ctx context.Context, blockHash BlockHash) ([]*ContractAddress, error)

	// GetAccountInfo get the state of an account in the given block.
	GetAccountInfo(ctx context.Context, blockHash BlockHash, accountAddress AccountAddress) (*AccountInfo, error)

	// GetInstanceInfo get information about the given smart contract instance in the given block.
	GetInstanceInfo(ctx context.Context, blockHash BlockHash, contractAddress *ContractAddress) (*InstanceInfo, error)

	// InvokeContract invokes a smart contract instance and view its results as if it had been
	// updated at the end of the given block. Please note that this is not a transaction, so it
	// won’t affect the contract on chain. It only simulates the invocation.
	InvokeContract(ctx context.Context, blockHash BlockHash, context *ContractContext) (*InvokeContractResult, error)

	// GetRewardStatus get an overview of the balance of special accounts in the given block.
	GetRewardStatus(ctx context.Context, blockHash BlockHash) (*RewardStatus, error)

	// GetBirkParameters get an overview of the parameters used for baking.
	GetBirkParameters(ctx context.Context, blockHash BlockHash) (*BirkParameters, error)

	// GetModuleList get a list of all smart contract modules that exist in the state at the end of the given block.
	GetModuleList(ctx context.Context, blockHash BlockHash) ([]ModuleRef, error)

	// GetModuleSource get the binary source of a smart contract module.
	GetModuleSource(ctx context.Context, blockHash BlockHash, moduleRef ModuleRef) (io.Reader, error)

	// GetIdentityProviders get a list of all identity providers that exist in the state at the end of the given block.
	GetIdentityProviders(ctx context.Context, blockHash BlockHash) ([]*IdentityProvider, error)

	// GetAnonymityRevokers get a list of all anonymity revokers that exist in the state at the end of the given block.
	GetAnonymityRevokers(ctx context.Context, blockHash BlockHash) ([]*AnonymityRevoker, error)

	// GetCryptographicParameters get the cryptographic parameters used in the given block.
	GetCryptographicParameters(ctx context.Context, blockHash BlockHash) (*CryptographicParameters, error)

	// GetBannedPeers get a list of banned peers.
	GetBannedPeers(ctx context.Context) (*PeerList, error)

	// Shutdown shut down the node.
	Shutdown(ctx context.Context) (bool, error)

	// DumpStart start dumping packages into the specified file. Only available
	// on a node built with the network_dump feature.
	DumpStart(ctx context.Context, file string, raw bool) (bool, error)

	// DumpStop stop dumping packages. Only available on a node built with the network_dump feature.
	DumpStop(ctx context.Context) (bool, error)

	// GetTransactionStatus
	GetTransactionStatus(ctx context.Context, hash TransactionHash) (*TransactionSummary, error)

	// GetTransactionStatusInBlock
	GetTransactionStatusInBlock(ctx context.Context, hash TransactionHash, blockHash BlockHash) (*TransactionSummary, error)

	// GetAccountNonFinalizedTransactions get a list of non-finalized
	// transactions present on an account.
	GetAccountNonFinalizedTransactions(ctx context.Context, accountAddress AccountAddress) ([]TransactionHash, error)

	// GetBlockSummary get a summary of the transactions and data in a given block.
	GetBlockSummary(ctx context.Context, blockHash BlockHash) (*BlockSummary, error)

	// GetNextAccountNonce returns the next available nonce for this account.
	GetNextAccountNonce(ctx context.Context, accountAddress AccountAddress) (*NextAccountNonce, error)
}

type perRPCCredentials string

func (c perRPCCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authentication": string(c),
	}, nil
}

func (c perRPCCredentials) RequireTransportSecurity() bool {
	return false
}

type client struct {
	grpc grpc_api.P2PClient
}

func NewClient(ctx context.Context, target, token string) (Client, error) {
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithPerRPCCredentials(perRPCCredentials(token)))
	if err != nil {
		return nil, err
	}
	cli := &client{
		grpc: grpc_api.NewP2PClient(conn),
	}
	return cli, nil
}

func (c *client) PeerConnect(ctx context.Context, ip string, port int32) (bool, error) {
	res, err := c.grpc.PeerConnect(ctx, &grpc_api.PeerConnectRequest{
		Ip: &wrapperspb.StringValue{
			Value: ip,
		},
		Port: &wrapperspb.Int32Value{
			Value: port,
		},
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) PeerDisconnect(ctx context.Context, ip string, port int32) (bool, error) {
	res, err := c.grpc.PeerDisconnect(ctx, &grpc_api.PeerConnectRequest{
		Ip: &wrapperspb.StringValue{
			Value: ip,
		},
		Port: &wrapperspb.Int32Value{
			Value: port,
		},
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) PeerUptime(ctx context.Context) (uint64, error) {
	res, err := c.grpc.PeerUptime(ctx, &grpc_api.Empty{})
	if err != nil {
		return 0, err
	}
	return res.GetValue(), nil
}

func (c *client) PeerTotalSent(ctx context.Context) (uint64, error) {
	res, err := c.grpc.PeerTotalSent(ctx, &grpc_api.Empty{})
	if err != nil {
		return 0, err
	}
	return res.GetValue(), nil
}

func (c *client) PeerTotalReceived(ctx context.Context) (uint64, error) {
	res, err := c.grpc.PeerTotalReceived(ctx, &grpc_api.Empty{})
	if err != nil {
		return 0, err
	}
	return res.GetValue(), nil
}

func (c *client) PeerVersion(ctx context.Context) (string, error) {
	res, err := c.grpc.PeerVersion(ctx, &grpc_api.Empty{})
	if err != nil {
		return "", err
	}
	return res.GetValue(), nil
}

func (c *client) PeerStats(ctx context.Context, includeBootstrappers bool) (*PeerStats, error) {
	res, err := c.grpc.PeerStats(ctx, &grpc_api.PeersRequest{
		IncludeBootstrappers: includeBootstrappers,
	})
	if err != nil {
		return nil, err
	}
	s := &PeerStats{
		Peers:     make([]*PeerStatsElement, len(res.Peerstats)),
		AvgBpsIn:  res.AvgBpsIn,
		AvgBpsOut: res.AvgBpsOut,
	}
	for i, e := range res.Peerstats {
		s.Peers[i] = &PeerStatsElement{
			NodeId:          NodeId(e.NodeId),
			PacketsSent:     e.PacketsSent,
			PacketsReceived: e.PacketsReceived,
			Latency:         e.Latency,
		}
	}
	return s, nil
}

func (c *client) PeerList(ctx context.Context, includeBootstrappers bool) (*PeerList, error) {
	res, err := c.grpc.PeerList(ctx, &grpc_api.PeersRequest{
		IncludeBootstrappers: includeBootstrappers,
	})
	if err != nil {
		return nil, err
	}
	l := &PeerList{
		Type:  PeerType(res.PeerType),
		Peers: make([]*PeerElement, len(res.Peers)),
	}
	for i, e := range res.Peers {
		l.Peers[i] = &PeerElement{
			NodeId:        NodeId(e.NodeId.GetValue()),
			Ip:            e.Ip.GetValue(),
			Port:          e.Port.GetValue(),
			CatchupStatus: PeerElementCatchupStatus(e.CatchupStatus),
		}
	}
	return l, nil
}

func (c *client) BanNode(ctx context.Context, nodeId NodeId, ip string) (bool, error) {
	req := &grpc_api.PeerElement{}
	if nodeId != "" {
		req.NodeId = &wrapperspb.StringValue{
			Value: string(nodeId),
		}
	}
	if ip != "" {
		req.Ip = &wrapperspb.StringValue{
			Value: ip,
		}
	}
	res, err := c.grpc.BanNode(ctx, req)
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) UnbanNode(ctx context.Context, nodeId NodeId, ip string) (bool, error) {
	req := &grpc_api.PeerElement{}
	if nodeId != "" {
		req.NodeId = &wrapperspb.StringValue{
			Value: string(nodeId),
		}
	}
	if ip != "" {
		req.Ip = &wrapperspb.StringValue{
			Value: ip,
		}
	}
	res, err := c.grpc.UnbanNode(ctx, req)
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) JoinNetwork(ctx context.Context, networkId NetworkId) (bool, error) {
	res, err := c.grpc.JoinNetwork(ctx, &grpc_api.NetworkChangeRequest{
		NetworkId: &wrapperspb.Int32Value{
			Value: int32(networkId),
		},
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) LeaveNetwork(ctx context.Context, networkId NetworkId) (bool, error) {
	res, err := c.grpc.LeaveNetwork(ctx, &grpc_api.NetworkChangeRequest{
		NetworkId: &wrapperspb.Int32Value{
			Value: int32(networkId),
		},
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) NodeInfo(ctx context.Context) (*NodeInfo, error) {
	res, err := c.grpc.NodeInfo(ctx, &grpc_api.Empty{})
	if err != nil {
		return nil, err
	}
	i := &NodeInfo{
		NodeId:                      NodeId(res.NodeId.GetValue()),
		CurrentLocaltime:            res.CurrentLocaltime,
		PeerType:                    PeerType(res.PeerType),
		ConsensusBakerRunning:       res.ConsensusBakerRunning,
		ConsensusRunning:            res.ConsensusRunning,
		ConsensusType:               ConsensusType(res.ConsensusType),
		ConsensusBakerCommittee:     NodeInfoIsInBakingCommittee(res.ConsensusBakerCommittee),
		ConsensusFinalizerCommittee: res.ConsensusFinalizerCommittee,
		ConsensusBakerId:            BakerId(res.ConsensusBakerId.GetValue()),
	}
	return i, nil
}

func (c *client) GetConsensusStatus(ctx context.Context) (*ConsensusStatus, error) {
	res, err := c.grpc.GetConsensusStatus(ctx, &grpc_api.Empty{})
	if err != nil {
		return nil, err
	}
	s := &ConsensusStatus{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, err
}

func (c *client) GetBlockInfo(ctx context.Context, b BlockHash) (*BlockInfo, error) {
	res, err := c.grpc.GetBlockInfo(ctx, &grpc_api.BlockHash{
		BlockHash: b.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, fmt.Errorf("not found")
	}
	i := &BlockInfo{}
	err = json.Unmarshal([]byte(res.GetValue()), i)
	return i, err
}

func (c *client) GetAncestors(ctx context.Context, blockHash BlockHash, amount int) ([]BlockHash, error) {
	res, err := c.grpc.GetAncestors(ctx, &grpc_api.BlockHashAndAmount{
		BlockHash: blockHash.String(),
		Amount:    uint64(amount),
	})
	if err != nil {
		return nil, err
	}
	var s []BlockHash
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetBranches(ctx context.Context) (*Branch, error) {
	res, err := c.grpc.GetBranches(ctx, &grpc_api.Empty{})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, fmt.Errorf("not found")
	}
	b := &Branch{}
	err = json.Unmarshal([]byte(res.GetValue()), b)
	return b, err
}

func (c *client) GetBlocksAtHeight(ctx context.Context, blockHeight BlockHeight) ([]BlockHash, error) {
	res, err := c.grpc.GetBlocksAtHeight(ctx, &grpc_api.BlockHeight{
		BlockHeight: uint64(blockHeight),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, fmt.Errorf("not found")
	}
	var s []BlockHash
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) SendTransactionAsync(ctx context.Context, networkId NetworkId, request TransactionRequest) (TransactionHash, error) {
	b, err := request.Serialize()
	if err != nil {
		return "", fmt.Errorf("unable to serialize request: %w", err)
	}
	p := make([]byte, 2+len(b))
	p[0] = uint8(request.Version())
	p[1] = uint8(request.Kind())
	copy(p[2:], b)
	res, err := c.grpc.SendTransaction(ctx, &grpc_api.SendTransactionRequest{
		NetworkId: uint32(networkId),
		Payload:   p,
	})
	if err != nil {
		return "", fmt.Errorf("unable to send request: %w", err)
	}
	if !res.Value {
		return "", fmt.Errorf("transaction was rejected")
	}
	return newTransactionHash(request.Kind(), b), nil
}

func (c *client) SendTransactionAwait(ctx context.Context, networkId NetworkId, request TransactionRequest) (*TransactionSummary, TransactionHash, error) {
	hash, err := c.SendTransactionAsync(ctx, networkId, request)
	if err != nil {
		return nil, hash, fmt.Errorf("unable to send transaction: %w", err)
	}
	var s *TransactionSummary
	t := time.NewTicker(time.Second)
	for range t.C {
		if request.ExpiredAt().Add(time.Minute).Before(time.Now()) {
			return nil, hash, fmt.Errorf("transaction %q timed out", hash)
		}
		s, err = c.GetTransactionStatus(ctx, hash)
		if err != nil {
			return nil, hash, fmt.Errorf("unable to get status of transaction %q: %w", hash, err)
		}
		if s.Status != TransactionStatusFinalized {
			continue
		}
		t.Stop()
		break
	}
	return s, hash, nil
}

func (c *client) StartBaker(ctx context.Context) (bool, error) {
	res, err := c.grpc.StartBaker(ctx, &grpc_api.Empty{})
	return res.GetValue(), err
}

func (c *client) StopBaker(ctx context.Context) (bool, error) {
	res, err := c.grpc.StopBaker(ctx, &grpc_api.Empty{})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) GetAccountList(ctx context.Context, blockHash BlockHash) ([]AccountAddress, error) {
	res, err := c.grpc.GetAccountList(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	var s []AccountAddress
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetInstances(ctx context.Context, blockHash BlockHash) ([]*ContractAddress, error) {
	res, err := c.grpc.GetInstances(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	var s []*ContractAddress
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetAccountInfo(ctx context.Context, blockHash BlockHash, accountAddress AccountAddress) (*AccountInfo, error) {
	res, err := c.grpc.GetAccountInfo(ctx, &grpc_api.GetAddressInfoRequest{
		BlockHash: blockHash.String(),
		Address:   accountAddress.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	i := &AccountInfo{}
	err = json.Unmarshal([]byte(res.GetValue()), i)
	return i, err
}

func (c *client) GetInstanceInfo(ctx context.Context, blockHash BlockHash, contractAddress *ContractAddress) (*InstanceInfo, error) {
	b, err := json.Marshal(contractAddress)
	if err != nil {
		return nil, err
	}
	res, err := c.grpc.GetInstanceInfo(ctx, &grpc_api.GetAddressInfoRequest{
		BlockHash: blockHash.String(),
		Address:   string(b),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	i := &InstanceInfo{}
	err = json.Unmarshal([]byte(res.GetValue()), i)
	return i, err
}

func (c *client) InvokeContract(ctx context.Context, blockHash BlockHash, context *ContractContext) (*InvokeContractResult, error) {
	b, err := json.Marshal(context)
	if err != nil {
		return nil, err
	}
	res, err := c.grpc.InvokeContract(ctx, &grpc_api.InvokeContractRequest{
		BlockHash: blockHash.String(),
		Context:   string(b),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	r := &InvokeContractResult{}
	err = json.Unmarshal([]byte(res.GetValue()), r)
	return r, nil
}

func (c *client) GetRewardStatus(ctx context.Context, blockHash BlockHash) (*RewardStatus, error) {
	res, err := c.grpc.GetRewardStatus(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	s := &RewardStatus{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, err
}

func (c *client) GetBirkParameters(ctx context.Context, blockHash BlockHash) (*BirkParameters, error) {
	res, err := c.grpc.GetBirkParameters(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	p := &BirkParameters{}
	err = json.Unmarshal([]byte(res.GetValue()), p)
	return p, err
}

func (c *client) GetModuleList(ctx context.Context, b BlockHash) ([]ModuleRef, error) {
	res, err := c.grpc.GetModuleList(ctx, &grpc_api.BlockHash{
		BlockHash: b.String(),
	})
	if err != nil {
		return nil, err
	}
	var s []ModuleRef
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetModuleSource(ctx context.Context, blockHash BlockHash, moduleRef ModuleRef) (io.Reader, error) {
	res, err := c.grpc.GetModuleSource(ctx, &grpc_api.GetModuleSourceRequest{
		BlockHash: blockHash.String(),
		ModuleRef: moduleRef.String(),
	})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(res.GetValue()), nil
}

func (c *client) GetIdentityProviders(ctx context.Context, blockHash BlockHash) ([]*IdentityProvider, error) {
	res, err := c.grpc.GetIdentityProviders(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	var s []*IdentityProvider
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetAnonymityRevokers(ctx context.Context, blockHash BlockHash) ([]*AnonymityRevoker, error) {
	res, err := c.grpc.GetAnonymityRevokers(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	var s []*AnonymityRevoker
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetCryptographicParameters(ctx context.Context, blockHash BlockHash) (*CryptographicParameters, error) {
	res, err := c.grpc.GetCryptographicParameters(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	p := &CryptographicParameters{}
	err = json.Unmarshal([]byte(res.GetValue()), p)
	return p, err
}

func (c *client) GetBannedPeers(ctx context.Context) (*PeerList, error) {
	res, err := c.grpc.GetBannedPeers(ctx, &grpc_api.Empty{})
	if err != nil {
		return nil, err
	}
	l := &PeerList{
		Type:  PeerType(res.PeerType),
		Peers: make([]*PeerElement, len(res.Peers)),
	}
	for i, e := range res.Peers {
		l.Peers[i] = &PeerElement{
			NodeId:        NodeId(e.NodeId.GetValue()),
			Ip:            e.Ip.GetValue(),
			Port:          e.Port.GetValue(),
			CatchupStatus: PeerElementCatchupStatus(e.CatchupStatus),
		}
	}
	return l, nil
}

func (c *client) Shutdown(ctx context.Context) (bool, error) {
	res, err := c.grpc.Shutdown(ctx, &grpc_api.Empty{})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) DumpStart(ctx context.Context, file string, raw bool) (bool, error) {
	res, err := c.grpc.DumpStart(ctx, &grpc_api.DumpRequest{
		File: file,
		Raw:  raw,
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) DumpStop(ctx context.Context) (bool, error) {
	res, err := c.grpc.DumpStop(ctx, &grpc_api.Empty{})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) GetTransactionStatus(ctx context.Context, hash TransactionHash) (*TransactionSummary, error) {
	res, err := c.grpc.GetTransactionStatus(ctx, &grpc_api.TransactionHash{
		TransactionHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, fmt.Errorf("not found")
	}
	s := &TransactionSummary{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, err
}

func (c *client) GetTransactionStatusInBlock(ctx context.Context, hash TransactionHash, blockHash BlockHash) (*TransactionSummary, error) {
	res, err := c.grpc.GetTransactionStatusInBlock(ctx, &grpc_api.GetTransactionStatusInBlockRequest{
		TransactionHash: string(hash),
		BlockHash:       blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, fmt.Errorf("not found")
	}
	s := &TransactionSummary{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, err
}

func (c *client) GetAccountNonFinalizedTransactions(ctx context.Context, accountAddress AccountAddress) ([]TransactionHash, error) {
	res, err := c.grpc.GetAccountNonFinalizedTransactions(ctx, &grpc_api.AccountAddress{
		AccountAddress: accountAddress.String(),
	})
	if err != nil {
		return nil, err
	}
	var s []TransactionHash
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, nil
}

func (c *client) GetBlockSummary(ctx context.Context, blockHash BlockHash) (*BlockSummary, error) {
	res, err := c.grpc.GetBlockSummary(ctx, &grpc_api.BlockHash{
		BlockHash: blockHash.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	s := &BlockSummary{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, nil
}

func (c *client) GetNextAccountNonce(ctx context.Context, accountAddress AccountAddress) (*NextAccountNonce, error) {
	res, err := c.grpc.GetNextAccountNonce(ctx, &grpc_api.AccountAddress{
		AccountAddress: accountAddress.String(),
	})
	if err != nil {
		return nil, err
	}
	if res.GetValue() == "null" {
		return nil, nil
	}
	n := &NextAccountNonce{}
	err = json.Unmarshal([]byte(res.GetValue()), n)
	return n, err
}
