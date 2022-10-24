# Concordium GoLang Client
___

## Installation

```shell
go get -u github.com/Concordium/concordium-go-sdk
```

## Packages

* [Account Transaction Client](transaction/account)
* [Credential Deployment Transaction Client](transaction/credential)
* [Update Transaction Client](transaction/update)

## Usage

### Client instantiation

```go
// token is optional
cli, err := concordium.NewClient(context.Background(), "127.0.0.1:10000", concordium.WithToken("token")
if err != nil {
    panic(err)
}
// the client is ready to use
```

### PeerConnect

Suggest the node to connect to the submitted peer. If successful, this adds the peer to the list of peers.

```go
// instantiate client
ok, err := cli.PeerConnect(context.Background(), "127.0.0.1", 10001)
if err != nil {
    panic(err)
}
fmt.Printf("peer connect status: %t\n", ok)
```

### PeerDisconnect

Disconnect from the peer and remove them from the given addresses list if they are on it.

```go
// instantiate client
ok, err := cli.PeerDisconnect(context.Background(), "127.0.0.1", 10001)
if err != nil {
    panic(err)
}
fmt.Printf("peer disconnect status: %t\n", ok)
```

### PeerUptime

Get the uptime of the node.

```go
// instantiate client
t, err := cli.PeerUptime(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("peer uptime: %s\n", t)
```

### PeerTotalSent

Get the total number of packets sent by the node.

```go
// instantiate client
c, err := cli.PeerTotalSent(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("sent count: %d\n", c)
```

### PeerTotalReceived

Get the total number of packets received by the node.

```go
// instantiate client
c, err := cli.PeerTotalReceived(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("received count: %d\n", c)
```

### PeerVersion

Get the version of the node software.

```go
// instantiate client
v, err := cli.PeerVersion(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("version: %s\n", v)
```

### PeerStats

Get information on the peers that the node is connected to.

```go
// instantiate client
s, err := cli.PeerStats(context.Background(), true)
if err != nil {
    panic(err)
}
fmt.Printf("peer statistics: %#v\n", s)
```

### PeerList

Get a list of the peers that the node is connected to.

```go
// instantiate client
s, err := cli.PeerList(context.Background(), true)
if err != nil {
    panic(err)
}
fmt.Printf("peer list: %#v\n", s)
```

### BanNode

Ban a node from being a peer. Note that you should provide a node_id OR an ip, but not both. Use empty value for the option not chosen.

```go
// instantiate client
ok, err := cli.BanNode(context.Background(), "node_id", "127.0.0.1")
if err != nil {
    panic(err)
}
fmt.Printf("ban node status: %t\n", ok)
```

### UnbanNode

Unban a previously banned node. Note that you should provide a node_id OR an ip, but not both. Use empty value for the option not chosen.

```go
// instantiate client
ok, err := cli.UnbanNode(context.Background(), "node_id", "127.0.0.1")
if err != nil {
    panic(err)
}
fmt.Printf("unban node status: %t\n", ok)
```

### JoinNetwork

Attempt to join the specified network.

```go
// instantiate client
ok, err := cli.JoinNetwork(context.Background(), concordium.DefaultNetworkId)
if err != nil {
    panic(err)
}
fmt.Printf("join network status: %t\n", ok)
```

### LeaveNetwork

Attempt to leave the specified network.

```go
// instantiate client
ok, err := cli.LeaveNetwork(context.Background(), concordium.DefaultNetworkId)
if err != nil {
    panic(err)
}
fmt.Printf("leave network status: %t\n", ok)
```

### NodeInfo

Get information about the running node.

```go
// instantiate client
i, err := cli.NodeInfo(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("node info: %#v\n", i)
```

### GetConsensusStatus

Get the information about the consensus.

```go
// instantiate client
s, err := cli.GetConsensusStatus(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("consensus status: %#v\n", s)
```

### GetBlockInfo

Get information, such as height, timings, and transaction counts for the given block.

