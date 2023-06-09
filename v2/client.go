package v2

import (
	"github.com/zeebo/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Error represents error received while communication with node.
var Error = errs.Class("grpc client")

// Config contains Concordium configurable values.
type Config struct {
	NodeAddress string `env:"NODE_ADDRESS"`
	IsTestnet   bool   `env:"IS_TESTNET"`
}

// Client provides grpc connection with node.
type Client struct {
	grpcClient QueriesClient
	config     Config
}

// NewClient creates new concordium grpc client.
func NewClient(config Config) (_ *Client, err error) {
	conn, err := grpc.Dial(
		config.NodeAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := NewQueriesClient(conn)

	return &Client{grpcClient: client, config: config}, nil
}
