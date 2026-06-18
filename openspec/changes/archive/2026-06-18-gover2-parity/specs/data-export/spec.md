## ADDED Requirements

### Requirement: Export active view to CSV

The user SHALL be able to export the currently active asset view to a CSV file chosen through a native save dialog.

#### Scenario: Export with a chosen path

- **WHEN** the user triggers Export CSV on a category view and selects a destination path
- **THEN** a CSV file is written at that path containing the view's columns as the header row and one row per record in the view's column order

#### Scenario: Cancel the save dialog

- **WHEN** the user triggers Export CSV and cancels the native save dialog
- **THEN** no file is written and no error is reported
