package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Transactor struct {
	Config
	conn *pgx.Conn
	tx   pgx.Tx
	ctx  context.Context
}

func NewTransactor(config Config, ctx context.Context) *Transactor {
	return &Transactor{Config: config, ctx: ctx}
}

func (t *Transactor) Begin(ctx context.Context) error {
	conn, err := pgx.Connect(t.ctx, t.Config.GetConnectionString())
	if err != nil {
		return err
	}
	t.conn = conn
	tx, err := t.conn.BeginTx(t.ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	t.tx = tx
	return nil
}

func (t *Transactor) Commit() error {
	return t.tx.Commit(t.ctx)
}

func (t *Transactor) Rollback() error {
	return t.tx.Rollback(t.ctx)
}

func (t *Transactor) Close() error {
	err := t.conn.Close(t.ctx)
	if err != nil {
		return err
	}
	return nil
}
