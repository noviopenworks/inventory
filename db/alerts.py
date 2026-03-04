"""
db.alerts – Warranty / licence expiry alert queries.
"""

from datetime import date, timedelta

from db.config import EXPIRY_WARNING_DAYS, TABLE_CONFIG
from db.connection import get_connection


def get_expiring_soon(days: int = EXPIRY_WARNING_DAYS) -> list[dict]:
    today = date.today().isoformat()
    threshold = (date.today() + timedelta(days=days)).isoformat()
    results: list[dict] = []
    with get_connection() as conn:
        for cat, cfg in TABLE_CONFIG.items():
            if cat == "Users":
                continue
            tbl = cfg["table"]
            exp_col = "expiry_date" if "expiry_date" in cfg["db_cols"] else "warranty_expiry"
            if exp_col not in cfg["db_cols"]:
                continue
            rows = conn.execute(
                f"SELECT *, '{cat}' AS type FROM {tbl} "
                f"WHERE {exp_col} BETWEEN ? AND ? AND status != 'Retired' "
                f"ORDER BY {exp_col}",
                (today, threshold),
            ).fetchall()
            results.extend(dict(r) for r in rows)
    return results


def get_expired() -> list[dict]:
    today = date.today().isoformat()
    results: list[dict] = []
    with get_connection() as conn:
        for cat, cfg in TABLE_CONFIG.items():
            if cat == "Users":
                continue
            tbl = cfg["table"]
            exp_col = "expiry_date" if "expiry_date" in cfg["db_cols"] else "warranty_expiry"
            if exp_col not in cfg["db_cols"]:
                continue
            rows = conn.execute(
                f"SELECT *, '{cat}' AS type FROM {tbl} "
                f"WHERE {exp_col} < ? AND status != 'Retired' "
                f"ORDER BY {exp_col}",
                (today,),
            ).fetchall()
            results.extend(dict(r) for r in rows)
    return results
