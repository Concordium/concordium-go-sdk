package concordium

import "time"

type NextAccountNonce struct {
	AllFinal bool   `json:"allFinal"`
	Nonce    uint64 `json:"nonce"`
}

// AccountInfo
// Messy documentations:
// https://github.com/Concordium/concordium-node/blob/main/docs/grpc-for-smart-contracts.md#getaccountinfo
// https://github.com/Concordium/concordium-node/blob/main/docs/grpc.md#getaccountinfo--blockhash---accountaddress---accountinfo
type AccountInfo struct {
	AccountNonce           uint64                  `json:"accountNonce"`
	AccountAmount          *Amount                 `json:"accountAmount"`
	AccountReleaseSchedule *AccountReleaseSchedule `json:"accountReleaseSchedule"`
	AccountCredentials     any                     `json:"accountCredentials"` // TODO
	AccountThreshold       int                     `json:"accountThreshold"`
	AccountEncryptedAmount *AccountEncryptedAmount `json:"accountEncryptedAmount"`
	AccountEncryptionKey   EncryptionKey           `json:"accountEncryptionKey"`
	AccountIndex           uint64                  `json:"accountIndex"`
	AccountAddress         AccountAddress          `json:"accountAddress"`
}

type AccountReleaseSchedule struct {
	Schedule []*Release `json:"schedule"`
	Total    *Amount    `json:"total"`
}

type Release struct {
	Timestamp    time.Time         `json:"timestamp"`
	Amount       *Amount           `json:"amount"`
	Transactions []TransactionHash `json:"transactions"`
}

type AccountEncryptedAmount struct {
	IncomingAmounts []EncryptedAmount `json:"incomingAmounts"`
	SelfAmount      EncryptedAmount   `json:"selfAmount"`
	StartIndex      int               `json:"startIndex"`
}

// EncryptionKey base-16 encoded string (192 characters)
type EncryptionKey string

// EncryptedAmount base-16 encoded string (384 characters)
type EncryptedAmount string
