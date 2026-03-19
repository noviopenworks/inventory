"""
db.config – Per-table configuration, column metadata, and application constants.
"""

from typing import TypedDict


# ── Per-table configuration ────────────────────────────────────────────────────
class TableSpec(TypedDict):
    """Schema for each entry in TABLE_CONFIG."""

    table: str
    db_cols: list[str]
    display_cols: list[str]
    has_user_fk: bool


TABLE_CONFIG: dict[str, TableSpec] = {
    "Computer": {
        "table": "computers",
        "db_cols": [
            "name",
            "model",
            "user_id",
            "status",
            "purchase_date",
            "warranty_expiry",
            "notes",
        ],
        "display_cols": [
            "id",
            "name",
            "model",
            "user_name",
            "status",
            "purchase_date",
            "warranty_expiry",
            "notes",
        ],
        "has_user_fk": True,
    },
    "Smartphone": {
        "table": "smartphones",
        "db_cols": [
            "name",
            "model",
            "user_id",
            "status",
            "purchase_date",
            "warranty_expiry",
            "notes",
        ],
        "display_cols": [
            "id",
            "name",
            "model",
            "user_name",
            "status",
            "purchase_date",
            "warranty_expiry",
            "notes",
        ],
        "has_user_fk": True,
    },
    "Tablet": {
        "table": "tablets",
        "db_cols": [
            "name",
            "model",
            "user_id",
            "status",
            "purchase_date",
            "warranty_expiry",
            "notes",
        ],
        "display_cols": [
            "id",
            "name",
            "model",
            "user_name",
            "status",
            "purchase_date",
            "warranty_expiry",
            "notes",
        ],
        "has_user_fk": True,
    },
    "Windows Key": {
        "table": "windows_keys",
        "db_cols": [
            "license_key",
            "computer_id",
            "status",
            "notes",
        ],
        "display_cols": [
            "id",
            "license_key",
            "device_name",
            "status",
            "notes",
        ],
        "has_user_fk": False,
    },
    "Antivirus": {
        "table": "antivirus",
        "db_cols": [
            "name",
            "license_key",
            "computer_id",
            "smartphone_id",
            "tablet_id",
            "status",
            "expiry_date",
            "notes",
        ],
        "display_cols": [
            "id",
            "name",
            "license_key",
            "device_name",
            "status",
            "expiry_date",
            "notes",
        ],
        "has_user_fk": False,
    },
    "Other Software": {
        "table": "other_software",
        "db_cols": [
            "name",
            "license_key",
            "computer_id",
            "smartphone_id",
            "tablet_id",
            "status",
            "expiry_date",
            "notes",
        ],
        "display_cols": [
            "id",
            "name",
            "license_key",
            "device_name",
            "status",
            "expiry_date",
            "notes",
        ],
        "has_user_fk": False,
    },
    "Users": {
        "table": "users",
        "db_cols": ["name", "surname", "status", "notes"],
        "display_cols": ["id", "name", "surname", "status", "notes"],
        "has_user_fk": False,
    },
}

CATEGORIES: list[str] = list(TABLE_CONFIG.keys())

# Columns visible in the "All" overview tab
ALL_DISPLAY_COLS = [
    "type",
    "id",
    "name",
    "model",
    "user_name",
    "status",
    "purchase_date",
    "warranty_expiry",
    "notes",
]

COLUMN_LABELS: dict[str, str] = {
    "id": "ID",
    "type": "Type",
    "name": "Name",
    "brand": "Brand / Vendor",
    "model": "Model",
    "serial_number": "Serial #",
    "email": "Email",
    "surname": "Surname",
    "user_name": "User",
    "user_id": "User",
    "device_name": "Device",
    "status": "Status",
    "purchase_date": "Purchase Date",
    "warranty_expiry": "Warranty Expiry",
    "expiry_date": "Expiry Date",
    "license_key": "License Key",
    "notes": "Notes",
    "created_at": "Created",
    "updated_at": "Updated",
}

ALL_TAB_CATS: set[str] = {"Computer", "Smartphone", "Tablet"}

STATUS_OPTIONS = ["Active", "In Repair", "Spare", "Retired", "Missing"]
USER_STATUS_OPTIONS = ["Active", "On Leave", "Inactive", "Terminated"]
EXPIRY_WARNING_DAYS = 30
