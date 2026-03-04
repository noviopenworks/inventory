"""
tests/helpers.py – Lightweight factory helpers shared across DB test modules.

These are plain functions (not fixtures) that insert common records into
whichever DB is currently active (i.e. the one monkeypatched by tmp_db).
"""

from datetime import date, timedelta

import db as db_module


def _add_user(name: str = "Alice", surname: str = "Smith") -> int:
    return db_module.insert_record("Users", {"name": name, "surname": surname, "status": "Active"})


def _add_computer(name: str = "PC-01", user_id: int | None = None) -> int:
    return db_module.insert_record(
        "Computer",
        {"name": name, "model": "ThinkPad X1", "user_id": user_id, "status": "Active"},
    )


def _add_phone(name: str = "iPhone-01", user_id: int | None = None) -> int:
    return db_module.insert_record(
        "Smartphone",
        {"name": name, "model": "iPhone 15", "user_id": user_id, "status": "Active"},
    )


def _days(delta: int) -> str:
    """Return an ISO date string relative to today."""
    return (date.today() + timedelta(days=delta)).isoformat()
