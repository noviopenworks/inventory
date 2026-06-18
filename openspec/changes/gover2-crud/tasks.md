## 1. Models — Input Structs + Dropdown Types

- [x] 1.1 Add `ComputerInput`, `SmartphoneInput`, `TabletInput` structs to `internal/models/models.go`
- [x] 1.2 Add `WindowsKeyInput`, `AntivirusInput`, `OtherSoftwareInput`, `UserInput` structs
- [x] 1.3 Add `DropdownItem` and `DeviceDropdownItem` structs

## 2. Services — Write Operations

- [x] 2.1 Implement `InsertComputer`, `UpdateComputer`, `DeleteComputer` in `internal/services/`
- [x] 2.2 Implement `InsertSmartphone`, `UpdateSmartphone`, `DeleteSmartphone`
- [x] 2.3 Implement `InsertTablet`, `UpdateTablet`, `DeleteTablet`
- [x] 2.4 Implement `InsertWindowsKey`, `UpdateWindowsKey`, `DeleteWindowsKey`
- [x] 2.5 Implement `InsertAntivirus`, `UpdateAntivirus`, `DeleteAntivirus`
- [x] 2.6 Implement `InsertOtherSoftware`, `UpdateOtherSoftware`, `DeleteOtherSoftware`
- [x] 2.7 Implement `InsertUser`, `UpdateUser`, `DeleteUser`
- [x] 2.8 Implement `ListUsersForDropdown` and `ListDevicesForDropdown`
- [x] 2.9 Write table-driven tests: insert→list, update→verify field, delete→list empty (all 7 categories)
- [x] 2.10 Run `go test -cover ./internal/services/...` — must stay ≥ 70%

## 3. Bridge — Write Bridge Methods

- [x] 3.1 Add `AddComputer`, `UpdateComputer`, `DeleteComputer` to `internal/bridge/bridge.go`
- [x] 3.2 Add `AddSmartphone`, `UpdateSmartphone`, `DeleteSmartphone`
- [x] 3.3 Add `AddTablet`, `UpdateTablet`, `DeleteTablet`
- [x] 3.4 Add `AddWindowsKey`, `UpdateWindowsKey`, `DeleteWindowsKey`
- [x] 3.5 Add `AddAntivirus`, `UpdateAntivirus`, `DeleteAntivirus`
- [x] 3.6 Add `AddOtherSoftware`, `UpdateOtherSoftware`, `DeleteOtherSoftware`
- [x] 3.7 Add `AddUser`, `UpdateUser`, `DeleteUser`
- [x] 3.8 Add `ListUsersForDropdown` and `ListDevicesForDropdown` bridge methods
- [x] 3.9 Write bridge tests: round-trip Add→List, Delete removes record, nil-db guard
- [x] 3.10 Run `go test -cover ./internal/bridge/...` — must stay ≥ 50%

## 4. Go Build Verification

- [ ] 4.1 Run `go build ./...` — exit 0
- [ ] 4.2 Run `wails build` — exit 0, binary produced, wailsjs regenerated

## 5. Frontend — api/index.ts

- [x] 5.1 Add `addComputer`, `updateComputer`, `deleteComputer` to `src/lib/api/index.ts`
- [x] 5.2 Add same for Smartphone, Tablet, WindowsKey, Antivirus, OtherSoftware, User (6 × 3 = 18 functions)
- [x] 5.3 Add `listUsersForDropdown` and `listDevicesForDropdown`

## 6. Frontend — EditPanel.vue

- [x] 6.1 Create `src/components/EditPanel.vue` with slide-in animation (300 px, right edge)
- [x] 6.2 Implement `FIELD_CONFIGS` map: fields per category with labels, types, required flags
- [x] 6.3 Implement user dropdown (blank + all users from `listUsersForDropdown`)
- [x] 6.4 Implement device dropdown (mixed, `"Name (Type)"` labels from `listDevicesForDropdown`)
- [x] 6.5 Implement status dropdown (device statuses vs user statuses)
- [x] 6.6 Implement date field with `YYYY-MM-DD` placeholder + save-time format validation
- [x] 6.7 Implement required-field validation (inline error message on Save)
- [x] 6.8 Implement Escape key → close panel
- [x] 6.9 Write `EditPanel.test.ts`: edit mode pre-fills fields; create mode fields empty; Cancel emits cancelled; Escape emits cancelled

## 7. Frontend — ConfirmDialog.vue

- [x] 7.1 Create `src/components/ConfirmDialog.vue` (centered overlay modal)
- [x] 7.2 Write `ConfirmDialog.test.ts`: confirm emits confirmed; cancel emits cancelled

## 8. Frontend — Wire All 7 Views

- [ ] 8.1 Update `ComputersView.vue`: row click → edit panel; Add button → create panel; Delete → confirm → delete
- [ ] 8.2 Update `SmartphonesView.vue` (same pattern)
- [ ] 8.3 Update `TabletsView.vue` (same pattern)
- [ ] 8.4 Update `WindowsKeysView.vue` (same pattern)
- [ ] 8.5 Update `AntivirusView.vue` (same pattern)
- [ ] 8.6 Update `OtherSoftwareView.vue` (same pattern)
- [ ] 8.7 Update `UsersView.vue` (same pattern, user status options)

## 9. Frontend Verification

- [ ] 9.1 Run `pnpm run typecheck` — exit 0
- [ ] 9.2 Run `pnpm run test` — all tests pass

## 10. Integration Smoke Test

- [ ] 10.1 Add a new Computer record — appears in table
- [ ] 10.2 Edit an existing record — changes reflected in table
- [ ] 10.3 Delete a record — row removed
- [ ] 10.4 Cancel edit — no data changed, panel closes
