"""
app.style – Theme helpers for qt-material integration.
"""

import db as db_module

# Theme names passed to apply_stylesheet()
THEME_DARK = "dark_blue.xml"
THEME_LIGHT = "light_blue.xml"

# Small app-level tweaks applied on top of qt-material via setStyleSheet.
# Keep this minimal – qt-material already handles inputs, buttons, combos, etc.
EXTRA_STYLESHEET = """
    QTabBar::tab {
        padding: 5px 20px;
        font-weight: bold;
        font-size: 12px;
        border-bottom: none;
        border-radius: 4px 4px 0 0;
        margin-right: 2px;
    }
    QStatusBar { font-size: 11px; }
"""


def get_legend_items(dark: bool = False) -> list[tuple[str, str]]:
    if dark:
        return [
            ("#b71c1c", "Expired"),
            ("#f57f17", f"Expires \u2264 {db_module.EXPIRY_WARNING_DAYS} days"),
            ("#2a2d4a", "Retired"),
        ]
    return [
        ("#ffcdd2", "Expired"),
        ("#fff9c4", f"Expires \u2264 {db_module.EXPIRY_WARNING_DAYS} days"),
        ("#d0d0d0", "Retired"),
    ]


LEGEND_ITEMS = get_legend_items(dark=False)

# Backwards-compat aliases
LIGHT_STYLESHEET = ""
DARK_STYLESHEET = ""
STYLESHEET = ""
