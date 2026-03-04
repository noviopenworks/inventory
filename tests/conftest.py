"""
tests/conftest.py – Fixtures shared across the entire test suite.
"""

import pytest

import db as db_module


@pytest.fixture()
def tmp_db(tmp_path, monkeypatch):
    """Point DB_PATH at a fresh temp file and initialise the schema."""
    db_file = tmp_path / "test.db"
    monkeypatch.setattr(db_module, "DB_PATH", db_file)
    db_module.init_db()
    yield db_file


@pytest.fixture()
def win(tmp_db, qtbot):
    """A fully constructed MainWindow backed by a temp DB."""
    from app.main_window import MainWindow

    w = MainWindow()
    qtbot.addWidget(w)
    return w
