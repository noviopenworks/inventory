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
