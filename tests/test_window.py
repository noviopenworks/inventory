"""
tests/test_window.py – Tests for AssetDialog (dialogs.py) and
                       MainWindow CRUD operations (main_window.py).

Uses pytest-qt's `qtbot` fixture for a QApplication context.
Dialog exec() is monkeypatched so no real modal loops open.
"""

import sqlite3

import pytest
from PyQt6.QtWidgets import QDialog, QMessageBox

import db as db_module
from app.dialogs import AssetDialog
from app.main_window import MainWindow

# ── Fixtures ───────────────────────────────────────────────────────────────────


@pytest.fixture()
def tmp_db(tmp_path, monkeypatch):
    db_file = tmp_path / "test.db"
    monkeypatch.setattr(db_module, "DB_PATH", db_file)
    db_module.init_db()
    yield db_file


@pytest.fixture()
def win(tmp_db, qtbot):
    """A fully constructed MainWindow backed by a temp DB."""
    w = MainWindow()
    qtbot.addWidget(w)
    return w


def _tab_index(name: str) -> int:
    from app.main_window import _MAIN_TABS

    return _MAIN_TABS.index(name)


def _switch_to(win: object, name: str) -> None:
    """Switch *win* to the given tab, handling the separate Users button."""
    if name == "Users":
        win._users_btn.click()  # type: ignore[attr-defined]
    else:
        win._tab_bar.setCurrentIndex(_tab_index(name))  # type: ignore[attr-defined]


# ══════════════════════════════════════════════════════════════════════════════
# AssetDialog – get_data() and _populate()
# ══════════════════════════════════════════════════════════════════════════════


