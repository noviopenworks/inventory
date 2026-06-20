package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS _meta (
			key   TEXT PRIMARY KEY,
			value TEXT
		);

		CREATE TABLE IF NOT EXISTS users (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			surname    TEXT,
			status     TEXT DEFAULT 'Active',
			notes      TEXT,
			created_at TEXT DEFAULT (datetime('now','localtime')),
			updated_at TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS computers (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			name            TEXT,
			model           TEXT,
			user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
			status          TEXT DEFAULT 'Active',
			purchase_date   TEXT,
			warranty_expiry TEXT,
			notes           TEXT,
			created_at      TEXT DEFAULT (datetime('now','localtime')),
			updated_at      TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS smartphones (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			name            TEXT,
			model           TEXT,
			user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
			status          TEXT DEFAULT 'Active',
			purchase_date   TEXT,
			warranty_expiry TEXT,
			notes           TEXT,
			created_at      TEXT DEFAULT (datetime('now','localtime')),
			updated_at      TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS tablets (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			name            TEXT,
			model           TEXT,
			user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
			status          TEXT DEFAULT 'Active',
			purchase_date   TEXT,
			warranty_expiry TEXT,
			notes           TEXT,
			created_at      TEXT DEFAULT (datetime('now','localtime')),
			updated_at      TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS windows_keys (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			license_key   TEXT,
			computer_id   INTEGER REFERENCES computers(id) ON DELETE SET NULL,
			status        TEXT DEFAULT 'Active',
			notes         TEXT,
			created_at    TEXT DEFAULT (datetime('now','localtime')),
			updated_at    TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS other_software (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT,
			license_key   TEXT,
			computer_id   INTEGER REFERENCES computers(id) ON DELETE SET NULL,
			smartphone_id INTEGER REFERENCES smartphones(id) ON DELETE SET NULL,
			tablet_id     INTEGER REFERENCES tablets(id) ON DELETE SET NULL,
			status        TEXT DEFAULT 'Active',
			expiry_date   TEXT,
			notes         TEXT,
			created_at    TEXT DEFAULT (datetime('now','localtime')),
			updated_at    TEXT DEFAULT (datetime('now','localtime'))
		);

		CREATE TABLE IF NOT EXISTS antivirus (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT,
			license_key   TEXT,
			computer_id   INTEGER REFERENCES computers(id) ON DELETE SET NULL,
			smartphone_id INTEGER REFERENCES smartphones(id) ON DELETE SET NULL,
			tablet_id     INTEGER REFERENCES tablets(id) ON DELETE SET NULL,
			status        TEXT DEFAULT 'Active',
			expiry_date   TEXT,
			notes         TEXT,
			created_at    TEXT DEFAULT (datetime('now','localtime')),
			updated_at    TEXT DEFAULT (datetime('now','localtime'))
		);
	`)
	return err
}
