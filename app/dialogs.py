"""
app.dialogs – Add / Edit record dialog (relational schema edition).
"""

import re

from PyQt6.QtCore import Qt, QTimer
from PyQt6.QtWidgets import (
    QComboBox,
    QDialog,
    QDialogButtonBox,
    QFormLayout,
    QFrame,
    QHBoxLayout,
    QLabel,
    QLineEdit,
    QMessageBox,
    QTextEdit,
    QVBoxLayout,
)

import db as db_module


def _close_popup(combo: QComboBox) -> None:
    """Schedule hidePopup after the current event cycle (fixes Linux/Wayland)."""
    combo.activated.connect(lambda: QTimer.singleShot(0, combo.hidePopup))


_DATE_HINT = "YYYY-MM-DD"
_DATE_COLS = {"purchase_date", "warranty_expiry", "expiry_date"}
_DATE_RE = re.compile(r"^\d{4}-\d{2}-\d{2}$")
_NA = "N/A"
# Columns that get an editable combobox with an N/A option
_NA_COMBO_COLS = {"license_key", "expiry_date"}
# columns that are auto-managed and never shown in the form
_SKIP_COLS = {"id", "created_at", "updated_at"}
# FK columns replaced by a single "Device" dropdown in the form
_DEVICE_FK_COLS = {"computer_id", "smartphone_id", "tablet_id"}


