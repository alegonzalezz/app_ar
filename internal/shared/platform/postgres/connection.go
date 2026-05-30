package postgres

import (
	"context"
	"database/sql"
)

// SQLQueryer define la interfaz común para sql.DB y sql.Tx.
type SQLQueryer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type ctxKey string

const txKey ctxKey = "db_tx"

// GetQueryer retorna el *sql.Tx del contexto si existe, o el *sql.DB.
func GetQueryer(ctx context.Context, db *sql.DB) SQLQueryer {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return db
}

// ContextWithTx inyecta un *sql.Tx en el contexto.
func ContextWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
