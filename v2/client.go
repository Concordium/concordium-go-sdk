package v2

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Config contains Concordium configurable values.
type Config struct {
	NodeAddress    string `env:"NODE_ADDRESS"`
	TlsCredentials credentials.TransportCredentials
}

// Client provides grpc connection with node.
type Client struct {
	GrpcClient pb.QueriesClient
	ClientConn *grpc.ClientConn
	config     Config
}

// NewClient creates new concordium grpc client.
func NewClient(config Config) (_ *Client, err error) {
	if config.TlsCredentials != nil {
		conn, err := grpc.Dial(
			config.NodeAddress,
			grpc.WithTransportCredentials(config.TlsCredentials),
		)
		if err != nil {
			return nil, err
		}
		client := pb.NewQueriesClient(conn)

		return &Client{GrpcClient: client, ClientConn: conn, config: config}, nil
	} else {
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
}

func HostTLSRoots() (_ credentials.TransportCredentials, err error) {
	RootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		RootCAs: RootCAs,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// Close closes client connection.
func (c *Client) Close() error {
	return c.ClientConn.Close()
}
