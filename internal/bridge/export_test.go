package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"inventory/internal/bridge"
)

func TestExportCSV_NoDB(t *testing.T) {
	app := bridge.NewApp()
	err := app.ExportCSV("computers")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database")
}
