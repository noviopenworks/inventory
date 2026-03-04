"""
db – Relational SQLite layer for IT Asset Inventory.

This package re-exports the full public API so that existing
``import db as db_module`` statements continue to work unchanged.

Schema
------
  users                                 <- people
  computers            user_id -> users
  smartphones          user_id -> users
  tablets              user_id -> users
  windows_keys         computer_id -> computers
  antivirus            computer_id -> computers, smartphone_id -> smartphones, tablet_id -> tablets
  other_software       computer_id -> computers, smartphone_id -> smartphones, tablet_id -> tablets
"""

# ── configuration & constants ──────────────────────────────────────────────────
# ── expiry alerts ──────────────────────────────────────────────────────────────
from db.alerts import (
    get_expired as get_expired,
)
from db.alerts import (
    get_expiring_soon as get_expiring_soon,
)
from db.config import (
    _ALL_TAB_CATS as _ALL_TAB_CATS,
)
from db.config import (
    ALL_DISPLAY_COLS as ALL_DISPLAY_COLS,
)
from db.config import (
    CATEGORIES as CATEGORIES,
)
from db.config import (
    COLUMN_LABELS as COLUMN_LABELS,
)
from db.config import (
    EXPIRY_WARNING_DAYS as EXPIRY_WARNING_DAYS,
)
from db.config import (
    STATUS_OPTIONS as STATUS_OPTIONS,
)
from db.config import (
    TABLE_CONFIG as TABLE_CONFIG,
)
from db.config import (
    USER_STATUS_OPTIONS as USER_STATUS_OPTIONS,
)
from db.config import (
    TableSpec as TableSpec,
)

# ── connection ─────────────────────────────────────────────────────────────────
from db.connection import (
    DB_PATH as DB_PATH,
)
from db.connection import (
    get_connection as get_connection,
)

# ── CRUD & fetch helpers ───────────────────────────────────────────────────────
from db.queries import (
    delete_record as delete_record,
)
from db.queries import (
    fetch_computers as fetch_computers,
)
from db.queries import (
    fetch_records as fetch_records,
)
from db.queries import (
    fetch_smartphones as fetch_smartphones,
)
from db.queries import (
    fetch_tablets as fetch_tablets,
)
from db.queries import (
    fetch_users as fetch_users,
)
from db.queries import (
    insert_record as insert_record,
)
from db.queries import (
    update_record as update_record,
)

# ── schema / migrations ────────────────────────────────────────────────────────
from db.schema import (
    SCHEMA_VERSION as SCHEMA_VERSION,
)
from db.schema import (
    _get_schema_version as _get_schema_version,
)
from db.schema import (
    _migrate as _migrate,
)
from db.schema import (
    _set_schema_version as _set_schema_version,
)
from db.schema import (
    init_db as init_db,
)

__all__ = [
    # config
    "ALL_DISPLAY_COLS",
    "CATEGORIES",
    "COLUMN_LABELS",
    "EXPIRY_WARNING_DAYS",
    "STATUS_OPTIONS",
    "TABLE_CONFIG",
    "TableSpec",
    # connection
    "DB_PATH",
    "get_connection",
    # schema
    "SCHEMA_VERSION",
    "init_db",
    # queries
    "delete_record",
    "fetch_computers",
    "fetch_records",
    "fetch_smartphones",
    "fetch_tablets",
    "fetch_users",
    "insert_record",
    "update_record",
    # alerts
    "get_expired",
    "get_expiring_soon",
]
