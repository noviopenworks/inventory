"""
db.connection – SQLite connection factory and DB_PATH.
"""

import sqlite3
from pathlib import Path

DB_PATH = Path(__file__).resolve().parent.parent / "inventory.db"


def get_connection() -> sqlite3.Connection:
    # Read DB_PATH from the *package* so that monkeypatch / runtime overrides
    # (which set ``db.DB_PATH``) are picked up.
    import db as _db_pkg  # noqa: E402  (deferred to avoid circular import at load time)

    conn = sqlite3.connect(str(_db_pkg.DB_PATH))
    conn.row_factory = sqlite3.Row
    conn.execute("PRAGMA foreign_keys = ON")
    conn.execute("PRAGMA journal_mode = WAL")
    return conn
