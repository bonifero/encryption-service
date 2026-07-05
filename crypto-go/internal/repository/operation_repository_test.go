package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"crypto-go/internal/model"
)

func newTestRepository(t *testing.T) *OperationRepository {
	t.Helper()

	db, err := sql.Open("pgx", "host=localhost port=5432 dbname=crypto_db user=postgres password=postgres sslmode=disable")
	if err != nil {
		t.Skipf("не удалось открыть подключение к БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Skipf("PostgreSQL недоступен на localhost:5432/crypto_db, пропускаем интеграционный тест: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	repo, err := NewOperationRepository(db)
	if err != nil {
		t.Fatalf("не удалось инициализировать репозиторий: %v", err)
	}

	if _, err := db.Exec("DELETE FROM crypto_operation_go"); err != nil {
		t.Fatalf("не удалось очистить таблицу перед тестом: %v", err)
	}

	return repo
}

func TestSave_ThenFindByID_ReturnsPersistedOperation(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	saved, err := repo.Save(ctx, model.CryptoOperation{
		OperationType: model.OperationHash,
		InputData:     `"hello"`,
		OutputData:    `{"hashHex":"abc"}`,
		CreatedAt:     time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("Save вернул ошибку: %v", err)
	}
	if saved.ID == 0 {
		t.Fatalf("ожидался присвоенный id")
	}

	found, err := repo.FindByID(ctx, saved.ID)
	if err != nil {
		t.Fatalf("FindByID вернул ошибку: %v", err)
	}
	if found.OperationType != model.OperationHash || found.InputData != `"hello"` {
		t.Fatalf("неожиданные данные записи: %+v", found)
	}
}

func TestFindByID_ReturnsErrNotFoundForMissingID(t *testing.T) {
	repo := newTestRepository(t)

	_, err := repo.FindByID(context.Background(), 999999999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("ожидалась ошибка ErrNotFound, получено: %v", err)
	}
}

func TestFindPage_FiltersByOperationType(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	mustSave(t, repo, model.OperationHash, "in1", "out1")
	mustSave(t, repo, model.OperationEncrypt, "in2", "out2")
	mustSave(t, repo, model.OperationHash, "in3", "out3")

	hashType := model.OperationHash
	operations, total, err := repo.FindPage(ctx, &hashType, 0, 10)
	if err != nil {
		t.Fatalf("FindPage вернул ошибку: %v", err)
	}
	if total != 2 {
		t.Fatalf("ожидалось 2 записи типа HASH, получено %d", total)
	}
	for _, op := range operations {
		if op.OperationType != model.OperationHash {
			t.Fatalf("фильтр вернул запись другого типа: %+v", op)
		}
	}
}

func TestFindPage_RespectsPageSize(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		mustSave(t, repo, model.OperationHash, "in", "out")
	}

	operations, total, err := repo.FindPage(ctx, nil, 0, 2)
	if err != nil {
		t.Fatalf("FindPage вернул ошибку: %v", err)
	}
	if total != 5 {
		t.Fatalf("ожидалось всего 5 записей, получено %d", total)
	}
	if len(operations) != 2 {
		t.Fatalf("ожидалось 2 записи на странице, получено %d", len(operations))
	}
}

func mustSave(t *testing.T, repo *OperationRepository, opType model.OperationType, input, output string) {
	t.Helper()
	_, err := repo.Save(context.Background(), model.CryptoOperation{
		OperationType: opType,
		InputData:     input,
		OutputData:    output,
		CreatedAt:     time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("Save вернул ошибку: %v", err)
	}
}
