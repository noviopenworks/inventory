package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/models"
	"inventory/internal/services"
)

func TestListTablets_WithData(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Exec(`INSERT INTO tablets (name, model, status) VALUES (?, ?, ?)`, "iPad", "Air 5", "Active")
	require.NoError(t, err)
	result, err := services.ListTablets(db)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "iPad", result[0].Name)
}

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
