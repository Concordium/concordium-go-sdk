package concordium

const (
	InvokeContractResultSuccess InvokeContractResultTag = "success"
	InvokeContractResultFailure InvokeContractResultTag = "failure"
)

type InvokeContractResultTag string

// ModuleRef base-16 encoded module reference (64 characters)
type ModuleRef string

type InitName string

type ReceiveName string

type ContractContext struct {
	Invoker   *Address         `json:"invoker"`
	Contract  *ContractAddress `json:"contract"`
	Amount    *Amount          `json:"amount"`
	Method    ReceiveName      `json:"method"`
	Parameter Model            `json:"parameter"`
	Energy    int              `json:"energy"`
}

type InvokeContractResult struct {
	Tag         InvokeContractResultTag `json:"tag"`
	UsedEnergy  int                     `json:"usedEnergy"`
	ReturnValue Model                   `json:"returnValue"`
	Events      []Model                 `json:"events"`
	Reason      any                     `json:"reason"` // TODO
}