```go
// instantiate client
i, err := cli.GetBlockInfo(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("block info: %#v\n", i)
```

### GetAncestors

Get a list of the blocks preceding the given block. The list will contain at most amount blocks.

```go
// instantiate client
s, err := cli.GetAncestors(context.Background(), concordium.MustNewBlockHash("block_hash"), 10)
if err != nil {
    panic(err)
}
fmt.Printf("ancestors: %#v\n", s)
```

### GetBranches

Get the branches of the tree. This is the part of the tree above the last finalized block.

```go
// instantiate client
s, err := cli.GetBranches(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("branches: %#v\n", s)
```

### GetBlocksAtHeight

Get a list of the blocks at the given height.

```go
// instantiate client
b, err := cli.GetBlocksAtHeight(context.Background(), 0)
if err != nil {
    panic(err)
}
fmt.Printf("block: %#v\n", b)
```

### SendTransactionAsync

Send a transaction to the given network. The node will do basic
transaction validation, such as signature checks and account nonce checks, and if these
fail, the call will return an error.

```go
// instantiate client
// instantiate request
h, err := cli.SendTransactionAsync(context.Background(), concordium.DefaultNetworkId, request)
if err != nil {
    panic(err)
}
fmt.Printf("transaction hash: %s\n", h)
```

### SendTransactionAwait

Send a transaction to the given network. and await its finalization.

```go
// instantiate client
// instantiate request
s, h, err := cli.SendTransactionAwait(context.Background(), concordium.DefaultNetworkId, request)
if err != nil {
    panic(err)
}
fmt.Printf("transaction hash: %s\n", h)
fmt.Printf("transaction status: %#v\n", s)
```

### StartBaker

Start the baker.

```go
// instantiate client
ok, err := cli.StartBaker(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("start backer status: %t\n", ok)
```

### StopBaker

Stop the baker.

```go
// instantiate client
ok, err := cli.StopBaker(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("stop backer status: %t\n", ok)
```

### GetAccountList

Get a list of all accounts that exist in the state at the end of the given block.

```go
// instantiate client
s, err := cli.GetAccountList(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("accounts: %#v\n", s)
```

### GetInstances

Get a list of all smart contract instances that exist in the state at the end of the given block.

```go
// instantiate client
s, err := cli.GetInstances(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("instances: %#v\n", s)
```

### GetAccountInfo

Get the state of an account in the given block.

```go
// instantiate client
i, err := cli.GetAccountInfo(context.Background(),
    concordium.MustNewBlockHash("block_hash"),
    concordium.MustNewAccountAddress("account_address"),
)
if err != nil {
    panic(err)
}
fmt.Printf("account: %#v\n", i)
```

### GetInstanceInfo

Get information about the given smart contract instance in the given block.

```go
// instantiate client
i, err := cli.GetInstanceInfo(context.Background(),
    concordium.MustNewBlockHash("block_hash"),
    &concordium.ContractAddress{Index: 0, SubIndex: 0},
)
if err != nil {
    panic(err)
}
fmt.Printf("instance: %#v\n", i)
```

### InvokeContract

Invoke a smart contract instance and view its results as if it had been updated at the end of the 
given block. Please note that this is not a transaction, so it won’t affect the contract on chain. 
It only simulates the invocation.

```go
// instantiate client
r, err := cli.InvokeContract(context.Background(), 
    concordium.MustNewBlockHash("block_hash"), 
    &concordium.ContractContext{
        Invoker:   concordium.WrapAccountAddress(concordium.MustNewAccountAddress("account_address")),
        Contract:  &concordium.ContractAddress{Index: 5129, SubIndex: 0},
        Amount:    concordium.NewAmountZero(),
        Method:    concordium.NewReceiveName("contract", "receiver"),
        Parameter: "",
        Energy:    10000000,
    },
)
if err != nil {
    panic(err)
}
fmt.Printf("invoke result: %#v\n", r)
```

