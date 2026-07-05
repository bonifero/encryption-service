package cryptoutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

const certValidity = 10 * 365 * 24 * time.Hour

func LoadOrGenerateKeyPair(path, password, alias string, keyBits int) (*rsa.PrivateKey, error) {
	if _, err := os.Stat(path); err == nil {
		return loadKeyPair(path, password)
	}
	return generateAndStoreKeyPair(path, password, alias, keyBits)
}

func loadKeyPair(path, password string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать хранилище ключей %s: %w", path, err)
	}

	key, _, err := pkcs12.Decode(data, password)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть хранилище ключей (неверный пароль?): %w", err)
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("хранилище ключей не содержит RSA приватный ключ")
	}
	return privateKey, nil
}

func generateAndStoreKeyPair(path, password, alias string, keyBits int) (*rsa.PrivateKey, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("не удалось создать каталог для хранилища ключей: %w", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return nil, fmt.Errorf("не удалось сгенерировать RSA ключ: %w", err)
	}

	cert, err := selfSignedCertificate(privateKey, alias)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать самоподписанный сертификат: %w", err)
	}

	pfxData, err := pkcs12.Encode(rand.Reader, privateKey, cert, nil, password)
	if err != nil {
		return nil, fmt.Errorf("не удалось закодировать PKCS12 хранилище: %w", err)
	}

	if err := os.WriteFile(path, pfxData, 0o600); err != nil {
		return nil, fmt.Errorf("не удалось сохранить хранилище ключей: %w", err)
	}

	return privateKey, nil
}

func selfSignedCertificate(privateKey *rsa.PrivateKey, commonName string) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:         commonName,
			OrganizationalUnit: []string{"Cybersecurity"},
			Organization:       []string{"Bank"},
			Country:            []string{"RU"},
		},
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(certValidity),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(der)
}