class TestAssetDialogGetData:
    def test_locked_add_returns_fixed_category(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        data = dlg.get_data()
        assert data["type"] == "Computer"

    def test_unlocked_add_returns_combo_category(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=False)
        qtbot.addWidget(dlg)
        # Switch the category combo to Antivirus
        dlg._cat_combo.setCurrentText("Antivirus")
        data = dlg.get_data()
        assert data["type"] == "Antivirus"

    def test_name_field_value_returned(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["name"].setText("My PC")
        data = dlg.get_data()
        assert data["name"] == "My PC"

    def test_empty_name_returns_none(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["name"].setText("")
        data = dlg.get_data()
        assert data["name"] is None

    def test_status_field_returns_combo_text(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["status"].setCurrentText("Retired")
        data = dlg.get_data()
        assert data["status"] == "Retired"

    def test_notes_field_returned_as_plaintext(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["notes"].setPlainText("some note")
        data = dlg.get_data()
        assert data["notes"] == "some note"

    def test_empty_notes_returns_none(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["notes"].setPlainText("")
        data = dlg.get_data()
        assert data["notes"] is None

    def test_device_resolved_from_map(self, tmp_db, qtbot):
        cid = db_module.insert_record("Computer", {"name": "PC-01", "status": "Active"})
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        # _device_map is populated in _build_device_row
        assert "PC-01" in dlg._device_map
        dlg._fields["_device"].setCurrentText("PC-01")
        data = dlg.get_data()
        assert data["computer_id"] == cid

    def test_device_blank_returns_none(self, tmp_db, qtbot):
        dlg = AssetDialog(category="Windows Key", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["_device"].setCurrentText("")
        data = dlg.get_data()
        assert data["computer_id"] is None

    def test_user_id_resolved_from_map(self, tmp_db, qtbot):
        uid = db_module.insert_record("Users", {"name": "Alice", "status": "Active"})
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        assert "Alice" in dlg._user_map
        dlg._fields["user_id"].setCurrentText("Alice")
        data = dlg.get_data()
        assert data["user_id"] == uid

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


class TestAssetDialogDateValidation:
    """Verify _validate_and_accept blocks bad dates."""

    def test_valid_date_accepted(self, tmp_db, qtbot, monkeypatch):
        dlg = AssetDialog(category="Computer", locked=True)
        qtbot.addWidget(dlg)
        dlg._fields["purchase_date"].setText("2025-01-15")
        # Patch accept so we can detect it was called
        accepted = []
        monkeypatch.setattr(QDialog, "accept", lambda self: accepted.append(True))
        dlg._validate_and_accept()
        assert accepted  # accept() was called

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
        assert warned  # warning was shown

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


# ══════════════════════════════════════════════════════════════════════════════
# MainWindow CRUD
# ══════════════════════════════════════════════════════════════════════════════


class TestMainWindowAdd:
    def _accept_with(self, monkeypatch, record: dict):
        """Patch AssetDialog so exec() returns Accepted and get_data() returns record."""
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Accepted)
        monkeypatch.setattr(AssetDialog, "get_data", lambda self: record)

    def test_add_inserts_record_to_db(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        self._accept_with(monkeypatch, {"type": "Computer", "name": "New PC", "status": "Active"})
        win._add_asset()
        rows = db_module.fetch_records("Computer")
        assert len(rows) == 1
        assert rows[0]["name"] == "New PC"

    def test_add_refreshes_model(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        self._accept_with(monkeypatch, {"type": "Computer", "name": "PC", "status": "Active"})
        win._add_asset()
        assert win._models["Computer"].rowCount() == 1

    def test_add_on_all_tab_does_nothing(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("All"))
        called = []
        monkeypatch.setattr(
            AssetDialog, "exec", lambda self: called.append(1) or QDialog.DialogCode.Accepted
        )
        win._add_asset()
        assert called == []  # dialog never opened

    def test_add_shows_error_on_integrity_error(self, win, monkeypatch, qtbot):
        _switch_to(win, "Users")
        # Return empty name which will trigger NOT NULL constraint
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Accepted)
        monkeypatch.setattr(AssetDialog, "get_data", lambda self: {"type": "Users", "name": None})
        shown = []
        monkeypatch.setattr(
            QMessageBox, "critical", staticmethod(lambda *a, **kw: shown.append(a[2]))
        )
        win._add_asset()
        assert shown  # error dialog was shown
        assert db_module.fetch_records("Users") == []  # nothing written

    def test_add_cancel_does_not_insert(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Rejected)
        win._add_asset()
        assert db_module.fetch_records("Computer") == []


class TestMainWindowEdit:
    def _setup_computer(self) -> dict:
        db_module.insert_record("Computer", {"name": "Original", "status": "Active"})
        rows = db_module.fetch_records("Computer")
        return rows[0]

    def _accept_with(self, monkeypatch, record: dict):
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Accepted)
        monkeypatch.setattr(AssetDialog, "get_data", lambda self: record)

    def test_edit_updates_record_in_db(self, win, monkeypatch):
        row = self._setup_computer()
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        win._selected_data = row
        self._accept_with(monkeypatch, {"type": "Computer", "name": "Updated", "status": "Active"})
        win._edit_asset()
        assert db_module.fetch_records("Computer")[0]["name"] == "Updated"

    def test_edit_no_selection_shows_message(self, win, monkeypatch):
        win._selected_data = None
        shown = []
        monkeypatch.setattr(
            QMessageBox, "information", staticmethod(lambda *a, **kw: shown.append(a[2]))
        )
        win._edit_asset()
        assert shown

    def test_edit_cancel_does_not_change_db(self, win, monkeypatch):
        row = self._setup_computer()
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        win._selected_data = row
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Rejected)
        win._edit_asset()
        assert db_module.fetch_records("Computer")[0]["name"] == "Original"

    def test_edit_shows_error_on_db_exception(self, win, monkeypatch):
        row = self._setup_computer()
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        win._selected_data = row
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Accepted)
        monkeypatch.setattr(AssetDialog, "get_data", lambda self: {"type": "Computer", "name": "X"})
        monkeypatch.setattr(
            db_module,
            "update_record",
            lambda *a, **kw: (_ for _ in ()).throw(sqlite3.OperationalError("forced error")),
        )
        shown = []
        monkeypatch.setattr(
            QMessageBox, "critical", staticmethod(lambda *a, **kw: shown.append(a[2]))
        )
        win._edit_asset()
        assert shown


class TestMainWindowDelete:
    def _setup_computer(self) -> dict:
        db_module.insert_record("Computer", {"name": "To Delete", "status": "Active"})
        return db_module.fetch_records("Computer")[0]

    def test_delete_removes_record(self, win, monkeypatch):
        row = self._setup_computer()
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        win._selected_data = row
        monkeypatch.setattr(
            QMessageBox, "question", staticmethod(lambda *a, **kw: QMessageBox.StandardButton.Yes)
        )
        win._delete_asset()
        assert db_module.fetch_records("Computer") == []

    def test_delete_cancel_keeps_record(self, win, monkeypatch):
        row = self._setup_computer()
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        win._selected_data = row
        monkeypatch.setattr(
            QMessageBox, "question", staticmethod(lambda *a, **kw: QMessageBox.StandardButton.No)
        )
        win._delete_asset()
        assert len(db_module.fetch_records("Computer")) == 1

    def test_delete_no_selection_shows_message(self, win, monkeypatch):
        win._selected_data = None
        shown = []
        monkeypatch.setattr(
            QMessageBox, "information", staticmethod(lambda *a, **kw: shown.append(a[2]))
        )
        win._delete_asset()
        assert shown

    def test_delete_refreshes_model(self, win, monkeypatch):
        row = self._setup_computer()
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        win._selected_data = row
        monkeypatch.setattr(
            QMessageBox, "question", staticmethod(lambda *a, **kw: QMessageBox.StandardButton.Yes)
        )
        win._delete_asset()
        assert win._models["Computer"].rowCount() == 0


# ══════════════════════════════════════════════════════════════════════════════
# MainWindow – tab switching & Add button state
# ══════════════════════════════════════════════════════════════════════════════


class TestMainWindowTabs:
    def test_add_button_disabled_on_all_tab(self, win):
        win._tab_bar.setCurrentIndex(_tab_index("All"))
        assert not win._btn_add.isEnabled()

    def test_add_button_enabled_on_computer_tab(self, win):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        assert win._btn_add.isEnabled()

    def test_add_button_enabled_on_each_category_tab(self, win):
        from app.main_window import TABS

        for tab in TABS:
            _switch_to(win, tab)
            expected = tab != "All"
            assert win._btn_add.isEnabled() == expected, f"Wrong state for tab '{tab}'"

    def test_switching_tab_clears_selected_data(self, win):
        db_module.insert_record("Computer", {"name": "PC"})
        win._models["Computer"].refresh()
        _switch_to(win, "Computer")
        win._selected_data = {"id": 1, "name": "PC"}
        _switch_to(win, "Users")
        assert win._selected_data is None
