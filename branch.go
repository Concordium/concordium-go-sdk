package concordium

type Branch struct {
	BlockHash BlockHash `json:"blockHash"`
	Children  []*Branch `json:"children"`
}
