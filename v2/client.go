package v2

import (
	"github.com/zeebo/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"concordium-go-sdk/v2/pb"
)

// Error represents error received while communication with node.
var Error = errs.Class("grpc client")

// Config contains Concordium configurable values.
type Config struct {
	NodeAddress string `env:"NODE_ADDRESS"`
}

// Client provides grpc connection with node.
type Client struct {
	grpcClient pb.QueriesClient
	ClientConn *grpc.ClientConn
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

	client := pb.NewQueriesClient(conn)

	return &Client{grpcClient: client, ClientConn: conn, config: config}, nil
}

// Close closes client connection.
func (c *Client) Close() error {
	return Error.Wrap(c.ClientConn.Close())
}
