package database

import (
	"sort"
	"testing"
)

type migrationOne struct{}
type migrationTwo struct{}
type migrationThree struct{}

func (m migrationOne) Name() string   { return "20170501_test1" }
func (m migrationTwo) Name() string   { return "20170502_test2" }
func (m migrationThree) Name() string { return "20170601_test3" }

func (m migrationOne) Up() string   { return "" }
func (m migrationTwo) Up() string   { return "" }
func (m migrationThree) Up() string { return "" }

func (m migrationOne) Down() string   { return "" }
func (m migrationTwo) Down() string   { return "" }
func (m migrationThree) Down() string { return "" }

func (m migrationOne) String() string   { return m.Name() }
func (m migrationTwo) String() string   { return m.Name() }
func (m migrationThree) String() string { return m.Name() }

func TestByVersion(t *testing.T) {
	migrations := []AppliedMigration{
		{Name: "v1", Version: 1},
		{Name: "v4", Version: 4},
		{Name: "v2", Version: 2},
		{Name: "v3", Version: 3},
	}

	sort.Sort(sort.Reverse(byVersion(migrations)))

	curVersion := len(migrations)

	for _, v := range migrations {
		if v.Version > curVersion {
			t.Error("Version should be less than the previous one")
		}

		curVersion = v.Version
	}
}

func TestByName(t *testing.T) {
	migrations := []Migration{
		migrationThree{},
		migrationTwo{},
		migrationOne{},
	}

	sort.Sort(byName(migrations))

	if migrations[0].Name() != "20170501_test1" {
		t.Error("Invalid migration")
	}

	if migrations[1].Name() != "20170502_test2" {
		t.Error("Invalid migration")
	}

	if migrations[2].Name() != "20170601_test3" {
		t.Error("Invalid migration")
	}
}
