export interface FieldConfig {
  key: string
  label: string
  type: 'text' | 'textarea' | 'select-status' | 'select-user' | 'select-device' | 'select-computer' | 'date'
  required?: boolean
}

export const FIELD_CONFIGS: Record<string, FieldConfig[]> = {
  computers: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'model', label: 'Model', type: 'text' },
    { key: 'userId', label: 'User', type: 'select-user' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'purchaseDate', label: 'Purchase Date', type: 'date' },
    { key: 'warrantyExpiry', label: 'Warranty Expiry', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  smartphones: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'model', label: 'Model', type: 'text' },
    { key: 'userId', label: 'User', type: 'select-user' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'purchaseDate', label: 'Purchase Date', type: 'date' },
    { key: 'warrantyExpiry', label: 'Warranty Expiry', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  tablets: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'model', label: 'Model', type: 'text' },
    { key: 'userId', label: 'User', type: 'select-user' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'purchaseDate', label: 'Purchase Date', type: 'date' },
    { key: 'warrantyExpiry', label: 'Warranty Expiry', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  windowskeys: [
    { key: 'licenseKey', label: 'License Key', type: 'text', required: true },
    { key: '_computerSelect', label: 'Computer', type: 'select-computer' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  antivirus: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'licenseKey', label: 'License Key', type: 'text' },
    { key: '_deviceSelect', label: 'Device', type: 'select-device' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'expiryDate', label: 'Expiry Date', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  othersoftware: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'licenseKey', label: 'License Key', type: 'text' },
    { key: '_deviceSelect', label: 'Device', type: 'select-device' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'expiryDate', label: 'Expiry Date', type: 'date' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
  users: [
    { key: 'name', label: 'Name', type: 'text', required: true },
    { key: 'surname', label: 'Surname', type: 'text' },
    { key: 'status', label: 'Status', type: 'select-status' },
    { key: 'notes', label: 'Notes', type: 'textarea' },
  ],
}

const DEVICE_STATUSES = ['active', 'inactive', 'repair', 'decommissioned']
const LICENSE_STATUSES = ['active', 'inactive', 'expired']
const USER_STATUSES = ['active', 'inactive']

export const STATUS_MAP: Record<string, string[]> = {
  computers: DEVICE_STATUSES,
  smartphones: DEVICE_STATUSES,
  tablets: DEVICE_STATUSES,
  windowskeys: LICENSE_STATUSES,
  antivirus: LICENSE_STATUSES,
  othersoftware: LICENSE_STATUSES,
  users: USER_STATUSES,
}

export const CATEGORY_LABELS: Record<string, string> = {
  computers: 'Computer',
  smartphones: 'Smartphone',
  tablets: 'Tablet',
  windowskeys: 'Windows Key',
  antivirus: 'Antivirus',
  othersoftware: 'Other Software',
  users: 'User',
}
