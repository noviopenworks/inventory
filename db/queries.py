"""
db.queries – Generic CRUD operations, fetch helpers, and the All-tab UNION.
"""

from db.config import ALL_TAB_CATS, CATEGORIES, TABLE_CONFIG
from db.connection import get_connection

# ── Non-text columns (excluded from text search) ──────────────────────────────
_NON_TEXT = frozenset(
    {
        "user_id",
        "computer_id",
        "smartphone_id",
        "tablet_id",
        "status",
        "purchase_date",
        "warranty_expiry",
        "expiry_date",
    }
)


# ── Users helpers ──────────────────────────────────────────────────────────────
def fetch_users() -> list[dict]:
    """Return [{id, name, surname, status, notes}] ordered by name."""
    with get_connection() as conn:
        rows = conn.execute(
            "SELECT id, name, surname, status, notes FROM users ORDER BY name"
        ).fetchall()
    return [dict(r) for r in rows]


def fetch_computers() -> list[dict]:
    """Return [{id, name}] ordered by name."""
    with get_connection() as conn:
        rows = conn.execute("SELECT id, name FROM computers ORDER BY name").fetchall()
    return [dict(r) for r in rows]


def fetch_smartphones() -> list[dict]:
    """Return [{id, name}] ordered by name."""
    with get_connection() as conn:
        rows = conn.execute("SELECT id, name FROM smartphones ORDER BY name").fetchall()
    return [dict(r) for r in rows]


def fetch_tablets() -> list[dict]:
    """Return [{id, name}] ordered by name."""
    with get_connection() as conn:
        rows = conn.execute("SELECT id, name FROM tablets ORDER BY name").fetchall()
    return [dict(r) for r in rows]


# ── Generic SELECT builder ─────────────────────────────────────────────────────
def _base_select(category: str) -> str:
    cfg = TABLE_CONFIG[category]
    tbl = cfg["table"]
    db_cols = cfg["db_cols"]
    selects = ["t.*"]
    joins: list[str] = []
    if "user_id" in db_cols:
        selects.append("u.name AS user_name")
        joins.append("LEFT JOIN users u ON t.user_id = u.id")
    has_comp = "computer_id" in db_cols
    has_phone = "smartphone_id" in db_cols
    has_tablet = "tablet_id" in db_cols
    if has_comp:
        joins.append("LEFT JOIN computers c ON t.computer_id = c.id")
    if has_phone:
        joins.append("LEFT JOIN smartphones s ON t.smartphone_id = s.id")
    if has_tablet:
        joins.append("LEFT JOIN tablets tb ON t.tablet_id = tb.id")
    # Build device_name expression
    device_parts: list[str] = []
    if has_comp:
        device_parts.append("c.name")
    if has_phone:
        device_parts.append("s.name")
    if has_tablet:
        device_parts.append("tb.name")
    if len(device_parts) > 1:
        selects.append(f"COALESCE({', '.join(device_parts)}) AS device_name")
    elif device_parts:
        selects.append(f"{device_parts[0]} AS device_name")
    return f"SELECT {', '.join(selects)} FROM {tbl} t {' '.join(joins)}"


# ── CRUD ───────────────────────────────────────────────────────────────────────
def fetch_records(category: str, search: str | None = None) -> list[dict]:
    if category == "All":
        return _fetch_all(search)

    cfg = TABLE_CONFIG[category]
    base = _base_select(category)
    conditions: list[str] = []
    params: list = []

    if search:
        text_cols = [c for c in cfg["db_cols"] if c not in _NON_TEXT]
        or_parts = [f"t.{c} LIKE ?" for c in text_cols]
        if "user_id" in cfg["db_cols"]:
            or_parts.append("u.name LIKE ?")
        if "computer_id" in cfg["db_cols"]:
            or_parts.append("c.name LIKE ?")
        if "smartphone_id" in cfg["db_cols"]:
            or_parts.append("s.name LIKE ?")
        if "tablet_id" in cfg["db_cols"]:
            or_parts.append("tb.name LIKE ?")
        conditions.append(f"({' OR '.join(or_parts)})")
        params.extend([f"%{search}%"] * len(or_parts))

    where = f" WHERE {' AND '.join(conditions)}" if conditions else ""
    with get_connection() as conn:
        rows = conn.execute(f"{base}{where} ORDER BY t.id DESC", params).fetchall()
    return [dict(r) for r in rows]


