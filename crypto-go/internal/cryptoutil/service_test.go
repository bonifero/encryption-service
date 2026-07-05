package cryptoutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"testing"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("не удалось сгенерировать тестовый RSA ключ: %v", err)
	}
	return NewService(key)
}

func TestEncryptThenDecrypt_ReturnsOriginalMessage(t *testing.T) {
	service := newTestService(t)
	message := "Тестовое сообщение для шифрования 123!"

	encrypted, err := service.Encrypt(message)
	if err != nil {
		t.Fatalf("Encrypt вернул ошибку: %v", err)
	}
	if encrypted.EncryptedKey == "" || encrypted.IV == "" || encrypted.CipherText == "" {
		t.Fatalf("ожидались непустые поля результата шифрования: %+v", encrypted)
	}
	if encrypted.CipherText == message {
		t.Fatalf("шифротекст не должен совпадать с открытым текстом")
	}

	decrypted, err := service.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt вернул ошибку: %v", err)
	}
	if decrypted != message {
		t.Fatalf("ожидалось %q, получено %q", message, decrypted)
	}
}

func TestEncrypt_ProducesDifferentCipherTextForSameMessage(t *testing.T) {
	service := newTestService(t)
	message := "Repeated message"

	first, err := service.Encrypt(message)
	if err != nil {
		t.Fatalf("Encrypt вернул ошибку: %v", err)
	}
	second, err := service.Encrypt(message)
	if err != nil {
		t.Fatalf("Encrypt вернул ошибку: %v", err)
	}

	if first.IV == second.IV {
		t.Fatalf("IV должен быть случайным для каждого вызова")
	}
	if first.CipherText == second.CipherText {
		t.Fatalf("разный IV/AES-ключ должны давать разный шифротекст")
	}
}

func TestDecrypt_RejectsTamperedCipherText(t *testing.T) {
	service := newTestService(t)

	encrypted, err := service.Encrypt("original message")
	if err != nil {
		t.Fatalf("Encrypt вернул ошибку: %v", err)
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(encrypted.CipherText)
	if err != nil {
		t.Fatalf("не удалось декодировать шифротекст: %v", err)
	}
	cipherBytes[0] ^= 0x01
	encrypted.CipherText = base64.StdEncoding.EncodeToString(cipherBytes)

	if _, err := service.Decrypt(encrypted); err == nil {
		t.Fatalf("ожидалась ошибка при дешифровании подделанного шифротекста")
	}
}

func TestDecrypt_RejectsInvalidBase64(t *testing.T) {
	service := newTestService(t)

	invalid := HybridEncryptionResult{
		EncryptedKey: "not-base64!!",
		IV:           "not-base64!!",
		CipherText:   "not-base64!!",
	}

	if _, err := service.Decrypt(invalid); err == nil {
		t.Fatalf("ожидалась ошибка при дешифровании невалидных base64-данных")
	}
}

func TestSHA256Hash_IsDeterministic(t *testing.T) {
	service := newTestService(t)
	message := "Hello, world!"

	first, err := service.Hash(message, AlgorithmSHA256)
	if err != nil {
		t.Fatalf("Hash вернул ошибку: %v", err)
	}
	second, err := service.Hash(message, AlgorithmSHA256)
	if err != nil {
		t.Fatalf("Hash вернул ошибку: %v", err)
	}

	if first.Algorithm != "SHA-256" {
		t.Fatalf("ожидался алгоритм SHA-256, получено %q", first.Algorithm)
	}
	if first.HashHex != second.HashHex {
		t.Fatalf("SHA-256 должен быть детерминированным: %q != %q", first.HashHex, second.HashHex)
	}
	if len(first.HashHex) != 64 {
		t.Fatalf("ожидалось 64 hex-символа, получено %d", len(first.HashHex))
	}
}

func TestPBKDF2Hash_UsesRandomSaltAndIsVerifiable(t *testing.T) {
	service := newTestService(t)
	message := "Hello, world!"

	first, err := service.Hash(message, AlgorithmPBKDF2)
	if err != nil {
		t.Fatalf("Hash вернул ошибку: %v", err)
	}
	second, err := service.Hash(message, AlgorithmPBKDF2)
	if err != nil {
		t.Fatalf("Hash вернул ошибку: %v", err)
	}

	if first.SaltHex == nil || second.SaltHex == nil {
		t.Fatalf("ожидалась соль для PBKDF2")
	}
	if *first.SaltHex == *second.SaltHex {
		t.Fatalf("соль должна быть случайной для каждого вызова")
	}
	if first.HashHex == second.HashHex {
		t.Fatalf("разная соль должна давать разный хеш")
	}
	if first.Iterations == nil || *first.Iterations != 100_000 {
		t.Fatalf("ожидалось 100000 итераций, получено %v", first.Iterations)
	}
}
