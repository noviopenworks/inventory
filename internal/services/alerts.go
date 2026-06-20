// Package services provides database access functions for the inventory application.
package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"inventory/internal/models"
)

// alertSources lists every (table, date column) pair that contributes to
// expiry/warranty alerts. Mirrors db/alerts.py: Windows Key is intentionally
// omitted (no expiry column); Users is irrelevant.
type alertSource struct {
	table    string
	category string
	dateCol  string
}

var alertSources = []alertSource{
	{"computers", "Computer", "warranty_expiry"},
	{"smartphones", "Smartphone", "warranty_expiry"},
	{"tablets", "Tablet", "warranty_expiry"},
	{"antivirus", "Antivirus", "expiry_date"},
	{"other_software", "Other Software", "expiry_date"},
}

// GetAlerts returns expiry and warranty alerts for items expiring within warningDays.
func GetAlerts(db *sql.DB, warningDays int) ([]models.Alert, error) {
	threshold := fmt.Sprintf("+%d days", warningDays)
	parts := make([]string, 0, len(alertSources))
	params := make([]any, 0, len(alertSources))
	for _, src := range alertSources {
		parts = append(parts, fmt.Sprintf(`
		SELECT '%s' AS category, id, COALESCE(name, '') AS name, %s AS expiry_date,
		       CAST(julianday(%s) - julianday('now') AS INTEGER) AS days_remaining
		FROM %s
		WHERE %s IS NOT NULL
		  AND date(%s) <= date('now', ?)
		  AND status != 'Retired'`,
			src.category, src.dateCol, src.dateCol, src.table, src.dateCol, src.dateCol))
		params = append(params, threshold)
	}
	query := strings.Join(parts, "\n UNION ALL \n") + "\n ORDER BY days_remaining ASC"

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	today := time.Now().Format("2006-01-02")
	var out []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.Category, &a.ID, &a.Name, &a.ExpiryDate, &a.DaysRemaining); err != nil {
			return nil, err
		}
		if a.ExpiryDate < today {
			a.Severity = "expired"
		} else {
			a.Severity = "expiring"
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
