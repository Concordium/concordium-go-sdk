# Concordium GoLang Account Transaction Client
___

## Installation

```shell
go get -u github.com/Concordium/concordium-go-sdk/tranasction/account
```

## Usage

### Client instantiation

```go
// instantiate main client
cli := account.NewClient(mainClient)
// the client is ready to use
```

### DeployModule 

Send DeployModule account transactions and awaits for its finalization.

```go
// instantiate client
// instantiate credentials
w, err := os.Open("/path/to/wasm")
if err != nil {
    panic(err)
}
e, err := cli.DeployModule(
    &account.Context{
        Context: ctx,
        Credentials: credentials,
        Sender: concordium.MustNewAccountAddress("account_address"),
    },
    w,
)
if err != nil {
    panic(err)
}
fmt.Printf("event: %s\n", e)
```

### InitContract 

Send InitContract account transactions and awaits for its finalization.

```go
// instantiate client
// instantiate credentials
e, err := cli.InitContract(
    &account.Context{
        Context: ctx,
        Credentials: credentials,
        Sender: concordium.MustNewAccountAddress("account_address"),
    },
    &account.InitContractParams{
        ModuleRef: concordium.MustNewModuleRef("module_ref"),
        Name: "contact",
    },
)
if err != nil {
    panic(err)
}
fmt.Printf("event: %#v\n", e)
```

### UpdateContract 

Send UpdateContract account transactions and awaits for its finalization.

```go
// instantiate client
// instantiate credentials
s, err := cli.UpdateContract(
    &account.Context{
        Context: ctx,
        Credentials: credentials,
        Sender: concordium.MustNewAccountAddress("account_address"),
    },
    &account.UpdateContractParams{
        ContractAddress: &concordium.ContractAddress{
            Index: 0,
        },
        ReceiveName: concordium.NewReceiveName("contract", "receiver"),
        Amount:      concordium.NewAmountZero(),
    },
)
if err != nil {
    panic(err)
}
fmt.Printf("events: %#v\n", s)
```

### SimpleTransfer 

Send SimpleTransfer account transactions and awaits for its finalization.

```go
// instantiate client
// instantiate sign key
e, err := cli.SimpleTransfer(
    &account.Context{
        Context: context.Background(),
        Credentials: credentials,
        Sender: concordium.MustNewAccountAddress("account_address"),
    },
    &account.SimpleTransferParams{
        To:     concordium.MustNewAccountAddress("account_address"),
        Amount: concordium.NewAmountFromMicroCCD(1000000),
    },
)
if err != nil {
    panic(err)
}
fmt.Printf("event: %#v\n", e)
```