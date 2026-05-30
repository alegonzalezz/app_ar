package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

// TxManager gestiona transacciones SQL.
type TxManager struct {
	db *sql.DB
}

// NewTxManager crea una nueva instancia del gestor de transacciones.
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// RunInTx ejecuta una función dentro de una transacción SQL.
func (m *TxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", err)
	}

	txCtx := ContextWithTx(ctx, tx)

	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error en la lógica (err: %v) y falló rollback (rbErr: %v)", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al hacer commit de la transacción: %w", err)
	}
	return nil
}
