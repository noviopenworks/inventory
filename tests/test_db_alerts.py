"""
tests/test_db_alerts.py – Expiry / warranty alert queries.
"""

import db as db_module
from tests.helpers import _add_computer, _days


class TestExpiryAlerts:
    def test_expired_item_appears_in_get_expired(self, tmp_db):
        db_module.insert_record(
            "Antivirus", {"name": "OldAV", "expiry_date": _days(-10), "status": "Active"}
        )
        assert any(r["name"] == "OldAV" for r in db_module.get_expired())

    def test_future_item_not_in_expired(self, tmp_db):
        db_module.insert_record(
            "Antivirus", {"name": "FutureAV", "expiry_date": _days(60), "status": "Active"}
        )
        assert db_module.get_expired() == []

    def test_expiring_soon_default_window(self, tmp_db):
        db_module.insert_record(
            "Antivirus", {"name": "NearExpiry", "expiry_date": _days(15), "status": "Active"}
        )
        assert any(r["name"] == "NearExpiry" for r in db_module.get_expiring_soon())

    def test_not_expiring_soon_outside_window(self, tmp_db):
        db_module.insert_record(
            "Antivirus", {"name": "FarFuture", "expiry_date": _days(90), "status": "Active"}
        )
        assert not any(r["name"] == "FarFuture" for r in db_module.get_expiring_soon(days=30))

    def test_expiring_soon_excludes_retired(self, tmp_db):
        db_module.insert_record(
            "Antivirus", {"name": "RetiredAV", "expiry_date": _days(5), "status": "Retired"}
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
        assert any(r["type"] == "Computer" for r in db_module.get_expired())

    def test_warranty_expiring_soon_for_computer(self, tmp_db):
        cid = _add_computer()
        db_module.update_record("Computer", cid, {"warranty_expiry": _days(10)})
        assert any(r["type"] == "Computer" for r in db_module.get_expiring_soon(days=30))

    def test_multiple_categories_in_expired(self, tmp_db):
        cid = _add_computer()
        db_module.update_record("Computer", cid, {"warranty_expiry": _days(-1)})
        db_module.insert_record("Antivirus", {"name": "AV", "expiry_date": _days(-2)})
        db_module.insert_record("Other Software", {"name": "SW", "expiry_date": _days(-3)})
        types = {r["type"] for r in db_module.get_expired()}
        assert {"Computer", "Antivirus", "Other Software"}.issubset(types)

    def test_custom_expiry_window(self, tmp_db):
        db_module.insert_record(
            "Antivirus", {"name": "Tight", "expiry_date": _days(5), "status": "Active"}
        )
        assert any(r["name"] == "Tight" for r in db_module.get_expiring_soon(days=7))
        assert not any(r["name"] == "Tight" for r in db_module.get_expiring_soon(days=3))
