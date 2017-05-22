package database

import (
	"database/sql"

	"fmt"

	"github.com/jmoiron/sqlx"
)

// SQLXAdapter represents an adapter for the migrator tied to the SQLX library
type SQLXAdapter struct {
	tableName string
	db        *sqlx.DB
	tx        *sql.Tx
}

// NewSQLXMigrator instantiates a new migrator with the given migration table name.
func NewSQLXMigrator(db *sqlx.DB, tableName string) *Migrator {
	return &Migrator{
		Adapter: &SQLXAdapter{
			db:        db,
			tableName: tableName,
		},
	}
}

func (s *SQLXAdapter) Begin() {
	s.tx, _ = s.db.Begin()
}

func (s *SQLXAdapter) Commit() error {
	return s.tx.Commit()
}

func (s *SQLXAdapter) CreateMigrationsTable() error {
	s.db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
	name VARCHAR(255),
	appliedat TIMESTAMP DEFAULT NOW(),
	version INT,
	CONSTRAINT migrations_pkey PRIMARY KEY (name)
);
`, s.tableName))
	return nil
}

func (s *SQLXAdapter) DropMigrationsTable() error {
	s.tx.Exec(fmt.Sprintf("DROP TABLE %s CASCADE;", s.tableName))
	return nil
}

func (s *SQLXAdapter) ExecDown(m Migration) error {
	_, err := s.tx.Exec(m.Down())
	return err
}

func (s *SQLXAdapter) ExecUp(m Migration) error {
	_, err := s.tx.Exec(m.Up())
	return err
}

func (s *SQLXAdapter) MigrationInserted(name string, version int) {
	s.tx.Exec(s.db.Rebind(fmt.Sprintf("INSERT INTO %s (name, version) VALUES (?, ?)", s.tableName)), name, version)
}

func (s *SQLXAdapter) MigrationRemoved(name string) {
	s.tx.Exec(s.db.Rebind(fmt.Sprintf("DELETE FROM %s WHERE name = ?", s.tableName)), name)
}

func (s *SQLXAdapter) SelectMigrations(migrations *[]AppliedMigration) error {
	return s.db.Select(migrations, fmt.Sprintf("SELECT * FROM %s ORDER BY version DESC", s.tableName))
}
