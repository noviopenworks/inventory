package services_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gover/internal/database"
	"gover/internal/models"
	"gover/internal/services"
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

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "expiring", result[0].Severity)
}

func TestGetAlerts_NotExpiringSoon(t *testing.T) {
	db := openTestDB(t)
	// expiry in 45 days — should NOT appear with 30-day window
	notExpiring := time.Now().AddDate(0, 0, 45).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", notExpiring)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Empty(t, result)

	// Also test other_software expiring soon
	expiringSoon := fmt.Sprintf("%s", time.Now().AddDate(0, 0, 10).Format("2006-01-02"))
	_, err = db.Exec(`INSERT INTO other_software (name, status, expiry_date) VALUES (?, ?, ?)`, "Slack", "Active", expiringSoon)
	require.NoError(t, err)
	result, err = services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetAlerts_AlreadyExpired(t *testing.T) {
	db := openTestDB(t)
	// expired yesterday — should appear with severity "expired"
	expired := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "AVG", "Active", expired)
	require.NoError(t, err)
	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "expired", result[0].Severity)
}

func TestGetAlerts_ComputerWarrantyExpiring(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, 15).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO computers (name, status, warranty_expiry) VALUES (?, ?, ?)`, "Dell XPS", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Computer", result[0].Category)
	assert.Equal(t, "expiring", result[0].Severity)
	assert.Equal(t, "Dell XPS", result[0].Name)
}

func TestGetAlerts_SmartphoneWarrantyExpired(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO smartphones (name, status, warranty_expiry) VALUES (?, ?, ?)`, "iPhone", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Smartphone", result[0].Category)
	assert.Equal(t, "expired", result[0].Severity)
}

