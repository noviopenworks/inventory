package services

import (
	"database/sql"

	"gover2/internal/models"
)

// ListUsersForDropdown returns id and name for all users, ordered by name.
func ListUsersForDropdown(db *sql.DB) ([]models.DropdownItem, error) {
	rows, err := db.Query(`SELECT id, COALESCE(name,'') FROM users ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.DropdownItem
	for rows.Next() {
		var item models.DropdownItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

// ListDevicesForDropdown returns id, name and kind for all computers,
// smartphones and tablets, ordered by kind then name.
func ListDevicesForDropdown(db *sql.DB) ([]models.DeviceDropdownItem, error) {
	rows, err := db.Query(`
		SELECT id, COALESCE(name,'') AS name, 'computer' AS kind FROM computers
		UNION ALL
		SELECT id, COALESCE(name,''), 'smartphone' FROM smartphones
		UNION ALL
		SELECT id, COALESCE(name,''), 'tablet'     FROM tablets
		ORDER BY kind, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.DeviceDropdownItem
	for rows.Next() {
		var item models.DeviceDropdownItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Kind); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
