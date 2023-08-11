package v2

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
)

// Config contains Concordium configurable values.
type Config struct {
	NodeAddress string `env:"NODE_ADDRESS"`
	CA          struct {
		Enabled               bool   `env:"CA_ENABLED" envDefault:"false"`
		InsecureSkipVerify    bool   `env:"CA_INSECURE_SKIP_VERIFY" envDefault:"false"`
		CACertificatePath     string `env:"CA_CA_CERTIFICATE_PATH" envDefault:""`
		ClientCertificatePath string `env:"CA_CLIENT_CERTIFICATE_PATH" envDefault:""`
		ClientCertificateKey  string `env:"CA_CLIENT_CERTIFICATE_KEY_PATH" envDefault:""`
	}
}

// Client provides grpc connection with node.
type Client struct {
	GrpcClient pb.QueriesClient
	ClientConn *grpc.ClientConn
	config     Config
}

// NewClient creates new concordium grpc client.
func NewClient(config Config) (_ *Client, err error) {
	creds, err := loadTLSCredentials(config)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(
		config.NodeAddress,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewQueriesClient(conn)

	return &Client{GrpcClient: client, ClientConn: conn, config: config}, nil
}

// loadTLSCredentials load the certificate of the CA who signed the server’s certificate.
func loadTLSCredentials(config Config) (credentials.TransportCredentials, error) {
	if !config.CA.Enabled {
		return insecure.NewCredentials(), nil
	}

	if config.CA.InsecureSkipVerify {
		return credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: config.CA.InsecureSkipVerify,
		}), nil
	}

	pemServerCA, err := os.ReadFile(config.CA.CACertificatePath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, errors.New("failed to add server CA's certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(config.CA.ClientCertificatePath, config.CA.ClientCertificateKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// Close closes client connection.
func (c *Client) Close() error {
	return c.ClientConn.Close()
}