### GetRewardStatus

Get an overview of the balance of special accounts in the given block.

```go
// instantiate client
s, err := cli.GetRewardStatus(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("reward status: %#v\n", s)
```

### GetBirkParameters

Get an overview of the parameters used for baking.

```go
// instantiate client
p, err := cli.GetBirkParameters(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("birk parameters: %#v\n", p)
```

### GetModuleList

Get a list of all smart contract modules that exist in the state at the end of the given block.

```go
// instantiate client
s, err := cli.GetModuleList(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("modules: %#v\n", s)
```

### GetModuleSource

Get the binary source of a smart contract module.

```go
// instantiate client
m, err := cli.GetModuleSource(context.Background(),
    concordium.MustNewBlockHash("block_hash"),
    concordium.MustNewModuleRef("module_ref"),
)
if err != nil {
    panic(err)
}
fmt.Printf("module source: %#v\n", m)
```

### GetIdentityProviders

Get a list of all identity providers that exist in the state at the end of the given block.

```go
// instantiate client
s, err := cli.GetIdentityProviders(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("identity providers: %#v\n", s)
```

### GetAnonymityRevokers

Get a list of all anonymity revokers that exist in the state at the end of the given block.

```go
// instantiate client
s, err := cli.GetAnonymityRevokers(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("anonymity revokers: %#v\n", s)
```

### GetCryptographicParameters

Get the cryptographic parameters used in the given block.

```go
// instantiate client
p, err := cli.GetCryptographicParameters(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("cryptographic parameters: %s\n", p)
```

### GetBannedPeers

Get a list of banned peers.

```go
// instantiate client
s, err := cli.GetBannedPeers(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("banned peers: %#v\n", s)
```

### Shutdown

Shut down the node.

```go
// instantiate client
ok, err := cli.Shutdown(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("shutdown status: %t\n", ok)
```

### DumpStart

Start dumping packages into the specified file. Only available on a node built with the network_dump feature.

```go
// instantiate client
ok, err := cli.DumpStart(context.Background(), "path/to/file", true)
if err != nil {
    panic(err)
}
fmt.Printf("dump status status: %t\n", ok)
```

### DumpStop

Stop dumping packages. Only available on a node built with the network_dump feature.

```go
// instantiate client
ok, err := cli.DumpStop(context.Background())
if err != nil {
    panic(err)
}
fmt.Printf("dump stop status: %t\n", ok)
```

### GetTransactionStatus

Get the status of a given transaction.

```go
// instantiate client
s, err := cli.GetTransactionStatus(context.Background(), "transaction_hash")
if err != nil {
    panic(err)
}
fmt.Printf("transaction status: %#v\n", s)
```

### GetTransactionStatusInBlock

Get the status of a given transaction in a given block.

```go
// instantiate client
s, err := cli.GetTransactionStatusInBlock(context.Background(), "transaction_hash", concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("transaction status: %#v\n", s)
```

### GetAccountNonFinalizedTransactions

Get a list of non-finalized transactions present on an account.

```go
// instantiate client
s, err := cli.GetAccountNonFinalizedTransactions(context.Background(), concordium.MustNewAccountAddress("account_address"))
if err != nil {
    panic(err)
}
fmt.Printf("account non-finalized transactions: %#v\n", s)
```

### GetBlockSummary

Get a summary of the transactions and data in a given block.

```go
// instantiate client
s, err := cli.GetBlockSummary(context.Background(), concordium.MustNewBlockHash("block_hash"))
if err != nil {
    panic(err)
}
fmt.Printf("block summary: %#v\n", s)
```

### GetNextAccountNonce

Returns the next available nonce for this account.

```go
// instantiate client
r, err := cli.GetNextAccountNonce(context.Background(), concordium.MustNewAccountAddress("account_address"))
if err != nil {
    panic(err)
}
fmt.Printf("next account nonce: %#v\n", r)
```
