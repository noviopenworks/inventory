package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/database"
	"inventory/internal/models"
	"inventory/internal/services"
)

func newBridgeWithDB(t *testing.T) *bridge.App {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })
	return bridge.NewAppWithDB(db)
}

func TestListComputers_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.ListComputers()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestListComputers_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListComputers()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListSmartphones_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListSmartphones()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListTablets_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListTablets()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListWindowsKeys_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListWindowsKeys()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListAntivirus_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListAntivirus()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListOtherSoftware_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListOtherSoftware()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestListUsers_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	result, err := app.ListUsers()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetDatabasePath_NoConfig(t *testing.T) {
	app := bridge.NewApp()
	path, err := app.GetDatabasePath()
	require.NoError(t, err)
	assert.Equal(t, "", path)
}

func TestGetConfig_Default(t *testing.T) {
	app := bridge.NewApp()
	cfg, err := app.GetConfig()
	require.NoError(t, err)
	assert.Equal(t, "comfortable", cfg.Density)
}

// ---------------------------------------------------------------------------
// Nil-DB guard tests for write methods
// ---------------------------------------------------------------------------

func TestAddComputer_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddComputer(models.ComputerInput{Name: "X", Model: "Y", Status: "active"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestAddSmartphone_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddSmartphone(models.SmartphoneInput{Name: "X", Model: "Y", Status: "active"})
	assert.Error(t, err)
}

func TestAddTablet_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddTablet(models.TabletInput{Name: "X", Model: "Y", Status: "active"})
	assert.Error(t, err)
}

func TestAddWindowsKey_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "KEY", Status: "active"})
	assert.Error(t, err)
}

func TestAddAntivirus_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddAntivirus(models.AntivirusInput{Name: "AV", LicenseKey: "KEY", Status: "active"})
	assert.Error(t, err)
}

func TestAddOtherSoftware_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddOtherSoftware(models.OtherSoftwareInput{Name: "SW", LicenseKey: "KEY", Status: "active"})
	assert.Error(t, err)
}

func TestAddUser_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.AddUser(models.UserInput{Name: "Alice", Status: "active"})
	assert.Error(t, err)
}

func TestListUsersForDropdown_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.ListUsersForDropdown()
	assert.Error(t, err)
}

func TestListDevicesForDropdown_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.ListDevicesForDropdown()
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// Computer CRUD round-trip
// ---------------------------------------------------------------------------

func TestAddComputer_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	input := models.ComputerInput{Name: "ThinkPad", Model: "X1 Carbon", Status: "active"}
	err := app.AddComputer(input)
	require.NoError(t, err)

	list, err := app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "ThinkPad", list[0].Name)
	assert.Equal(t, "X1 Carbon", list[0].Model)
}

func TestUpdateComputer(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddComputer(models.ComputerInput{Name: "OldName", Model: "M1", Status: "active"}))
	list, err := app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	id := list[0].ID

	require.NoError(t, app.UpdateComputer(id, models.ComputerInput{Name: "NewName", Model: "M1", Status: "active"}))

	list, err = app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "NewName", list[0].Name)
}

func TestDeleteComputer(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddComputer(models.ComputerInput{Name: "ToDelete", Model: "M2", Status: "active"}))
	list, err := app.ListComputers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	id := list[0].ID

	require.NoError(t, app.DeleteComputer(id))

	list, err = app.ListComputers()
	require.NoError(t, err)
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Smartphone CRUD round-trip
// ---------------------------------------------------------------------------

func TestAddSmartphone_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "Pixel", Model: "9 Pro", Status: "active"}))
	list, err := app.ListSmartphones()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Pixel", list[0].Name)
}

func TestUpdateSmartphone(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "Old", Model: "M", Status: "active"}))
	list, _ := app.ListSmartphones()
	id := list[0].ID

	require.NoError(t, app.UpdateSmartphone(id, models.SmartphoneInput{Name: "New", Model: "M", Status: "active"}))
	list, _ = app.ListSmartphones()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteSmartphone(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "S", Model: "M", Status: "active"}))
	list, _ := app.ListSmartphones()
	require.NoError(t, app.DeleteSmartphone(list[0].ID))
	list, _ = app.ListSmartphones()
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Tablet CRUD round-trip
// ---------------------------------------------------------------------------

func TestAddTablet_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddTablet(models.TabletInput{Name: "iPad", Model: "Pro 13", Status: "active"}))
	list, err := app.ListTablets()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "iPad", list[0].Name)
}

func TestUpdateTablet(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddTablet(models.TabletInput{Name: "Old", Model: "M", Status: "active"}))
	list, _ := app.ListTablets()
	id := list[0].ID

	require.NoError(t, app.UpdateTablet(id, models.TabletInput{Name: "New", Model: "M", Status: "active"}))
	list, _ = app.ListTablets()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteTablet(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddTablet(models.TabletInput{Name: "T", Model: "M", Status: "active"}))
	list, _ := app.ListTablets()
	require.NoError(t, app.DeleteTablet(list[0].ID))
	list, _ = app.ListTablets()
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// WindowsKey CRUD round-trip
// ---------------------------------------------------------------------------

func TestAddWindowsKey_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "ABCD-1234", Status: "active"}))
	list, err := app.ListWindowsKeys()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "ABCD-1234", list[0].LicenseKey)
}

