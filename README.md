<h1 align="center">Concordium Go SDK</h1>
<div align="center">
	<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/BoostyLabs/concordium-go-sdk">
</div>



## Getting Started

### Installing

```sh
go get -v github.com/BoostyLabs/concordium-go-sdk
```

### Basic usage

The core structure of the SDK is the Client which maintains a connection to the node and supports querying the node and sending messages to it. 
This client is cheaply clonable.

The Client is constructed using the new method.

```go
// NewClient creates new concordium grpc client.
func NewClient(config Config) (_ *Client, err error) {
	// creating grpc connection, using node address
	conn, err := grpc.Dial(
		config.NodeAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	
	client := pb.NewQueriesClient(conn)

	return &Client{grpcClient: client, ClientConn: conn, config: config}, nil
}
```

### Example

#### Get Node Info

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
)

// in this example we receive and print node info.
func main() {
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})

	// sending empty context, can also use any other context instead.
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

for more examples, follow `v2/examples/` folder, or in `v2/tests/`