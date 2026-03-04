"""
tests/test_db.py – Unit tests for db.py CRUD, helpers, FK behaviour,
                   search, the All-tab UNION, and expiry alerts.

Each test receives an isolated temporary SQLite file via the `tmp_db`
fixture so the production inventory.db is never touched.
"""

import sqlite3
from datetime import date, timedelta

import pytest

import db as db_module

# ── Fixture ────────────────────────────────────────────────────────────────────


@pytest.fixture()
def tmp_db(tmp_path, monkeypatch):
    """Point DB_PATH at a fresh temp file and initialise the schema."""
    db_file = tmp_path / "test_inventory.db"
    monkeypatch.setattr(db_module, "DB_PATH", db_file)
    db_module.init_db()
    yield db_file


# ── Convenience helpers (keep tests DRY) ──────────────────────────────────────


def _add_user(name: str = "Alice", surname: str = "Smith") -> int:
    return db_module.insert_record("Users", {"name": name, "surname": surname, "status": "Active"})


def _add_computer(name: str = "PC-01", user_id: int | None = None) -> int:
    return db_module.insert_record(
        "Computer",
        {"name": name, "model": "ThinkPad X1", "user_id": user_id, "status": "Active"},
    )


def _add_phone(name: str = "iPhone-01", user_id: int | None = None) -> int:
    return db_module.insert_record(
        "Smartphone",
        {
            "name": name,
            "model": "iPhone 15",
            "user_id": user_id,
            "status": "Active",
        },
    )


def _days(delta: int) -> str:
    """ISO date relative to today."""
    return (date.today() + timedelta(days=delta)).isoformat()


# ── Schema ─────────────────────────────────────────────────────────────────────


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


class TestMigration:
    def test_migrate_from_zero(self, tmp_db):
        """If schema_version is missing, _migrate sets it to SCHEMA_VERSION."""
        with sqlite3.connect(str(tmp_db)) as conn:
            conn.execute("DELETE FROM _meta WHERE key = 'schema_version'")
            conn.commit()
        # Re-run migration
        with db_module.get_connection() as conn:
            db_module._migrate(conn)
        with sqlite3.connect(str(tmp_db)) as conn:
            row = conn.execute("SELECT value FROM _meta WHERE key = 'schema_version'").fetchone()
        assert int(row[0]) == db_module.SCHEMA_VERSION

    def test_migrate_skips_when_current(self, tmp_db):
        """_migrate is a no-op when version is already at SCHEMA_VERSION."""
        with db_module.get_connection() as conn:
            db_module._migrate(conn)  # should not raise
            row = conn.execute("SELECT value FROM _meta WHERE key = 'schema_version'").fetchone()
        assert int(row[0]) == db_module.SCHEMA_VERSION


class TestTableConfig:
    def test_all_categories_have_required_keys(self, tmp_db):
        for cat, spec in db_module.TABLE_CONFIG.items():
            assert "table" in spec, f"{cat} missing 'table'"
            assert "db_cols" in spec, f"{cat} missing 'db_cols'"
            assert "display_cols" in spec, f"{cat} missing 'display_cols'"
            assert "has_user_fk" in spec, f"{cat} missing 'has_user_fk'"


# ── Users ──────────────────────────────────────────────────────────────────────


class TestUsers:
    def test_insert_returns_int_id(self, tmp_db):
        uid = _add_user()
        assert isinstance(uid, int) and uid > 0

    def test_fetch_contains_inserted_record(self, tmp_db):
        _add_user("Bob", "Jones")
        users = db_module.fetch_users()
        assert len(users) == 1
        assert users[0]["name"] == "Bob"
        assert users[0]["surname"] == "Jones"

    def test_fetch_ordered_by_name(self, tmp_db):
        _add_user("Zara")
        _add_user("Alice")
        names = [u["name"] for u in db_module.fetch_users()]
        assert names == sorted(names)

    def test_update_surname(self, tmp_db):
        uid = _add_user("Charlie")
        db_module.update_record("Users", uid, {"surname": "Brown"})
        row = db_module.fetch_records("Users")[0]
        assert row["surname"] == "Brown"

    def test_update_status(self, tmp_db):
        uid = _add_user("Dave")
        db_module.update_record("Users", uid, {"status": "Retired"})
        row = db_module.fetch_records("Users")[0]
        assert row["status"] == "Retired"

    def test_delete_removes_record(self, tmp_db):
        uid = _add_user()
        db_module.delete_record("Users", uid)
        assert db_module.fetch_users() == []

    def test_insert_empty_data_does_not_raise(self, tmp_db):
        """DEFAULT VALUES path: name has NOT NULL but uv run will hit the constraint;
        this confirms the real guard is the DB constraint, not our code."""
        with pytest.raises(Exception):
            db_module.insert_record("Users", {})


