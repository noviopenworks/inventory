## ADDED Requirements

### Requirement: Startup expiry alerts

On application startup, the app SHALL query for assets whose expiry or warranty date falls within the configured warning window and, if any exist, present them in an alerts modal automatically.

#### Scenario: Assets expiring within the window on launch

- **WHEN** the app starts and at least one asset expires within `ExpiryWarningDays`
- **THEN** the alerts modal opens automatically listing each expiring asset with its category, name, expiry date, days remaining, and severity

#### Scenario: No expiring assets on launch

- **WHEN** the app starts and no asset expires within the warning window
- **THEN** no alerts modal is shown

### Requirement: On-demand alerts

The user SHALL be able to reopen the alerts modal at any time from a topbar action.

#### Scenario: Reopen alerts from the topbar

- **WHEN** the user clicks the ⚠ Alerts action in the topbar
- **THEN** the alerts modal opens showing the current expiring assets
