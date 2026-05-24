package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed *.sql
var fs embed.FS

// Run ejecuta todas las migraciones SQL pendientes en la carpeta migrations
func Run(ctx context.Context, db *sql.DB) error {
	// Crear tabla de schema_migrations si no existe
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("error al crear tabla schema_migrations: %w", err)
	}

	// Leer archivos sql
	entries, err := fs.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error al leer directorio de migraciones: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// Obtener migraciones ya aplicadas
	rows, err := db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("error al consultar schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		applied[version] = true
	}

	// Ejecutar migraciones pendientes en orden
	for _, file := range files {
		if applied[file] {
			continue
		}

		log.Printf("[MIGRATION] Aplicando migración: %s", file)
		content, err := fs.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error al leer archivo de migración %s: %w", file, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("error ejecutando migración %s: %w", file, err)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", file); err != nil {
			tx.Rollback()
			return fmt.Errorf("error registrando migración %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("[MIGRATION] Migración aplicada exitosamente: %s", file)
	}

	return nil
}
