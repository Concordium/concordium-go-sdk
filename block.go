package concordium

import "time"

// BlockHash base-16 encoded hash of a block (64 characters)
type BlockHash string

type BlockHeight uint64

type BlockInfo struct {
	BlockHash             BlockHash   `json:"blockHash"`
	BlockParent           BlockHash   `json:"blockParent"`
	BlockLastFinalized    BlockHash   `json:"blockLastFinalized"`
	BlockHeight           BlockHeight `json:"blockHeight"`
	EraBlockHeight        BlockHeight `json:"eraBlockHeight"`
	GenesisIndex          int         `json:"genesisIndex"`
	BlockReceiveTime      time.Time   `json:"blockReceiveTime"`
	BlockArriveTime       time.Time   `json:"blockArriveTime"`
	BlockSlot             int         `json:"blockSlot"`
	BlockSlotTime         time.Time   `json:"blockSlotTime"`
	BlockBaker            BakerId     `json:"blockBaker"`
	Finalized             bool        `json:"finalized"`
	TransactionCount      int         `json:"transactionCount"`
	TransactionEnergyCost int         `json:"transactionEnergyCost"`
	// undocumented but returned in fact
	TransactionsSize int       `json:"transactionsSize"`
	BlockStateHash   BlockHash `json:"blockStateHash"`
}
