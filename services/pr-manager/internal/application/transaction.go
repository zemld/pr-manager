package application

import (
	"context"
	"os"

	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
)

var config = db.Config{
	Host:     os.Getenv("POSTGRES_HOST"),
	User:     os.Getenv("POSTGRES_USER"),
	Password: os.Getenv("POSTGRES_PASSWORD"),
	Db:       os.Getenv("POSTGRES_DB"),
	Port:     os.Getenv("POSTGRES_PORT"),
}

type TransactionExecutor struct {
	config db.Config
}

func NewTransactionExecutor(config db.Config) *TransactionExecutor {
	return &TransactionExecutor{config: config}
}

var executor = NewTransactionExecutor(config)

func (e *TransactionExecutor) withTransaction(ctx context.Context, fn func(*db.Transactor) error) error {
	transactor := db.NewTransactor(e.config, ctx)
	if err := transactor.Begin(ctx); err != nil {
		return err
	}
	defer transactor.Rollback()

	if err := fn(transactor); err != nil {
		return err
	}

	return transactor.Commit()
}
