package examples_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

	"concordium-go-sdk/v2/pb"
)

func TestGetNodeInfo(t *testing.T) {
	conn, err := grpc.Dial(
		"node.testnet.concordium.com:20000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	require.NotNil(t, conn)

	defer conn.Close()

	client := pb.NewQueriesClient(conn)

	var info *pb.NodeInfo

	info, err = client.GetNodeInfo(context.Background(), new(pb.Empty))
	require.NoError(t, err)
	require.NotNil(t, info)

	blocks, err := client.GetBlocks(context.Background(), new(pb.Empty))
	require.NoError(t, err)
	require.NotNil(t, blocks)

	finalizedblocks, err := client.GetFinalizedBlocks(context.Background(), new(pb.Empty))
	require.NoError(t, err)
	require.NotNil(t, finalizedblocks)

	bb, err := finalizedblocks.Recv()
	require.NoError(t, err)
	require.NotNil(t, bb)

	stream, err := client.GetAccountList(context.Background(), &pb.BlockHashInput{BlockHashInput: &pb.BlockHashInput_Given{Given: &pb.BlockHash{Value: bb.Hash.Value}}})
	require.NoError(t, err)

	acInfo, err := stream.Recv()
	require.NoError(t, err)
	require.NotNil(t, acInfo)

	accInfo, err := client.GetAccountInfo(context.Background(), &pb.AccountInfoRequest{
		BlockHash:         &pb.BlockHashInput{BlockHashInput: nil},
		AccountIdentifier: nil,
	})
	require.NoError(t, err)
	require.NotNil(t, accInfo)

	asd, err := blocks.Recv()
	require.NoError(t, err)
	require.NotNil(t, asd)

	dd := []byte("initial")

	var tx pb.CredentialDeployment
	tx.Payload = &pb.CredentialDeployment_RawPayload{RawPayload: dd}
	tm := pb.TransactionTime{Value: uint64(12345)}
	tx.MessageExpiry = &tm

	asdd, err := client.SendBlockItem(context.Background(), &pb.SendBlockItemRequest{BlockItem: &pb.SendBlockItemRequest_CredentialDeployment{CredentialDeployment: &tx}})
	require.NoError(t, err)
	require.NotNil(t, asdd)
}
