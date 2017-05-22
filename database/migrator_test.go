package database

import (
	"sort"
	"testing"
)

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
