package v2_test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/BoostyLabs/concordium-go-sdk/v2"
	"github.com/BoostyLabs/concordium-go-sdk/v2/pb"
	"github.com/stretchr/testify/require"
)

func TestConvertBlockItems(t *testing.T) {
	randomBytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		randomBytes[i] = byte(rand.Intn(10))
	}

	signatures := make(map[uint32]*pb.AccountSignatureMap)
	signs := make(map[uint32]*pb.Signature)
	signs[1] = &pb.Signature{
		Value: randomBytes,
	}
	signatures[1] = &pb.AccountSignatureMap{Signatures: signs}

	testCases := [...]struct {
		blockItems []*pb.BlockItem
	}{
		{[]*pb.BlockItem{{
			Hash: &pb.TransactionHash{
				Value: randomBytes,
			},
			BlockItem: &pb.BlockItem_AccountTransaction{
				AccountTransaction: &pb.AccountTransaction{
					Signature: &pb.AccountTransactionSignature{
						Signatures: signatures,
					},
					Header: &pb.AccountTransactionHeader{
						Sender: &pb.AccountAddress{
							Value: randomBytes,
						},
						SequenceNumber: &pb.SequenceNumber{
							Value: rand.Uint64(),
						},
						EnergyAmount: &pb.Energy{
							Value: rand.Uint64(),
						},
						Expiry: &pb.TransactionTime{
							Value: rand.Uint64(),
						},
					},
					Payload: &pb.AccountTransactionPayload{
						Payload: &pb.AccountTransactionPayload_Transfer{
							Transfer: &pb.TransferPayload{
								Amount: &pb.Amount{
									Value: rand.Uint64(),
								},
								Receiver: &pb.AccountAddress{
									Value: randomBytes,
								},
							},
						}},
				},
			}}},
		},
		{[]*pb.BlockItem{{
			Hash: &pb.TransactionHash{
				Value: randomBytes,
			},
			BlockItem: &pb.BlockItem_CredentialDeployment{
				CredentialDeployment: &pb.CredentialDeployment{
					MessageExpiry: &pb.TransactionTime{
						Value: rand.Uint64(),
					},
					Payload: &pb.CredentialDeployment_RawPayload{
						RawPayload: randomBytes,
					},
				},
			},
		}},
		},
		{[]*pb.BlockItem{{
			Hash: &pb.TransactionHash{
				Value: randomBytes,
			},
			BlockItem: &pb.BlockItem_UpdateInstruction{
				UpdateInstruction: &pb.UpdateInstruction{
					Signatures: &pb.SignatureMap{
						Signatures: signs,
					},
					Header: &pb.UpdateInstructionHeader{
						SequenceNumber: &pb.UpdateSequenceNumber{
							Value: rand.Uint64(),
						},
						EffectiveTime: &pb.TransactionTime{
							Value: rand.Uint64(),
						},
						Timeout: &pb.TransactionTime{
							Value: rand.Uint64(),
						},
					},
					Payload: &pb.UpdateInstructionPayload{
						Payload: &pb.UpdateInstructionPayload_RawPayload{
							RawPayload: randomBytes,
						}},
				},
			},
		}},
		},
	}

	for _, tc := range testCases {
		convert := v2.ConvertBlockItems(tc.blockItems)
		require.True(t, bytes.Equal(tc.blockItems[0].Hash.Value, convert[0].Hash.Value[:]))

		switch item := tc.blockItems[0].BlockItem.(type) {
		case *pb.BlockItem_AccountTransaction:
			tx, ok := convert[0].BlockItem.(*v2.AccountTransaction)
			if ok {
				require.Equal(t, item.AccountTransaction.Signature.Signatures[1].Signatures[1].Value, tx.Signature.Signatures[1].Signatures[1].Value)
				require.Equal(t, item.AccountTransaction.Header.SequenceNumber.Value, tx.Header.SequenceNumber.Value)
				require.Equal(t, item.AccountTransaction.Header.Expiry.Value, tx.Header.Expiry.Value)
				require.Equal(t, item.AccountTransaction.Header.EnergyAmount.Value, tx.Header.EnergyAmount.Value)
				require.True(t, bytes.Equal(item.AccountTransaction.Header.Sender.Value, tx.Header.Sender.Value[:]))

				payload1, ok1 := item.AccountTransaction.Payload.Payload.(*pb.AccountTransactionPayload_Transfer)
				payload2, ok2 := tx.Payload.Payload.(*v2.Transfer)
				if ok1 {
					if ok2 {
						require.Equal(t, payload1.Transfer.Amount.Value, payload2.Payload.Amount.Value)
						require.True(t, bytes.Equal(payload1.Transfer.Receiver.Value, payload2.Payload.Receiver.Value[:]))
					}
				}
			}

		case *pb.BlockItem_CredentialDeployment:
			tx, ok := convert[0].BlockItem.(*v2.CredentialDeployment)
			if ok {
				require.Equal(t, item.CredentialDeployment.MessageExpiry.Value, tx.MessageExpiry.Value)

				payload1, ok1 := item.CredentialDeployment.Payload.(*pb.CredentialDeployment_RawPayload)
				payload2, ok2 := tx.Payload.(*v2.RawPayload)
				if ok1 {
					if ok2 {
						require.Equal(t, payload1.RawPayload, payload2.Value)
					}
				}
			}
		case *pb.BlockItem_UpdateInstruction:
			tx, ok := convert[0].BlockItem.(*v2.UpdateInstruction)
			if ok {
				require.Equal(t, item.UpdateInstruction.Header.Timeout.Value, tx.Header.Timeout.Value)
				require.Equal(t, item.UpdateInstruction.Header.SequenceNumber.Value, tx.Header.SequenceNumber.Value)
				require.Equal(t, item.UpdateInstruction.Header.EffectiveTime.Value, tx.Header.EffectiveTime.Value)
				require.Equal(t, item.UpdateInstruction.Signatures.Signatures[1].Value, tx.Signatures.Signatures[1].Value)

				payload1, ok1 := item.UpdateInstruction.Payload.Payload.(*pb.UpdateInstructionPayload_RawPayload)
				payload2, ok2 := tx.Payload.Payload.(*v2.RawPayload)

				if ok1 {
					if ok2 {
						require.Equal(t, payload1.RawPayload, payload2.Value)
					}
				}
			}
		}
	}
}
