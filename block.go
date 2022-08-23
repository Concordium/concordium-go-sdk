package concordium

import (
	"encoding/hex"
	"fmt"
	"time"
)

const blockHashSize = 32

// BlockHash base-16 encoded hash of a block (64 characters)
type BlockHash [blockHashSize]byte

func NewBlockHashFromString(s string) (BlockHash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return BlockHash{}, fmt.Errorf("hex decode: %w", err)
	}
	if len(b) != blockHashSize {
		return BlockHash{}, fmt.Errorf("expect %d bytes but %d given", blockHashSize, len(b))
	}
	var h BlockHash
	copy(h[:], b)
	return h, nil
}

func MustNewBlockHashFromString(s string) BlockHash {
	h, err := NewBlockHashFromString(s)
	if err != nil {
		panic("MustNewBlockHashFromString: " + err.Error())
	}
	return h
}

func (h *BlockHash) String() string {
	return hex.EncodeToString((*h)[:])
}

func (h BlockHash) MarshalJSON() ([]byte, error) {
	b, err := hexMarshalJSON(h[:])
	if err != nil {
		return nil, fmt.Errorf("%T: %w", h, err)
	}
	return b, nil
}

func (h *BlockHash) UnmarshalJSON(b []byte) error {
	v, err := hexUnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("%T: %w", *h, err)
	}
	var x BlockHash
	copy(x[:], v)
	*h = x
	return nil
}

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
