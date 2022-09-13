package account

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Concordium/concordium-go-sdk"
)

// Context contains main parameters for all the account transaction.
type Context struct {
	Context context.Context
	// Concordium network identifier
	NetworkId concordium.NetworkId
	// Sender's credentials
	Credentials concordium.Credentials
	// The transaction sender
	Sender concordium.AccountAddress
	// Expiry date for the transaction
	Expiry time.Time
}

// Client is a helper to work with account transactions
type Client interface {
	// DeployModule sends DeployModule account transactions and awaits for its finalization.
	DeployModule(ctx *Context, wasm io.Reader) (*concordium.EventModuleDeployed, error)

	// InitContract sends InitContract account transactions and awaits for its finalization.
	InitContract(ctx *Context, params *InitContractParams) (*concordium.EventContractInitialized, error)

	// UpdateContract sends UpdateContract account transactions and awaits for its finalization.
	UpdateContract(ctx *Context, params *UpdateContractParams) (concordium.Events, error)

	// SimpleTransfer sends SimpleTransfer account transactions and awaits for its finalization.
	SimpleTransfer(ctx *Context, params *SimpleTransferParams) (*concordium.EventTransferred, error)
}

// NewClient is a Client constructor.
func NewClient(cli concordium.Client) Client {
	return &client{
		cli: cli,
	}
}

type client struct {
	cli concordium.Client
}

// DeployModule implements Client.DeployModule method.
func (c *client) DeployModule(ctx *Context, wasm io.Reader) (*concordium.EventModuleDeployed, error) {
	b, err := io.ReadAll(wasm)
	if err != nil {
		return nil, fmt.Errorf("unable to read wasm: %w", err)
	}
	s, err := c.sendRequestAwait(ctx, newDeployModuleBody(b))
	if err != nil {
		return nil, err
	}
	if s.Result.Outcome != concordium.BlockItemResultOutcomeSuccess {
		return nil, s.Result.RejectReason.Error()
	}
	if len(s.Result.Events) != 1 {
		return nil, fmt.Errorf("unexpected events count in transaction %q", s.Hash)
	}
	return s.Result.Events[0].ModuleDeployed, nil
}

// InitContract implements Client.InitContract method.
func (c *client) InitContract(ctx *Context, params *InitContractParams) (*concordium.EventContractInitialized, error) {
	s, err := c.sendRequestAwait(ctx, newInitContractBody(concordium.NewAmountZero(), params.ModuleRef, params.Name, params.Params...))
	if err != nil {
		return nil, err
	}
	if s.Result.Outcome != concordium.BlockItemResultOutcomeSuccess {
		return nil, s.Result.RejectReason.Error()
	}
	if len(s.Result.Events) != 1 {
		return nil, fmt.Errorf("unexpected events count in transaction %q", s.Hash)
	}
	return s.Result.Events[0].ContractInitialized, nil
}

// UpdateContract implements Client.UpdateContract method.
func (c *client) UpdateContract(ctx *Context, params *UpdateContractParams) (concordium.Events, error) {
	s, err := c.sendRequestAwait(ctx, newUpdateContractBody(params.Amount, params.ContractAddress, params.ReceiveName, params.Params...))
	if err != nil {
		return nil, err
	}
	if s.Result.Outcome != concordium.BlockItemResultOutcomeSuccess {
		return nil, s.Result.RejectReason.Error()
	}
	return s.Result.Events, nil
}

// SimpleTransfer implements Client.SimpleTransfer method.
func (c *client) SimpleTransfer(ctx *Context, params *SimpleTransferParams) (*concordium.EventTransferred, error) {
	s, err := c.sendRequestAwait(ctx, newSimpleTransferBody(params.To, params.Amount))
	if err != nil {
		return nil, err
	}
	if s.Result.Outcome != concordium.BlockItemResultOutcomeSuccess {
		return nil, s.Result.RejectReason.Error()
	}
	if len(s.Result.Events) != 1 {
		return nil, fmt.Errorf("unexpected events count in transaction %q", s.Hash)
	}
	return s.Result.Events[0].Transferred, nil
}

func (c *client) sendRequestAwait(ctx *Context, body body) (*concordium.BlockItemSummary, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context not allowed")
	}
	if ctx.Sender.IsZero() {
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
		net = concordium.DefaultNetworkId
	}
	expiry := ctx.Expiry
	if ctx.Expiry.IsZero() {
		expiry = time.Now().Add(concordium.DefaultExpiry)
	}
	nonce, err := c.cli.GetNextAccountNonce(ct, ctx.Sender)
	if err != nil {
		return nil, fmt.Errorf("unable to get next nonce: %w", err)
	}
	req := newRequest(ctx.Credentials, ctx.Sender, nonce.Nonce, expiry, body)

	s, hash, err := c.cli.SendTransactionAwait(ct, net, req)
	if err != nil {
		return nil, fmt.Errorf("unable to await for transaction: %w", err)
	}
	for _, v := range s.Outcomes {
		return v, nil
	}
	return nil, fmt.Errorf("%q has no outcomes", hash)
}
