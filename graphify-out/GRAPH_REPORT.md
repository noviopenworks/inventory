# Graph Report - .  (2026-06-17)

## Corpus Check
- Corpus is ~12,588 words - fits in a single context window. You may not need a graph.

## Summary
- 459 nodes · 916 edges · 30 communities (22 shown, 8 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 70 edges (avg confidence: 0.63)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CRUD Test Suite|CRUD Test Suite]]
- [[_COMMUNITY_Asset Dialog Layer|Asset Dialog Layer]]
- [[_COMMUNITY_Main Window UI|Main Window UI]]
- [[_COMMUNITY_Main Window Tests|Main Window Tests]]
- [[_COMMUNITY_Window Integration Tests|Window Integration Tests]]
- [[_COMMUNITY_DB Schema & Fixtures|DB Schema & Fixtures]]
- [[_COMMUNITY_Expiry Alert Engine|Expiry Alert Engine]]
- [[_COMMUNITY_Query Config & Builder|Query Config & Builder]]
- [[_COMMUNITY_Alert Test Coverage|Alert Test Coverage]]
- [[_COMMUNITY_Theme & Style System|Theme & Style System]]
- [[_COMMUNITY_Asset Table Model|Asset Table Model]]
- [[_COMMUNITY_Model Unit Tests|Model Unit Tests]]
- [[_COMMUNITY_Dialog Form Builder|Dialog Form Builder]]
- [[_COMMUNITY_DB Migration System|DB Migration System]]
- [[_COMMUNITY_Model Data Access|Model Data Access]]
- [[_COMMUNITY_Model Row Tests|Model Row Tests]]
- [[_COMMUNITY_App Init & Tab Bar|App Init & Tab Bar]]
- [[_COMMUNITY_Model Column Tests|Model Column Tests]]
- [[_COMMUNITY_DB Config Schema|DB Config Schema]]
- [[_COMMUNITY_Dev Toolchain|Dev Toolchain]]
- [[_COMMUNITY_Model Dark Mode|Model Dark Mode]]
- [[_COMMUNITY_Model SetData Tests|Model SetData Tests]]
- [[_COMMUNITY_Version String|Version String]]
- [[_COMMUNITY_All Display Cols|All Display Cols]]
- [[_COMMUNITY_Column Labels|Column Labels]]
- [[_COMMUNITY_Status Options|Status Options]]
- [[_COMMUNITY_User Status Options|User Status Options]]

