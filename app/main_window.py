"""
app.main_window – Main application window.
"""

import csv
from datetime import date
from pathlib import Path

from PyQt6.QtCore import Qt
from PyQt6.QtGui import QAction, QKeySequence, QShortcut
from PyQt6.QtWidgets import (
    QAbstractItemView,
    QDialog,
    QFileDialog,
    QFrame,
    QHBoxLayout,
    QHeaderView,
    QLabel,
    QLineEdit,
    QMainWindow,
    QMessageBox,
    QPushButton,
    QStackedWidget,
    QStatusBar,
    QTabBar,
    QTableView,
    QVBoxLayout,
    QWidget,
)

import db as db_module
from app.dialogs import AssetDialog
from app.models import AssetTableModel
from app.style import LEGEND_ITEMS, STYLESHEET

_MAIN_TABS: list[str] = ["All"] + [c for c in db_module.CATEGORIES if c != "Users"]
TABS: list[str] = _MAIN_TABS + ["Users"]


class MainWindow(QMainWindow):
    def __init__(self) -> None:
        super().__init__()
        from app._version import __version__
        self.setWindowTitle(f"IT Asset Inventory v{__version__}")
        self.setMinimumSize(1280, 720)

        self._models: dict[str, AssetTableModel] = {}
        self._tables: dict[str, QTableView] = {}
        self._current_tab: str = "All"
        self._selected_data: dict | None = None  # last selected row, survives focus loss

        self._build_ui()
        self._apply_style()
        self._check_alerts()

    # ══════════════════════════════════════════════════════════════════════════
    # UI construction
    # ══════════════════════════════════════════════════════════════════════════

    def _build_ui(self) -> None:
        central = QWidget()
        self.setCentralWidget(central)
        root = QVBoxLayout(central)
        root.setContentsMargins(8, 8, 8, 4)
        root.setSpacing(4)

        root.addLayout(self._build_toolbar())
        root.addWidget(self._build_legend())
        root.addWidget(self._build_tabs())
        self._btn_add.setEnabled(self._current_tab != "All")

        _status_bar = self.statusBar()
        assert _status_bar is not None
        self._status_bar: QStatusBar = _status_bar
        self._build_menu_bar()
        self._refresh_status()

    # ── Menu bar ────────────────────────────────────────────────────────────────
    def _build_menu_bar(self) -> None:
        mb = self.menuBar()
        assert mb is not None

        # ── File ────────────────────────────────────────────────────────────
        file_menu = mb.addMenu("&File")
        assert file_menu is not None

        act_new = QAction("&New Database…", self)
        act_new.setShortcut(QKeySequence("Ctrl+N"))
        act_new.setStatusTip("Create a new empty database file")
        act_new.triggered.connect(self._new_db)

        act_open = QAction("&Open Database…", self)
        act_open.setShortcut(QKeySequence("Ctrl+O"))
        act_open.setStatusTip("Open an existing database file")
        act_open.triggered.connect(self._open_db)

        act_export = QAction("&Export current tab to CSV…", self)
        act_export.setShortcut(QKeySequence("Ctrl+E"))
        act_export.setStatusTip("Export the currently visible tab to a CSV file")
        act_export.triggered.connect(self._export_csv)

        act_quit = QAction("&Quit", self)
        act_quit.setShortcut(QKeySequence("Ctrl+Q"))
        act_quit.setStatusTip("Exit the application")
        act_quit.triggered.connect(self.close)

        file_menu.addAction(act_new)
        file_menu.addAction(act_open)
        file_menu.addSeparator()
        file_menu.addAction(act_export)
        file_menu.addSeparator()
        file_menu.addAction(act_quit)

        # ── Help ────────────────────────────────────────────────────────────
        help_menu = mb.addMenu("&Help")
        assert help_menu is not None

        act_about = QAction("&About", self)
        act_about.setStatusTip("About this application")
        act_about.triggered.connect(self._show_about)

        help_menu.addAction(act_about)

    # ── Toolbar ───────────────────────────────────────────────────────────────
    def _build_toolbar(self) -> QHBoxLayout:
        row = QHBoxLayout()
        row.setSpacing(6)

        def btn(text: str, tip: str) -> QPushButton:
            b = QPushButton(text)
            b.setToolTip(tip)
            b.setFixedHeight(34)
            return b

        self._btn_add = btn("➕  Add", "Add new asset  (Ins)")
        self._btn_edit = btn("✏️  Edit", "Edit selected asset  (Enter)")
        self._btn_del = btn("🗑  Delete", "Delete selected asset  (Del)")
        self._btn_export = btn("⬇  Export CSV", "Export current view to CSV")
        self._btn_alerts = btn("🔔  Alerts", "Show expiry alerts")

        self._btn_add.clicked.connect(self._add_asset)
        self._btn_edit.clicked.connect(self._edit_asset)
        self._btn_del.clicked.connect(self._delete_asset)
        self._btn_export.clicked.connect(self._export_csv)
        self._btn_alerts.clicked.connect(self._check_alerts)

        sep = QFrame()
        sep.setFrameShape(QFrame.Shape.VLine)
        sep.setFixedHeight(28)

        self._search = QLineEdit()
        self._search.setPlaceholderText("🔍  Search across all fields…")
        self._search.setMinimumWidth(280)
        self._search.setFixedHeight(34)
        self._search.textChanged.connect(self._on_search)

        btn_clear = QPushButton("✕")
        btn_clear.setFixedSize(34, 34)
        btn_clear.setToolTip("Clear search")
        btn_clear.clicked.connect(self._search.clear)

        for w in (
            self._btn_add,
            self._btn_edit,
            self._btn_del,
            self._btn_export,
            self._btn_alerts,
            sep,
            self._search,
            btn_clear,
        ):
            row.addWidget(w)
        row.addStretch()

        # keyboard shortcuts
        QShortcut(QKeySequence(Qt.Key.Key_Insert), self, self._add_asset)
        QShortcut(QKeySequence(Qt.Key.Key_Return), self, self._edit_asset)
        QShortcut(QKeySequence(Qt.Key.Key_Delete), self, self._delete_asset)
        QShortcut(QKeySequence("Ctrl+F"), self, self._search.setFocus)

        return row

    # ── Legend ────────────────────────────────────────────────────────────────
    def _build_legend(self) -> QWidget:
        bar = QWidget()
        bar.setFixedHeight(22)
        row = QHBoxLayout(bar)
        row.setContentsMargins(4, 0, 4, 0)
        row.setSpacing(16)

        for color, label in LEGEND_ITEMS:
            box = QFrame()
            box.setFixedSize(14, 14)
            box.setStyleSheet(f"background:{color}; border:1px solid #aaa; border-radius:2px;")
            lbl = QLabel(label)
            lbl.setStyleSheet("color:#555; font-size:11px;")
            row.addWidget(box)
            row.addWidget(lbl)

        row.addStretch()
        return bar

    # ── Tab bar + stacked tables ───────────────────────────────────────────────
    def _build_tabs(self) -> QWidget:
        container = QWidget()
        layout = QVBoxLayout(container)
        layout.setContentsMargins(0, 0, 0, 0)
        layout.setSpacing(0)

        tab_row = QHBoxLayout()
        tab_row.setContentsMargins(0, 0, 0, 0)
        tab_row.setSpacing(0)

        self._tab_bar = QTabBar()
        self._tab_bar.setExpanding(False)
        self._tab_bar.setDocumentMode(True)
        for tab in _MAIN_TABS:
            self._tab_bar.addTab(tab)
        self._tab_bar.currentChanged.connect(self._on_tab_changed)
        tab_row.addWidget(self._tab_bar)

        tab_row.addStretch()

        self._users_btn = QPushButton("Users")
        self._users_btn.setObjectName("usersTabBtn")
        self._users_btn.setCheckable(True)
        self._users_btn.clicked.connect(self._on_users_clicked)
        tab_row.addWidget(self._users_btn)

        layout.addLayout(tab_row)

        self._stack = QStackedWidget()
        for cat in TABS:
            model = AssetTableModel(category=cat)
            self._models[cat] = model

            table = QTableView()
            table.setModel(model)
            table.setSelectionBehavior(QAbstractItemView.SelectionBehavior.SelectRows)
            table.setSelectionMode(QAbstractItemView.SelectionMode.SingleSelection)
            table.setAlternatingRowColors(True)
            table.setEditTriggers(
                QAbstractItemView.EditTrigger.DoubleClicked
                | QAbstractItemView.EditTrigger.AnyKeyPressed
            )
            h_header = table.horizontalHeader()
            assert h_header is not None
            h_header.setSectionResizeMode(QHeaderView.ResizeMode.Interactive)
            h_header.setStretchLastSection(True)
            h_header.setMinimumSectionSize(60)
            v_header = table.verticalHeader()
            assert v_header is not None
            v_header.setDefaultSectionSize(24)
            v_header.hide()
            table.setShowGrid(True)
            table.doubleClicked.connect(lambda _idx, t=cat: self._edit_asset())

            sel_model = table.selectionModel()
            assert sel_model is not None
            sel_model.selectionChanged.connect(
                lambda _sel, _desel, t=cat: self._on_selection_changed(t)
            )

            self._tables[cat] = table
            self._stack.addWidget(table)

        layout.addWidget(self._stack)
        return container

    # ── Style ─────────────────────────────────────────────────────────────────
    def _apply_style(self) -> None:
        self.setStyleSheet(STYLESHEET)

    # ══════════════════════════════════════════════════════════════════════════
    # Helpers
    # ══════════════════════════════════════════════════════════════════════════

    def _current_model(self) -> AssetTableModel:
        return self._models[self._current_tab]

    def _current_table(self) -> QTableView:
        return self._tables[self._current_tab]

    def _refresh_all(self) -> None:
        for m in self._models.values():
            m.refresh()
        self._selected_data = None
        self._refresh_status()

    def _refresh_status(self) -> None:
        count = self._current_model().rowCount()
        self._status_bar.showMessage(
            f"  {count} record(s) shown    |    "
            f"Tab: {self._current_tab}    |    "
            f"Today: {date.today().isoformat()}"
        )

    # ══════════════════════════════════════════════════════════════════════════
    # Slots
    # ══════════════════════════════════════════════════════════════════════════

    def _on_tab_changed(self, idx: int) -> None:
        if idx < 0:
            return
        self._users_btn.setChecked(False)
        self._current_tab = _MAIN_TABS[idx]
        self._stack.setCurrentIndex(idx)
        self._selected_data = None
        self._btn_add.setEnabled(self._current_tab != "All")
        self._refresh_status()

    def _on_users_clicked(self) -> None:
        self._users_btn.setChecked(True)
        self._tab_bar.setCurrentIndex(-1)
        self._current_tab = "Users"
        self._stack.setCurrentIndex(len(_MAIN_TABS))  # Users is the last page
        self._selected_data = None
        self._btn_add.setEnabled(True)
        self._refresh_status()

    def _on_selection_changed(self, tab: str) -> None:
        """Cache the selected row for the active tab so button clicks see it."""
        if tab != self._current_tab:
            return
        sel_model = self._current_table().selectionModel()
        assert sel_model is not None
        indexes = sel_model.selectedRows()
        if indexes:
            self._selected_data = self._current_model().get_row_data(indexes[0].row())
        else:
            self._selected_data = None

    def _on_search(self, text: str) -> None:
        for m in self._models.values():
            m.set_search(text)
        self._refresh_status()

    # ── CRUD ──────────────────────────────────────────────────────────────────
    def _add_asset(self) -> None:
        if self._current_tab == "All":
            return
        dlg = AssetDialog(category=self._current_tab, locked=True, parent=self)
        if dlg.exec() == QDialog.DialogCode.Accepted:
            record = dlg.get_data()
            try:
                db_module.insert_record(record["type"], record)
            except Exception as exc:
                QMessageBox.critical(self, "Save failed", str(exc))
                return
            self._refresh_all()

    def _edit_asset(self) -> None:
        data = self._selected_data
        if data is None:
            QMessageBox.information(self, "No selection", "Select a row first.")
            return
        # "type" is present on All-tab rows; fall back to current tab
        category = data.get("type") or self._current_tab
        dlg = AssetDialog(category=category, data=data, parent=self)
        if dlg.exec() == QDialog.DialogCode.Accepted:
            record = dlg.get_data()
            try:
                db_module.update_record(category, data["id"], record)
            except Exception as exc:
                QMessageBox.critical(self, "Save failed", str(exc))
                return
            self._refresh_all()

    def _delete_asset(self) -> None:
        data = self._selected_data
        if data is None:
            QMessageBox.information(self, "No selection", "Select a row first.")
            return
        category = data.get("type") or self._current_tab
        name = data.get("name") or f"ID {data['id']}"
        reply = QMessageBox.question(
            self,
            "Confirm Delete",
            f'Permanently delete  "{name}" ?',
            QMessageBox.StandardButton.Yes | QMessageBox.StandardButton.No,
        )
        if reply == QMessageBox.StandardButton.Yes:
            db_module.delete_record(category, data["id"])
            self._refresh_all()

    # ── File menu actions ─────────────────────────────────────────────────────
    def _new_db(self) -> None:
        path, _ = QFileDialog.getSaveFileName(
            self, "New Database", "inventory.db", "SQLite Database (*.db)"
        )
        if not path:
            return
        db_module.DB_PATH = Path(path)
        db_module.init_db()
        self._refresh_all()
        self.setWindowTitle(f"IT Asset Inventory — {Path(path).name}")
        self._status_bar.showMessage(f"  New database created: {path}", 5000)

    def _open_db(self) -> None:
        path, _ = QFileDialog.getOpenFileName(
            self, "Open Database", "", "SQLite Database (*.db);;All Files (*)"
        )
        if not path:
            return
        db_module.DB_PATH = Path(path)
        db_module.init_db()  # creates missing tables if opening older file
        self._refresh_all()
        self.setWindowTitle(f"IT Asset Inventory — {Path(path).name}")
        self._status_bar.showMessage(f"  Opened: {path}", 5000)

    def _show_about(self) -> None:
        QMessageBox.about(
            self,
            "About IT Asset Inventory",
            "<h3>IT Asset Inventory</h3>"
            "<p>A desktop tool to track computers, devices, licenses "
            "and software subscriptions.</p>"
            "<p><b>Stack:</b> Python 3.14 · PyQt6 · SQLite</p>"
            "<p><b>License:</b> GNU GPL v3 or later</p>"
            "<p><b>DB file:</b> " + str(db_module.DB_PATH) + "</p>",
        )

    # ── Export ────────────────────────────────────────────────────────────────
    def _export_csv(self) -> None:
        model = self._current_model()
        if model.rowCount() == 0:
            QMessageBox.information(self, "Export", "Nothing to export.")
            return
        default_name = f"{self._current_tab.replace(' ', '_')}_assets.csv"
        path, _ = QFileDialog.getSaveFileName(self, "Save CSV", default_name, "CSV Files (*.csv)")
        if not path:
            return
        cols = model.column_keys()
        with open(path, "w", newline="", encoding="utf-8") as fh:
            writer = csv.writer(fh)
            writer.writerow([db_module.COLUMN_LABELS.get(c, c) for c in cols])
            for r in range(model.rowCount()):
                writer.writerow([model.data(model.index(r, c)) for c in range(model.columnCount())])
        QMessageBox.information(
            self,
            "Export complete",
            f"Exported {model.rowCount()} record(s) to:\n{path}",
        )

    # ── Alerts ────────────────────────────────────────────────────────────────
    def _check_alerts(self) -> None:
        expired = db_module.get_expired()
        expiring = db_module.get_expiring_soon(db_module.EXPIRY_WARNING_DAYS)

        if not expired and not expiring:
            # Only show positive message when triggered manually via button
            sender = self.sender()
            if sender is self._btn_alerts:
                QMessageBox.information(self, "Alerts", "No expiry issues found. ✔")
            return

        lines: list[str] = []
        if expired:
            lines.append(f"🔴  EXPIRED ({len(expired)} item(s)):")
            for r in expired[:10]:
                exp = r.get("expiry_date") or r.get("warranty_expiry") or "?"
                lines.append(f"   \u2022 {r.get('name', '?')}  [{r.get('type', '')}]  \u2014 {exp}")
            if len(expired) > 10:
                lines.append(f"   \u2026 and {len(expired) - 10} more")

        if expiring:
            lines.append("")
            lines.append(
                f"\U0001f7e1  Expiring within {db_module.EXPIRY_WARNING_DAYS} days ({len(expiring)} item(s)):"
            )
            for r in expiring[:10]:
                exp = r.get("expiry_date") or r.get("warranty_expiry") or "?"
                lines.append(f"   \u2022 {r.get('name', '?')}  [{r.get('type', '')}]  \u2014 {exp}")
            if len(expiring) > 10:
                lines.append(f"   … and {len(expiring) - 10} more")

        QMessageBox.warning(self, "Expiry Alerts", "\n".join(lines))
