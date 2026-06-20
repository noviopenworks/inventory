package services

import (
	"database/sql"

	"inventory/internal/models"
)

// ListUsers returns all users ordered by name.
func ListUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query(
		`SELECT id, COALESCE(name,''), surname, COALESCE(status,''), notes, created_at, updated_at
		 FROM users ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
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

// InsertUser inserts a new user record and returns its id.
func InsertUser(db *sql.DB, d models.UserInput) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO users (name, surname, status, notes)
		 VALUES (?, ?, ?, ?)`,
		d.Name, d.Surname, d.Status, d.Notes,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateUser updates the user identified by id.
func UpdateUser(db *sql.DB, id int, d models.UserInput) error {
	_, err := db.Exec(
		`UPDATE users SET name=?, surname=?, status=?, notes=?,
		 updated_at=datetime('now','localtime') WHERE id=?`,
		d.Name, d.Surname, d.Status, d.Notes, id,
	)
	return err
}

// DeleteUser removes the user identified by id.
func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE id=?`, id)
	return err
}
