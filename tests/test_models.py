"""
tests/test_models.py – Tests for AssetTableModel (models.py).

All tests use an isolated temp DB via the shared tmp_db fixture.
"""

import pytest
from PyQt6.QtCore import QModelIndex, Qt

import db as db_module
from app.models import AssetTableModel

# ── Fixture ────────────────────────────────────────────────────────────────────


@pytest.fixture()
def tmp_db(tmp_path, monkeypatch):
    db_file = tmp_path / "test.db"
    monkeypatch.setattr(db_module, "DB_PATH", db_file)
    db_module.init_db()
    yield db_file


def _insert(category: str, data: dict) -> int:
    return db_module.insert_record(category, data)


# ── Construction & columns ─────────────────────────────────────────────────────


class TestModelColumns:
    def test_all_tab_uses_all_display_cols(self, tmp_db, qtbot):
        m = AssetTableModel(category="All")
        assert m.column_keys() == list(db_module.ALL_DISPLAY_COLS)

    def test_computer_tab_uses_correct_cols(self, tmp_db, qtbot):
        m = AssetTableModel(category="Computer")
        assert "user_name" in m.column_keys()
        assert "model" in m.column_keys()

    def test_windows_key_tab_has_device_name(self, tmp_db, qtbot):
        m = AssetTableModel(category="Windows Key")
        assert "device_name" in m.column_keys()
        assert "user_name" not in m.column_keys()

    def test_antivirus_has_device_name(self, tmp_db, qtbot):
        m = AssetTableModel(category="Antivirus")
        assert "device_name" in m.column_keys()

    def test_column_count_matches_display_cols(self, tmp_db, qtbot):
        for cat in db_module.CATEGORIES:
            m = AssetTableModel(category=cat)
            expected = len(db_module.TABLE_CONFIG[cat]["display_cols"])
            assert m.columnCount() == expected


# ── Row data ───────────────────────────────────────────────────────────────────


