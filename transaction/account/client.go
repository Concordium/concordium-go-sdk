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

	UpdateContract(ctx *Context, contractAddress *cc.ContractAddress, receiveName cc.ReceiveName, amount *cc.Amount, params ...any) ([]*UpdateContractResultEvent, error)

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
	b, err := io.ReadAll(wasm)
	if err != nil {
		return "", fmt.Errorf("unable to read wasm: %w", err)
	}
	s := &cc.TransactionSummary[*DeployModuleResultEvent, *DeployModuleRejectReason]{}
	h, err := c.sendRequestAwait(ctx, newDeployModuleBody(b), s)
	if err != nil {
		return "", err
	}
	_, o, ok := s.Outcomes.First()
	if !ok {
		return "", fmt.Errorf("%q has no outcomes", h)
	}
	err = o.Error()
	if err != nil {
		return "", err
	}
	if len(o.Result.Events) != 1 {
		return "", fmt.Errorf("unexpected events count in transaction %q", o.Hash)
	}
	return o.Result.Events[0].Contents, nil
}

func (c *client) InitContract(ctx *Context, moduleRef cc.ModuleRef, name cc.ContractName, params ...any) (*cc.ContractAddress, error) {
	s := &cc.TransactionSummary[*InitContractResultEvent, *InitContractRejectReason]{}
	h, err := c.sendRequestAwait(ctx, newInitContractBody(cc.NewAmountZero(), moduleRef, name, params...), s)
	if err != nil {
		return nil, err
	}
	_, o, ok := s.Outcomes.First()
	if !ok {
		return nil, fmt.Errorf("%q has no outcomes", h)
	}
	err = o.Error()
	if err != nil {
		return nil, err
	}
	if len(o.Result.Events) != 1 {
		return nil, fmt.Errorf("unexpected events count in transaction %q", o.Hash)
	}
	return o.Result.Events[0].Address, nil
}

func (c *client) UpdateContract(ctx *Context, contractAddress *cc.ContractAddress, receiveName cc.ReceiveName, amount *cc.Amount, params ...any) ([]*UpdateContractResultEvent, error) {
	s := &cc.TransactionSummary[*UpdateContractResultEvent, *UpdateContractRejectReason]{}
	h, err := c.sendRequestAwait(ctx, newUpdateContractBody(amount, contractAddress, receiveName, params...), s)
	if err != nil {
		return nil, err
	}
	_, o, ok := s.Outcomes.First()
	if !ok {
		return nil, fmt.Errorf("%q has no outcomes", h)
	}
	err = o.Error()
	if err != nil {
		return nil, err
	}
	return o.Result.Events, nil
}

func (c *client) SimpleTransfer(ctx *Context, to cc.AccountAddress, amount *cc.Amount) error {
	s := &cc.TransactionSummary[*SimpleTransferResultEvent, *SimpleTransferRejectReason]{}
	h, err := c.sendRequestAwait(ctx, newSimpleTransferBody(to, amount), s)
	if err != nil {
		return err
	}
	_, o, ok := s.Outcomes.First()
	if !ok {
		return fmt.Errorf("%q has no outcomes", h)
	}
	return o.Error()
}

func (c *client) sendRequestAwait(ctx *Context, body body, summary cc.ITransactionSummary[cc.ITransactionResultEvent, cc.ITransactionRejectReason]) (cc.TransactionHash, error) {
	if ctx == nil {
		return "", fmt.Errorf("nil context not allowed")
	}
	if ctx.Sender == "" {
		return "", fmt.Errorf("empty sender not allowed")
	}
	if len(ctx.Credentials) == 0 {
		return "", fmt.Errorf("empty credentials not allowed")
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
		return "", fmt.Errorf("unable to get next nonce: %w", err)
	}
	// TODO what to do in case of false
	if !nonce.AllFinal {
		return "", fmt.Errorf("account nonce is not reliable")
	}
	req := newRequest(ctx.Credentials, ctx.Sender, nonce.Nonce, expiry, body)

	hash, err := c.cli.SendTransactionAwait(ct, net, req, summary)
	if err != nil {
		return hash, fmt.Errorf("unable to await for transaction: %w", err)
	}
	return hash, nil
}
