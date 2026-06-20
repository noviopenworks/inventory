package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
	"inventory/internal/database"
)

func newBridgeWithDB(t *testing.T) *bridge.App {
	t.Helper()
	db, err := database.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, database.InitSchema(db))
	t.Cleanup(func() { _ = db.Close() })
	return bridge.NewAppWithDB(db)
}
