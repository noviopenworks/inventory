"""
tests/test_db_schema.py – Schema initialisation, migration, and TABLE_CONFIG.
"""

import sqlite3

import db as db_module

# ── Init ───────────────────────────────────────────────────────────────────────


class TestInitDb:
    def test_all_tables_created(self, tmp_db):
        with sqlite3.connect(str(tmp_db)) as conn:
            tables = {
                r[0]
                for r in conn.execute(
                    "SELECT name FROM sqlite_master WHERE type='table'"
                ).fetchall()
            }
        expected = {
            "users",
            "computers",
            "smartphones",
            "tablets",
            "windows_keys",
            "other_software",
            "antivirus",
            "_meta",
        }
        assert expected.issubset(tables)

    def test_no_subscriptions_table(self, tmp_db):
        """The old 'subscriptions' table must not be created by init_db."""
        with sqlite3.connect(str(tmp_db)) as conn:
            tables = {
                r[0]
                for r in conn.execute(
                    "SELECT name FROM sqlite_master WHERE type='table'"
                ).fetchall()
            }
        assert "subscriptions" not in tables

    def test_idempotent(self, tmp_db):
        """Calling init_db() twice must not raise."""
        db_module.init_db()

    def test_meta_table_has_schema_version(self, tmp_db):
        """After init_db the _meta table stores the current schema version."""
        with sqlite3.connect(str(tmp_db)) as conn:
            row = conn.execute("SELECT value FROM _meta WHERE key = 'schema_version'").fetchone()
        assert row is not None
        assert int(row[0]) == db_module.SCHEMA_VERSION

    def test_wal_mode_enabled(self, tmp_db):
        """get_connection() should activate WAL journal mode."""
        conn = db_module.get_connection()
        mode = conn.execute("PRAGMA journal_mode").fetchone()[0]
        conn.close()
        assert mode == "wal"


# ── Migration ──────────────────────────────────────────────────────────────────


class TestMigration:
    def test_migrate_from_zero(self, tmp_db):
        """If schema_version is missing, _migrate sets it to SCHEMA_VERSION."""
        with sqlite3.connect(str(tmp_db)) as conn:
            conn.execute("DELETE FROM _meta WHERE key = 'schema_version'")
            conn.commit()
        with db_module.get_connection() as conn:
            db_module._migrate(conn)
        with sqlite3.connect(str(tmp_db)) as conn:
            row = conn.execute("SELECT value FROM _meta WHERE key = 'schema_version'").fetchone()
        assert int(row[0]) == db_module.SCHEMA_VERSION

    def test_migrate_skips_when_current(self, tmp_db):
        """_migrate is a no-op when version is already at SCHEMA_VERSION."""
        with db_module.get_connection() as conn:
            db_module._migrate(conn)
            row = conn.execute("SELECT value FROM _meta WHERE key = 'schema_version'").fetchone()
        assert int(row[0]) == db_module.SCHEMA_VERSION


# ── TABLE_CONFIG ───────────────────────────────────────────────────────────────


class TestTableConfig:
    def test_all_categories_have_required_keys(self, tmp_db):
        for cat, spec in db_module.TABLE_CONFIG.items():
            assert "table" in spec, f"{cat} missing 'table'"
            assert "db_cols" in spec, f"{cat} missing 'db_cols'"
            assert "display_cols" in spec, f"{cat} missing 'display_cols'"
            assert "has_user_fk" in spec, f"{cat} missing 'has_user_fk'"