class AssetDialog(QDialog):
    """Modal form to create or edit a single record."""

    def __init__(
        self,
        category: str = "Computer",
        data: dict | None = None,
        locked: bool = False,
        parent=None,
    ) -> None:
        super().__init__(parent)
        self._editing = data is not None
        self._locked = locked  # suppress category selector when adding to a specific tab
        self._data = data or {}
        self._category: str = data.get("type", category) if data else category
        # FK lookup maps: display name -> id  (populated in _build_form_rows)
        self._user_map: dict[str, int] = {}
        self._device_map: dict[str, tuple[str, int]] = {}  # label -> (fk_col, fk_id)
        self._fields: dict[str, QComboBox | QTextEdit | QLineEdit] = {}

        self.setWindowTitle("Edit" if self._editing else "Add")
        self.setMinimumWidth(520)
        self.setModal(True)
        self._build_ui()
        if self._editing:
            self._populate()

    # ── UI ────────────────────────────────────────────────────────────────────
    def _build_ui(self) -> None:
        root = QVBoxLayout(self)
        root.setSpacing(8)

        if not self._editing and not self._locked:
            cat_row = QHBoxLayout()
            cat_row.addWidget(QLabel("<b>Category:</b>"))
            self._cat_combo = QComboBox()
            self._cat_combo.addItems(db_module.CATEGORIES)
            self._cat_combo.setCurrentText(self._category)
            _close_popup(self._cat_combo)
            self._cat_combo.currentTextChanged.connect(self._rebuild_form)
            cat_row.addWidget(self._cat_combo)
            cat_row.addStretch()
            root.addLayout(cat_row)
            self._add_separator(root)
        else:
            root.addWidget(QLabel(f"<b>Category:</b>  {self._category}"))

        self._form_widget = QFrame()
        self._form_layout = QFormLayout(self._form_widget)
        self._form_layout.setRowWrapPolicy(QFormLayout.RowWrapPolicy.DontWrapRows)
        self._form_layout.setLabelAlignment(Qt.AlignmentFlag.AlignRight)
        self._form_layout.setSpacing(6)
        root.addWidget(self._form_widget)
        self._build_form_rows()

        self._add_separator(root)
        buttons = QDialogButtonBox(
            QDialogButtonBox.StandardButton.Ok | QDialogButtonBox.StandardButton.Cancel
        )
        buttons.accepted.connect(self._validate_and_accept)
        buttons.rejected.connect(self.reject)
        root.addWidget(buttons)

    @staticmethod
    def _add_separator(layout: QVBoxLayout) -> None:
        sep = QFrame()
        sep.setFrameShape(QFrame.Shape.HLine)
        sep.setFrameShadow(QFrame.Shadow.Sunken)
        layout.addWidget(sep)

    def _clear_form(self) -> None:
        self._fields.clear()
        self._user_map.clear()
        self._device_map.clear()
        while self._form_layout.rowCount():
            self._form_layout.removeRow(0)

    def _build_form_rows(self) -> None:
        spec = db_module.TABLE_CONFIG.get(self._category)
        db_cols: list[str] = list(spec["db_cols"]) if spec else []

        device_built = False
        for col in db_cols:
            if col in _SKIP_COLS:
                continue
            if col in _DEVICE_FK_COLS:
                if not device_built:
                    self._build_device_row(
                        "computer_id" in db_cols,
                        "smartphone_id" in db_cols,
                        "tablet_id" in db_cols,
                    )
                    device_built = True
                continue

            label = db_module.COLUMN_LABELS.get(col, col)
            widget: QComboBox | QTextEdit | QLineEdit

            if col == "user_id":
                widget = QComboBox()
                _close_popup(widget)
                widget.addItem("")  # blank / unassigned
                for u in db_module.fetch_users():
                    widget.addItem(u["name"])
                    self._user_map[u["name"]] = u["id"]
            elif col == "status":
                widget = QComboBox()
                _close_popup(widget)
                opts = (
                    db_module.USER_STATUS_OPTIONS
                    if self._category == "Users"
                    else db_module.STATUS_OPTIONS
                )
                widget.addItems(opts)
            elif col == "notes":
                widget = QTextEdit()
                widget.setFixedHeight(72)
                widget.setAcceptRichText(False)
            elif col in _NA_COMBO_COLS:
                widget = QComboBox()
                widget.setEditable(True)
                widget.setInsertPolicy(QComboBox.InsertPolicy.NoInsert)
                widget.setCompleter(None)
                _close_popup(widget)
                widget.addItems(["", _NA])
                le = widget.lineEdit()
                if le is not None and col == "license_key":
                    le.setFont(self._mono_font())
                    le.setPlaceholderText("XXXXX-XXXXX-XXXXX-XXXXX-XXXXX")
                elif le is not None and col in _DATE_COLS:
                    le.setPlaceholderText(_DATE_HINT)
            elif col in _DATE_COLS:
                widget = QLineEdit()
                widget.setPlaceholderText(_DATE_HINT)
            else:
                widget = QLineEdit()

            self._form_layout.addRow(f"{label}:", widget)
            self._fields[col] = widget

    def _build_device_row(self, computers: bool, phones: bool, tablets: bool) -> None:
        """Add a single 'Device' dropdown combining computers, phones, and/or tablets."""
        widget = QComboBox()
        _close_popup(widget)
        widget.addItem("")
        kind_count = sum([computers, phones, tablets])
        mixed = kind_count > 1
        if computers:
            for c in db_module.fetch_computers():
                lbl = f"{c['name']} (Computer)" if mixed else c["name"]
                widget.addItem(lbl)
                self._device_map[lbl] = ("computer_id", c["id"])
        if phones:
            for s in db_module.fetch_smartphones():
                lbl = f"{s['name']} (Smartphone)" if mixed else s["name"]
                widget.addItem(lbl)
                self._device_map[lbl] = ("smartphone_id", s["id"])
        if tablets:
            for t in db_module.fetch_tablets():
                lbl = f"{t['name']} (Tablet)" if mixed else t["name"]
                widget.addItem(lbl)
                self._device_map[lbl] = ("tablet_id", t["id"])
        self._form_layout.addRow("Device:", widget)
        self._fields["_device"] = widget

    def _rebuild_form(self, category: str) -> None:
        self._category = category
        self._clear_form()
        self._build_form_rows()
        self.adjustSize()

    def _populate(self) -> None:
        for col, widget in self._fields.items():
            if col == "_device":
                assert isinstance(widget, QComboBox)
                comp_id = self._data.get("computer_id")
                phone_id = self._data.get("smartphone_id")
                tablet_id = self._data.get("tablet_id")
                target = ""
                for lbl, (fk_col, fk_id) in self._device_map.items():
                    if fk_col == "computer_id" and fk_id == comp_id:
                        target = lbl
                        break
                    if fk_col == "smartphone_id" and fk_id == phone_id:
                        target = lbl
                        break
                    if fk_col == "tablet_id" and fk_id == tablet_id:
                        target = lbl
                        break
                idx = widget.findText(target)
                widget.setCurrentIndex(max(idx, 0))
            elif col == "user_id":
                assert isinstance(widget, QComboBox)
                user_name = str(self._data.get("user_name") or "")
                idx = widget.findText(user_name)
                widget.setCurrentIndex(max(idx, 0))
            elif isinstance(widget, QComboBox):
                value = str(self._data.get(col) or "")
                idx = widget.findText(value)
                widget.setCurrentIndex(max(idx, 0))
            elif isinstance(widget, QTextEdit):
                widget.setPlainText(str(self._data.get(col) or ""))
            else:
                widget.setText(str(self._data.get(col) or ""))

    # ── Validation ─────────────────────────────────────────────────────────────
    def _validate_and_accept(self) -> None:
        """Check date fields before accepting the dialog."""
        errors: list[str] = []
        for col, widget in self._fields.items():
            if col not in _DATE_COLS:
                continue
            # Support both QLineEdit and editable QComboBox
            if isinstance(widget, QComboBox):
                text = widget.currentText().strip()
            elif isinstance(widget, QLineEdit):
                text = widget.text().strip()
            else:
                continue
            if not text or text == _NA:
                continue  # blank and N/A are allowed
            if not _DATE_RE.match(text):
                label = db_module.COLUMN_LABELS.get(col, col)
                errors.append(f"{label}: must be YYYY-MM-DD (got '{text}')")
        if errors:
            QMessageBox.warning(self, "Invalid date(s)", "\n".join(errors))
            return
        self.accept()

    # ── Data extraction ────────────────────────────────────────────────────────
    def get_data(self) -> dict:
        category = (
            self._category if (self._editing or self._locked) else self._cat_combo.currentText()
        )
        result: dict = {"type": category}
        spec = db_module.TABLE_CONFIG.get(category)
        db_cols = spec["db_cols"] if spec else []

        for col, widget in self._fields.items():
            if col == "_device":
                assert isinstance(widget, QComboBox)
                text = widget.currentText()
                entry = self._device_map.get(text)
                if entry:
                    fk_col, fk_id = entry
                    result[fk_col] = fk_id
                    # Clear the other FK if present
                    if fk_col != "computer_id" and "computer_id" in db_cols:
                        result["computer_id"] = None
                    if fk_col != "smartphone_id" and "smartphone_id" in db_cols:
                        result["smartphone_id"] = None
                    if fk_col != "tablet_id" and "tablet_id" in db_cols:
                        result["tablet_id"] = None
                else:
                    # Blank / unset — clear all device FKs
                    if "computer_id" in db_cols:
                        result["computer_id"] = None
                    if "smartphone_id" in db_cols:
                        result["smartphone_id"] = None
                    if "tablet_id" in db_cols:
                        result["tablet_id"] = None
            elif col == "user_id":
                assert isinstance(widget, QComboBox)
                n = widget.currentText()
                result["user_id"] = self._user_map.get(n) if n else None
            elif col in _NA_COMBO_COLS:
                assert isinstance(widget, QComboBox)
                val = widget.currentText().strip()
                result[col] = val or None
            elif isinstance(widget, QComboBox):
                result[col] = widget.currentText()
            elif isinstance(widget, QTextEdit):
                result[col] = widget.toPlainText().strip() or None
            else:
                result[col] = widget.text().strip() or None
        return result

    # ── Helpers ────────────────────────────────────────────────────────────────
    @staticmethod
    def _mono_font():
        from PyQt6.QtGui import QFont

        f = QFont("Consolas, Courier New, monospace")
        f.setPointSize(10)
        return f
