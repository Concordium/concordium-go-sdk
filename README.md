<h1 align="center">Concordium Go SDK</h1>


## Getting Started

### Installing

```sh
go get -v github.com/BoostyLabs/concordium-go-sdk
```

### Example

#### Get Node Info

```go
package main

import (
	"context"
	"fmt"
	"log"

	"concordium-go-sdk/v2"
)

func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	resp, err := client.GetNodeInfo(context.TODO())
	if err != nil {
		log.Fatalf("failed to get node info, err: %v", err)
	}

	fmt.Println("node info: ", resp.String())
}

```

## RPC

All interfaces of rpc follow the [concordium protocol docs](https://developer.concordium.software/concordium-grpc-api/#v2%2fconcordium%2fservice.proto).

### More Example

for more examples, follow `examples/` folder, or in `tests/`