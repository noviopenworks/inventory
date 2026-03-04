"""
db.schema – DDL initialisation and forward-only migrations.
"""

import sqlite3

from db.connection import get_connection

# Current schema version – bump when adding migrations.
SCHEMA_VERSION = 1


def init_db() -> None:
    with get_connection() as conn:
        conn.executescript("""
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
        """)
        conn.commit()
        _migrate(conn)


def _get_schema_version(conn: sqlite3.Connection) -> int:
    """Read the stored schema version (0 if never set)."""
    try:
        row = conn.execute("SELECT value FROM _meta WHERE key = 'schema_version'").fetchone()
        return int(row[0]) if row else 0
    except sqlite3.OperationalError:
        return 0


def _set_schema_version(conn: sqlite3.Connection, version: int) -> None:
    conn.execute(
        "INSERT OR REPLACE INTO _meta (key, value) VALUES ('schema_version', ?)",
        (str(version),),
    )


def _migrate(conn: sqlite3.Connection) -> None:
    """Run forward-only migrations up to SCHEMA_VERSION."""
    current = _get_schema_version(conn)
    if current >= SCHEMA_VERSION:
        return

    _set_schema_version(conn, SCHEMA_VERSION)
    conn.commit()