## God Nodes (most connected - your core abstractions)
1. `AssetDialog` - 92 edges
2. `MainWindow` - 44 edges
3. `_add_computer()` - 42 edges
4. `AssetTableModel` - 39 edges
5. `_add_user()` - 24 edges
6. `get_connection()` - 22 edges
7. `_add_phone()` - 21 edges
8. `_tab_index()` - 16 edges
9. `TestAssetDialogGetData` - 16 edges
10. `_tab_index()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `IT Asset Inventory (README)` --references--> `db Package (Public API Facade)`  [INFERRED]
  README.md → db/__init__.py
- `WAL Journal Mode (SQLite Write-Ahead Logging)` --rationale_for--> `IT Asset Inventory (README)`  [INFERRED]
  db/connection.py → README.md
- `int` --uses--> `AssetDialog`  [INFERRED]
  tests/test_window.py → app/dialogs.py
- `object` --uses--> `AssetDialog`  [INFERRED]
  tests/test_window.py → app/dialogs.py
- `int` --uses--> `AssetDialog`  [INFERRED]
  tests/test_main_window.py → app/dialogs.py

## Hyperedges (group relationships)
- **MainWindow CRUD: add/edit/delete via AssetDialog + AssetTableModel** — app_main_window_mainwindow_add_asset, app_main_window_mainwindow_edit_asset, app_main_window_mainwindow_delete_asset, app_dialogs_assetdialog, app_models_assettablemodel [EXTRACTED 1.00]
- **Expiry alerting: get_expired / get_expiring_soon drives check_alerts and row background colouring** — app_main_window_mainwindow_check_alerts, app_main_window_mainwindow_show_alerts_dialog, app_models_assettablemodel_row_bg, tests_test_db_alerts_testexpiryalerts [INFERRED 0.85]
- **Duplicated test coverage in test_window.py and test_main_window.py / test_dialogs.py** — tests_test_window_testmainwindowadd, tests_test_main_window_testmainwindowadd, tests_test_window_testassetdialoggetdata, tests_test_dialogs_testassetdialoggetdata [INFERRED 0.85]
- **DB Layer CRUD Pattern: TABLE_CONFIG + get_connection + query functions** — db_config_table_config, db_connection_get_connection, db_queries_insert_record, db_queries_update_record, db_queries_delete_record, db_queries_fetch_records [INFERRED 0.95]
- **Test Isolation Pattern: tmp_db fixture + monkeypatch DB_PATH + init_db** — tests_conftest_tmp_db, db_connection_db_path, db_schema_init_db [EXTRACTED 1.00]
- **CI/Release Pipeline: Taskfile tasks orchestrated by GitHub Actions workflows** — taskfile_taskfile, workflows_ci_yml, workflows_release_yml, precommit_config [INFERRED 0.85]

## Communities (30 total, 8 thin omitted)

### Community 0 - "CRUD Test Suite"
Cohesion: 0.05
Nodes (18): _add_computer(), _add_phone(), _add_user(), int, str, tests/helpers.py – Lightweight factory helpers shared across DB test modules.  T, tests/test_db_crud.py – CRUD operations for every entity type., DEFAULT VALUES path: name has NOT NULL – confirms the guard is the DB constraint (+10 more)

### Community 1 - "Asset Dialog Layer"
Cohesion: 0.06
Nodes (14): AssetDialog, bool, str, Check date fields before accepting the dialog., Modal form to create or edit a single record., str, tests/test_dialogs.py – Tests for AssetDialog: get_data(), populate, and date va, TestAssetDialogDateValidation (+6 more)

### Community 2 - "Main Window UI"
Cohesion: 0.06
Nodes (24): MainWindow, bool, int, Cache the selected row for the active tab so button clicks see it., Load persisted column widths from QSettings and apply to current tab., Select a row by index, scroll to it, and update button states., Find the row whose 'id' matches row_id and select it., Double-clicking empty table space opens the Add dialog. (+16 more)

### Community 3 - "Main Window Tests"
Cohesion: 0.11
Nodes (13): _accept_with(), int, object, str, tests/test_main_window.py – MainWindow CRUD operations and tab behaviour., Patch AssetDialog so exec() returns Accepted and get_data() returns record., Switch win to the given tab., _switch_to() (+5 more)

### Community 4 - "Window Integration Tests"
Cohesion: 0.11
Nodes (13): _accept_with(), int, object, str, tests/test_window.py – Tests for AssetDialog (dialogs.py) and, Patch AssetDialog so exec() returns Accepted and get_data() returns record., Switch *win* to the given tab., _switch_to() (+5 more)

### Community 5 - "DB Schema & Fixtures"
Cohesion: 0.10
Nodes (14): DB_PATH, SCHEMA_VERSION, Point DB_PATH at a fresh temp file and initialise the schema., tmp_db(), tests/test_db_schema.py – Schema initialisation, migration, and TABLE_CONFIG., The old 'subscriptions' table must not be created by init_db., Calling init_db() twice must not raise., After init_db the _meta table stores the current schema version. (+6 more)

### Community 6 - "Expiry Alert Engine"
Cohesion: 0.16
Nodes (18): WAL Journal Mode (SQLite Write-Ahead Logging), get_expired(), get_expiring_soon(), int, db.alerts – Warranty / licence expiry alert queries., EXPIRY_WARNING_DAYS, get_connection(), Connection (+10 more)

### Community 7 - "Query Config & Builder"
Cohesion: 0.25
Nodes (16): UNION ALL All-tab Overview Pattern, ALL_TAB_CATS, CATEGORIES, TABLE_CONFIG, db Package (Public API Facade), _base_select(), delete_record(), _fetch_all() (+8 more)

### Community 8 - "Alert Test Coverage"
Cohesion: 0.24
Nodes (4): _days(), Return an ISO date string relative to today., tests/test_db_alerts.py – Expiry / warranty alert queries., TestExpiryAlerts

### Community 9 - "Theme & Style System"
Cohesion: 0.19
Nodes (8): AssetTableModel._row_bg(), get_legend_items(), bool, str, app.style – Theme helpers for qt-material integration., THEME_DARK constant, object, TestModelBackground

### Community 10 - "Asset Table Model"
Cohesion: 0.28
Nodes (6): AssetTableModel, _cols_for(), str, Flat table model for one asset category (or 'All')., Orientation, QAbstractTableModel

### Community 11 - "Model Unit Tests"
Cohesion: 0.17
Nodes (4): tests/test_models.py – Tests for AssetTableModel (models.py).  All tests use an, TestModelFlags, TestModelHeader, TestModelSearch

### Community 12 - "Dialog Form Builder"
Cohesion: 0.31
Nodes (7): _add_separator(), _close_popup(), _mono_font(), app.dialogs – Add / Edit record dialog (relational schema edition)., Add a single 'Device' dropdown combining computers, phones, and/or tablets., Schedule hidePopup after the current event cycle (fixes Linux/Wayland)., QComboBox

### Community 13 - "DB Migration System"
Cohesion: 0.29
Nodes (10): Forward-Only Schema Migration Pattern, _get_schema_version(), init_db(), _migrate(), Connection, int, db.schema – DDL initialisation and forward-only migrations., Read the stored schema version (0 if never set). (+2 more)

### Community 14 - "Model Data Access"
Cohesion: 0.24
Nodes (4): int, _row_bg(), ItemFlag, QModelIndex

### Community 16 - "App Init & Tab Bar"
Cohesion: 0.22
Nodes (6): app – GUI layer for IT Asset Inventory.  Re-exports the main public classes so t, _NameSizedTabBar, app.main_window – Main application window., Tab bar where each tab is sized to fit its label text., QSize, QTabBar

### Community 18 - "DB Config Schema"
Cohesion: 0.40
Nodes (4): db.config – Per-table configuration, column metadata, and application constants., Schema for each entry in TABLE_CONFIG., TableSpec, TypedDict

### Community 19 - "Dev Toolchain"
Cohesion: 0.40
Nodes (5): Pre-commit Config (ruff, mypy, version-sync), IT Asset Inventory (README), Taskfile (Build & Dev Tasks), CI Workflow (GitHub Actions), Release Workflow (GitHub Actions)

### Community 20 - "Model Dark Mode"
Cohesion: 0.50
Nodes (3): bool, app.models – QAbstractTableModel backed by the SQLite db., set_dark_mode()

## Knowledge Gaps
- **20 isolated node(s):** `Orientation`, `ItemFlag`, `bool`, `str`, `Connection` (+15 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `tmp_db()` connect `DB Schema & Fixtures` to `CRUD Test Suite`, `Main Window UI`, `DB Migration System`?**
  _High betweenness centrality (0.469) - this node is a cross-community bridge._
- **Why does `AssetDialog` connect `Asset Dialog Layer` to `Main Window UI`, `Main Window Tests`, `Window Integration Tests`, `Asset Table Model`, `Dialog Form Builder`, `App Init & Tab Bar`?**
  _High betweenness centrality (0.450) - this node is a cross-community bridge._
- **Why does `AssetTableModel` connect `Asset Table Model` to `Asset Dialog Layer`, `Main Window UI`, `Theme & Style System`, `Model Unit Tests`, `Model Data Access`, `Model Row Tests`, `App Init & Tab Bar`, `Model Column Tests`, `Model Dark Mode`, `Model SetData Tests`?**
  _High betweenness centrality (0.380) - this node is a cross-community bridge._
- **Are the 29 inferred relationships involving `AssetDialog` (e.g. with `_NameSizedTabBar` and `int`) actually correct?**
  _`AssetDialog` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Are the 17 inferred relationships involving `AssetTableModel` (e.g. with `_NameSizedTabBar` and `int`) actually correct?**
  _`AssetTableModel` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `main.py – Entry point for IT Asset Inventory.`, `app.dialogs – Add / Edit record dialog (relational schema edition).`, `Schedule hidePopup after the current event cycle (fixes Linux/Wayland).` to the rest of the system?**
  _80 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CRUD Test Suite` be split into smaller, more focused modules?**
  _Cohesion score 0.0547945205479452 - nodes in this community are weakly interconnected._
