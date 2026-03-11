"""
tests/conftest.py – Fixtures shared across the entire test suite.
"""

import pytest

import db as db_module
from app.models import AssetTableModel


@pytest.fixture(autouse=True)
def reset_dark_mode():
    """Ensure AssetTableModel always starts in light mode for every test."""
    AssetTableModel.set_dark_mode(False)
    yield
    AssetTableModel.set_dark_mode(False)


@pytest.fixture()
def tmp_db(tmp_path, monkeypatch):
    """Point DB_PATH at a fresh temp file and initialise the schema."""
    db_file = tmp_path / "test.db"
    monkeypatch.setattr(db_module, "DB_PATH", db_file)
    db_module.init_db()
    yield db_file


@pytest.fixture()
def win(tmp_db, qtbot, monkeypatch):
    """A fully constructed MainWindow backed by a temp DB, always in light mode."""
    from PyQt6.QtCore import QSettings

    # Prevent QSettings from reading real user preferences (e.g. dark_mode=True)
    monkeypatch.setattr(QSettings, "value", lambda self, key, default=None, type=None: default)
    monkeypatch.setattr(QSettings, "setValue", lambda self, key, val: None)

    from app.main_window import MainWindow

    w = MainWindow()
    qtbot.addWidget(w)
    return w
