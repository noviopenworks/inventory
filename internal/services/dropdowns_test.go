package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/services"
)

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
