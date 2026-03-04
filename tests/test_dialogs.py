"""
tests/test_dialogs.py – Tests for AssetDialog: get_data(), populate, and date validation.
"""

from PyQt6.QtWidgets import QDialog, QMessageBox

import db as db_module
from app.dialogs import AssetDialog

# ══════════════════════════════════════════════════════════════════════════════
# get_data()
# ══════════════════════════════════════════════════════════════════════════════


class TestAssetDialogGetData:
    def test_locked_add_returns_fixed_category(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        assert dlg.get_data()["type"] == "Computer"

    def test_unlocked_add_returns_combo_category(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=False)
        qtbot.addWidget(dlg)
        dlg._cat_combo.setCurrentText("Antivirus")
        assert dlg.get_data()["type"] == "Antivirus"

    def test_name_field_value_returned(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["name"].setText("My PC")
        assert dlg.get_data()["name"] == "My PC"

    def test_empty_name_returns_none(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["name"].setText("")
        assert dlg.get_data()["name"] is None

    def test_status_field_returns_combo_text(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["status"].setCurrentText("Retired")
        assert dlg.get_data()["status"] == "Retired"

    def test_notes_field_returned_as_plaintext(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["notes"].setPlainText("some note")
        assert dlg.get_data()["notes"] == "some note"

    def test_empty_notes_returns_none(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["notes"].setPlainText("")
        assert dlg.get_data()["notes"] is None

    def test_device_resolved_from_map(self, tmp_db, qtbot):
        cid = db_module.insert_record("Computer", {"name": "PC-01", "status": "Active"})
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        assert "PC-01" in dlg._device_map
        dlg._fields["_device"].setCurrentText("PC-01")
        assert dlg.get_data()["computer_id"] == cid

    def test_device_blank_returns_none(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["_device"].setCurrentText("")
        assert dlg.get_data()["computer_id"] is None

    def test_user_id_resolved_from_map(self, tmp_db, qtbot):
        uid = db_module.insert_record("Users", {"name": "Alice", "status": "Active"})
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        assert "Alice" in dlg._user_map
        dlg._fields["user_id"].setCurrentText("Alice")
        assert dlg.get_data()["user_id"] == uid

    def test_device_smartphone_resolved(self, tmp_db, qtbot):
        sid = db_module.insert_record("Smartphone", {"name": "Phone-01", "status": "Active"})
        dlg = AssetDialog(category="Antivirus", locked=True)
        qtbot.addWidget(dlg)
        label = "Phone-01 (Smartphone)"
        assert label in dlg._device_map
        dlg._fields["_device"].setCurrentText(label)
        data = dlg.get_data()
        assert data["smartphone_id"] == sid
        assert data["computer_id"] is None

    def test_all_db_cols_present_in_result(self, tmp_db, qtbot):
        for cat in db_module.CATEGORIES:
            dlg = AssetDialog(category=cat, locked=True)
            qtbot.addWidget(dlg)
            data = dlg.get_data()
            for col in db_module.TABLE_CONFIG[cat]["db_cols"]:
                assert col in data, f"Missing {col} in {cat} get_data() result"


# ══════════════════════════════════════════════════════════════════════════════
# Date validation
# ══════════════════════════════════════════════════════════════════════════════


class TestAssetDialogDateValidation:
    def test_valid_date_accepted(self, tmp_db, qtbot, monkeypatch):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["purchase_date"].setText("2025-01-15")
        accepted = []
        monkeypatch.setattr(QDialog, "accept", lambda self: accepted.append(True))
        dlg._validate_and_accept()
        assert accepted

    def test_empty_date_accepted(self, tmp_db, qtbot, monkeypatch):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["purchase_date"].setText("")
        accepted = []
        monkeypatch.setattr(QDialog, "accept", lambda self: accepted.append(True))
        dlg._validate_and_accept()
        assert accepted

    def test_bad_date_rejected(self, tmp_db, qtbot, monkeypatch):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["purchase_date"].setText("not-a-date")
        warned = []
        monkeypatch.setattr(
            QMessageBox, "warning", staticmethod(lambda *a, **kw: warned.append(a[2]))
        )
        dlg._validate_and_accept()
        assert warned

    def test_partial_date_rejected(self, tmp_db, qtbot, monkeypatch):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["purchase_date"].setText("2025-1-5")
        warned = []
        monkeypatch.setattr(
            QMessageBox, "warning", staticmethod(lambda *a, **kw: warned.append(a[2]))
        )
        dlg._validate_and_accept()
        assert warned

    def test_bad_warranty_expiry_rejected(self, tmp_db, qtbot, monkeypatch):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["warranty_expiry"].setText("31/12/2025")
        warned = []
        monkeypatch.setattr(
            QMessageBox, "warning", staticmethod(lambda *a, **kw: warned.append(a[2]))
        )
        dlg._validate_and_accept()
        assert warned


# ══════════════════════════════════════════════════════════════════════════════
# _populate() (edit mode)
# ══════════════════════════════════════════════════════════════════════════════


class TestAssetDialogPopulate:
    def test_populate_fills_name_field(self, tmp_db, qtbot):
        row = {
            "type": "Computer",
            "id": 1,
            "name": "Old PC",
            "status": "Active",
            "model": None,
            "user_id": None,
            "user_name": None,
            "purchase_date": None,
            "warranty_expiry": None,
            "notes": None,
        }
        dlg = AssetDialog(category="Computer", data=row)
        qtbot.addWidget(dlg)
        assert dlg._fields["name"].text() == "Old PC"

    def test_populate_selects_correct_status(self, tmp_db, qtbot):
        row = {
            "type": "Computer",
            "id": 1,
            "name": "PC",
            "status": "Retired",
            "model": None,
            "user_id": None,
            "user_name": None,
            "purchase_date": None,
            "warranty_expiry": None,
            "notes": None,
        }
        dlg = AssetDialog(category="Computer", data=row)
        qtbot.addWidget(dlg)
        assert dlg._fields["status"].currentText() == "Retired"

    def test_populate_selects_user_by_name(self, tmp_db, qtbot):
        uid = db_module.insert_record("Users", {"name": "Alice", "status": "Active"})
        row = {
            "type": "Computer",
            "id": 1,
            "name": "PC",
            "status": "Active",
            "model": None,
            "user_id": uid,
            "user_name": "Alice",
            "purchase_date": None,
            "warranty_expiry": None,
            "notes": None,
        }
        dlg = AssetDialog(category="Computer", data=row)
        qtbot.addWidget(dlg)
        assert dlg._fields["user_id"].currentText() == "Alice"

    def test_populate_selects_device(self, tmp_db, qtbot):
        cid = db_module.insert_record("Computer", {"name": "PC-01", "status": "Active"})
        row = {
            "type": "Windows Key",
            "id": 1,
            "license_key": "X",
            "status": "Active",
            "computer_id": cid,
            "device_name": "PC-01",
            "notes": None,
        }
        dlg = AssetDialog(category="Windows Key", data=row)
        qtbot.addWidget(dlg)
        assert dlg._fields["_device"].currentText() == "PC-01"
