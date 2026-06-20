package services_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"inventory/internal/database"
)

// openTestDB opens an in-memory SQLite database with schema initialized.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { _ = db.Close() })
	return db
}
