package services

import (
	"database/sql"

	"inventory/internal/models"
)

// ListSmartphones returns all smartphones ordered by name.
func ListSmartphones(db *sql.DB) ([]models.Smartphone, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM smartphones ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []models.Smartphone
	for rows.Next() {
		var s models.Smartphone
		if err := rows.Scan(&s.ID, &s.Name, &s.Model, &s.UserID, &s.Status,
			&s.PurchaseDate, &s.WarrantyExpiry, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// InsertSmartphone inserts a new smartphone record and returns its id.
func InsertSmartphone(db *sql.DB, d models.SmartphoneInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO smartphones (name, model, user_id, status, purchase_date, warranty_expiry, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateSmartphone updates the smartphone identified by id.
func UpdateSmartphone(db *sql.DB, id int, d models.SmartphoneInput) error {
	_, err := db.Exec(
		`UPDATE smartphones SET name=?, model=?, user_id=?, status=?, purchase_date=?,
		 warranty_expiry=?, notes=?, updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Model, d.UserID, d.Status, d.PurchaseDate, d.WarrantyExpiry, d.Notes, id,
	)
	return err
}

// DeleteSmartphone removes the smartphone identified by id.
func DeleteSmartphone(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM smartphones WHERE id=?`, id)
	return err
}
