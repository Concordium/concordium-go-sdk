package concordium

// AccountAddress base-58 check with version byte 1 encoded address (with Bitcoin mapping table)
type AccountAddress string

// ContractAddress is a JSON record with two fields {index : Int, subindex : Int}
type ContractAddress struct {
	Index    uint64 `json:"index"`
	SubIndex uint64 `json:"subindex"`
}
