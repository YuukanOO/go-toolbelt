// Package database exposes common tools when dealing with a database. The primary struct
// being exposed is the Migrator. It's a configurable struct to help you manage your
// database easily.
package database

import (
	"fmt"
	"sort"
	"time"
)

// MigratorAdapter represents the interface to implement for the Migrator struct.
type MigratorAdapter interface {
	CreateMigrationsTable() error
	DropMigrationsTable() error
	SelectMigrations(migrations *[]AppliedMigration) error
	MigrationInserted(name string, version int)
	MigrationRemoved(name string)
	Begin()
	Commit() error
	ExecUp(migration Migration) error
	ExecDown(migration Migration) error
}

// AppliedMigration represents a migration already applied to the database. Name is
// the name of the migration and version is the version of the database after it was applied.
// The database version is just the number of migrations applied.
type AppliedMigration struct {
	Name      string
	AppliedAt time.Time
	Version   int
}

type byVersion []AppliedMigration

func (s byVersion) Len() int {
	return len(s)
}

func (s byVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byVersion) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}

// Migration is the primary interface to implement for your migrations.
type Migration interface {
	Name() string
	Up() string
	Down() string
}

type byName []Migration

func (s byName) Len() int {
	return len(s)
}

func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byName) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}

// Migrator offers a configurable interface for your migration needs.
type Migrator struct {
	Adapter    MigratorAdapter
	Migrations []Migration
}

// Register given migrations into this migrator. Order does not matter since they will be
// sorted by name when applying migrations.
func (m *Migrator) Register(migrations ...Migration) {
	m.Migrations = append(m.Migrations, migrations...)
}

// Migrate the database to the latest version applying needed migrations and returns the
// new database version.
func (m *Migrator) Migrate() (int, error) {
	if err := m.Adapter.CreateMigrationsTable(); err != nil {
		return -1, err
	}

	var appliedMigrations []AppliedMigration

	if err := m.Adapter.SelectMigrations(&appliedMigrations); err != nil {
		return -1, err
	}

	version := len(appliedMigrations)

	m.Adapter.Begin()

	sort.Sort(byName(m.Migrations))

	for _, mig := range m.Migrations {
		name := mig.Name()
		applied := false

		for _, am := range appliedMigrations {
			if am.Name == name {
				applied = true
				break
			}
		}

		if !applied {
			version++

			if err := m.Adapter.ExecUp(mig); err != nil {
				return -1, err
			}

			m.Adapter.MigrationInserted(name, version)
		}
	}

	return version, m.Adapter.Commit()
}

// RollBackToVersion rolls back the database to a given version.
func (m *Migrator) RollBackToVersion(version int) error {
	var appliedMigrations []AppliedMigration

	if err := m.Adapter.SelectMigrations(&appliedMigrations); err != nil {
		return err
	}

	// Constructs a map to ease the process of retrieving the migration
	migrationsByName := map[string]Migration{}

	for _, rm := range m.Migrations {
		migrationsByName[rm.Name()] = rm
	}

	sort.Sort(sort.Reverse(byVersion(appliedMigrations)))

	m.Adapter.Begin()

	for _, v := range appliedMigrations {
		if v.Version > version {
			curMigration := migrationsByName[v.Name]

			// Rollback it
			if err := m.Adapter.ExecDown(curMigration); err != nil {
				return err
			}

			m.Adapter.MigrationRemoved(v.Name)
		}
	}

	if version == 0 {
		if err := m.Adapter.DropMigrationsTable(); err != nil {
			return err
		}
	}

	return m.Adapter.Commit()
}

// RollBackToName rolls back the database to the given migration name.
func (m *Migrator) RollBackToName(name string) error {
	var migrations []AppliedMigration

	if err := m.Adapter.SelectMigrations(&migrations); err != nil {
		return err
	}

	for _, v := range migrations {
		if v.Name == name {
			return m.RollBackToVersion(v.Version)
		}
	}

	return fmt.Errorf("Could not find migration %s", name)
}
