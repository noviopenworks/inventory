## ADDED Requirements

### Requirement: Search the active view

The user SHALL be able to filter the currently active view's rows in real time by typing in a topbar search box, matching across the view's visible fields.

#### Scenario: Filter rows by query

- **WHEN** the user types text in the search box on a view
- **THEN** the table shows only rows whose visible field values contain the query, updating as the user types

#### Scenario: Clear the query

- **WHEN** the user clears the search box
- **THEN** the table shows all rows of the active view again
