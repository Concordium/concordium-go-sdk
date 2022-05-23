package concordium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	grpc_api "github.com/Concordium/concordium-go-sdk/grpc-api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"time"
)

type Client interface {
	// PeerConnect Suggest to w peer to connect to the submitted peer details.
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

	// SendTransactionAsync Submit a local transaction
	SendTransactionAsync(ctx context.Context, id NetworkId, request TransactionRequest) (TransactionHash, error)

	SendTransactionAwait(ctx context.Context, id NetworkId, request TransactionRequest) (*TransactionOutcome, error)

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
	InvokeContract(ctx context.Context, hash BlockHash, context *ContractContext) (*InvokeContractResult, error)

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
	GetTransactionStatus(ctx context.Context, hash TransactionHash) (*TransactionSummary, error)

	// GetTransactionStatusInBlock Query for transactions in a block by its hash
	GetTransactionStatusInBlock(ctx context.Context, hash TransactionHash, blockHash BlockHash) (*TransactionSummary, error)

	// GetAccountNonFinalizedTransactions Query for non-finalized
	// transactions present on an account by the account address
	GetAccountNonFinalizedTransactions(ctx context.Context, address AccountAddress) ([]byte, error) // TODO

	// GetBlockSummary Request a summary for a block by its hash
	GetBlockSummary(ctx context.Context, hash BlockHash) ([]byte, error) // TODO

	// GetNextAccountNonce Request next nonce information for an account
	GetNextAccountNonce(ctx context.Context, address AccountAddress) (*NextAccountNonce, error)
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

func (c *client) PeerConnect(ctx context.Context, ip string, port int) (bool, error) {
	res, err := c.grpc.PeerConnect(ctx, &grpc_api.PeerConnectRequest{
		Ip: &wrapperspb.StringValue{
			Value: ip,
		},
		Port: &wrapperspb.Int32Value{
			Value: int32(port),
		},
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) PeerDisconnect(ctx context.Context, ip string, port int) (bool, error) {
	res, err := c.grpc.PeerDisconnect(ctx, &grpc_api.PeerConnectRequest{
		Ip: &wrapperspb.StringValue{
			Value: ip,
		},
		Port: &wrapperspb.Int32Value{
			Value: int32(port),
		},
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) PeerUptime(ctx context.Context) (int64, error) {
	res, err := c.grpc.PeerUptime(ctx, &grpc_api.Empty{})
	if err != nil {
		return 0, err
	}
	return int64(res.GetValue()), nil
}

func (c *client) PeerTotalSent(ctx context.Context) (int, error) {
	res, err := c.grpc.PeerTotalSent(ctx, &grpc_api.Empty{})
	if err != nil {
		return 0, err
	}
	return int(res.GetValue()), nil
}

func (c *client) PeerTotalReceived(ctx context.Context) (int, error) {
	res, err := c.grpc.PeerTotalReceived(ctx, &grpc_api.Empty{})
	if err != nil {
		return 0, err
	}
	return int(res.GetValue()), nil
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
		AvgBpsIn:  int(res.AvgBpsIn),
		AvgBpsOut: int(res.AvgBpsOut),
	}
	for i, e := range res.Peerstats {
		s.Peers[i] = &PeerStatsElement{
			NodeId:          NodeId(e.NodeId),
			PacketsSent:     int(e.PacketsSent),
			PacketsReceived: int(e.PacketsReceived),
			Latency:         int(e.Latency),
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
			Port:          int(e.Port.GetValue()),
			CatchupStatus: PeerElementCatchupStatus(e.CatchupStatus),
		}
	}
	return l, nil
}

func (c *client) BanNode(ctx context.Context, element PeerElement) (bool, error) {
	res, err := c.grpc.BanNode(ctx, &grpc_api.PeerElement{
		NodeId: &wrapperspb.StringValue{
			Value: string(element.NodeId),
		},
		Port: &wrapperspb.UInt32Value{
			Value: uint32(element.Port),
		},
		Ip: &wrapperspb.StringValue{
			Value: element.Ip,
		},
		CatchupStatus: grpc_api.PeerElement_CatchupStatus(element.CatchupStatus),
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) UnbanNode(ctx context.Context, element PeerElement) (bool, error) {
	res, err := c.grpc.UnbanNode(ctx, &grpc_api.PeerElement{
		NodeId: &wrapperspb.StringValue{
			Value: string(element.NodeId),
		},
		Port: &wrapperspb.UInt32Value{
			Value: uint32(element.Port),
		},
		Ip: &wrapperspb.StringValue{
			Value: element.Ip,
		},
		CatchupStatus: grpc_api.PeerElement_CatchupStatus(element.CatchupStatus),
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) JoinNetwork(ctx context.Context, id NetworkId) (bool, error) {
	res, err := c.grpc.JoinNetwork(ctx, &grpc_api.NetworkChangeRequest{
		NetworkId: &wrapperspb.Int32Value{
			Value: int32(id),
		},
	})
	if err != nil {
		return false, err
	}
	return res.GetValue(), nil
}

func (c *client) LeaveNetwork(ctx context.Context, id NetworkId) (bool, error) {
	res, err := c.grpc.LeaveNetwork(ctx, &grpc_api.NetworkChangeRequest{
		NetworkId: &wrapperspb.Int32Value{
			Value: int32(id),
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
		CurrentLocaltime:            int64(res.CurrentLocaltime),
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

func (c *client) GetBlockInfo(ctx context.Context, hash BlockHash) (*BlockInfo, error) {
	res, err := c.grpc.GetBlockInfo(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	i := &BlockInfo{}
	err = json.Unmarshal([]byte(res.GetValue()), i)
	return i, err
}

func (c *client) GetAncestors(ctx context.Context, hash BlockHash, amount int) ([]BlockHash, error) {
	res, err := c.grpc.GetAncestors(ctx, &grpc_api.BlockHashAndAmount{
		BlockHash: string(hash),
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
	b := &Branch{}
	err = json.Unmarshal([]byte(res.GetValue()), b)
	return b, err
}

func (c *client) GetBlocksAtHeight(ctx context.Context, height BlockHeight) ([]BlockHash, error) {
	res, err := c.grpc.GetBlocksAtHeight(ctx, &grpc_api.BlockHeight{
		BlockHeight: uint64(height),
	})
	if err != nil {
		return nil, err
	}
	var s []BlockHash
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) SendTransactionAsync(ctx context.Context, id NetworkId, request TransactionRequest) (TransactionHash, error) {
	b, err := request.Serialize()
	if err != nil {
		return "", fmt.Errorf("unable to serialize request: %w", err)
	}
	p := make([]byte, 2+len(b))
	p[0] = uint8(request.Version())
	p[1] = uint8(request.Kind())
	copy(p[2:], b)
	res, err := c.grpc.SendTransaction(ctx, &grpc_api.SendTransactionRequest{
		NetworkId: uint32(id),
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

func (c *client) SendTransactionAwait(ctx context.Context, id NetworkId, request TransactionRequest) (*TransactionOutcome, error) {
	hash, err := c.SendTransactionAsync(ctx, id, request)
	if err != nil {
		return nil, fmt.Errorf("unable to send transaction: %w", err)
	}
	outcome := TransactionOutcome{}
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if request.ExpiredAt().Add(time.Minute).Before(time.Now()) {
			return nil, fmt.Errorf("transaction %q timed out", hash)
		}
		status, err := c.GetTransactionStatus(ctx, hash)
		if err != nil {
			return nil, fmt.Errorf("unable to get status of transaction %q: %w", hash, err)
		}
		if status.Status != TransactionStatusFinalized {
			continue
		}
		if len(status.Outcomes) != 1 {
			return nil, fmt.Errorf("unexpected outcomes count in transaction %q", hash)
		}
		for _, outcome = range status.Outcomes {
		}
		ticker.Stop()
	}
	return &outcome, nil
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

func (c *client) GetAccountList(ctx context.Context, hash BlockHash) ([]AccountAddress, error) {
	res, err := c.grpc.GetAccountList(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	var s []AccountAddress
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetInstances(ctx context.Context, hash BlockHash) ([]*ContractAddress, error) {
	res, err := c.grpc.GetInstances(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	var s []*ContractAddress
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetAccountInfo(ctx context.Context, hash BlockHash, address AccountAddress) (*AccountInfo, error) {
	res, err := c.grpc.GetAccountInfo(ctx, &grpc_api.GetAddressInfoRequest{
		BlockHash: string(hash),
		Address:   string(address),
	})
	if err != nil {
		return nil, err
	}
	i := &AccountInfo{}
	err = json.Unmarshal([]byte(res.GetValue()), i)
	return i, err
}

func (c *client) GetInstanceInfo(ctx context.Context, hash BlockHash, address *ContractAddress) (*InstanceInfo, error) {
	b, err := json.Marshal(address)
	if err != nil {
		return nil, err
	}
	res, err := c.grpc.GetInstanceInfo(ctx, &grpc_api.GetAddressInfoRequest{
		BlockHash: string(hash),
		Address:   string(b),
	})
	if err != nil {
		return nil, err
	}
	i := &InstanceInfo{}
	err = json.Unmarshal([]byte(res.GetValue()), i)
	return i, err
}

func (c *client) InvokeContract(ctx context.Context, hash BlockHash, context *ContractContext) (*InvokeContractResult, error) {
	b, err := json.Marshal(context)
	if err != nil {
		return nil, err
	}
	res, err := c.grpc.InvokeContract(ctx, &grpc_api.InvokeContractRequest{
		BlockHash: string(hash),
		Context:   string(b),
	})
	if err != nil {
		return nil, err
	}
	r := &InvokeContractResult{}
	err = json.Unmarshal([]byte(res.GetValue()), r)
	return r, nil
}

func (c *client) GetRewardStatus(ctx context.Context, hash BlockHash) (*RewardStatus, error) {
	res, err := c.grpc.GetRewardStatus(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	s := &RewardStatus{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, err
}

func (c *client) GetBirkParameters(ctx context.Context, hash BlockHash) (*BirkParameters, error) {
	res, err := c.grpc.GetBirkParameters(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	p := &BirkParameters{}
	err = json.Unmarshal([]byte(res.GetValue()), p)
	return p, err
}

func (c *client) GetModuleList(ctx context.Context, hash BlockHash) ([]ModuleRef, error) {
	res, err := c.grpc.GetModuleList(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	var s []ModuleRef
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetModuleSource(ctx context.Context, hash BlockHash, ref ModuleRef) (io.Reader, error) {
	res, err := c.grpc.GetModuleSource(ctx, &grpc_api.GetModuleSourceRequest{
		BlockHash: string(hash),
		ModuleRef: string(ref),
	})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(res.GetValue()), nil
}

func (c *client) GetIdentityProviders(ctx context.Context, hash BlockHash) ([]*IdentityProvider, error) {
	res, err := c.grpc.GetIdentityProviders(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	var s []*IdentityProvider
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetAnonymityRevokers(ctx context.Context, hash BlockHash) ([]*AnonymityRevoker, error) {
	res, err := c.grpc.GetAnonymityRevokers(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	var s []*AnonymityRevoker
	err = json.Unmarshal([]byte(res.GetValue()), &s)
	return s, err
}

func (c *client) GetCryptographicParameters(ctx context.Context, hash BlockHash) ([]byte, error) {
	res, err := c.grpc.GetCryptographicParameters(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	return []byte(res.GetValue()), nil
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
			Port:          int(e.Port.GetValue()),
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
	s := &TransactionSummary{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, err
}

func (c *client) GetTransactionStatusInBlock(ctx context.Context, hash TransactionHash, blockHash BlockHash) (*TransactionSummary, error) {
	res, err := c.grpc.GetTransactionStatusInBlock(ctx, &grpc_api.GetTransactionStatusInBlockRequest{
		TransactionHash: string(hash),
		BlockHash:       string(blockHash),
	})
	if err != nil {
		return nil, err
	}
	s := &TransactionSummary{}
	err = json.Unmarshal([]byte(res.GetValue()), s)
	return s, err
}

func (c *client) GetAccountNonFinalizedTransactions(ctx context.Context, address AccountAddress) ([]byte, error) {
	res, err := c.grpc.GetAccountNonFinalizedTransactions(ctx, &grpc_api.AccountAddress{
		AccountAddress: string(address),
	})
	if err != nil {
		return nil, err
	}
	return []byte(res.GetValue()), nil
}

func (c *client) GetBlockSummary(ctx context.Context, hash BlockHash) ([]byte, error) {
	res, err := c.grpc.GetBlockSummary(ctx, &grpc_api.BlockHash{
		BlockHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	return []byte(res.GetValue()), nil
}

func (c *client) GetNextAccountNonce(ctx context.Context, address AccountAddress) (*NextAccountNonce, error) {
	res, err := c.grpc.GetNextAccountNonce(ctx, &grpc_api.AccountAddress{
		AccountAddress: string(address),
	})
	if err != nil {
		return nil, err
	}
	n := &NextAccountNonce{}
	err = json.Unmarshal([]byte(res.GetValue()), n)
	return n, err
}
