# database-management Specification

## Purpose
TBD - created by archiving change gover2-parity. Update Purpose after archive.
## Requirements
### Requirement: Create and open databases from the UI

The user SHALL be able to create a new database file or open an existing one through native dialogs reachable from the topbar.

#### Scenario: Create a new database

- **WHEN** the user chooses New Database and selects a destination path
- **THEN** a new database file with the initialized schema is created and becomes the active database

#### Scenario: Open an existing database

- **WHEN** the user chooses Open Database and selects an existing database file
- **THEN** that database becomes the active database and the views reflect its contents
