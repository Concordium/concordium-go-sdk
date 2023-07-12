package v2

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// Config contains Concordium configurable values.
type Config struct {
	NodeAddress string `env:"NODE_ADDRESS"`
}

// Client provides grpc connection with node.
type Client struct {
	GrpcClient pb.QueriesClient
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

	return &Client{GrpcClient: client, ClientConn: conn, config: config}, nil
}

// Close closes client connection.
func (c *Client) Close() error {
	return c.ClientConn.Close()
}
