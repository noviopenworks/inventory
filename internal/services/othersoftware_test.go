package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/models"
	"inventory/internal/services"
)

func TestListOtherSoftware_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO other_software (name, status) VALUES (?, ?)`, "Slack", "Active")
	require.NoError(t, err)
	result, err := services.ListOtherSoftware(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Slack", result[0].Name)
}

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
