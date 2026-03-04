"""
tests/test_main_window.py – MainWindow CRUD operations and tab behaviour.
"""

import sqlite3

from PyQt6.QtWidgets import QDialog, QMessageBox

import db as db_module
from app.dialogs import AssetDialog
from app.main_window import _MAIN_TABS, TABS

# ── Helpers ────────────────────────────────────────────────────────────────────


def _accept_with(monkeypatch, record: dict) -> None:
    """Patch AssetDialog so exec() returns Accepted and get_data() returns record."""
    monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Accepted)
    monkeypatch.setattr(AssetDialog, "get_data", lambda self: record)


def _tab_index(name: str) -> int:
    return _MAIN_TABS.index(name)


def _switch_to(win: object, name: str) -> None:
    """Switch win to the given tab."""
    win._tab_bar.setCurrentIndex(_tab_index(name))  # type: ignore[attr-defined]


# ══════════════════════════════════════════════════════════════════════════════
# Add
# ══════════════════════════════════════════════════════════════════════════════


class TestMainWindowAdd:
    def test_add_inserts_record_to_db(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        _accept_with(monkeypatch, {"type": "Computer", "name": "New PC", "status": "Active"})
        win._add_asset()
        rows = db_module.fetch_records("Computer")
        assert len(rows) == 1
        assert rows[0]["name"] == "New PC"

    def test_add_refreshes_model(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        _accept_with(monkeypatch, {"type": "Computer", "name": "PC", "status": "Active"})
        win._add_asset()
        assert win._models["Computer"].rowCount() == 1

    def test_add_on_all_tab_does_nothing(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("All"))
        called = []
        monkeypatch.setattr(
            AssetDialog, "exec", lambda self: called.append(1) or QDialog.DialogCode.Accepted
        )
        win._add_asset()
        assert called == []

    def test_add_shows_error_on_integrity_error(self, win, monkeypatch, qtbot):
        _switch_to(win, "Users")
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Accepted)
        monkeypatch.setattr(AssetDialog, "get_data", lambda self: {"type": "Users", "name": None})
        shown = []
        monkeypatch.setattr(
            QMessageBox, "critical", staticmethod(lambda *a, **kw: shown.append(a[2]))
        )
        win._add_asset()
        assert shown
        assert db_module.fetch_records("Users") == []

    def test_add_cancel_does_not_insert(self, win, monkeypatch):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        monkeypatch.setattr(AssetDialog, "exec", lambda self: QDialog.DialogCode.Rejected)
        win._add_asset()
        assert db_module.fetch_records("Computer") == []


# ══════════════════════════════════════════════════════════════════════════════
# Edit
# ══════════════════════════════════════════════════════════════════════════════


class TestMainWindowEdit:
    def _setup_computer(self) -> dict:
        db_module.insert_record("Computer", {"name": "Original", "status": "Active"})
        return db_module.fetch_records("Computer")[0]

    def test_edit_updates_record_in_db(self, win, monkeypatch):
        row = self._setup_computer()
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        win._selected_data = row
        _accept_with(monkeypatch, {"type": "Computer", "name": "Updated", "status": "Active"})
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
        _accept_with(monkeypatch, {"type": "Computer", "name": "X"})
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


# ══════════════════════════════════════════════════════════════════════════════
# Delete
# ══════════════════════════════════════════════════════════════════════════════


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
# Tab switching & Add button state
# ══════════════════════════════════════════════════════════════════════════════


class TestMainWindowTabs:
    def test_add_button_disabled_on_all_tab(self, win):
        win._tab_bar.setCurrentIndex(_tab_index("All"))
        assert not win._btn_add.isEnabled()

    def test_add_button_enabled_on_computer_tab(self, win):
        win._tab_bar.setCurrentIndex(_tab_index("Computer"))
        assert win._btn_add.isEnabled()

    def test_add_button_enabled_on_each_category_tab(self, win):
        for tab in TABS:
            _switch_to(win, tab)
            assert win._btn_add.isEnabled() == (tab != "All"), f"Wrong state for tab '{tab}'"

    def test_switching_tab_clears_selected_data(self, win):
        db_module.insert_record("Computer", {"name": "PC"})
        win._models["Computer"].refresh()
        _switch_to(win, "Computer")
        win._selected_data = {"id": 1, "name": "PC"}
        _switch_to(win, "Users")
        assert win._selected_data is None
