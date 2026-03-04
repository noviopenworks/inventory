"""
app.style – Application-wide Qt stylesheet.
"""

import db as db_module

STYLESHEET = """
    QMainWindow, QWidget { background: #f4f6f9; }

    QTabBar::tab {
        padding: 5px 20px;
        font-weight: bold;
        font-size: 12px;
        background: #dce3ed;
        color: #444;
        border: 1px solid #b0bec5;
        border-bottom: none;
        border-radius: 4px 4px 0 0;
        margin-right: 2px;
    }
    QTabBar::tab:selected  { background: #1565c0; color: white; }
    QTabBar::tab:hover:!selected { background: #bbdefb; }

    QTableView {
        background: white;
        alternate-background-color: #f1f5fb;
        gridline-color: #d0d8e4;
        font-size: 12px;
        selection-background-color: #bbdefb;
        selection-color: #0d47a1;
        border: 1px solid #b0bec5;
    }
    QHeaderView::section {
        background: #e3eaf4;
        padding: 4px 8px;
        border: 1px solid #c5cdd8;
        font-weight: bold;
        font-size: 11px;
    }

    QPushButton {
        padding: 4px 14px;
        border-radius: 4px;
        background: #e8edf4;
        border: 1px solid #b0bec5;
        font-size: 12px;
    }
    QPushButton:hover   { background: #bbdefb; border-color: #1976d2; }
    QPushButton:pressed { background: #90caf9; }

    QPushButton#usersTabBtn {
        padding: 5px 20px;
        font-weight: bold;
        font-size: 12px;
        background: #dce3ed;
        color: #444;
        border: 1px solid #b0bec5;
        border-bottom: none;
        border-radius: 4px 4px 0 0;
        margin-bottom: 0px;
    }
    QPushButton#usersTabBtn:checked { background: #1565c0; color: white; }
    QPushButton#usersTabBtn:hover:!checked { background: #bbdefb; }

    QLineEdit {
        border: 1px solid #b0bec5;
        border-radius: 4px;
        padding: 4px 8px;
        background: white;
        font-size: 12px;
    }
    QLineEdit:focus { border-color: #1976d2; }

    QStatusBar { font-size: 11px; color: #555; }
"""

LEGEND_ITEMS: list[tuple[str, str]] = [
    ("#ffcdd2", "Expired"),
    ("#fff9c4", f"Expires ≤ {db_module.EXPIRY_WARNING_DAYS} days"),
    ("#eeeeee", "Retired"),
]
