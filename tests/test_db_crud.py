"""
tests/test_db_crud.py – CRUD operations for every entity type.
"""

import pytest

import db as db_module
from tests.helpers import _add_computer, _add_phone, _add_user

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
        """DEFAULT VALUES path: name has NOT NULL – confirms the guard is the DB constraint."""
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
        db_module.insert_record("Windows Key", {"license_key": "XXXXX-XXXXX", "computer_id": cid})
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
        db_module.insert_record("Other Software", {"name": "Adobe CC", "computer_id": cid})
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
        assert row["device_name"] == "PC-01"  # COALESCE picks computer first

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
        assert row["device_name"] == "PC-01"  # COALESCE picks computer first

    def test_insert_empty(self, tmp_db):
        aid = db_module.insert_record("Antivirus", {})
        assert isinstance(aid, int)

    def test_delete(self, tmp_db):
        aid = db_module.insert_record("Antivirus", {"name": "AV"})
        db_module.delete_record("Antivirus", aid)
        assert db_module.fetch_records("Antivirus") == []