# ── Computers ──────────────────────────────────────────────────────────────────


class TestComputers:
    def test_insert_with_user_fk_shows_user_name(self, tmp_db):
        uid = _add_user()
        cid = _add_computer(user_id=uid)
        rows = db_module.fetch_records("Computer")
        assert len(rows) == 1
        assert rows[0]["user_name"] == "Alice"
        assert rows[0]["id"] == cid

    def test_insert_without_user_returns_null_user_name(self, tmp_db):
        _add_computer()
        rows = db_module.fetch_records("Computer")
        assert rows[0]["user_name"] is None

    def test_fetch_computers_helper_returns_all(self, tmp_db):
        _add_computer("PC-A")
        _add_computer("PC-B")
        computers = db_module.fetch_computers()
        assert len(computers) == 2
        names = {c["name"] for c in computers}
        assert names == {"PC-A", "PC-B"}

    def test_fetch_computers_ordered_by_name(self, tmp_db):
        _add_computer("Zeta")
        _add_computer("Alpha")
        names = [c["name"] for c in db_module.fetch_computers()]
        assert names == sorted(names)

    def test_update_status(self, tmp_db):
        cid = _add_computer()
        db_module.update_record("Computer", cid, {"status": "Retired"})
        row = db_module.fetch_records("Computer")[0]
        assert row["status"] == "Retired"

    def test_update_noop_when_no_matching_cols(self, tmp_db):
        """update_record with no recognised columns must not raise."""
        cid = _add_computer()
        db_module.update_record("Computer", cid, {"nonexistent_col": "x"})

    def test_delete(self, tmp_db):
        cid = _add_computer()
        db_module.delete_record("Computer", cid)
        assert db_module.fetch_records("Computer") == []

    def test_user_deleted_sets_computer_user_id_null(self, tmp_db):
        uid = _add_user()
        _add_computer(user_id=uid)
        db_module.delete_record("Users", uid)
        row = db_module.fetch_records("Computer")[0]
        assert row["user_id"] is None

    def test_multiple_computers_ordered_by_id_desc(self, tmp_db):
        _add_computer("First")
        id2 = _add_computer("Second")
        rows = db_module.fetch_records("Computer")
        assert rows[0]["id"] == id2  # newest first


# ── Smartphones ────────────────────────────────────────────────────────────────


class TestSmartphones:
    def test_insert_with_user_fk(self, tmp_db):
        uid = _add_user()
        _add_phone(user_id=uid)
        rows = db_module.fetch_records("Smartphone")
        assert len(rows) == 1
        assert rows[0]["user_name"] == "Alice"

    def test_fetch_smartphones_helper(self, tmp_db):
        _add_phone("Phone-A")
        _add_phone("Phone-B")
        phones = db_module.fetch_smartphones()
        assert len(phones) == 2

    def test_user_deleted_sets_smartphone_user_id_null(self, tmp_db):
        uid = _add_user()
        _add_phone(user_id=uid)
        db_module.delete_record("Users", uid)
        row = db_module.fetch_records("Smartphone")[0]
        assert row["user_id"] is None

    def test_delete(self, tmp_db):
        sid = _add_phone()
        db_module.delete_record("Smartphone", sid)
        assert db_module.fetch_records("Smartphone") == []


# ── Windows Keys ───────────────────────────────────────────────────────────────


