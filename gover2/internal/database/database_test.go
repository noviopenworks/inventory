package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover2/internal/database"
)

func TestOpen_InMemory(t *testing.T) {
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()
}

func TestInitSchema_CreatesAllTables(t *testing.T) {
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	defer db.Close()

	err = database.InitSchema(db)
	require.NoError(t, err)

	tables := []string{"_meta", "users", "computers", "smartphones", "tablets", "windows_keys", "antivirus", "other_software"}
	for _, table := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "table %s should exist", table)
	}
}

func TestInitSchema_Idempotent(t *testing.T) {
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, database.InitSchema(db))
	// Second call must not fail (IF NOT EXISTS)
	assert.NoError(t, database.InitSchema(db))
}
