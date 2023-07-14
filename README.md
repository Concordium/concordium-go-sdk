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

### Transaction constructor

There are two helping packets with high level wrappers for making transactions with minimal user input.
These wrappers handle encoding, setting energy costs when those are fixed for transaction.

#### Construct

Functions from `github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct` packet build PreAccountTransaction.

PreAccountTransaction is an AccountTransaction without signatures. To get signed transaction you can use `PreAccountTransaction.Sign(TransactionSigner)`
with transmitted implementation of `TransactionSigner` interface.

```go
// TransactionSigner is an interface for signing transactions.
type TransactionSigner interface {
    // SignTransactionHash signs transaction hash and returns signatures in TransactionSignature type.
    SignTransactionHash(hashToSign *TransactionHash) (*AccountTransactionSignature, error)
}
```

#### Send

Functions from `github.com/BoostyLabs/concordium-go-sdk/v2/transactions/send` packet build AccountTransaction.

AccountTransaction already have signatures, but the first parameter in constructing functions from `send` packet must implement
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

For signature with only one private key you can use `SimpleSigner` struct.
To create `SimpleSigner` just call `NewSimpleSigner` from
`github.com/BoostyLabs/concordium-go-sdk/v2` packet with your private key.

### Examples

#### Transfer transaction with transfer payload (construct)

```go
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcutil/base58"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/construct"
)

func main() {
	privateKey, err := hex.DecodeString("YOUR_PRIVATE_KEY_IN_HEX_ENCODING")
	if err != nil {
		log.Fatalf("couldn't get private key, err: %v", err)
	}

	signer := v2.NewSimpleSigner(privateKey)

	sender, _, err := base58.CheckDecode("ALICE_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode sender address, err: %v", err)
	}

	receiver, _, err := base58.CheckDecode("ROB_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode receiver address, err: %v", err)
	}

	nonce := v2.SequenceNumber{Value: 1}
	expiry := v2.TransactionTime{Value: uint64(time.Now().Add(time.Hour).UTC().Unix())}
	amount := v2.Amount{Value: 100} // 100 micro CCD.

	ctx := context.Background()
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("couldn't create node client, err: %v", err)
	}

	preAccountTransaction := construct.Transfer(
		signer.NumberOfKeys(),
		v2.AccountAddressFromBytes(sender),
		nonce,
		expiry,
		v2.AccountAddressFromBytes(receiver),
		amount,
	)

	accountTransaction, err := preAccountTransaction.Sign(signer)
	if err != nil {
		log.Fatalf("couldn't sign transction, err: %v", err)
	}

	txHash, err := client.SendBlockItem(ctx, &pb.SendBlockItemRequest{
		BlockItem: &pb.SendBlockItemRequest_AccountTransaction{
			AccountTransaction: &pb.AccountTransaction{
				Signature: &pb.AccountTransactionSignature{Signatures: map[uint32]*pb.AccountSignatureMap{
					0: {Signatures: map[uint32]*pb.Signature{
						0: {Value: accountTransaction.Signature.Signatures[0].Signatures[0].Value},
					}}}},
				Header: &pb.AccountTransactionHeader{
					Sender:         &pb.AccountAddress{Value: accountTransaction.Header.Sender.Value[:]},
					SequenceNumber: &pb.SequenceNumber{Value: accountTransaction.Header.SequenceNumber.Value},
					EnergyAmount:   &pb.Energy{Value: accountTransaction.Header.EnergyAmount.Value},
					Expiry:         &pb.TransactionTime{Value: accountTransaction.Header.Expiry.Value},
				},
				Payload: &pb.AccountTransactionPayload{Payload: &pb.AccountTransactionPayload_Transfer{
					Transfer: &pb.TransferPayload{
						Amount:   &pb.Amount{Value: amount.Value},
						Receiver: &pb.AccountAddress{Value: receiver},
					},
				}},
			},
		},
	})
	if err != nil {
		log.Fatalf("couldn't send block item, err: %v", err)
	}
	fmt.Printf("Success! Transaction hash: %s", txHash.Hex())
}
```

#### Transfer transaction with transfer payload (send)

