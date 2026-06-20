package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/models"
	"inventory/internal/services"
)

func TestListWindowsKeys_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO windows_keys (license_key, status) VALUES (?, ?)`, "XXXXX-XXXXX", "Active")
	require.NoError(t, err)
	result, err := services.ListWindowsKeys(db)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "XXXXX-XXXXX", result[0].LicenseKey)
}

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
