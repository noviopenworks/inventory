"""
main.py – Entry point for IT Asset Inventory.
"""

import sys

from PyQt6.QtGui import QFont
from PyQt6.QtWidgets import QApplication

import db as db_module
from app._version import __version__
from app.main_window import MainWindow


def main() -> None:
    db_module.init_db()

    app = QApplication(sys.argv)
    app.setApplicationName("IT Asset Inventory")
    app.setApplicationVersion(__version__)
    app.setOrganizationName("ITDept")

    font = QFont("Segoe UI", 10)
    app.setFont(font)

    window = MainWindow()
    window.show()
    sys.exit(app.exec())


if __name__ == "__main__":
    main()
