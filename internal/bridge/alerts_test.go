package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inventory/internal/bridge"
)

func TestApp_GetAlerts_NoDB(t *testing.T) {
	app := bridge.NewApp()
	_, err := app.GetAlerts()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}

func TestApp_GetAlerts_WithDB(t *testing.T) {
	app := newBridgeWithDB(t)
	alerts, err := app.GetAlerts()
	require.NoError(t, err)
	assert.NotNil(t, alerts)
}
