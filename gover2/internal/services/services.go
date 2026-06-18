package services

import (
	"database/sql"
	"time"

	"gover2/internal/models"
)

func ListComputers(db *sql.DB) ([]models.Computer, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM computers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Computer
	for rows.Next() {
		var c models.Computer
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.UserID, &c.Status,
			&c.PurchaseDate, &c.WarrantyExpiry, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func ListSmartphones(db *sql.DB) ([]models.Smartphone, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM smartphones ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

func ListTablets(db *sql.DB) ([]models.Tablet, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(model,''), user_id, COALESCE(status,''),
		        purchase_date, warranty_expiry, notes, created_at, updated_at
		 FROM tablets ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

func ListWindowsKeys(db *sql.DB) ([]models.WindowsKey, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(license_key,''), computer_id, COALESCE(status,''),
		        notes, created_at, updated_at
		 FROM windows_keys ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.WindowsKey
	for rows.Next() {
		var w models.WindowsKey
		if err := rows.Scan(&w.ID, &w.LicenseKey, &w.ComputerID, &w.Status,
			&w.Notes, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func ListAntivirus(db *sql.DB) ([]models.Antivirus, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(license_key,''), computer_id, smartphone_id, tablet_id,
		        COALESCE(status,''), expiry_date, notes, created_at, updated_at
		 FROM antivirus ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Antivirus
	for rows.Next() {
		var a models.Antivirus
		if err := rows.Scan(&a.ID, &a.Name, &a.LicenseKey, &a.ComputerID, &a.SmartphoneID,
			&a.TabletID, &a.Status, &a.ExpiryDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func ListOtherSoftware(db *sql.DB) ([]models.OtherSoftware, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), COALESCE(license_key,''), computer_id, smartphone_id, tablet_id,
		        COALESCE(status,''), expiry_date, notes, created_at, updated_at
		 FROM other_software ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.OtherSoftware
	for rows.Next() {
		var s models.OtherSoftware
		if err := rows.Scan(&s.ID, &s.Name, &s.LicenseKey, &s.ComputerID, &s.SmartphoneID,
			&s.TabletID, &s.Status, &s.ExpiryDate, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func ListUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), surname, COALESCE(status,''), notes, created_at, updated_at
		 FROM users ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Status,
			&u.Notes, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func GetAlerts(db *sql.DB) ([]models.Alert, error) {
	threshold := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	query := `
		SELECT 'antivirus' as category, id, COALESCE(name, '') as name, expiry_date,
		       CAST(julianday(expiry_date) - julianday('now') AS INTEGER) as days_remaining
		FROM antivirus
		WHERE expiry_date IS NOT NULL AND expiry_date <= ?
		UNION ALL
		SELECT 'other_software', id, COALESCE(name, ''), expiry_date,
		       CAST(julianday(expiry_date) - julianday('now') AS INTEGER)
		FROM other_software
		WHERE expiry_date IS NOT NULL AND expiry_date <= ?
		ORDER BY expiry_date`
	rows, err := db.Query(query, threshold, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Alert
	for rows.Next() {
		var a models.Alert
		var daysRemaining int
		if err := rows.Scan(&a.Category, &a.ID, &a.Name, &a.ExpiryDate, &daysRemaining); err != nil {
			return nil, err
		}
		a.DaysRemaining = daysRemaining
		if daysRemaining < 0 {
			a.Severity = "critical"
		} else if daysRemaining <= 7 {
			a.Severity = "high"
		} else {
			a.Severity = "warning"
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
