package services

import (
	"database/sql"

	"inventory/internal/models"
)

// ListTablets returns all tablets ordered by name.
func ListTablets(db *sql.DB) ([]models.Tablet, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM tablets ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.Tablet
	for rows.Next() {
		var tab models.Tablet
		if err := rows.Scan(&tab.ID, &tab.Name, &tab.Model, &tab.UserID, &tab.Status,
			&tab.PurchaseDate, &tab.WarrantyExpiry, &tab.Notes, &tab.CreatedAt, &tab.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, tab)
	}
	return out, rows.Err()
}

// InsertTablet inserts a new tablet record and returns its id.
func InsertTablet(db *sql.DB, d models.TabletInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO tablets (name, model, user_id, status, purchase_date, warranty_expiry, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateTablet updates the tablet identified by id.
func UpdateTablet(db *sql.DB, id int, d models.TabletInput) error {
	_, err := db.Exec(
		`UPDATE tablets SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

// DeleteTablet removes the tablet identified by id.
func DeleteTablet(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tablets WHERE id=?`, id)
	return err
}
