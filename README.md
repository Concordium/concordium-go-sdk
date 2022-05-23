# Concordium GoLang Client Library
___

## Main Client usage

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