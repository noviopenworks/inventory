"""
app – GUI layer for IT Asset Inventory.

Re-exports the main public classes so that
``from app.main_window import MainWindow`` etc. work.
"""

from app.dialogs import AssetDialog as AssetDialog
from app.main_window import _MAIN_TABS as _MAIN_TABS
from app.main_window import TABS as TABS
from app.main_window import MainWindow as MainWindow
from app.models import AssetTableModel as AssetTableModel