class TestModelRows:
    def test_empty_table_has_zero_rows(self, tmp_db, qtbot):
        m = AssetTableModel(category="Computer")
        assert m.rowCount() == 0

    def test_inserted_record_appears_in_model(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01", "status": "Active"})
        m = AssetTableModel(category="Computer")
        assert m.rowCount() == 1

    def test_get_row_data_returns_correct_dict(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-Alpha", "status": "Active"})
        m = AssetTableModel(category="Computer")
        row = m.get_row_data(0)
        assert row["name"] == "PC-Alpha"

    def test_get_row_data_out_of_range_returns_none(self, tmp_db, qtbot):
        m = AssetTableModel(category="Computer")
        assert m.get_row_data(99) is None

    def test_display_role_returns_string(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01", "status": "Active"})
        m = AssetTableModel(category="Computer")
        idx = m.index(0, m.column_keys().index("name"))
        assert m.data(idx, Qt.ItemDataRole.DisplayRole) == "PC-01"

    def test_none_value_displayed_as_empty_string(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01"})  # no model set
        m = AssetTableModel(category="Computer")
        idx = m.index(0, m.column_keys().index("model"))
        assert m.data(idx, Qt.ItemDataRole.DisplayRole) == ""

    def test_user_role_returns_dict(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01"})
        m = AssetTableModel(category="Computer")
        idx = m.index(0, 0)
        assert isinstance(m.data(idx, Qt.ItemDataRole.UserRole), dict)

    def test_invalid_index_returns_none(self, tmp_db, qtbot):
        m = AssetTableModel(category="Computer")
        assert m.data(QModelIndex()) is None

    def test_all_tab_shows_records_from_all_categories(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01"})
        _insert("Smartphone", {"name": "Phone-01"})
        _insert("Antivirus", {"name": "AV-01"})
        m = AssetTableModel(category="All")
        assert m.rowCount() == 2  # only hardware devices

    def test_all_tab_excludes_users(self, tmp_db, qtbot):
        _insert("Users", {"name": "Bob"})
        m = AssetTableModel(category="All")
        assert m.rowCount() == 0


# ── Header ─────────────────────────────────────────────────────────────────────


class TestModelHeader:
    def test_horizontal_header_returns_label(self, tmp_db, qtbot):
        m = AssetTableModel(category="Computer")
        name_col = m.column_keys().index("name")
        label = m.headerData(name_col, Qt.Orientation.Horizontal, Qt.ItemDataRole.DisplayRole)
        assert label == db_module.COLUMN_LABELS["name"]

    def test_vertical_header_returns_row_number(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01"})
        m = AssetTableModel(category="Computer")
        assert m.headerData(0, Qt.Orientation.Vertical, Qt.ItemDataRole.DisplayRole) == "1"


# ── Search ─────────────────────────────────────────────────────────────────────


class TestModelSearch:
    def test_set_search_filters_rows(self, tmp_db, qtbot):
        _insert("Computer", {"name": "Alpha"})
        _insert("Computer", {"name": "Beta"})
        m = AssetTableModel(category="Computer")
        m.set_search("Alpha")
        assert m.rowCount() == 1
        assert m.get_row_data(0)["name"] == "Alpha"

    def test_clear_search_restores_all_rows(self, tmp_db, qtbot):
        _insert("Computer", {"name": "Alpha"})
        _insert("Computer", {"name": "Beta"})
        m = AssetTableModel(category="Computer")
        m.set_search("Alpha")
        m.set_search("")
        assert m.rowCount() == 2


# ── Flags & editability ───────────────────────────────────────────────────────


class TestModelFlags:
    def test_id_column_not_editable(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01"})
        m = AssetTableModel(category="Computer")
        id_col = m.column_keys().index("id")
        idx = m.index(0, id_col)
        assert not (m.flags(idx) & Qt.ItemFlag.ItemIsEditable)

    def test_name_column_is_editable(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01"})
        m = AssetTableModel(category="Computer")
        name_col = m.column_keys().index("name")
        idx = m.index(0, name_col)
        assert m.flags(idx) & Qt.ItemFlag.ItemIsEditable

    def test_user_name_column_not_editable(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC-01"})
        m = AssetTableModel(category="Computer")
        un_col = m.column_keys().index("user_name")
        idx = m.index(0, un_col)
        assert not (m.flags(idx) & Qt.ItemFlag.ItemIsEditable)


# ── Background colour ──────────────────────────────────────────────────────────


class TestModelBackground:
    def _bg(self, tmp_db, qtbot, data: dict) -> object:
        _insert("Antivirus", data)
        m = AssetTableModel(category="Antivirus")
        idx = m.index(0, 0)
        return m.data(idx, Qt.ItemDataRole.BackgroundRole)

    def test_expired_row_has_red_background(self, tmp_db, qtbot):
        from datetime import date, timedelta

        past = (date.today() - timedelta(days=5)).isoformat()
        bg = self._bg(tmp_db, qtbot, {"name": "AV", "expiry_date": past, "status": "Active"})
        assert bg is not None
        assert bg.name().lower() == "#ffcdd2"

    def test_expiring_soon_row_has_yellow_background(self, tmp_db, qtbot):
        from datetime import date, timedelta

        soon = (date.today() + timedelta(days=5)).isoformat()
        bg = self._bg(tmp_db, qtbot, {"name": "AV", "expiry_date": soon, "status": "Active"})
        assert bg is not None
        assert bg.name().lower() == "#fff9c4"

    def test_retired_row_has_grey_background(self, tmp_db, qtbot):
        bg = self._bg(tmp_db, qtbot, {"name": "AV", "status": "Retired"})
        assert bg is not None
        assert bg.name().lower() == "#eeeeee"

    def test_active_valid_row_has_no_background(self, tmp_db, qtbot):
        from datetime import date, timedelta

        future = (date.today() + timedelta(days=365)).isoformat()
        bg = self._bg(tmp_db, qtbot, {"name": "AV", "expiry_date": future, "status": "Active"})
        assert bg is None


# ── setData (inline edit) ─────────────────────────────────────────────────────


class TestModelSetData:
    def test_setdata_updates_db_and_refreshes(self, tmp_db, qtbot):
        _insert("Computer", {"name": "Old Name"})
        m = AssetTableModel(category="Computer")
        name_col = m.column_keys().index("name")
        idx = m.index(0, name_col)
        result = m.setData(idx, "New Name")
        assert result is True
        assert m.get_row_data(0)["name"] == "New Name"

    def test_setdata_rejects_id_column(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC"})
        m = AssetTableModel(category="Computer")
        id_col = m.column_keys().index("id")
        idx = m.index(0, id_col)
        assert m.setData(idx, "999") is False

    def test_setdata_on_all_tab_is_rejected(self, tmp_db, qtbot):
        _insert("Computer", {"name": "PC"})
        m = AssetTableModel(category="All")
        idx = m.index(0, 0)
        # All-tab type column is non-editable; any editable col should bail on type=="All"
        assert m.setData(idx, "x") is False
