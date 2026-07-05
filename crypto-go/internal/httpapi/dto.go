package httpapi

import "time"

type MessageRequest struct {
	Message string `json:"message"`
}

type EncryptResponse struct {
	EncryptedKey string `json:"encryptedKey"`
	IV           string `json:"iv"`
	CipherText   string `json:"cipherText"`
}

type DecryptRequest struct {
	EncryptedKey string `json:"encryptedKey"`
	IV           string `json:"iv"`
	CipherText   string `json:"cipherText"`
}

type DecryptResponse struct {
	Message string `json:"message"`
}

type HashRequest struct {
	Message   string `json:"message"`
	Algorithm string `json:"algorithm"`
}

type HashResponse struct {
	Algorithm  string  `json:"algorithm"`
	HashHex    string  `json:"hashHex"`
	SaltHex    *string `json:"saltHex"`
	Iterations *int    `json:"iterations"`
}

type CryptoOperationResponse struct {
	ID            int64     `json:"id"`
	OperationType string    `json:"operationType"`
	InputData     string    `json:"inputData"`
	OutputData    string    `json:"outputData"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PageResponse struct {
	Content          []CryptoOperationResponse `json:"content"`
	TotalElements    int64                     `json:"totalElements"`
	TotalPages       int64                     `json:"totalPages"`
	Size             int                       `json:"size"`
	Number           int                       `json:"number"`
	NumberOfElements int                       `json:"numberOfElements"`
	First            bool                      `json:"first"`
	Last             bool                      `json:"last"`
	Empty            bool                      `json:"empty"`
}

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
}
