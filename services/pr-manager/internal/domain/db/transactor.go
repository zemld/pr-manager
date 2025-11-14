package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Transactor struct {
	Config
	isReadOnly bool
	conn       *pgx.Conn
	tx         pgx.Tx
	ctx        context.Context
}

func NewTransactor(config Config, ctx context.Context, isReadOnly bool) *Transactor {
	return &Transactor{Config: config, ctx: ctx, isReadOnly: isReadOnly}
}

func (t *Transactor) Begin(ctx context.Context) error {
	if t.isReadOnly {
		conn, err := pgx.Connect(t.ctx, t.Config.GetConnectionString())
		if err != nil {
			return err
		}
		t.conn = conn
		return nil
	}

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
	if t.isReadOnly {
		if t.conn != nil {
			return t.conn.Close(t.ctx)
		}
		return nil
	}
	if t.tx != nil {
		if err := t.tx.Commit(t.ctx); err != nil {
			return err
		}
	}
	if t.conn != nil {
		if err := t.conn.Close(t.ctx); err != nil {
			return err
		}
	}
	return nil
}

func (t *Transactor) Rollback() error {
	if t.isReadOnly {
		if t.conn != nil {
			return t.conn.Close(t.ctx)
		}
		return nil
	}
	if t.tx != nil {
		if err := t.tx.Rollback(t.ctx); err != nil {
			return err
		}
	}
	if t.conn != nil {
		if err := t.conn.Close(t.ctx); err != nil {
			return err
		}
	}
	return nil
}

func (t *Transactor) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if t.isReadOnly {
		return t.conn.Query(ctx, sql, args...)
	}
	return t.tx.Query(ctx, sql, args...)
}

func (t *Transactor) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if t.isReadOnly {
		return t.conn.Exec(ctx, sql, args...)
	}
	return t.tx.Exec(ctx, sql, args...)
}