def _fetch_all(search: str | None = None) -> list[dict]:
    """UNION ALL across hardware device tables (Computer, Smartphone, Tablet)."""
    asset_cats = [c for c in CATEGORIES if c in ALL_TAB_CATS]

    def row_expr(cat: str) -> tuple[str, list]:
        cfg = TABLE_CONFIG[cat]
        tbl = cfg["table"]
        db_cols = cfg["db_cols"]
        has_name = "name" in db_cols
        has_model = "model" in db_cols
        has_purchase = "purchase_date" in db_cols
        has_warranty = "warranty_expiry" in db_cols

        # All ALL_TAB_CATS are hardware with user_id
        user_expr = "u.name" if "user_id" in db_cols else "NULL"
        joins = "LEFT JOIN users u ON t.user_id = u.id" if "user_id" in db_cols else ""

        select = (
            f"SELECT t.id, '{cat}' AS type, "
            f"{'t.name' if has_name else 'NULL'} AS name, "
            f"{'t.model' if has_model else 'NULL'} AS model, "
            f"{user_expr} AS user_name, "
            f"t.status, "
            f"{'t.purchase_date' if has_purchase else 'NULL'} AS purchase_date, "
            f"{'t.warranty_expiry' if has_warranty else 'NULL'} AS warranty_expiry, "
            f"t.notes "
            f"FROM {tbl} t {joins}"
        )
        if not search:
            return select, []

        text_cols = [c for c in db_cols if c not in _NON_TEXT]
        or_parts = [f"t.{c} LIKE ?" for c in text_cols]
        if "user_id" in db_cols:
            or_parts.append("u.name LIKE ?")
        p = [f"%{search}%"] * len(or_parts)
        return f"{select} WHERE ({' OR '.join(or_parts)})", p

    parts, params = [], []
    for cat in asset_cats:
        sql, p = row_expr(cat)
        parts.append(sql)
        params.extend(p)

    query = f"SELECT * FROM ({' UNION ALL '.join(parts)}) ORDER BY type, id DESC"
    with get_connection() as conn:
        rows = conn.execute(query, params).fetchall()
    return [dict(r) for r in rows]


def insert_record(category: str, data: dict) -> int:
    cfg = TABLE_CONFIG[category]
    tbl = cfg["table"]
    cols = [c for c in cfg["db_cols"] if data.get(c) not in (None, "")]
    with get_connection() as conn:
        if cols:
            placeholders = ", ".join("?" * len(cols))
            values = [data[c] for c in cols]
            cur = conn.execute(
                f"INSERT INTO {tbl} ({', '.join(cols)}) VALUES ({placeholders})", values
            )
        else:
            cur = conn.execute(f"INSERT INTO {tbl} DEFAULT VALUES")
        conn.commit()
        return cur.lastrowid or 0


def update_record(category: str, record_id: int, data: dict) -> None:
    cfg = TABLE_CONFIG[category]
    tbl = cfg["table"]
    cols = [c for c in cfg["db_cols"] if c in data]
    if not cols:
        return
    set_clause = ", ".join(f"{c} = ?" for c in cols)
    set_clause += ", updated_at = datetime('now','localtime')"
    values = [data[c] for c in cols] + [record_id]
    with get_connection() as conn:
        conn.execute(f"UPDATE {tbl} SET {set_clause} WHERE id = ?", values)
        conn.commit()


def delete_record(category: str, record_id: int) -> None:
    tbl = TABLE_CONFIG[category]["table"]
    with get_connection() as conn:
        conn.execute(f"DELETE FROM {tbl} WHERE id = ?", (record_id,))
        conn.commit()
