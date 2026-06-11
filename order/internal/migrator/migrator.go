package migrator

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/Steadypim/rocket-factory/order/migrations"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Up(ctx context.Context) error {
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.UpContext(ctx, m.db, "."); err != nil {
		return fmt.Errorf("goose.UpContext: %w", err)
	}

	return nil
}
