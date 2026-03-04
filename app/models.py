"""
app.models – QAbstractTableModel backed by the SQLite db.
"""

from datetime import date, timedelta

from PyQt6.QtCore import QAbstractTableModel, QModelIndex, Qt
from PyQt6.QtGui import QColor, QFont

import db as db_module

_EXPIRED_COLOR = QColor("#ffcdd2")  # light red
_EXPIRING_COLOR = QColor("#fff9c4")  # light yellow
_RETIRED_COLOR = QColor("#eeeeee")  # grey


class AssetTableModel(QAbstractTableModel):
    """Flat table model for one asset category (or 'All')."""

    def __init__(self, category: str = "All", search: str = "", parent=None) -> None:
        super().__init__(parent)
        self._category = category
        self._search = search
        self._columns = self._cols_for(category)
        self._rows: list[dict] = []
        self.refresh()

    @staticmethod
    def _cols_for(category: str) -> list[str]:
        if category == "All":
            return list(db_module.ALL_DISPLAY_COLS)
        cfg = db_module.TABLE_CONFIG.get(category)
        return list(cfg["display_cols"]) if cfg else list(db_module.ALL_DISPLAY_COLS)

    # ── Public API ─────────────────────────────────────────────────────────────
    def set_category(self, category: str) -> None:
        self._category = category
        self._columns = self._cols_for(category)
        self.refresh()

    def set_search(self, search: str) -> None:
        self._search = search.strip()
        self.refresh()

    def refresh(self) -> None:
        self.beginResetModel()
        self._rows = db_module.fetch_records(
            self._category,
            search=self._search or None,
        )
        self.endResetModel()

    def get_row_data(self, row: int) -> dict | None:
        if 0 <= row < len(self._rows):
            return self._rows[row]
        return None

    def column_keys(self) -> list[str]:
        return list(self._columns)

    # ── QAbstractTableModel interface ──────────────────────────────────────────
    def rowCount(self, parent: QModelIndex = QModelIndex()) -> int:
        return len(self._rows)

    def columnCount(self, parent: QModelIndex = QModelIndex()) -> int:
        return len(self._columns)

    def headerData(
        self, section: int, orientation: Qt.Orientation, role: int = Qt.ItemDataRole.DisplayRole
    ):
        if orientation == Qt.Orientation.Horizontal:
            if role == Qt.ItemDataRole.DisplayRole:
                return db_module.COLUMN_LABELS.get(self._columns[section], self._columns[section])
            if role == Qt.ItemDataRole.FontRole:
                f = QFont()
                f.setBold(True)
                return f
        else:
            if role == Qt.ItemDataRole.DisplayRole:
                return str(section + 1)
        return None

    def data(self, index: QModelIndex, role: int = Qt.ItemDataRole.DisplayRole):
        if not index.isValid():
            return None

        row_data = self._rows[index.row()]
        col_key = self._columns[index.column()]
        value = row_data.get(col_key)

        if role in (Qt.ItemDataRole.DisplayRole, Qt.ItemDataRole.EditRole):
            return str(value) if value is not None else ""

        if role == Qt.ItemDataRole.BackgroundRole:
            return self._row_bg(row_data)

        if role == Qt.ItemDataRole.ForegroundRole:
            if row_data.get("status") == "Retired":
                return QColor("#9e9e9e")
            if row_data.get("status") == "Missing":
                return QColor("#b71c1c")

        if role == Qt.ItemDataRole.ToolTipRole:
            tip_parts = [f"ID: {row_data.get('id')}"]
            if row_data.get("expiry_date"):
                tip_parts.append(f"Expires: {row_data['expiry_date']}")
            if row_data.get("warranty_expiry"):
                tip_parts.append(f"Warranty: {row_data['warranty_expiry']}")
            if row_data.get("notes"):
                tip_parts.append(f"Notes: {row_data['notes']}")
            return "  |  ".join(tip_parts)

        if role == Qt.ItemDataRole.UserRole:
            return row_data

        return None

    def flags(self, index: QModelIndex) -> Qt.ItemFlag:
        if not index.isValid():
            return Qt.ItemFlag.NoItemFlags
        base = Qt.ItemFlag.ItemIsEnabled | Qt.ItemFlag.ItemIsSelectable
        non_editable = ("id", "created_at", "updated_at", "user_name", "device_name", "type")
        if self._columns[index.column()] not in non_editable:
            base |= Qt.ItemFlag.ItemIsEditable
        return base

    def setData(self, index: QModelIndex, value, role: int = Qt.ItemDataRole.EditRole) -> bool:
        if not index.isValid() or role != Qt.ItemDataRole.EditRole:
            return False
        col_key = self._columns[index.column()]
        row_data = self._rows[index.row()]
        non_editable = ("id", "created_at", "updated_at", "user_name", "device_name", "type")
        if col_key in non_editable:
            return False
        # For "All" tab rows, derive the category from the 'type' field
        category = row_data.get("type") or self._category
        if category == "All":
            return False
        db_module.update_record(category, row_data["id"], {col_key: value.strip() or None})
        self.refresh()
        return True

    # ── Helpers ────────────────────────────────────────────────────────────────
    @staticmethod
    def _row_bg(row_data: dict):
        status = row_data.get("status", "")
        if status == "Retired":
            return _RETIRED_COLOR

        today = date.today().isoformat()
        threshold = (date.today() + timedelta(days=db_module.EXPIRY_WARNING_DAYS)).isoformat()
        expiry = row_data.get("expiry_date") or ""
        warranty = row_data.get("warranty_expiry") or ""

        if (expiry and expiry < today) or (warranty and warranty < today):
            return _EXPIRED_COLOR
        if (expiry and expiry <= threshold) or (warranty and warranty <= threshold):
            return _EXPIRING_COLOR
        return None
