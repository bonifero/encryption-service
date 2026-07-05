package model

import "time"

type OperationType string

const (
	OperationEncrypt OperationType = "ENCRYPT"
	OperationDecrypt OperationType = "DECRYPT"
	OperationHash    OperationType = "HASH"
)

func (t OperationType) Valid() bool {
	switch t {
	case OperationEncrypt, OperationDecrypt, OperationHash:
		return true
	default:
		return false
	}
}

type CryptoOperation struct {
	ID            int64
	OperationType OperationType
	InputData     string
	OutputData    string
	CreatedAt     time.Time
}
