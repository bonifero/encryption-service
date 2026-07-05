package httpapi

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"crypto-go/internal/cryptoutil"
	"crypto-go/internal/repository"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	db, err := sql.Open("pgx", "host=localhost port=5432 dbname=crypto_db user=postgres password=postgres sslmode=disable")
	if err != nil {
		t.Skipf("не удалось открыть подключение к БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Skipf("PostgreSQL недоступен на localhost:5432/crypto_db, пропускаем интеграционный тест: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	repo, err := repository.NewOperationRepository(db)
	if err != nil {
		t.Fatalf("не удалось инициализировать репозиторий: %v", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("не удалось сгенерировать тестовый RSA ключ: %v", err)
	}

	handler := NewHandler(cryptoutil.NewService(key), repo)
	return NewRouter(handler)
}

func doJSON(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("не удалось сериализовать тело запроса: %v", err)
		}
		reader = bytes.NewReader(data)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestHash_ReturnsSHA256HexDigest(t *testing.T) {
	router := newTestRouter(t)

	rec := doJSON(t, router, http.MethodPost, "/api/crypto/hash", HashRequest{Message: "Hello, world!", Algorithm: "SHA_256"})
	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200, получено %d: %s", rec.Code, rec.Body.String())
	}

	var response HashResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("не удалось разобрать ответ: %v", err)
	}
	if response.Algorithm != "SHA-256" {
		t.Fatalf("ожидался алгоритм SHA-256, получено %q", response.Algorithm)
	}
	if response.HashHex == "" {
		t.Fatalf("ожидался непустой hashHex")
	}
}

func TestEncryptThenDecrypt_ViaHTTP_RoundTrips(t *testing.T) {
	router := newTestRouter(t)

	encRec := doJSON(t, router, http.MethodPost, "/api/crypto/encrypt", MessageRequest{Message: "secret payload"})
	if encRec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200 при шифровании, получено %d: %s", encRec.Code, encRec.Body.String())
	}
	var encrypted EncryptResponse
	if err := json.Unmarshal(encRec.Body.Bytes(), &encrypted); err != nil {
		t.Fatalf("не удалось разобрать ответ шифрования: %v", err)
	}

	decRec := doJSON(t, router, http.MethodPost, "/api/crypto/decrypt", DecryptRequest{
		EncryptedKey: encrypted.EncryptedKey, IV: encrypted.IV, CipherText: encrypted.CipherText,
	})
	if decRec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200 при дешифровании, получено %d: %s", decRec.Code, decRec.Body.String())
	}
	var decrypted DecryptResponse
	if err := json.Unmarshal(decRec.Body.Bytes(), &decrypted); err != nil {
		t.Fatalf("не удалось разобрать ответ дешифрования: %v", err)
	}
	if decrypted.Message != "secret payload" {
		t.Fatalf("ожидалось %q, получено %q", "secret payload", decrypted.Message)
	}
}

func TestEncrypt_RejectsBlankMessage(t *testing.T) {
	router := newTestRouter(t)

	rec := doJSON(t, router, http.MethodPost, "/api/crypto/encrypt", MessageRequest{Message: ""})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидался статус 400, получено %d: %s", rec.Code, rec.Body.String())
	}
}

func TestLogs_ReturnsPersistedOperations(t *testing.T) {
	router := newTestRouter(t)

	doJSON(t, router, http.MethodPost, "/api/crypto/hash", HashRequest{Message: "for logs test", Algorithm: "SHA_256"})

	rec := doJSON(t, router, http.MethodGet, "/api/crypto/logs?type=HASH&page=0&size=5", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200, получено %d: %s", rec.Code, rec.Body.String())
	}

	var page PageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("не удалось разобрать ответ: %v", err)
	}
	if len(page.Content) == 0 {
		t.Fatalf("ожидалась хотя бы одна запись в журнале")
	}
	for _, entry := range page.Content {
		if entry.OperationType != "HASH" {
			t.Fatalf("фильтр type=HASH вернул запись другого типа: %q", entry.OperationType)
		}
	}
}

func TestLogByID_ReturnsNotFoundForMissingID(t *testing.T) {
	router := newTestRouter(t)

	rec := doJSON(t, router, http.MethodGet, "/api/crypto/logs/999999999", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("ожидался статус 404, получено %d: %s", rec.Code, rec.Body.String())
	}
}

func TestDecrypt_RejectsInvalidData(t *testing.T) {
	router := newTestRouter(t)

	rec := doJSON(t, router, http.MethodPost, "/api/crypto/decrypt", DecryptRequest{
		EncryptedKey: "not-valid", IV: "not-valid", CipherText: "not-valid",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидался статус 400, получено %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHash_RejectsUnknownAlgorithm(t *testing.T) {
	router := newTestRouter(t)

	rec := doJSON(t, router, http.MethodPost, "/api/crypto/hash", HashRequest{Message: "test", Algorithm: "MD5"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидался статус 400, получено %d: %s", rec.Code, rec.Body.String())
	}
}

func TestLogs_RespectsPageSize(t *testing.T) {
	router := newTestRouter(t)

	for i := 0; i < 3; i++ {
		doJSON(t, router, http.MethodPost, "/api/crypto/hash", HashRequest{Message: "page test", Algorithm: "SHA_256"})
	}

	rec := doJSON(t, router, http.MethodGet, "/api/crypto/logs?page=0&size=2", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200, получено %d: %s", rec.Code, rec.Body.String())
	}

	var page PageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("не удалось разобрать ответ: %v", err)
	}
	if len(page.Content) != 2 {
		t.Fatalf("ожидалось 2 записи на странице, получено %d", len(page.Content))
	}
	if page.Size != 2 {
		t.Fatalf("ожидался size=2 в ответе, получено %d", page.Size)
	}
}