func TestUpdateWindowsKey(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "OLD-KEY", Status: "active"}))
	list, _ := app.ListWindowsKeys()
	id := list[0].ID

	require.NoError(t, app.UpdateWindowsKey(id, models.WindowsKeyInput{LicenseKey: "NEW-KEY", Status: "active"}))
	list, _ = app.ListWindowsKeys()
	assert.Equal(t, "NEW-KEY", list[0].LicenseKey)
}

func TestDeleteWindowsKey(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddWindowsKey(models.WindowsKeyInput{LicenseKey: "KEY", Status: "active"}))
	list, _ := app.ListWindowsKeys()
	require.NoError(t, app.DeleteWindowsKey(list[0].ID))
	list, _ = app.ListWindowsKeys()
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Antivirus CRUD round-trip
// ---------------------------------------------------------------------------

func TestAddAntivirus_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddAntivirus(models.AntivirusInput{Name: "Norton", LicenseKey: "NRT-001", Status: "active"}))
	list, err := app.ListAntivirus()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Norton", list[0].Name)
}

func TestUpdateAntivirus(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddAntivirus(models.AntivirusInput{Name: "Old", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListAntivirus()
	id := list[0].ID

	require.NoError(t, app.UpdateAntivirus(id, models.AntivirusInput{Name: "New", LicenseKey: "K", Status: "active"}))
	list, _ = app.ListAntivirus()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteAntivirus(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddAntivirus(models.AntivirusInput{Name: "AV", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListAntivirus()
	require.NoError(t, app.DeleteAntivirus(list[0].ID))
	list, _ = app.ListAntivirus()
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// OtherSoftware CRUD round-trip
// ---------------------------------------------------------------------------

func TestAddOtherSoftware_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddOtherSoftware(models.OtherSoftwareInput{Name: "Office", LicenseKey: "OFF-001", Status: "active"}))
	list, err := app.ListOtherSoftware()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Office", list[0].Name)
}

func TestUpdateOtherSoftware(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddOtherSoftware(models.OtherSoftwareInput{Name: "Old", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListOtherSoftware()
	id := list[0].ID

	require.NoError(t, app.UpdateOtherSoftware(id, models.OtherSoftwareInput{Name: "New", LicenseKey: "K", Status: "active"}))
	list, _ = app.ListOtherSoftware()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteOtherSoftware(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddOtherSoftware(models.OtherSoftwareInput{Name: "SW", LicenseKey: "K", Status: "active"}))
	list, _ := app.ListOtherSoftware()
	require.NoError(t, app.DeleteOtherSoftware(list[0].ID))
	list, _ = app.ListOtherSoftware()
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// User CRUD round-trip
// ---------------------------------------------------------------------------

func TestAddUser_RoundTrip(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Alice", Status: "active"}))
	list, err := app.ListUsers()
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Alice", list[0].Name)
}

func TestUpdateUser(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Old", Status: "active"}))
	list, _ := app.ListUsers()
	id := list[0].ID

	require.NoError(t, app.UpdateUser(id, models.UserInput{Name: "New", Status: "active"}))
	list, _ = app.ListUsers()
	assert.Equal(t, "New", list[0].Name)
}

func TestDeleteUser(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Bob", Status: "active"}))
	list, _ := app.ListUsers()
	require.NoError(t, app.DeleteUser(list[0].ID))
	list, _ = app.ListUsers()
	assert.Empty(t, list)
}

// ---------------------------------------------------------------------------
// Dropdown helpers
// ---------------------------------------------------------------------------

func TestListUsersForDropdown(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddUser(models.UserInput{Name: "Alice", Status: "active"}))
	require.NoError(t, app.AddUser(models.UserInput{Name: "Bob", Status: "active"}))

	items, err := app.ListUsersForDropdown()
	require.NoError(t, err)
	require.Len(t, items, 2)
	// Ordered by name: Alice, Bob
	assert.Equal(t, "Alice", items[0].Name)
	assert.Equal(t, "Bob", items[1].Name)
}

func TestListDevicesForDropdown(t *testing.T) {
	app := newBridgeWithDB(t)

	require.NoError(t, app.AddComputer(models.ComputerInput{Name: "PC1", Model: "M", Status: "active"}))
	require.NoError(t, app.AddSmartphone(models.SmartphoneInput{Name: "Phone1", Model: "M", Status: "active"}))
	require.NoError(t, app.AddTablet(models.TabletInput{Name: "Tablet1", Model: "M", Status: "active"}))

	items, err := app.ListDevicesForDropdown()
	require.NoError(t, err)
	require.Len(t, items, 3)

	// Verify all three kinds are present
	kinds := make(map[string]bool)
	for _, item := range items {
		kinds[item.Kind] = true
	}
	assert.True(t, kinds["computer"])
	assert.True(t, kinds["smartphone"])
	assert.True(t, kinds["tablet"])
}

// ---------------------------------------------------------------------------
// ExportCSV tests
// ---------------------------------------------------------------------------

func TestExportCSV_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.ExportCSV("computers")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestExportCSV_BuildRows_Computers(t *testing.T) {
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { db.Close() })

	_, err = services.InsertComputer(db, models.ComputerInput{Name: "PC1", Model: "Dell", Status: "active"})
	require.NoError(t, err)

	rows, err := services.BuildCSV(db, "computers")
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, "PC1", rows[1][0])
}
