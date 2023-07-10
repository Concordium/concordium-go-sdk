package construct

import (
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/costs"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/payloads"
	"github.com/BoostyLabs/concordium-go-sdk/v2/transactions/types"
)

// TransferWithMemo constructs a transfer transaction with a memo.
func TransferWithMemo(numSigs uint32, sender types.AccountAddress, nonce types.Nonce, expiry types.TransactionTime,
	receiver types.AccountAddress, amount types.Amount, memo types.Memo) *transactions.PreAccountTransaction {
	payload := payloads.TransferWithMemoPayload{
		ToAddress: receiver,
		Memo:      memo,
		Amount:    amount,
	}
	energy := types.AddEnergy{
		NumSigs: numSigs,
		Energy:  costs.SimpleTransfer,
	}
	return makeTransaction(sender, nonce, expiry, &energy, &payload)
}
