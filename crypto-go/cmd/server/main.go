package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"crypto-go/internal/config"
	"crypto-go/internal/cryptoutil"
	"crypto-go/internal/httpapi"
	"crypto-go/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	privateKey, err := cryptoutil.LoadOrGenerateKeyPair(
		cfg.KeystorePath, cfg.KeystorePassword, cfg.KeystoreAlias, cfg.KeystoreKeyBits)
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища ключей: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}

	repo, err := repository.NewOperationRepository(db)
	if err != nil {
		log.Fatalf("Ошибка инициализации таблицы операций: %v", err)
	}

	cryptoService := cryptoutil.NewService(privateKey)
	handler := httpapi.NewHandler(cryptoService, repo)
	router := httpapi.NewRouter(handler)

	log.Printf("crypto-go запущен на :%s (Swagger UI: http://localhost:%s/swagger-ui/)",
		cfg.ServerPort, cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
