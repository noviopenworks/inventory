package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/models"
	"inventory/internal/services"
)

func TestListAntivirus_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO antivirus (name, status) VALUES (?, ?)`, "AVG", "Active")
	require.NoError(t, err)
	result, err := services.ListAntivirus(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "AVG", result[0].Name)
}

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
