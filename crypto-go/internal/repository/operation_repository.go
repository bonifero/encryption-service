package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"crypto-go/internal/model"
)

var ErrNotFound = errors.New("запись не найдена")

const createTableSQL = `
CREATE TABLE IF NOT EXISTS crypto_operation_go (
    id             BIGSERIAL PRIMARY KEY,
    operation_type TEXT NOT NULL,
    input_data     TEXT NOT NULL,
    output_data    TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL
)`

type OperationRepository struct {
	db *sql.DB
}

func NewOperationRepository(db *sql.DB) (*OperationRepository, error) {
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("не удалось создать таблицу crypto_operation_go: %w", err)
	}
	return &OperationRepository{db: db}, nil
}

func (r *OperationRepository) Save(ctx context.Context, op model.CryptoOperation) (model.CryptoOperation, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO crypto_operation_go (operation_type, input_data, output_data, created_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		op.OperationType, op.InputData, op.OutputData, op.CreatedAt)

	if err := row.Scan(&op.ID); err != nil {
		return model.CryptoOperation{}, fmt.Errorf("ошибка при сохранении результата операции в БД: %w", err)
	}
	return op, nil
}

func (r *OperationRepository) FindByID(ctx context.Context, id int64) (model.CryptoOperation, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, operation_type, input_data, output_data, created_at
		 FROM crypto_operation_go WHERE id = $1`, id)

	var op model.CryptoOperation
	if err := row.Scan(&op.ID, &op.OperationType, &op.InputData, &op.OutputData, &op.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CryptoOperation{}, ErrNotFound
		}
		return model.CryptoOperation{}, fmt.Errorf("ошибка при чтении записи из БД: %w", err)
	}
	return op, nil
}

func (r *OperationRepository) FindPage(ctx context.Context, operationType *model.OperationType, page, size int) ([]model.CryptoOperation, int64, error) {
	offset := page * size

	var (
		rows  *sql.Rows
		total int64
		err   error
	)

	if operationType != nil {
		err = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM crypto_operation_go WHERE operation_type = $1`, *operationType).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при подсчете записей: %w", err)
		}
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, operation_type, input_data, output_data, created_at
			 FROM crypto_operation_go WHERE operation_type = $1
			 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, *operationType, size, offset)
	} else {
		err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM crypto_operation_go`).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("ошибка при подсчете записей: %w", err)
		}
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, operation_type, input_data, output_data, created_at
			 FROM crypto_operation_go ORDER BY created_at DESC LIMIT $1 OFFSET $2`, size, offset)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при чтении записей из БД: %w", err)
	}
	defer rows.Close()

	var result []model.CryptoOperation
	for rows.Next() {
		var op model.CryptoOperation
		if err := rows.Scan(&op.ID, &op.OperationType, &op.InputData, &op.OutputData, &op.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("ошибка при чтении записи из БД: %w", err)
		}
		result = append(result, op)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return result, total, nil
}
