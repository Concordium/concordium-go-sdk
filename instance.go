package concordium

type InstanceInfo struct {
	Amount       *Amount        `json:"amount"`
	SourceModule ModuleRef      `json:"sourceModule"`
	Owner        AccountAddress `json:"owner"`
	Name         InitName       `json:"name"`
	Model        Model          `json:"model"`
	Methods      []ReceiveName  `json:"methods"`
}
