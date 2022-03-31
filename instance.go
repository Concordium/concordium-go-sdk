package concordium

type InstanceInfo struct {
	Amount       *Amount        `json:"amount"`
	SourceModule ModuleRef      `json:"sourceModule"`
	Owner        AccountAddress `json:"owner"`
	Name         InitName       `json:"name"`
	Model        State          `json:"model"`
	Methods      []ReceiveName  `json:"methods"`
}

// ModuleRef base-16 encoded module reference (64 characters)
type ModuleRef string

type InitName string

type ReceiveName string
