<template>
  <div
    class="fixed inset-0 z-40"
    @click.self="emit('cancelled')"
    @keydown.esc="emit('cancelled')"
    tabindex="-1"
  >
    <aside
      class="absolute right-0 top-0 h-full w-80 bg-surface border-l border-border p-4 overflow-y-auto
             transform transition-transform duration-200 ease-in-out translate-x-0"
    >
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-sm font-semibold text-text-primary">
          {{ props.mode === 'create' ? 'Add' : 'Edit' }} {{ categoryLabel }}
        </h2>
        <button @click="emit('cancelled')" class="text-text-secondary hover:text-text-primary">✕</button>
      </div>

      <div class="flex flex-col gap-3">
        <template v-for="field in currentFields" :key="field.key">
          <FormField :label="field.label">
            <!-- Text input -->
            <input
              v-if="field.type === 'text'"
              v-model="form[field.key]"
              :data-testid="`field-${field.key}`"
              type="text"
              class="w-full text-sm border border-border rounded px-2 py-1 bg-surface text-text-primary focus:outline-none focus:border-accent"
            />

            <!-- Textarea -->
            <textarea
              v-else-if="field.type === 'textarea'"
              v-model="form[field.key]"
              :data-testid="`field-${field.key}`"
              rows="3"
              class="w-full text-sm border border-border rounded px-2 py-1 bg-surface text-text-primary focus:outline-none focus:border-accent resize-none"
            />

            <!-- Date input -->
            <input
              v-else-if="field.type === 'date'"
              v-model="form[field.key]"
              :data-testid="`field-${field.key}`"
              type="text"
              placeholder="YYYY-MM-DD"
              class="w-full text-sm border border-border rounded px-2 py-1 bg-surface text-text-primary focus:outline-none focus:border-accent"
            />

            <!-- Status select -->
            <select
              v-else-if="field.type === 'select-status'"
              v-model="form[field.key]"
              :data-testid="`field-${field.key}`"
              class="w-full text-sm border border-border rounded px-2 py-1 bg-surface text-text-primary focus:outline-none focus:border-accent"
            >
              <option v-for="s in currentStatuses" :key="s" :value="s">{{ s }}</option>
            </select>

            <!-- User select -->
            <select
              v-else-if="field.type === 'select-user'"
              v-model="form[field.key]"
              :data-testid="`field-${field.key}`"
              class="w-full text-sm border border-border rounded px-2 py-1 bg-surface text-text-primary focus:outline-none focus:border-accent"
            >
              <option value="">— none —</option>
              <option v-for="u in props.users" :key="u.id" :value="u.id">{{ u.name }}</option>
            </select>

            <!-- Device select (for antivirus / othersoftware) -->
            <select
              v-else-if="field.type === 'select-device'"
              v-model="form['_deviceSelect']"
              :data-testid="`field-${field.key}`"
              class="w-full text-sm border border-border rounded px-2 py-1 bg-surface text-text-primary focus:outline-none focus:border-accent"
            >
              <option value="">— none —</option>
              <option
                v-for="d in props.devices"
                :key="`${d.kind}:${d.id}`"
                :value="`${d.kind}:${d.id}`"
              >
                {{ d.name }} ({{ d.kind }})
              </option>
            </select>

            <p v-if="errors[field.key]" class="text-xs text-red-500 mt-1">{{ errors[field.key] }}</p>
          </FormField>
        </template>
      </div>

      <div class="flex gap-2 mt-4">
        <button
          @click="onSave"
          data-testid="save-btn"
          class="flex-1 bg-accent text-white px-3 py-2 rounded text-sm"
        >
          Save
        </button>
        <button
          @click="emit('cancelled')"
          data-testid="cancel-btn"
          class="px-3 py-2 border border-border rounded text-sm text-text-secondary"
        >
          Cancel
        </button>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, onMounted, onUnmounted } from 'vue'
import FormField from './FormField.vue'
import type { DropdownItem, DeviceDropdownItem } from '@/lib/types'
import {
  addComputer, updateComputer,
  addSmartphone, updateSmartphone,
  addTablet, updateTablet,
  addWindowsKey, updateWindowsKey,
  addAntivirus, updateAntivirus,
  addOtherSoftware, updateOtherSoftware,
  addUser, updateUser,
} from '@/lib/api'

interface FieldConfig {
  key: string
  label: string
  type: 'text' | 'textarea' | 'select-status' | 'select-user' | 'select-device' | 'date'
  required?: boolean
}

