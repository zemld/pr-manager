package application

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
)

var config = db.Config{
	Host:     os.Getenv("POSTGRES_HOST"),
	User:     os.Getenv("POSTGRES_USER"),
	Password: os.Getenv("POSTGRES_PASSWORD"),
	Db:       os.Getenv("POSTGRES_DB"),
	Port:     os.Getenv("POSTGRES_PORT"),
}

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

func getPool() (*pgxpool.Pool, error) {
	poolOnce.Do(func() {
		connString := config.GetConnectionString()
		pool, poolErr = pgxpool.New(context.Background(), connString)
	})
	return pool, poolErr
}

type TransactionExecutor struct {
	pool *pgxpool.Pool
}

func NewTransactionExecutor(pool *pgxpool.Pool) *TransactionExecutor {
	return &TransactionExecutor{pool: pool}
}

var executor *TransactionExecutor

func (e *TransactionExecutor) withTransaction(ctx context.Context, fn func(*db.Transactor) error, isReadOnly bool) error {
	if e == nil || e.pool == nil {
		p, err := getPool()
		if err != nil {
			return fmt.Errorf("failed to get database pool: %w", err)
		}
		if e == nil {
			executor = NewTransactionExecutor(p)
			e = executor
		} else {
			e.pool = p
		}
	}

	transactor := db.NewTransactor(e.pool, ctx, isReadOnly)
	if err := transactor.Begin(ctx); err != nil {
		return err
	}
	if isReadOnly {
		defer transactor.Commit()
	} else {
		defer transactor.Rollback()
	}

	if err := fn(transactor); err != nil {
		return err
	}
	if !isReadOnly {
		return transactor.Commit()
	}
	return nil
}
