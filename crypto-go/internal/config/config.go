package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort string

	DBHost     string
	DBPort     string
	DBName     string
	DBUsername string
	DBPassword string

	KeystorePath     string
	KeystorePassword string
	KeystoreAlias    string
	KeystoreKeyBits  int
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func Load() (Config, error) {
	cfg := Config{
		ServerPort: envOr("SERVER_PORT", "8081"),

		DBHost:     envOr("DB_HOST", "localhost"),
		DBPort:     envOr("DB_PORT", "5432"),
		DBName:     envOr("DB_NAME", "crypto_db"),
		DBUsername: envOr("DB_USERNAME", "postgres"),
		DBPassword: envOr("DB_PASSWORD", "postgres"),

		KeystorePath:     envOr("CRYPTO_KEYSTORE_PATH", "keystore/crypto-service.p12"),
		KeystorePassword: os.Getenv("CRYPTO_KEYSTORE_PASSWORD"),
		KeystoreAlias:    envOr("CRYPTO_KEYSTORE_ALIAS", "crypto-service"),
		KeystoreKeyBits:  2048,
	}

	if cfg.KeystorePassword == "" {
		return Config{}, fmt.Errorf(
			"пароль хранилища ключей не задан: установите переменную окружения CRYPTO_KEYSTORE_PASSWORD")
	}

	return cfg, nil
}

func (c Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBName, c.DBUsername, c.DBPassword)
}
