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

            <!-- Device select (for antivirus / othersoftware — all device kinds) -->
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

            <!-- Computer-only select (for windowskeys) -->
            <select
              v-else-if="field.type === 'select-computer'"
              v-model="form['_computerSelect']"
              :data-testid="`field-${field.key}`"
              class="w-full text-sm border border-border rounded px-2 py-1 bg-surface text-text-primary focus:outline-none focus:border-accent"
            >
              <option value="">— none —</option>
              <option
                v-for="d in computerDevices"
                :key="d.id"
                :value="`computer:${d.id}`"
              >
                {{ d.name }}
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
import { FIELD_CONFIGS, STATUS_MAP, CATEGORY_LABELS, type FieldConfig } from './editpanel/fieldConfigs'
import { initForm, buildPayload } from './editpanel/payload'
import { callApi } from './editpanel/save'

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
const computerDevices = computed(() => props.devices.filter(d => d.kind === 'computer'))

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const form = reactive<Record<string, any>>(initForm(props.category, props.row))
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
  for (const key of Object.keys(errors)) {
    delete errors[key]
  }

  const fields = currentFields.value
  let valid = true

  for (const field of fields) {
    const value = form[field.key]

    if (field.required) {
      if (value === '' || value === null || value === undefined) {
        errors[field.key] = 'This field is required'
        valid = false
      }
    }

    if (field.type === 'date' && value !== '' && value !== null && value !== undefined) {
      if (!/^\d{4}-\d{2}-\d{2}$/.test(String(value))) {
        errors[field.key] = 'Date must be in YYYY-MM-DD format'
        valid = false
      }
    }
  }

  return valid
}

async function onSave() {
  if (!validateForm()) return

  const payload = buildPayload(props.category, form)
  const id = props.mode === 'edit' ? (props.row!['id'] as number) : 0
  try {
    await callApi(props.category, props.mode === 'edit', id, payload)
    emit('saved')
  } catch (err) {
    console.error('[EditPanel] save failed:', err)
  }
}
</script>
