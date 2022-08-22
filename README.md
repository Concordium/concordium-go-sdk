# Concordium GoLang Client Library
___

## Installation

```shell
go get -u github.com/Concordium/concordium-go-sdk
```

## Main Client usage

* [PeerConnect](#PeerConnect)
* [PeerDisconnect](#PeerDisconnect)
* [PeerUptime](#PeerUptime)
* [PeerTotalSent](#PeerTotalSent)
* [PeerTotalReceived](#PeerTotalReceived)
* [PeerVersion](#PeerVersion)
* [PeerStats](#PeerStats)
* [PeerList](#PeerList)
* [BanNode](#BanNode)
* [UnbanNode](#UnbanNode)
* [JoinNetwork](#JoinNetwork)
* [LeaveNetwork](#LeaveNetwork)
* [NodeInfo](#NodeInfo)
* [GetConsensusStatus](#GetConsensusStatus)
* [GetBlockInfo](#GetBlockInfo)
* [GetAncestors](#GetAncestors)
* [GetBranches](#GetBranches)
* [GetBlocksAtHeight](#GetBlocksAtHeight)
* [StartBaker](#StartBaker)
* [StopBaker](#StopBaker)
* [GetAccountList](#GetAccountList)
* [GetInstances](#GetInstances)
* [GetAccountInfo](#GetAccountInfo)
* [GetInstanceInfo](#GetInstanceInfo)
* [InvokeContract](#InvokeContract)
* [GetRewardStatus](#GetRewardStatus)
* [GetBirkParameters](#GetBirkParameters)
* [GetModuleList](#GetModuleList)
* [GetModuleSource](#GetModuleSource)
* [GetIdentityProviders](#GetIdentityProviders)
* [GetAnonymityRevokers](#GetAnonymityRevokers)
* [GetCryptographicParameters](#GetCryptographicParameters)
* [GetBannedPeers](#GetBannedPeers)
* [Shutdown](#Shutdown)
* [DumpStart](#DumpStart)
* [DumpStop](#DumpStop)
* [GetTransactionStatus](#GetTransactionStatus)
* [GetTransactionStatusInBlock](#GetTransactionStatusInBlock)
* [GetAccountNonFinalizedTransactions](#GetAccountNonFinalizedTransactions)
* [GetBlockSummary](#GetBlockSummary)
* [GetNextAccountNonce](#GetNextAccountNonce)

### PeerConnect

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerConnect(ctx, "127.0.0.1", 10001)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### PeerDisconnect

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerDisconnect(ctx, "127.0.0.1", 10001)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### PeerUptime

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerUptime(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", r)
}
```

### PeerTotalSent

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerTotalSent(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", r)
}
```

### PeerTotalReceived

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerTotalReceived(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", r)
}
```

### PeerVersion

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerVersion(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", r)
}
```

### PeerStats

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerStats(ctx, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### PeerList

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.PeerList(ctx, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### BanNode

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.BanNode(ctx, &concordium.PeerElement{
		NodeId:        concordium.NodeId(0),
		Ip:            "127.0.0.1",
		Port:          10001,
		CatchupStatus: concordium.PeerElementCatchupStatus(0),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### UnbanNode

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.UnbanNode(ctx, &concordium.PeerElement{
		NodeId:        "85d72ab53b6cd728",
		Ip:            "127.0.0.1",
		Port:          10001,
		CatchupStatus: concordium.PeerElementCatchupStatus(0),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### JoinNetwork

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.JoinNetwork(ctx, concordium.DefaultNetworkId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### LeaveNetwork

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.LeaveNetwork(ctx, concordium.DefaultNetworkId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### NodeInfo

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.NodeInfo(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetConsensusStatus

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetConsensusStatus(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetBlockInfo

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetBlockInfo(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetAncestors

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetAncestors(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae", 10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetBranches

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetBranches(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetBlocksAtHeight

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetBlocksAtHeight(ctx, 3335584)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### StartBaker

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.StartBaker(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### StopBaker

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.StopBaker(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### GetAccountList

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetAccountList(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetInstances

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetInstances(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetAccountInfo

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetAccountInfo(ctx,
		"7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae",
		"4hvvPeHb9HY4Lur7eUZv4KfL3tYBug8DRc4X9cVU8mpJLa1f2X",
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetInstanceInfo

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetInstanceInfo(ctx,
		"7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae",
		&concordium.ContractAddress{Index: 5129, SubIndex: 0},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### InvokeContract

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.InvokeContract(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae", &concordium.ContractContext{
		Invoker:   concordium.WrapAccountAddress("4hvvPeHb9HY4Lur7eUZv4KfL3tYBug8DRc4X9cVU8mpJLa1f2X"),
		Contract:  &concordium.ContractAddress{Index: 5129, SubIndex: 0},
		Amount:    concordium.NewAmountZero(),
		Method:    "a.func",
		Parameter: "",
		Energy:    10000000,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetRewardStatus

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetRewardStatus(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetBirkParameters

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetBirkParameters(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetModuleList

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetModuleList(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetModuleSource

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetModuleSource(ctx,
		"7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae",
		"935d17711a4dea10ba5a851df4f19cfdd7cdbd79c8d6ec9abfe5cacff873f6d0",
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetIdentityProviders

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetIdentityProviders(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetAnonymityRevokers

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetAnonymityRevokers(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### GetCryptographicParameters

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetCryptographicParameters(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", r)
}
```

### GetBannedPeers

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetBannedPeers(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### Shutdown

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### DumpStart

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.DumpStart(ctx, "path/to/file", true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### DumpStop

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.DumpStop(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", r)
}
```

### GetTransactionStatus

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	s := &concordium.TransactionSummary[*concordium.TransactionResultEvent, *concordium.TransactionRejectReason]{}
	err = cli.GetTransactionStatus(ctx, "8af811b649875f09d6f5b7ebfcc1899cf0e58466f33f07f74daf073dc7aea17f", s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", s)
}
```

### GetTransactionStatusInBlock

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	s := &concordium.TransactionSummary[*concordium.TransactionResultEvent, *concordium.TransactionRejectReason]{}
	err := cli.GetTransactionStatusInBlock(ctx,
		"8af811b649875f09d6f5b7ebfcc1899cf0e58466f33f07f74daf073dc7aea17f",
		"7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae",
		s, 
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", s)
}
```

### GetAccountNonFinalizedTransactions

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetAccountNonFinalizedTransactions(ctx, "4hvvPeHb9HY4Lur7eUZv4KfL3tYBug8DRc4X9cVU8mpJLa1f2X")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", r)
}
```

### GetBlockSummary

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetBlockSummary(ctx, "7f25ab75a1045321220d6a54ef76d5cd1b107228046b8cc349c69d90f2bf7fae")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", r)
}
```

### GetNextAccountNonce

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	r, err := cli.GetNextAccountNonce(ctx, "4hvvPeHb9HY4Lur7eUZv4KfL3tYBug8DRc4X9cVU8mpJLa1f2X")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", r)
}
```

## Account Transaction Client usage

```go
package main

import (
	"context"
	"github.com/Concordium/concordium-go-sdk"
	"github.com/Concordium/concordium-go-sdk/transaction/account"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "host", "token")
	if err != nil {
		panic(err)
	}

	acc := account.NewClient(cli)

	err = acc.SimpleTransfer(
		&account.Context{
			Context: ctx,
			Credentials: concordium.Credentials{}, // add your credentials
			Sender: "source-address",
		},
		"destination-address",
		concordium.NewAmountFromCCD(10),
	)
	if err != nil {
		panic(err)
	}
}
```