package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	aesKeyBytes    = 32
	gcmNonceBytes  = 12
	pbkdf2SaltLen  = 16
	pbkdf2Iters    = 100_000
	pbkdf2KeyBytes = 32
)

type HashAlgorithm string

const (
	AlgorithmSHA256 HashAlgorithm = "SHA_256"
	AlgorithmPBKDF2 HashAlgorithm = "PBKDF2"
)

type HybridEncryptionResult struct {
	EncryptedKey string
	IV           string
	CipherText   string
}

type HashResult struct {
	Algorithm  string
	HashHex    string
	SaltHex    *string
	Iterations *int
}

type Service struct {
	privateKey *rsa.PrivateKey
}

func NewService(privateKey *rsa.PrivateKey) *Service {
	return &Service{privateKey: privateKey}
}

func (s *Service) Encrypt(message string) (HybridEncryptionResult, error) {
	aesKey := make([]byte, aesKeyBytes)
	if _, err := rand.Read(aesKey); err != nil {
		return HybridEncryptionResult{}, fmt.Errorf("ошибка при генерации AES-ключа: %w", err)
	}

	nonce := make([]byte, gcmNonceBytes)
	if _, err := rand.Read(nonce); err != nil {
		return HybridEncryptionResult{}, fmt.Errorf("ошибка при генерации IV: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return HybridEncryptionResult{}, fmt.Errorf("ошибка при шифровании сообщения: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return HybridEncryptionResult{}, fmt.Errorf("ошибка при шифровании сообщения: %w", err)
	}
	cipherText := gcm.Seal(nil, nonce, []byte(message), nil)

	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &s.privateKey.PublicKey, aesKey, nil)
	if err != nil {
		return HybridEncryptionResult{}, fmt.Errorf("ошибка при шифровании сообщения: %w", err)
	}

	return HybridEncryptionResult{
		EncryptedKey: base64.StdEncoding.EncodeToString(encryptedKey),
		IV:           base64.StdEncoding.EncodeToString(nonce),
		CipherText:   base64.StdEncoding.EncodeToString(cipherText),
	}, nil
}

func (s *Service) Decrypt(encrypted HybridEncryptionResult) (string, error) {
	encryptedKey, err := base64.StdEncoding.DecodeString(encrypted.EncryptedKey)
	if err != nil {
		return "", fmt.Errorf("ошибка при дешифровании сообщения. Проверьте корректность данных: %w", err)
	}
	nonce, err := base64.StdEncoding.DecodeString(encrypted.IV)
	if err != nil {
		return "", fmt.Errorf("ошибка при дешифровании сообщения. Проверьте корректность данных: %w", err)
	}
	cipherText, err := base64.StdEncoding.DecodeString(encrypted.CipherText)
	if err != nil {
		return "", fmt.Errorf("ошибка при дешифровании сообщения. Проверьте корректность данных: %w", err)
	}

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, s.privateKey, encryptedKey, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при дешифровании сообщения. Проверьте корректность данных: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("ошибка при дешифровании сообщения. Проверьте корректность данных: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("ошибка при дешифровании сообщения. Проверьте корректность данных: %w", err)
	}

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при дешифровании сообщения. Проверьте корректность данных: %w", err)
	}
	return string(plainText), nil
}

func (s *Service) Hash(message string, algorithm HashAlgorithm) (HashResult, error) {
	switch algorithm {
	case AlgorithmSHA256:
		return sha256Hash(message), nil
	case AlgorithmPBKDF2:
		return pbkdf2Hash(message)
	default:
		return HashResult{}, fmt.Errorf("неизвестный алгоритм хеширования: %s", algorithm)
	}
}

func sha256Hash(message string) HashResult {
	sum := sha256.Sum256([]byte(message))
	return HashResult{
		Algorithm: "SHA-256",
		HashHex:   hex.EncodeToString(sum[:]),
	}
}

func pbkdf2Hash(message string) (HashResult, error) {
	salt := make([]byte, pbkdf2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return HashResult{}, fmt.Errorf("ошибка при расчете PBKDF2 хеша: %w", err)
	}

	derived := pbkdf2.Key([]byte(message), salt, pbkdf2Iters, pbkdf2KeyBytes, sha256.New)

	saltHex := hex.EncodeToString(salt)
	iterations := pbkdf2Iters
	return HashResult{
		Algorithm:  "PBKDF2WithHmacSHA256",
		HashHex:    hex.EncodeToString(derived),
		SaltHex:    &saltHex,
		Iterations: &iterations,
	}, nil
}
