export interface Computer {
  id: number
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface Smartphone {
  id: number
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface Tablet {
  id: number
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface WindowsKey {
  id: number
  licenseKey: string
  computerId: number | null
  status: string
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface Antivirus {
  id: number
  name: string
  licenseKey: string
  computerId: number | null
  smartphoneId: number | null
  tabletId: number | null
  status: string
  expiryDate: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface OtherSoftware {
  id: number
  name: string
  licenseKey: string
  computerId: number | null
  smartphoneId: number | null
  tabletId: number | null
  status: string
  expiryDate: string | null
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface User {
  id: number
  name: string
  surname: string | null
  status: string
  notes: string | null
  createdAt: string
  updatedAt: string
}

export interface AppConfig {
  dbPath: string
  density: 'comfortable' | 'compact'
  darkMode: boolean
  expiryWarningDays: number
}

export interface Alert {
  category: string
  id: number
  name: string
  expiryDate: string
  daysRemaining: number
  severity: 'expired' | 'expiring'
}

export interface ComputerInput {
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
}

export interface SmartphoneInput {
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
}

export interface TabletInput {
  name: string
  model: string
  userId: number | null
  status: string
  purchaseDate: string | null
  warrantyExpiry: string | null
  notes: string | null
}

export interface WindowsKeyInput {
  licenseKey: string
  computerId: number | null
  status: string
  notes: string | null
}

export interface AntivirusInput {
  name: string
  licenseKey: string
  computerId: number | null
  smartphoneId: number | null
  tabletId: number | null
  status: string
  expiryDate: string | null
  notes: string | null
}

export interface OtherSoftwareInput {
  name: string
  licenseKey: string
  computerId: number | null
  smartphoneId: number | null
  tabletId: number | null
  status: string
  expiryDate: string | null
  notes: string | null
}

export interface UserInput {
  name: string
  surname: string | null
  status: string
  notes: string | null
}

export interface DropdownItem {
  id: number
  name: string
}

export interface DeviceDropdownItem {
  id: number
  name: string
  kind: 'computer' | 'smartphone' | 'tablet'
}
