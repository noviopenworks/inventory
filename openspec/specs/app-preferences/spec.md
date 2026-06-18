# app-preferences Specification

## Purpose
TBD - created by archiving change gover2-parity. Update Purpose after archive.
## Requirements
### Requirement: Dark mode toggle and persistence

The user SHALL be able to toggle the UI between light and dark themes, and the chosen theme SHALL persist across application restarts.

#### Scenario: Toggle dark mode

- **WHEN** the user clicks the 🌙 dark-mode action
- **THEN** the entire UI switches between light and dark themes

#### Scenario: Theme persists across restart

- **WHEN** the user has selected dark mode and restarts the app
- **THEN** the app launches in dark mode

### Requirement: Table density

The tables SHALL render according to the configured density (comfortable or compact) with alternating row colors.

#### Scenario: Apply density to tables

- **WHEN** the density preference is set to compact or comfortable
- **THEN** all data tables render with the corresponding row spacing and theme-aware alternating row colors