class TestWindowsKeys:
    def test_insert_linked_to_computer_shows_name(self, tmp_db):
        cid = _add_computer()
        db_module.insert_record(
            "Windows Key",
            {"license_key": "XXXXX-XXXXX", "computer_id": cid},
        )
        row = db_module.fetch_records("Windows Key")[0]
        assert row["device_name"] == "PC-01"

    def test_insert_without_computer_returns_null(self, tmp_db):
        db_module.insert_record("Windows Key", {"license_key": "KEY-1"})
        row = db_module.fetch_records("Windows Key")[0]
        assert row["device_name"] is None

    def test_computer_deleted_sets_fk_null(self, tmp_db):
        cid = _add_computer()
        db_module.insert_record("Windows Key", {"license_key": "X", "computer_id": cid})
        db_module.delete_record("Computer", cid)
        row = db_module.fetch_records("Windows Key")[0]
        assert row["computer_id"] is None

    def test_insert_empty_data_does_not_raise(self, tmp_db):
        kid = db_module.insert_record("Windows Key", {})
        assert isinstance(kid, int)

    def test_update_license_key(self, tmp_db):
        kid = db_module.insert_record("Windows Key", {"license_key": "OLD"})
        db_module.update_record("Windows Key", kid, {"license_key": "NEW-KEY"})
        row = db_module.fetch_records("Windows Key")[0]
        assert row["license_key"] == "NEW-KEY"

    def test_delete(self, tmp_db):
        kid = db_module.insert_record("Windows Key", {"license_key": "K"})
        db_module.delete_record("Windows Key", kid)
        assert db_module.fetch_records("Windows Key") == []


# ── Other Software ─────────────────────────────────────────────────────────────


class TestOtherSoftware:
    def test_linked_to_computer(self, tmp_db):
        cid = _add_computer()
        db_module.insert_record(
            "Other Software",
            {"name": "Adobe CC", "computer_id": cid},
        )
        row = db_module.fetch_records("Other Software")[0]
        assert row["device_name"] == "PC-01"

    def test_linked_to_smartphone(self, tmp_db):
        sid = _add_phone()
        db_module.insert_record("Other Software", {"name": "App", "smartphone_id": sid})
        row = db_module.fetch_records("Other Software")[0]
        assert row["device_name"] == "iPhone-01"

    def test_linked_to_both(self, tmp_db):
        cid = _add_computer()
        sid = _add_phone()
        db_module.insert_record(
            "Other Software", {"name": "Multi", "computer_id": cid, "smartphone_id": sid}
        )
        row = db_module.fetch_records("Other Software")[0]
        # COALESCE picks computer first
        assert row["device_name"] == "PC-01"

    def test_fks_set_null_after_device_delete(self, tmp_db):
        cid = _add_computer()
        sid = _add_phone()
        db_module.insert_record(
            "Other Software", {"name": "App", "computer_id": cid, "smartphone_id": sid}
        )
        db_module.delete_record("Computer", cid)
        db_module.delete_record("Smartphone", sid)
        row = db_module.fetch_records("Other Software")[0]
        assert row["computer_id"] is None
        assert row["smartphone_id"] is None

    def test_delete(self, tmp_db):
        oid = db_module.insert_record("Other Software", {"name": "X"})
        db_module.delete_record("Other Software", oid)
        assert db_module.fetch_records("Other Software") == []


# ── Antivirus ──────────────────────────────────────────────────────────────────


class TestAntivirus:
    def test_linked_to_computer(self, tmp_db):
        cid = _add_computer()
        db_module.insert_record("Antivirus", {"name": "Norton 360", "computer_id": cid})
        row = db_module.fetch_records("Antivirus")[0]
        assert row["device_name"] == "PC-01"

    def test_linked_to_smartphone(self, tmp_db):
        sid = _add_phone()
        db_module.insert_record("Antivirus", {"name": "Lookout", "smartphone_id": sid})
        row = db_module.fetch_records("Antivirus")[0]
        assert row["device_name"] == "iPhone-01"

    def test_linked_to_both(self, tmp_db):
        cid = _add_computer()
        sid = _add_phone()
        db_module.insert_record(
            "Antivirus", {"name": "MB", "computer_id": cid, "smartphone_id": sid}
        )
        row = db_module.fetch_records("Antivirus")[0]
        # COALESCE picks computer first
        assert row["device_name"] == "PC-01"

    def test_insert_empty(self, tmp_db):
        aid = db_module.insert_record("Antivirus", {})
        assert isinstance(aid, int)

    def test_delete(self, tmp_db):
        aid = db_module.insert_record("Antivirus", {"name": "AV"})
        db_module.delete_record("Antivirus", aid)
        assert db_module.fetch_records("Antivirus") == []


