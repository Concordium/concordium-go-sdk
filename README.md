<h1 align="center">Concordium Go SDK</h1>
<div align="center">
	<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/Concordium/concordium-go-sdk">
</div>

## Getting Started

### Installing

```sh
go get -v github.com/Concordium/concordium-go-sdk
```

### Basic usage

The core structure of the SDK is the Client which maintains a connection to the
node and supports querying the node and sending messages to it.

See examples in [`v2/examples/`](v2/examples) for example usage.

### Transaction constructor

There are two helping packets with high-level wrappers for making transactions with minimal user input.
These wrappers handle encoding, setting energy costs when those are fixed for the transaction.

#### Construct

Functions from `github.com/Concordium/concordium-go-sdk/v2/transactions/construct` packet build PreAccountTransaction.

PreAccountTransaction is an AccountTransaction without signatures. To get a signed transaction you can use `PreAccountTransaction.Sign(TransactionSigner)`
with the transmitted implementation of the `TransactionSigner` interface.

```go
// TransactionSigner is an interface for signing transactions.
type TransactionSigner interface {
    // SignTransactionHash signs transaction hash and returns signatures in TransactionSignature type.
    SignTransactionHash(hashToSign *TransactionHash) (*AccountTransactionSignature, error)
}
```

#### Send

Functions from `github.com/Concordium/concordium-go-sdk/v2/transactions/send` packet build AccountTransaction.

AccountTransaction already has signatures, but the first parameter in constructing functions from `send` packet must implement the
`ExactSizeTransactionSigner` interface.

```go
// ExactSizeTransactionSigner describes TransactionSigner with ability to return number of signers.
type ExactSizeTransactionSigner interface {
    TransactionSigner
    // NumberOfKeys returns number of signers.
    NumberOfKeys() uint32
}
```

#### One key signature

For signature with only one private key, you can use the `WalletAccount` struct.
To create `WalletAccount` with one private key just call `NewWalletAccount`
with the account address and key pair. To create `WalletAccount` with one or more
keys, you can add a key pair to the existing `WalletAccount` using the method `AddKeyPair`
of call `NewWalletAccountFromFile` to read credentials from the `<account_address>.export` file.

### Examples

#### Transfer transaction with transfer payload (using construct)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/construct"
)

func main() {
	pathToExportFile := "PATH_TO_YOUR_EXPORT_FILE"
	walletAccount, err := v2.NewWalletAccountFromFile(pathToExportFile)
	if err != nil {
		log.Fatalf("couldn't create wallet account, err: %v", err)
	}

	receiver, err := v2.AccountAddressFromString("RECEIVER_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode receiver address, err: %v", err)
	}

	ctx := context.Background()
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("couldn't create node client, err: %v", err)
	}

	sequenceNumber, err := client.GetNextAccountSequenceNumber(ctx, walletAccount.Address)
	if err != nil {
		log.Fatalf("failed to get next sender sequnce number, err: %v", err)
	}

	nonce := v2.SequenceNumber{Value: sequenceNumber.SequenceNumber.Value}
	expiry := v2.TransactionTime{Value: uint64(time.Now().Add(time.Hour).UTC().Unix())}
	amount := v2.Amount{Value: 100} // 100 micro CCD.

	preAccountTransaction := construct.Transfer(
		walletAccount.NumberOfKeys(),
		*walletAccount.Address,
		nonce,
		expiry,
		receiver,
		amount,
	)

	accountTransaction, err := preAccountTransaction.Sign(walletAccount)
	if err != nil {
		log.Fatalf("couldn't sign transction, err: %v", err)
	}

	txHash, err := accountTransaction.Send(ctx, client)
	if err != nil {
		log.Fatalf("couldn't send block item, err: %v", err)
	}
	fmt.Printf("Success! Transaction hash: %s", txHash.Hex())
}
```

#### Transfer transaction with transfer payload (using send)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Concordium/concordium-go-sdk/v2"
	"github.com/Concordium/concordium-go-sdk/v2/transactions/send"
)

func main() {
	pathToExportFile := "PATH_TO_YOUR_EXPORT_FILE"
	walletAccount, err := v2.NewWalletAccountFromFile(pathToExportFile)
	if err != nil {
		log.Fatalf("couldn't create wallet account, err: %v", err)
	}

	receiver, err := v2.AccountAddressFromString("RECEIVER_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode receiver address, err: %v", err)
	}

	ctx := context.Background()
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("couldn't create node client, err: %v", err)
	}

	sequenceNumber, err := client.GetNextAccountSequenceNumber(ctx, walletAccount.Address)
	if err != nil {
		log.Fatalf("failed to get next sender sequnce number, err: %v", err)
	}

	nonce := v2.SequenceNumber{Value: sequenceNumber.SequenceNumber.Value}
	expiry := v2.TransactionTime{Value: uint64(time.Now().Add(time.Hour).UTC().Unix())}
	amount := v2.Amount{Value: 100} // 100 micro CCD.

	accountTransaction, err := send.Transfer(
		walletAccount,
		*walletAccount.Address,
		nonce,
		expiry,
		receiver,
		amount,
	)
	if err != nil {
		log.Fatalf("couldn't build account transction, err: %v", err)
	}

	txHash, err := accountTransaction.Send(ctx, client)
	if err != nil {
		log.Fatalf("couldn't send block item, err: %v", err)
	}
	fmt.Printf("Success! Transaction hash: %s", txHash.Hex())
}

```

## RPC

All RPC interfaces follow the [concordium protocol docs](https://developer.concordium.software/concordium-grpc-api/#v2%2fconcordium%2fservice.proto).

### More Examples

For more examples, see the `v2/examples/` or `v2/tests/` folders.

## Building

To update the generated protobuf files, ensure that you have the GRPC API repository pulled down:
```bash
git submodule update --init --recursive
```
Then ensure that you have the prerequisite tooling installed for protobuf. First install the Protobuf compiler, then install the Go and Go GRPC protobuf plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
Be sure to add `~/go/bin` to your path, where these executables are installed.

Then, to generate the protobuf files, run the following command:

```
protoc --go_out=./v2 --go-grpc_out=./v2 --proto_path=concordium-grpc-api concordium-grpc-api/v2/concordium/*.proto
```
Now, you should be able to easily build the project:
```
go build ./v2
```