func TestGetAlerts_TabletWarrantyCovered(t *testing.T) {
	db := openTestDB(t)
	warranty := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO tablets (name, status, warranty_expiry) VALUES (?, ?, ?)`, "iPad", "Active", warranty)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Tablet", result[0].Category)
	assert.Equal(t, "expiring", result[0].Severity)
}

func TestGetAlerts_RetiredExcluded(t *testing.T) {
	db := openTestDB(t)
	expired := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := db.Exec(`INSERT INTO antivirus (name, status, expiry_date) VALUES (?, ?, ?)`, "Old AV", "Retired", expired)
	require.NoError(t, err)

	result, err := services.GetAlerts(db, 30)
	require.NoError(t, err)
	assert.Empty(t, result)
}

// ---------------------------------------------------------------------------
// Write ops — Computers
// ---------------------------------------------------------------------------

func TestComputer_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	// Insert
	id, err := services.InsertComputer(db, models.ComputerInput{
		Name:   "TestPC",
		Model:  "Dell XPS",
		Status: "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Verify appears in list
	list, err := services.ListComputers(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "TestPC", list[0].Name)

	// Update
	err = services.UpdateComputer(db, int(id), models.ComputerInput{
		Name:   "UpdatedPC",
		Model:  "HP EliteBook",
		Status: "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListComputers(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "UpdatedPC", list[0].Name)
	assert.Equal(t, "inactive", list[0].Status)

	// Delete
	err = services.DeleteComputer(db, int(id))
	require.NoError(t, err)

	list, err = services.ListComputers(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Write ops — Smartphones
// ---------------------------------------------------------------------------

func TestSmartphone_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertSmartphone(db, models.SmartphoneInput{
		Name:   "iPhone",
		Model:  "15 Pro",
		Status: "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListSmartphones(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "iPhone", list[0].Name)

	err = services.UpdateSmartphone(db, int(id), models.SmartphoneInput{
		Name:   "Samsung",
		Model:  "Galaxy S24",
		Status: "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListSmartphones(db)
	require.NoError(t, err)
	assert.Equal(t, "Samsung", list[0].Name)

	err = services.DeleteSmartphone(db, int(id))
	require.NoError(t, err)

	list, err = services.ListSmartphones(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Write ops — Tablets
// ---------------------------------------------------------------------------

func TestTablet_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertTablet(db, models.TabletInput{
		Name:   "iPad",
		Model:  "Air 5",
		Status: "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListTablets(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "iPad", list[0].Name)

	err = services.UpdateTablet(db, int(id), models.TabletInput{
		Name:   "Samsung Tab",
		Model:  "S9",
		Status: "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListTablets(db)
	require.NoError(t, err)
	assert.Equal(t, "Samsung Tab", list[0].Name)

	err = services.DeleteTablet(db, int(id))
	require.NoError(t, err)

	list, err = services.ListTablets(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Write ops — Windows Keys
// ---------------------------------------------------------------------------

func TestWindowsKey_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertWindowsKey(db, models.WindowsKeyInput{
		LicenseKey: "AAAAA-BBBBB-CCCCC",
		Status:     "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListWindowsKeys(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "AAAAA-BBBBB-CCCCC", list[0].LicenseKey)

	err = services.UpdateWindowsKey(db, int(id), models.WindowsKeyInput{
		LicenseKey: "DDDDD-EEEEE-FFFFF",
		Status:     "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListWindowsKeys(db)
	require.NoError(t, err)
	assert.Equal(t, "DDDDD-EEEEE-FFFFF", list[0].LicenseKey)

	err = services.DeleteWindowsKey(db, int(id))
	require.NoError(t, err)

	list, err = services.ListWindowsKeys(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Write ops — Antivirus
// ---------------------------------------------------------------------------

func TestAntivirus_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertAntivirus(db, models.AntivirusInput{
		Name:       "Kaspersky",
		LicenseKey: "K-123",
		Status:     "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListAntivirus(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Kaspersky", list[0].Name)

	err = services.UpdateAntivirus(db, int(id), models.AntivirusInput{
		Name:       "Norton",
		LicenseKey: "N-456",
		Status:     "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListAntivirus(db)
	require.NoError(t, err)
	assert.Equal(t, "Norton", list[0].Name)

	err = services.DeleteAntivirus(db, int(id))
	require.NoError(t, err)

	list, err = services.ListAntivirus(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Write ops — OtherSoftware
// ---------------------------------------------------------------------------

func TestOtherSoftware_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertOtherSoftware(db, models.OtherSoftwareInput{
		Name:       "Slack",
		LicenseKey: "SL-001",
		Status:     "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListOtherSoftware(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Slack", list[0].Name)

	err = services.UpdateOtherSoftware(db, int(id), models.OtherSoftwareInput{
		Name:       "Teams",
		LicenseKey: "TM-002",
		Status:     "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListOtherSoftware(db)
	require.NoError(t, err)
	assert.Equal(t, "Teams", list[0].Name)

	err = services.DeleteOtherSoftware(db, int(id))
	require.NoError(t, err)

	list, err = services.ListOtherSoftware(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Write ops — Users
// ---------------------------------------------------------------------------

func TestUser_InsertUpdateDelete(t *testing.T) {
	db := openTestDB(t)

	id, err := services.InsertUser(db, models.UserInput{
		Name:   "Alice",
		Status: "active",
	})
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	list, err := services.ListUsers(db)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Alice", list[0].Name)

	err = services.UpdateUser(db, int(id), models.UserInput{
		Name:   "Bob",
		Status: "inactive",
	})
	require.NoError(t, err)

	list, err = services.ListUsers(db)
	require.NoError(t, err)
	assert.Equal(t, "Bob", list[0].Name)

	err = services.DeleteUser(db, int(id))
	require.NoError(t, err)

	list, err = services.ListUsers(db)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Dropdown helpers
// ---------------------------------------------------------------------------

func TestListUsersForDropdown(t *testing.T) {
	db := openTestDB(t)

	_, err := db.Exec(`INSERT INTO users (name, status) VALUES (?, ?)`, "Alice", "active")
	require.NoError(t, err)

	items, err := services.ListUsersForDropdown(db)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Alice", items[0].Name)
	assert.Greater(t, items[0].ID, 0)
}

func TestListDevicesForDropdown(t *testing.T) {
	db := openTestDB(t)

	_, err := db.Exec(`INSERT INTO computers (name, status) VALUES (?, ?)`, "PC1", "active")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO smartphones (name, status) VALUES (?, ?)`, "Phone1", "active")
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO tablets (name, status) VALUES (?, ?)`, "Tablet1", "active")
	require.NoError(t, err)

	items, err := services.ListDevicesForDropdown(db)
	require.NoError(t, err)
	require.Len(t, items, 3)

	// ordered by kind, name: computer < smartphone < tablet
	assert.Equal(t, "computer", items[0].Kind)
	assert.Equal(t, "PC1", items[0].Name)
	assert.Equal(t, "smartphone", items[1].Kind)
	assert.Equal(t, "Phone1", items[1].Name)
	assert.Equal(t, "tablet", items[2].Kind)
	assert.Equal(t, "Tablet1", items[2].Name)
}