# ── Search ─────────────────────────────────────────────────────────────────────


class TestSearch:
    def test_search_by_name(self, tmp_db):
        _add_computer("Workstation-Alpha")
        _add_computer("Laptop-Beta")
        results = db_module.fetch_records("Computer", search="Alpha")
        assert len(results) == 1
        assert results[0]["name"] == "Workstation-Alpha"

    def test_search_is_case_insensitive(self, tmp_db):
        _add_computer("MyLaptop")
        assert len(db_module.fetch_records("Computer", search="mylaptop")) == 1

    def test_search_partial_match(self, tmp_db):
        _add_computer("MyLaptop-001")
        assert len(db_module.fetch_records("Computer", search="Laptop")) == 1

    def test_search_by_user_name_via_join(self, tmp_db):
        uid = _add_user("FindMe")
        _add_computer(user_id=uid)
        _add_computer("OtherPC")
        results = db_module.fetch_records("Computer", search="FindMe")
        assert len(results) == 1
        assert results[0]["user_name"] == "FindMe"

    def test_search_by_computer_name_in_windows_key(self, tmp_db):
        cid = _add_computer("TargetPC")
        db_module.insert_record("Windows Key", {"license_key": "Key1", "computer_id": cid})
        db_module.insert_record("Windows Key", {"license_key": "Key2"})
        results = db_module.fetch_records("Windows Key", search="TargetPC")
        assert len(results) == 1
        assert results[0]["device_name"] == "TargetPC"

    def test_search_by_smartphone_name_in_antivirus(self, tmp_db):
        sid = _add_phone("TargetPhone")
        db_module.insert_record("Antivirus", {"name": "AV1", "smartphone_id": sid})
        db_module.insert_record("Antivirus", {"name": "AV2"})
        results = db_module.fetch_records("Antivirus", search="TargetPhone")
        assert len(results) == 1

    def test_none_search_returns_all(self, tmp_db):
        _add_computer("PC1")
        _add_computer("PC2")
        assert len(db_module.fetch_records("Computer", search=None)) == 2

    def test_empty_string_search_returns_all(self, tmp_db):
        _add_computer("PC1")
        _add_computer("PC2")
        assert len(db_module.fetch_records("Computer", search="")) == 2

    def test_no_results_for_nonexistent_term(self, tmp_db):
        _add_computer("PC1")
        assert db_module.fetch_records("Computer", search="ZZZNOMATCH") == []


# ── All-tab UNION ALL ──────────────────────────────────────────────────────────


class TestFetchAll:
    def test_all_asset_categories_present(self, tmp_db):
        _add_computer()
        _add_phone()
        db_module.insert_record("Tablet", {"name": "iPad", "status": "Active"})
        rows = db_module.fetch_records("All")
        types = {r["type"] for r in rows}
        assert types == {
            "Computer",
            "Smartphone",
            "Tablet",
        }

    def test_software_excluded_from_all(self, tmp_db):
        db_module.insert_record("Windows Key", {"license_key": "WK"})
        db_module.insert_record("Antivirus", {"name": "AV"})
        db_module.insert_record("Other Software", {"name": "SW"})
        rows = db_module.fetch_records("All")
        assert rows == []

    def test_users_excluded_from_all(self, tmp_db):
        _add_user()
        rows = db_module.fetch_records("All")
        assert all(r["type"] != "Users" for r in rows)

    def test_empty_db_returns_empty_list(self, tmp_db):
        assert db_module.fetch_records("All") == []

    def test_search_across_hardware_categories(self, tmp_db):
        _add_computer("UniqueComputer")
        _add_phone("UniquePhone")
        db_module.insert_record("Antivirus", {"name": "UniqueAV"})
        results = db_module.fetch_records("All", search="Unique")
        assert len(results) == 2
        types = {r["type"] for r in results}
        assert types == {"Computer", "Smartphone"}

    def test_device_rows_show_user_name(self, tmp_db):
        """Computer rows show user_name."""
        uid = _add_user("Bob")
        _add_computer(user_id=uid)
        rows = db_module.fetch_records("All")
        pc = next(r for r in rows if r["type"] == "Computer")
        assert pc["user_name"] == "Bob"

    def test_device_rows_show_model(self, tmp_db):
        """All tab includes model column."""
        _add_computer()
        rows = db_module.fetch_records("All")
        pc = next(r for r in rows if r["type"] == "Computer")
        assert pc["model"] == "ThinkPad X1"


