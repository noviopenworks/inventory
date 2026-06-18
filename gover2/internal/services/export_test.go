package services_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover2/internal/database"
	"gover2/internal/models"
	"gover2/internal/services"
)

func newExportDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })
	return db
}

func TestBuildCSV_Computers(t *testing.T) {
	db := newExportDB(t)
	_, err := services.InsertComputer(db, models.ComputerInput{Name: "PC1", Model: "Dell", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "computers")
	require.NoError(t, err)
	require.Len(t, rows, 2) // header + 1 data row
	assert.Equal(t, []string{"Name", "Model", "Status", "Purchase Date", "Warranty Expiry"}, rows[0])
	assert.Equal(t, "PC1", rows[1][0])
	assert.Equal(t, "Dell", rows[1][1])
	assert.Equal(t, "active", rows[1][2])
}

func TestBuildCSV_Users(t *testing.T) {
	db := newExportDB(t)
	_, err := services.InsertUser(db, models.UserInput{Name: "Alice", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "users")
	require.NoError(t, err)
	assert.Equal(t, []string{"Name", "Surname", "Status"}, rows[0])
	assert.Equal(t, "Alice", rows[1][0])
}

func TestBuildCSV_Antivirus(t *testing.T) {
	db := newExportDB(t)
	_, err := services.InsertAntivirus(db, models.AntivirusInput{Name: "Norton", LicenseKey: "K1", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "antivirus")
	require.NoError(t, err)
	assert.Equal(t, []string{"Name", "License Key", "Status", "Expiry Date"}, rows[0])
	assert.Equal(t, "Norton", rows[1][0])
	assert.Equal(t, "K1", rows[1][1])
}

func TestBuildCSV_UnknownCategory(t *testing.T) {
	db := newExportDB(t)
	_, err := services.BuildCSV(db, "nope")
	assert.Error(t, err)
}

func TestWriteCSV_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.csv")
	rows := [][]string{{"A", "B"}, {"1", "2"}}
	require.NoError(t, services.WriteCSV(path, rows))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "A,B\n1,2\n", string(data))
}
