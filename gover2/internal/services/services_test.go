package services_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover2/internal/database"
	"gover2/internal/services"
)

// openTestDB opens an in-memory SQLite database with schema initialized.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })
	return db
}

func TestListComputers_Empty(t *testing.T) {
	db := openTestDB(t)
	result, err := services.ListComputers(db)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestListComputers_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO computers (name, model, status) VALUES (?, ?, ?)`, "PC1", "Dell XPS", "Active")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO computers (name, model, status) VALUES (?, ?, ?)`, "PC2", "HP Elite", "Repair")
	require.NoError(t, err)

	result, err := services.ListComputers(db)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "PC1", result[0].Name)
	assert.Equal(t, "PC2", result[1].Name)
}

func TestListSmartphones_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO smartphones (name, model, status) VALUES (?, ?, ?)`, "iPhone", "15 Pro", "Active")
	require.NoError(t, err)
	result, err := services.ListSmartphones(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "iPhone", result[0].Name)
}

func TestListTablets_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO tablets (name, model, status) VALUES (?, ?, ?)`, "iPad", "Air 5", "Active")
	require.NoError(t, err)
	result, err := services.ListTablets(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "iPad", result[0].Name)
}

func TestListWindowsKeys_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO windows_keys (license_key, status) VALUES (?, ?)`, "XXXXX-XXXXX", "Active")
	require.NoError(t, err)
	result, err := services.ListWindowsKeys(db)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "XXXXX-XXXXX", result[0].LicenseKey)
}

func TestListAntivirus_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO antivirus (name, status) VALUES (?, ?)`, "AVG", "Active")
	require.NoError(t, err)
	result, err := services.ListAntivirus(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "AVG", result[0].Name)
}

func TestListOtherSoftware_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO other_software (name, status) VALUES (?, ?)`, "Slack", "Active")
	require.NoError(t, err)
	result, err := services.ListOtherSoftware(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Slack", result[0].Name)
}

func TestListUsers_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO users (name, surname) VALUES (?, ?)`, "Alice", "Smith")
	require.NoError(t, err)
	result, err := services.ListUsers(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Alice", result[0].Name)
}

func TestGetAlerts_ExpiringWithin30Days(t *testing.T) {
	db := openTestDB(t)
	// expiry in 15 days — should appear
	expiring := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", expiring)
	require.NoError(t, err)

	result, err := services.GetAlerts(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetAlerts_NotExpiringSoon(t *testing.T) {
	db := openTestDB(t)
	// expiry in 45 days — should NOT appear
	notExpiring := time.Now().AddDate(0, 0, 45).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", notExpiring)
	require.NoError(t, err)

	result, err := services.GetAlerts(db)
	require.NoError(t, err)
	assert.Empty(t, result)

	// Also test other_software
	expiringSoon := fmt.Sprintf("%s", time.Now().AddDate(0, 0, 10).Format("2006-01-02"))
	_, err = db.Exec(`INSERT INTO other_software (name, status, expiry_date) VALUES (?, ?, ?)`, "Slack", "Active", expiringSoon)
	require.NoError(t, err)
	result, err = services.GetAlerts(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetAlerts_AlreadyExpired(t *testing.T) {
	db := openTestDB(t)
	// expired yesterday — should appear (expired = also ≤30 days from now in the past)
	expired := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", expired)
	require.NoError(t, err)
	result, err := services.GetAlerts(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}