const FIELD_CONFIGS: Record<string, FieldConfig[]> = {
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
    { key: 'computerId', label: 'Computer ID', type: 'text' },
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

const STATUS_MAP: Record<string, string[]> = {
  computers: DEVICE_STATUSES,
  smartphones: DEVICE_STATUSES,
  tablets: DEVICE_STATUSES,
  windowskeys: LICENSE_STATUSES,
  antivirus: LICENSE_STATUSES,
  othersoftware: LICENSE_STATUSES,
  users: USER_STATUSES,
}

const CATEGORY_LABELS: Record<string, string> = {
  computers: 'Computer',
  smartphones: 'Smartphone',
  tablets: 'Tablet',
  windowskeys: 'Windows Key',
  antivirus: 'Antivirus',
  othersoftware: 'Other Software',
  users: 'User',
}

const props = defineProps<{
  category: string
  mode: 'create' | 'edit'
  row: Record<string, unknown> | null
  users: DropdownItem[]
  devices: DeviceDropdownItem[]
}>()

const emit = defineEmits<{ saved: []; cancelled: [] }>()

const currentFields = computed<FieldConfig[]>(() => FIELD_CONFIGS[props.category] ?? [])
const currentStatuses = computed<string[]>(() => STATUS_MAP[props.category] ?? [])
const categoryLabel = computed(() => CATEGORY_LABELS[props.category] ?? props.category)

/** Derive initial _deviceSelect value from antivirus/othersoftware row FKs */
function deriveDeviceSelect(row: Record<string, unknown> | null): string {
  if (!row) return ''
  if (row.computerId != null) return `computer:${row.computerId}`
  if (row.smartphoneId != null) return `smartphone:${row.smartphoneId}`
  if (row.tabletId != null) return `tablet:${row.tabletId}`
  return ''
}

function initForm(): Record<string, unknown> {
  const base: Record<string, unknown> = {}
  const fields = FIELD_CONFIGS[props.category] ?? []
  for (const field of fields) {
    if (field.key === '_deviceSelect') {
      base['_deviceSelect'] = deriveDeviceSelect(props.row)
    } else {
      base[field.key] = props.row != null ? (props.row[field.key] ?? '') : ''
    }
  }
  // Ensure status has a default
  if (!base['status']) {
    const statuses = STATUS_MAP[props.category]
    base['status'] = statuses?.[0] ?? ''
  }
  return base
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const form = reactive<Record<string, any>>(initForm())
const errors = reactive<Record<string, string>>({})

// Global Escape key handler (catches key events even when focus is on inputs)
function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('cancelled')
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeyDown)
})

function validateForm(): boolean {
  // Clear previous errors
  for (const key of Object.keys(errors)) {
    delete errors[key]
  }

  const fields = currentFields.value
  let valid = true

  for (const field of fields) {
    const value = form[field.key]

    // Required check
    if (field.required) {
      if (value === '' || value === null || value === undefined) {
        errors[field.key] = 'This field is required'
        valid = false
      }
    }

    // Date format check
    if (field.type === 'date' && value !== '' && value !== null && value !== undefined) {
      if (!/^\d{4}-\d{2}-\d{2}$/.test(String(value))) {
        errors[field.key] = 'Date must be in YYYY-MM-DD format'
        valid = false
      }
    }
  }

  return valid
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

function buildPayload(): unknown {
  const cat = props.category

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
    return {
      licenseKey: String(form['licenseKey'] ?? ''),
      computerId: nullOrNum(form['computerId']),
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

async function callApi(payload: unknown): Promise<void> {
  const cat = props.category
  const isEdit = props.mode === 'edit'
  const id = isEdit ? (props.row!['id'] as number) : 0

  if (cat === 'computers') {
    return isEdit
      ? updateComputer(id, payload as Parameters<typeof updateComputer>[1])
      : addComputer(payload as Parameters<typeof addComputer>[0])
  }
  if (cat === 'smartphones') {
    return isEdit
      ? updateSmartphone(id, payload as Parameters<typeof updateSmartphone>[1])
      : addSmartphone(payload as Parameters<typeof addSmartphone>[0])
  }
  if (cat === 'tablets') {
    return isEdit
      ? updateTablet(id, payload as Parameters<typeof updateTablet>[1])
      : addTablet(payload as Parameters<typeof addTablet>[0])
  }
  if (cat === 'windowskeys') {
    return isEdit
      ? updateWindowsKey(id, payload as Parameters<typeof updateWindowsKey>[1])
      : addWindowsKey(payload as Parameters<typeof addWindowsKey>[0])
  }
  if (cat === 'antivirus') {
    return isEdit
      ? updateAntivirus(id, payload as Parameters<typeof updateAntivirus>[1])
      : addAntivirus(payload as Parameters<typeof addAntivirus>[0])
  }
  if (cat === 'othersoftware') {
    return isEdit
      ? updateOtherSoftware(id, payload as Parameters<typeof updateOtherSoftware>[1])
      : addOtherSoftware(payload as Parameters<typeof addOtherSoftware>[0])
  }
  if (cat === 'users') {
    return isEdit
      ? updateUser(id, payload as Parameters<typeof updateUser>[1])
      : addUser(payload as Parameters<typeof addUser>[0])
  }
}

async function onSave() {
  if (!validateForm()) return

  const payload = buildPayload()
  try {
    await callApi(payload)
    emit('saved')
  } catch (err) {
    console.error('[EditPanel] save failed:', err)
  }
}
</script>