```go
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcutil/base58"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/send"
)

func main() {
	privateKey, err := hex.DecodeString("YOUR_PRIVATE_KEY_IN_HEX_ENCODING")
	if err != nil {
		log.Fatalf("couldn't get private key, err: %v", err)
	}

	signer := v2.NewSimpleSigner(privateKey)

	sender, _, err := base58.CheckDecode("ALICE_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode sender address, err: %v", err)
	}

	receiver, _, err := base58.CheckDecode("ROB_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode receiver address, err: %v", err)
	}

	nonce := v2.SequenceNumber{Value: 1}
	expiry := v2.TransactionTime{Value: uint64(time.Now().Add(time.Hour).UTC().Unix())}
	amount := v2.Amount{Value: 100} // 100 micro CCD.

	ctx := context.Background()
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("couldn't create node client, err: %v", err)
	}

	accountTransaction, err := send.Transfer(
		signer,
		v2.AccountAddressFromBytes(sender),
		nonce,
		expiry,
		v2.AccountAddressFromBytes(receiver),
		amount,
	)
	if err != nil {
		log.Fatalf("couldn't build account transction, err: %v", err)
	}

	txHash, err := client.SendBlockItem(ctx, &pb.SendBlockItemRequest{
		BlockItem: &pb.SendBlockItemRequest_AccountTransaction{
			AccountTransaction: &pb.AccountTransaction{
				Signature: &pb.AccountTransactionSignature{Signatures: map[uint32]*pb.AccountSignatureMap{
					0: {Signatures: map[uint32]*pb.Signature{
						0: {Value: accountTransaction.Signature.Signatures[0].Signatures[0].Value},
					}}}},
				Header: &pb.AccountTransactionHeader{
					Sender:         &pb.AccountAddress{Value: accountTransaction.Header.Sender.Value[:]},
					SequenceNumber: &pb.SequenceNumber{Value: accountTransaction.Header.SequenceNumber.Value},
					EnergyAmount:   &pb.Energy{Value: accountTransaction.Header.EnergyAmount.Value},
					Expiry:         &pb.TransactionTime{Value: accountTransaction.Header.Expiry.Value},
				},
				Payload: &pb.AccountTransactionPayload{Payload: &pb.AccountTransactionPayload_Transfer{
					Transfer: &pb.TransferPayload{
						Amount:   &pb.Amount{Value: amount.Value},
						Receiver: &pb.AccountAddress{Value: receiver},
					},
				}},
			},
		},
	})
	if err != nil {
		log.Fatalf("couldn't send block item, err: %v", err)
	}
	fmt.Printf("Success! Transaction hash: %s", txHash.Hex())
}
```

#### Transfer transaction with raw payload (send)

```go
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/btcsuite/btcutil/base58"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/send"
)

func main() {
	privateKey, err := hex.DecodeString("YOUR_PRIVATE_KEY_IN_HEX_ENCODING")
	if err != nil {
		log.Fatalf("couldn't get private key, err: %v", err)
	}

	signer := v2.NewSimpleSigner(privateKey)

	sender, _, err := base58.CheckDecode("ALICE_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode sender address, err: %v", err)
	}

	receiver, _, err := base58.CheckDecode("ROB_ADDRESS_IN_BASE58_CHECK")
	if err != nil {
		log.Fatalf("couldn't decode receiver address, err: %v", err)
	}

	nonce := v2.SequenceNumber{Value: 1}
	expiry := v2.TransactionTime{Value: uint64(time.Now().Add(time.Hour).UTC().Unix())}
	amount := v2.Amount{Value: 100} // 100 micro CCD.

	ctx := context.Background()
	client, err := v2.NewClient(v2.Config{NodeAddress: "node.testnet.concordium.com:20000"})
	if err != nil {
		log.Fatalf("couldn't create node client, err: %v", err)
	}

	accountTransaction, err := send.Transfer(
		signer,
		v2.AccountAddressFromBytes(sender),
		nonce,
		expiry,
		v2.AccountAddressFromBytes(receiver),
		amount,
	)
	if err != nil {
		log.Fatalf("couldn't build account transction, err: %v", err)
	}

	txHash, err := client.SendBlockItem(ctx, &pb.SendBlockItemRequest{
		BlockItem: &pb.SendBlockItemRequest_AccountTransaction{
			AccountTransaction: &pb.AccountTransaction{
				Signature: &pb.AccountTransactionSignature{Signatures: map[uint32]*pb.AccountSignatureMap{
					0: {Signatures: map[uint32]*pb.Signature{
						0: {Value: accountTransaction.Signature.Signatures[0].Signatures[0].Value},
					}}}},
				Header: &pb.AccountTransactionHeader{
					Sender:         &pb.AccountAddress{Value: accountTransaction.Header.Sender.Value[:]},
					SequenceNumber: &pb.SequenceNumber{Value: accountTransaction.Header.SequenceNumber.Value},
					EnergyAmount:   &pb.Energy{Value: accountTransaction.Header.EnergyAmount.Value},
					Expiry:         &pb.TransactionTime{Value: accountTransaction.Header.Expiry.Value},
				},
				Payload: &pb.AccountTransactionPayload{Payload: &pb.AccountTransactionPayload_RawPayload{
					RawPayload: accountTransaction.Payload.Payload.Encode().Value,
				}},
			},
		},
	})
	if err != nil {
		log.Fatalf("couldn't send block item, err: %v", err)
	}
	fmt.Printf("Success! Transaction hash: %s", txHash.Hex())
}
```

## RPC

All interfaces of rpc follow the [concordium protocol docs](https://developer.concordium.software/concordium-grpc-api/#v2%2fconcordium%2fservice.proto).

### More Example

for more examples, follow `v2/examples/` folder, or in `v2/tests/`