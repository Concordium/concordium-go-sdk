package account

import (
	"context"
	"fmt"
	"io"
	"time"

	cc "github.com/Concordium/concordium-go-sdk"
)

type Context struct {
	Context     context.Context
	NetworkId   cc.NetworkId
	Credentials cc.Credentials
	Sender      cc.AccountAddress
	Expiry      time.Time
}

type Client interface {
	DeployModule(ctx *Context, wasm io.Reader) (cc.ModuleRef, error)

	InitContract(ctx *Context, moduleRef cc.ModuleRef, name cc.ContractName, params ...any) (*cc.ContractAddress, error)

	UpdateContract(ctx *Context, contractAddress *cc.ContractAddress, receiveName cc.ReceiveName, amount *cc.Amount, params ...any) (any, error)

	SimpleTransfer(ctx *Context, to cc.AccountAddress, amount *cc.Amount) error
}

func NewClient(cli cc.Client) Client {
	return &client{
		cli: cli,
	}
}

type client struct {
	cli cc.Client
}

func (c *client) DeployModule(ctx *Context, wasm io.Reader) (cc.ModuleRef, error) {
	wb, err := io.ReadAll(wasm)
	if err != nil {
		return "", fmt.Errorf("unable to read wasm: %w", err)
	}
	outcome, err := c.sendRequestAwait(ctx, newDeployModuleBody(wb))
	if err != nil {
		return "", err
	}
	err = outcome.Error()
	if err != nil {
		return "", err
	}
	if len(outcome.Result.Events) != 1 {
		return "", fmt.Errorf("unexpected events count in transaction %q", outcome.Hash)
	}
	return cc.ModuleRef(outcome.Result.Events[0].Contents), nil
}

func (c *client) InitContract(ctx *Context, moduleRef cc.ModuleRef, name cc.ContractName, params ...any) (*cc.ContractAddress, error) {
	outcome, err := c.sendRequestAwait(ctx, newInitContractBody(cc.NewAmountZero(), moduleRef, name, params...))
	if err != nil {
		return nil, err
	}
	err = outcome.Error()
	if err != nil {
		return nil, err
	}
	if len(outcome.Result.Events) != 1 {
		return nil, fmt.Errorf("unexpected events count in transaction %q", outcome.Hash)
	}
	return &outcome.Result.Events[0].Address, nil
}

func (c *client) UpdateContract(ctx *Context, contractAddress *cc.ContractAddress, receiveName cc.ReceiveName, amount *cc.Amount, params ...any) (any, error) {
	outcome, err := c.sendRequestAwait(ctx, newUpdateContractBody(amount, contractAddress, receiveName, params...))
	if err != nil {
		return nil, err
	}
	err = outcome.Error()
	if err != nil {
		return nil, err
	}
	return outcome.Result.Events, nil
}

func (c *client) SimpleTransfer(ctx *Context, to cc.AccountAddress, amount *cc.Amount) error {
	outcome, err := c.sendRequestAwait(ctx, newSimpleTransferBody(to, amount))
	if err != nil {
		return err
	}
	return outcome.Error()
}

func (c *client) sendRequestAwait(ctx *Context, body body) (*cc.TransactionOutcome, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context not allowed")
	}
	if ctx.Sender == "" {
		return nil, fmt.Errorf("empty sender not allowed")
	}
	if len(ctx.Credentials) == 0 {
		return nil, fmt.Errorf("empty credentials not allowed")
	}
	ct := ctx.Context
	if ct == nil {
		ct = context.Background()
	}
	net := ctx.NetworkId
	if net == 0 {
		net = cc.DefaultNetworkId
	}
	expiry := ctx.Expiry
	if ctx.Expiry.IsZero() {
		expiry = time.Now().Add(cc.DefaultExpiry)
	}
	nonce, err := c.cli.GetNextAccountNonce(ct, ctx.Sender)
	if err != nil {
		return nil, fmt.Errorf("unable to get next nonce: %w", err)
	}
	// TODO what to do in case of false
	if !nonce.AllFinal {
		return nil, fmt.Errorf("account nonce is not reliable")
	}
	req := newRequest(ctx.Credentials, ctx.Sender, nonce.Nonce, expiry, body)

	s, err := c.cli.SendTransactionAwait(ct, net, req)
	if err != nil {
		return nil, fmt.Errorf("unable to await for transaction: %w", err)
	}
	return s, nil
}
