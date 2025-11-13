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
	if err := t.tx.Commit(t.ctx); err != nil {
		return err
	}
	if err := t.conn.Close(t.ctx); err != nil {
		return err
	}
	return nil
}

func (t *Transactor) Rollback() error {
	if err := t.tx.Rollback(t.ctx); err != nil {
		return err
	}
	if err := t.conn.Close(t.ctx); err != nil {
		return err
	}
	return nil
}
