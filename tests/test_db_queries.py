"""
tests/test_db_queries.py – fetch_records() search and the All-tab UNION query.
"""

import db as db_module
from tests.helpers import _add_computer, _add_phone, _add_user

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
        assert {r["type"] for r in rows} == {"Computer", "Smartphone", "Tablet"}

    def test_software_excluded_from_all(self, tmp_db):
        db_module.insert_record("Windows Key", {"license_key": "WK"})
        db_module.insert_record("Antivirus", {"name": "AV"})
        db_module.insert_record("Other Software", {"name": "SW"})
        assert db_module.fetch_records("All") == []

    def test_users_excluded_from_all(self, tmp_db):
        _add_user()
        assert all(r["type"] != "Users" for r in db_module.fetch_records("All"))

    def test_empty_db_returns_empty_list(self, tmp_db):
        assert db_module.fetch_records("All") == []

    def test_search_across_hardware_categories(self, tmp_db):
        _add_computer("UniqueComputer")
        _add_phone("UniquePhone")
        db_module.insert_record("Antivirus", {"name": "UniqueAV"})
        results = db_module.fetch_records("All", search="Unique")
        assert len(results) == 2
        assert {r["type"] for r in results} == {"Computer", "Smartphone"}

    def test_device_rows_show_user_name(self, tmp_db):
        uid = _add_user("Bob")
        _add_computer(user_id=uid)
        rows = db_module.fetch_records("All")
        pc = next(r for r in rows if r["type"] == "Computer")
        assert pc["user_name"] == "Bob"

    def test_device_rows_show_model(self, tmp_db):
        _add_computer()
        rows = db_module.fetch_records("All")
        pc = next(r for r in rows if r["type"] == "Computer")
        assert pc["model"] == "ThinkPad X1"
