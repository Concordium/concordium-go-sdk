package concordium

// Branch contains branches of the tree. This is the part of the tree
// above the last finalized block.
type Branch struct {
	// Root of the tree.
	BlockHash BlockHash `json:"blockHash"`
	// And children.
	Children []*Branch `json:"children"`
}