# ── Expiry alerts ──────────────────────────────────────────────────────────────


class TestExpiryAlerts:
    def test_expired_item_appears_in_get_expired(self, tmp_db):
        db_module.insert_record(
            "Antivirus",
            {"name": "OldAV", "expiry_date": _days(-10), "status": "Active"},
        )
        results = db_module.get_expired()
        assert any(r["name"] == "OldAV" for r in results)

    def test_future_item_not_in_expired(self, tmp_db):
        db_module.insert_record(
            "Antivirus",
            {"name": "FutureAV", "expiry_date": _days(60), "status": "Active"},
        )
        assert db_module.get_expired() == []

    def test_expiring_soon_default_window(self, tmp_db):
        db_module.insert_record(
            "Antivirus",
            {"name": "NearExpiry", "expiry_date": _days(15), "status": "Active"},
        )
        results = db_module.get_expiring_soon()
        assert any(r["name"] == "NearExpiry" for r in results)

    def test_not_expiring_soon_outside_window(self, tmp_db):
        db_module.insert_record(
            "Antivirus",
            {"name": "FarFuture", "expiry_date": _days(90), "status": "Active"},
        )
        assert not any(r["name"] == "FarFuture" for r in db_module.get_expiring_soon(days=30))

    def test_expiring_soon_excludes_retired(self, tmp_db):
        db_module.insert_record(
            "Antivirus",
            {"name": "RetiredAV", "expiry_date": _days(5), "status": "Retired"},
        )
        assert all(r["name"] != "RetiredAV" for r in db_module.get_expiring_soon())

    def test_expired_excludes_retired(self, tmp_db):
        db_module.insert_record(
            "Other Software",
            {"name": "RetiredSW", "expiry_date": _days(-5), "status": "Retired"},
        )
        assert db_module.get_expired() == []

    def test_warranty_expiry_triggers_get_expired_for_computer(self, tmp_db):
        cid = _add_computer()
        db_module.update_record("Computer", cid, {"warranty_expiry": _days(-1)})
        results = db_module.get_expired()
        assert any(r["type"] == "Computer" for r in results)

    def test_warranty_expiring_soon_for_computer(self, tmp_db):
        cid = _add_computer()
        db_module.update_record("Computer", cid, {"warranty_expiry": _days(10)})
        results = db_module.get_expiring_soon(days=30)
        assert any(r["type"] == "Computer" for r in results)

    def test_multiple_categories_in_expired(self, tmp_db):
        cid = _add_computer()
        db_module.update_record("Computer", cid, {"warranty_expiry": _days(-1)})
        db_module.insert_record("Antivirus", {"name": "AV", "expiry_date": _days(-2)})
        db_module.insert_record("Other Software", {"name": "SW", "expiry_date": _days(-3)})
        results = db_module.get_expired()
        types = {r["type"] for r in results}
        assert {"Computer", "Antivirus", "Other Software"}.issubset(types)

    def test_custom_expiry_window(self, tmp_db):
        db_module.insert_record(
            "Antivirus",
            {"name": "Tight", "expiry_date": _days(5), "status": "Active"},
        )
        assert any(r["name"] == "Tight" for r in db_module.get_expiring_soon(days=7))
        assert not any(r["name"] == "Tight" for r in db_module.get_expiring_soon(days=3))
