package backup_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"inventory/internal/backup"
	"inventory/internal/database"
)

func TestCreate_NilDB_ReturnsError(t *testing.T) {
	err := backup.Create(nil, filepath.Join(t.TempDir(), "out.db"))
	require.Error(t, err)
}

func TestCreate_ProducesOpenableCopyWithSameRows(t *testing.T) {
	src, err := database.Open(filepath.Join(t.TempDir(), "src.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = src.Close() })
	require.NoError(t, database.InitSchema(src))
	_, err = src.Exec(`INSERT INTO computers (name, model, status) VALUES (?, ?, ?)`, "PC1", "Dell", "Active")
	require.NoError(t, err)
	_, err = src.Exec(`INSERT INTO computers (name, model, status) VALUES (?, ?, ?)`, "PC2", "HP", "Active")
	require.NoError(t, err)

	dest := filepath.Join(t.TempDir(), "backup.db")
	require.NoError(t, backup.Create(src, dest))

	copyDB, err := database.Open(dest)
	require.NoError(t, err)
	t.Cleanup(func() { _ = copyDB.Close() })
	var count int
	require.NoError(t, copyDB.QueryRow(`SELECT COUNT(*) FROM computers`).Scan(&count))
	require.Equal(t, 2, count)
}
