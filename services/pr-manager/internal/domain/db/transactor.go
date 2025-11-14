package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactor struct {
	pool       *pgxpool.Pool
	isReadOnly bool
	conn       *pgxpool.Conn
	tx         pgx.Tx
	ctx        context.Context
}

func NewTransactor(pool *pgxpool.Pool, ctx context.Context, isReadOnly bool) *Transactor {
	return &Transactor{pool: pool, ctx: ctx, isReadOnly: isReadOnly}
}

func (t *Transactor) Begin(ctx context.Context) error {
	if t.isReadOnly {
		return nil
	}

	conn, err := t.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	t.conn = conn
	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return err
	}
	t.tx = tx
	return nil
}

func (t *Transactor) Commit() error {
	if t.isReadOnly {
		if t.conn != nil {
			t.conn.Release()
		}
		return nil
	}
	if t.tx != nil {
		if err := t.tx.Commit(t.ctx); err != nil {
			return err
		}
	}
	if t.conn != nil {
		t.conn.Release()
	}
	return nil
}

func (t *Transactor) Rollback() error {
	if t.isReadOnly {
		if t.conn != nil {
			t.conn.Release()
		}
		return nil
	}
	if t.tx != nil {
		if err := t.tx.Rollback(t.ctx); err != nil {
			return err
		}
	}
	if t.conn != nil {
		t.conn.Release()
	}
	return nil
}

func (t *Transactor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if t.isReadOnly {
		return t.pool.Query(ctx, sql, args...)
	}
	return t.tx.Query(ctx, sql, args...)
}

func (t *Transactor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if t.isReadOnly {
		return t.pool.Exec(ctx, sql, args...)
	}
	return t.tx.Exec(ctx, sql, args...)
}
