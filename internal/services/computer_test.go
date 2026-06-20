package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/models"
	"inventory/internal/services"
)

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
