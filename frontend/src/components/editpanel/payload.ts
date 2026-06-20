import { FIELD_CONFIGS, STATUS_MAP } from './fieldConfigs'

function deriveDeviceSelect(row: Record<string, unknown> | null): string {
  if (!row) return ''
  if (row.computerId != null) return `computer:${row.computerId}`
  if (row.smartphoneId != null) return `smartphone:${row.smartphoneId}`
  if (row.tabletId != null) return `tablet:${row.tabletId}`
  return ''
}

function deriveComputerSelect(row: Record<string, unknown> | null): string {
  if (!row || row.computerId == null) return ''
  return `computer:${row.computerId}`
}

function nullOrStr(val: unknown): string | null {
  if (val === '' || val === null || val === undefined) return null
  return String(val)
}

function nullOrNum(val: unknown): number | null {
  if (val === '' || val === null || val === undefined) return null
  const n = Number(val)
  return isNaN(n) ? null : n
}

export function initForm(category: string, row: Record<string, unknown> | null): Record<string, unknown> {
  const base: Record<string, unknown> = {}
  const fields = FIELD_CONFIGS[category] ?? []
  for (const field of fields) {
    if (field.key === '_deviceSelect') {
      base['_deviceSelect'] = deriveDeviceSelect(row)
    } else if (field.key === '_computerSelect') {
      base['_computerSelect'] = deriveComputerSelect(row)
    } else {
      base[field.key] = row != null ? (row[field.key] ?? '') : ''
    }
  }
  if (!base['status']) {
    const statuses = STATUS_MAP[category]
    base['status'] = statuses?.[0] ?? ''
  }
  return base
}

export function buildPayload(category: string, form: Record<string, unknown>): unknown {
  const cat = category

  if (cat === 'computers' || cat === 'smartphones' || cat === 'tablets') {
    return {
      name: String(form['name'] ?? ''),
      model: nullOrStr(form['model']) ?? '',
      userId: nullOrNum(form['userId']),
      status: String(form['status'] ?? ''),
      purchaseDate: nullOrStr(form['purchaseDate']),
      warrantyExpiry: nullOrStr(form['warrantyExpiry']),
      notes: nullOrStr(form['notes']),
    }
  }

  if (cat === 'windowskeys') {
    const computerSel = String(form['_computerSelect'] ?? '')
    const computerId = computerSel ? Number(computerSel.split(':')[1]) : null
    return {
      licenseKey: String(form['licenseKey'] ?? ''),
      computerId: isNaN(computerId as number) ? null : computerId,
      status: String(form['status'] ?? ''),
      notes: nullOrStr(form['notes']),
    }
  }

  if (cat === 'antivirus' || cat === 'othersoftware') {
    const deviceSel = String(form['_deviceSelect'] ?? '')
    let computerId: number | null = null
    let smartphoneId: number | null = null
    let tabletId: number | null = null
    if (deviceSel) {
      const [kind, idStr] = deviceSel.split(':')
      const id = Number(idStr)
      if (kind === 'computer') computerId = id
      else if (kind === 'smartphone') smartphoneId = id
      else if (kind === 'tablet') tabletId = id
    }
    return {
      name: String(form['name'] ?? ''),
      licenseKey: nullOrStr(form['licenseKey']) ?? '',
      computerId,
      smartphoneId,
      tabletId,
      status: String(form['status'] ?? ''),
      expiryDate: nullOrStr(form['expiryDate']),
      notes: nullOrStr(form['notes']),
    }
  }

  if (cat === 'users') {
    return {
      name: String(form['name'] ?? ''),
      surname: nullOrStr(form['surname']),
      status: String(form['status'] ?? ''),
      notes: nullOrStr(form['notes']),
    }
  }

  return {}
}
